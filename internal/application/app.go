package application

import (
	"context"
	"fmt"
	"github.com/GermanBogatov/auth-service/internal/config"
	httpHandler "github.com/GermanBogatov/auth-service/internal/handler/http"
	"github.com/GermanBogatov/auth-service/internal/repository/cache"
	"github.com/GermanBogatov/auth-service/internal/repository/postgres"
	"github.com/GermanBogatov/auth-service/internal/service"
	"github.com/GermanBogatov/auth-service/pkg/logging"
	"github.com/GermanBogatov/auth-service/pkg/postgresql"
	"github.com/GermanBogatov/auth-service/pkg/redis"
	"github.com/GermanBogatov/auth-service/pkg/sentry"
	"github.com/GermanBogatov/auth-service/pkg/tracer"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg          *config.Config
	httpServer   *http.Server
	router       *chi.Mux
	cancelTracer func(ctx context.Context)
}

// NewApplication - подключаем различные бд, инициализируем слои и роуты.
func NewApplication(ctx context.Context, cfg *config.Config) (App, error) {

	redisClient, err := redis.NewClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return App{}, errors.Wrap(err, "connection redis client")
	}

	logging.Info("connection postgresql...")
	pgClient, err := postgresql.NewPostgresqlClient(ctx, cfg.Postgres.URL, cfg.Postgres.MaxOpenConn,
		cfg.Postgres.ConnMaxLifetimeMinute, cfg.Postgres.ConnAttempts, cfg.Postgres.ConnTimeout)
	if err != nil {
		return App{}, errors.Wrap(err, "connection postgresql")
	}

	logging.Info("repo initializing...")
	userRepo := postgres.NewUser(pgClient)

	cacheRepo := cache.NewStorage(redisClient, cfg.Redis.UserTTL, cfg.Redis.RefreshTTL)
	logging.Info("cache initializing...")
	jwtService := service.NewJWT(userRepo, cacheRepo, config.JWTSecret, cfg.JwtTTL)

	logging.Info("service initializing...")
	userService := service.NewUser(userRepo)

	logging.Info("handler initializing...")
	appHandler := httpHandler.NewHandler(cfg, userService, jwtService)
	router := appHandler.InitRoutes()

	logging.Info("tracer initializing...")
	cancelTrace, err := tracer.New(&tracer.Config{
		ServiceName:        config.Namespace,
		Environment:        config.ServiceEnv,
		TraceRatioFraction: cfg.Tracer.TraceRatioFraction,
		Endpoint:           cfg.Tracer.Endpoint + ":" + cfg.Tracer.Port,
	})
	if err != nil {
		return App{}, errors.Wrap(err, "connection tracer")
	}

	if cfg.Sentry.DSN != "" && config.ServiceEnv != "" && config.SystemName != "" {
		err = sentry.InitSentry(config.SystemName, cfg.Sentry.DSN, config.ServiceEnv, cfg.Sentry.Debug)
		if err != nil {
			logging.Errorf("sentry.Init: %s", err.Error())
		}
	}

	return App{
		cfg:          cfg,
		router:       router,
		cancelTracer: cancelTrace,
	}, nil
}

// Start - старт сервера и хеслчеков
func (a *App) Start(ctx context.Context) error {
	go a.gracefulShutdown([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM})

	return a.startHttpServer()
}

// startHttpServer - старт http-сервера
func (a *App) startHttpServer() error {
	logging.Infof("http server started on :%v", a.cfg.Http.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", a.cfg.Http.Port))
	if err != nil {
		return errors.New("failed to create listener")
	}

	a.httpServer = &http.Server{
		Handler:      a.router,
		WriteTimeout: time.Second * time.Duration(a.cfg.Http.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(a.cfg.Http.ReadTimeout),
	}

	return a.httpServer.Serve(listener)
}

// gracefulShutdown - плавное завершение сервера
func (a *App) gracefulShutdown(signals []os.Signal) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)

	sig := <-sigc

	logging.Info("--- shutdown application ---")
	time.Sleep(time.Duration(a.cfg.ShutdownTimeoutSec) * time.Second)

	logging.Info("cancel tracer...")
	a.cancelTracer(context.Background())

	logging.Infof("Caught signal %s. Shutting down...", sig)
	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		logging.Errorf("failed to shutdown: %v", err)
	}

}
