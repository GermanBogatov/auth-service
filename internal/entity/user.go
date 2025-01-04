package entity

import (
	"github.com/google/uuid"
	"time"
)

type RoleType string

const (
	RoleUser       RoleType = "user"
	RoleAdmin      RoleType = "admin"
	RoleSuperAdmin RoleType = "super-admin"
)

// User - модель пользователя
type User struct {
	CreatedDate time.Time
	UpdatedDate *time.Time
	ID          string
	Name        string
	Surname     string
	Email       string
	Password    string
	Role        RoleType
	JWT         JWT
}

// UserUpdateBase - базовая модель пользователя для редактирования
type UserUpdateBase struct {
	Name    *string
	Surname *string
	Email   *string
	ID      string
}

// UserUpdate - модель обновления пользователя
type UserUpdate struct {
	Password *string
	UserUpdateBase
}

// UserUpdatePrivate - модель приватного обновления пользователя
type UserUpdatePrivate struct {
	Role *RoleType
	UserUpdateBase
}

// Filter - модель фильтра
type Filter struct {
	Role   *RoleType
	Sort   string
	Order  string
	Limit  int
	Offset int
}

func (u *User) GenerateID() {
	u.ID = uuid.New().String()
}

func (u *User) GenerateCreatedDate() {
	u.CreatedDate = time.Now().UTC()
}

func (u *User) AddRoleUser() {
	u.Role = RoleUser
}

func (u *User) SetPasswordHash(hash string) {
	u.Password = hash
}

func (u *User) SetJWT(token, refreshToken string) {
	u.JWT = JWT{
		Token:        token,
		RefreshToken: refreshToken,
	}
}

// JWT - модель токена с рефрешом
type JWT struct {
	Token        string
	RefreshToken string
}
