package service

import (
	"context"
	"github.com/GermanBogatov/auth-service/internal/entity"
	"github.com/GermanBogatov/auth-service/internal/repository/postgres"
	"github.com/pkg/errors"
)

var _ IUser = &User{}

type IUser interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByID(ctx context.Context, id string) (entity.User, error)
	GetUsers(ctx context.Context, filter entity.Filter) ([]entity.User, error)
	GetUserByEmailAndPassword(ctx context.Context, email, password string) (entity.User, error)
	DeleteUserByID(ctx context.Context, id string) error
	UpdateUserByID(ctx context.Context, userUpdate entity.UserUpdate) (entity.User, error)

	UpdatePrivateUserByID(ctx context.Context, userUpdate entity.UserUpdatePrivate) (entity.User, error)
}

type User struct {
	userRepo postgres.IUser
}

func NewUser(client postgres.IUser) IUser {
	return &User{
		userRepo: client,
	}
}

// CreateUser - создание пользователя
func (u *User) CreateUser(ctx context.Context, user entity.User) error {
	err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, "userRepo.CreateUser")
	}

	return nil
}

// GetUserByID - получение пользователя по идентификатору
func (u *User) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, errors.Wrap(err, "userRepo.GetUserByID")
	}

	return user, nil
}

// DeleteUserByID - удаление пользователя по идентификатору
func (u *User) DeleteUserByID(ctx context.Context, id string) error {
	err := u.userRepo.DeleteUserByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "userRepo.DeleteUserByID")
	}

	return nil
}

// GetUserByEmailAndPassword - получение пользователя по майлу и паролю
func (u *User) GetUserByEmailAndPassword(ctx context.Context, email, password string) (entity.User, error) {
	user, err := u.userRepo.GetUserByEmailAndPassword(ctx, email, password)
	if err != nil {
		return entity.User{}, errors.Wrap(err, "userRepo.GetUserByEmailAndPassword")
	}

	return user, nil
}

// UpdateUserByID - обновление пользователя
func (u *User) UpdateUserByID(ctx context.Context, userUpdate entity.UserUpdate) (entity.User, error) {
	user, err := u.userRepo.UpdateUserByID(ctx, userUpdate)
	if err != nil {
		return entity.User{}, errors.Wrap(err, "userRepo.UpdateUserByID")
	}

	return user, nil
}

// GetUsers - получение списка пользователей
func (u *User) GetUsers(ctx context.Context, filter entity.Filter) ([]entity.User, error) {
	users, err := u.userRepo.GetUsers(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.GetUsers")
	}

	return users, nil
}

// UpdatePrivateUserByID - приватное обновление пользователя
func (u *User) UpdatePrivateUserByID(ctx context.Context, userUpdate entity.UserUpdatePrivate) (entity.User, error) {
	user, err := u.userRepo.UpdatePrivateUserByID(ctx, userUpdate)
	if err != nil {
		return entity.User{}, errors.Wrap(err, "userRepo.UpdatePrivateUserByID")
	}

	return user, nil
}
