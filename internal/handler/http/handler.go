package http

import (
	_ "github.com/GermanBogatov/auth-service/docs"
	"github.com/GermanBogatov/auth-service/internal/config"
	"github.com/GermanBogatov/auth-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	metricPath    = "/metrics"
	healthPath    = "/health"
	publicV1      = "/public/v1"
	privateV1     = "/private/v1"
	integrationV1 = "/integration/v1"
	authV1        = "/public/v1/auth"

	livePath       = "/live"
	readinessPath  = "/readiness"
	swaggerPattern = "/swagger-ui/*"
)

type Handler struct {
	userService service.IUser
	jwtService  service.IJWT
	cfg         *config.Config
}

func NewHandler(cfg *config.Config, userService service.IUser, jwtService service.IJWT) *Handler {
	return &Handler{
		userService: userService,
		jwtService:  jwtService,
		cfg:         cfg,
	}
}

// InitRoutes - инициализация роутера приложения
func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Handle(metricPath, promhttp.Handler())

	r.Route(healthPath, func(r chi.Router) {
		r.Get(livePath, live)
		r.Get(readinessPath, readiness)
	})
	r.Get(swaggerPattern, httpSwagger.Handler())

	r.Route(authV1, func(r chi.Router) {
		r.Post("/sign-up", appMiddleware(h.SignUp))
		r.Post("/sign-in", appMiddleware(h.SignIn))
		r.Get("/refresh/{id}", appMiddleware(h.UpdateRefreshToken))
	})

	r.Route(publicV1, func(r chi.Router) {
		r.Get("/users", appMiddleware(h.GetUsers))
		r.Get("/users/{id}", appMiddleware(h.GetUserByID))
		r.Delete("/users/{id}", appMiddleware(h.DeleteUserByID))
		r.Patch("/users/{id}", appMiddleware(h.UpdateUserByID))
	})

	r.Route(privateV1, func(r chi.Router) {
		r.Patch("/users/{id}", appMiddleware(h.PrivateUpdateUser))
	})

	return r
}
