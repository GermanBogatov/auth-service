package cache

import (
	"context"
	"encoding/json"
	"github.com/GermanBogatov/auth-service/internal/common/apperror"
	"github.com/GermanBogatov/auth-service/internal/common/metrics"
	"github.com/GermanBogatov/auth-service/internal/entity"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type ICache interface {
	Get(ctx context.Context, key string) (string, error)
	GetUser(ctx context.Context, key string) (entity.User, error)
	Delete(ctx context.Context, key string) error
	SetUser(ctx context.Context, key string, user entity.User) error
	SetRefreshToken(ctx context.Context, key, userID string) error
}

var _ ICache = &Cache{}

type Cache struct {
	client     *redis.Client
	userTTL    time.Duration
	refreshTTL time.Duration
}

func NewStorage(client *redis.Client, userTTL, refreshTTL int) ICache {
	return &Cache{
		client:     client,
		userTTL:    time.Duration(userTTL) * time.Second,
		refreshTTL: time.Duration(refreshTTL) * time.Second,
	}
}

// Get - получение данных из кеша по ключу.
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	defer metrics.ObserveRequestDurationPerMethodDB(metrics.Cache, metrics.GetCache)()

	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		metrics.IncRequestTotalDB(metrics.GetCache, metrics.FailStatus)
		if errors.Is(err, redis.Nil) {
			return "", apperror.ErrRedisNil
		}
		return "", err
	}

	metrics.IncRequestTotalDB(metrics.GetCache, metrics.OkStatus)
	return value, nil
}

// Delete - удаление записи из кэша по ключу
func (c *Cache) Delete(ctx context.Context, key string) error {
	defer metrics.ObserveRequestDurationPerMethodDB(metrics.Cache, metrics.DeleteCache)()

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		metrics.IncRequestTotalDB(metrics.DeleteCache, metrics.FailStatus)
		return err
	}

	metrics.IncRequestTotalDB(metrics.DeleteCache, metrics.OkStatus)
	return nil
}

// GetUser - получение пользователя из кэша
func (c *Cache) GetUser(ctx context.Context, key string) (entity.User, error) {
	defer metrics.ObserveRequestDurationPerMethodDB(metrics.Cache, metrics.GetUserCache)()
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		metrics.IncRequestTotalDB(metrics.GetUserCache, metrics.OkStatus)
		if errors.Is(err, redis.Nil) {
			return entity.User{}, apperror.ErrRedisNil
		}
		return entity.User{}, err
	}

	var user entity.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		metrics.IncRequestTotalDB(metrics.GetUserCache, metrics.OkStatus)
		return entity.User{}, err
	}

	metrics.IncRequestTotalDB(metrics.GetUserCache, metrics.FailStatus)
	return user, nil
}

// SetUser - добавление пользователя в кеш.
func (c *Cache) SetUser(ctx context.Context, key string, user entity.User) error {
	defer metrics.ObserveRequestDurationPerMethodDB(metrics.Cache, metrics.SetUserCache)()
	data, errJson := json.Marshal(user)
	if errJson != nil {
		return errJson
	}

	_, err := c.client.Set(ctx, key, string(data), c.userTTL).Result()
	if err != nil {
		metrics.IncRequestTotalDB(metrics.SetUserCache, metrics.FailStatus)
		return err
	}

	metrics.IncRequestTotalDB(metrics.SetUserCache, metrics.OkStatus)
	return nil
}

// SetRefreshToken - добавление рефреш токена в кэш.
func (c *Cache) SetRefreshToken(ctx context.Context, key, userID string) error {
	defer metrics.ObserveRequestDurationPerMethodDB(metrics.Cache, metrics.SetRefreshTokenCache)()
	_, err := c.client.Set(ctx, key, userID, c.refreshTTL).Result()
	if err != nil {
		metrics.IncRequestTotalDB(metrics.SetRefreshTokenCache, metrics.FailStatus)
		return err
	}

	metrics.IncRequestTotalDB(metrics.SetRefreshTokenCache, metrics.OkStatus)
	return nil
}
