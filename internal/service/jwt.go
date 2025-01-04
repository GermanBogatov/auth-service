package service

import (
	"context"
	"github.com/GermanBogatov/auth-service/internal/common/apperror"
	"github.com/GermanBogatov/auth-service/internal/entity"
	"github.com/GermanBogatov/auth-service/internal/repository/cache"
	"github.com/GermanBogatov/auth-service/internal/repository/postgres"
	"github.com/GermanBogatov/auth-service/pkg/logging"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

var _ IJWT = &JWT{}

type JWT struct {
	userRepo postgres.IUser
	cache    cache.ICache
	secret   string
	jwtTTL   time.Duration
}

func NewJWT(userRepo postgres.IUser, cache cache.ICache, secret string, jwtTTL int) IJWT {
	return &JWT{
		userRepo: userRepo,
		cache:    cache,
		secret:   secret,
		jwtTTL:   time.Duration(jwtTTL) * time.Second,
	}
}

type IJWT interface {
	UpdateRefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	GenerateAccessAndRefreshTokens(user entity.User) (string, string, error)
}

// UpdateRefreshToken - обновление рефреш токена
func (j *JWT) UpdateRefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	userID, err := j.cache.Get(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", "", apperror.ErrRefreshTokenNotFound
		}
		return "", "", errors.Wrap(err, "cache.Get")
	}

	var (
		user    entity.User
		errUser error
	)
	user, errUser = j.cache.GetUser(ctx, userID)
	if errUser != nil {
		if errors.Is(errUser, redis.Nil) {
			user, errUser = j.userRepo.GetUserByID(ctx, userID)
			if errUser != nil {
				return "", "", errors.Wrap(errUser, "userRepo.GetUserByID")
			}
		} else {
			return "", "", errors.Wrap(errUser, "cache.GetUser")
		}
	}

	go func() {
		ctxDel, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
		defer cancel()
		errDelete := j.cache.Delete(ctxDel, refreshToken)
		if errDelete != nil {
			logging.Errorf("error deleting refresh token [%s]: %v", refreshToken, errDelete)
		}
	}()

	newAccessToken, newRefreshToken, err := j.GenerateAccessAndRefreshTokens(user)
	if err != nil {
		return "", "", errors.Wrap(err, "GenerateAccessToken")
	}

	return newAccessToken, newRefreshToken, nil
}

// GenerateAccessAndRefreshTokens - генерация токенов
func (j *JWT) GenerateAccessAndRefreshTokens(user entity.User) (string, string, error) {
	key := []byte(j.secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        user.ID,
			Audience:  jwt.ClaimStrings{"users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.jwtTTL)),
		},
		Email: user.Email,
		Role:  string(user.Role),
	})

	accessToken, err := token.SignedString(key)
	if err != nil {
		return "", "", err
	}

	refreshToken := uuid.New().String()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
		defer cancel()

		errSet := j.cache.SetRefreshToken(ctx, refreshToken, user.ID)
		if errSet != nil {
			logging.Errorf("error set refresh token [%s]: %v", refreshToken, errSet)
		}

		errSet = j.cache.SetUser(ctx, user.ID, user)
		if errSet != nil {
			logging.Errorf("error set user [%s]: %v", user.ID, errSet)
		}
	}()

	return accessToken, refreshToken, err
}
