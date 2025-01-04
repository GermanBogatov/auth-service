package mapper

import (
	"github.com/GermanBogatov/auth-service/internal/common/response"
	"github.com/GermanBogatov/auth-service/internal/config"
	"github.com/GermanBogatov/auth-service/internal/entity"
	"github.com/GermanBogatov/auth-service/internal/handler/http/model"
)

// MapToEntityUser - маппинг в модель пользователя
func MapToEntityUser(user model.SignUpRequest) entity.User {
	return entity.User{
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
	}
}

// MapToEntityUserUpdate - маппинг в модель редактирования пользователя
func MapToEntityUserUpdate(user model.UserUpdate) entity.UserUpdate {
	u := entity.UserUpdate{}
	u.Name = user.Name
	u.Surname = user.Surname
	u.Email = user.Email
	return u
}

// MapToEntityUserUpdatePrivate - маппинг в модель приватного редактирования пользователя
func MapToEntityUserUpdatePrivate(user model.UserUpdatePrivate) entity.UserUpdatePrivate {

	u := entity.UserUpdatePrivate{}
	u.Name = user.Name
	u.Surname = user.Surname
	u.Email = user.Email

	if user.Role != nil {
		role := entity.RoleType(*user.Role)
		u.Role = &role
		return u
	}

	return u
}

// MapToEntityFilter - маппинг в модель фильтра
func MapToEntityFilter(limit, offset int, sort, order string, role *string) entity.Filter {
	filter := entity.Filter{
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
		Order:  mapOrderType(order),
	}
	if role == nil {
		return filter
	}

	r := entity.RoleType(*role)
	filter.Role = &r
	return filter
}

// mapOrderType - маппинг поля сортировки
func mapOrderType(order string) string {
	switch order {
	case config.OrderCreatedDate:
		return config.OrderCreatedDateDB
	default:
		return order
	}
}

// MapToUserWithJWTResponse - маппинг пользователя с jwt в ответ
func MapToUserWithJWTResponse(code int, user entity.User) response.ViewResponse {
	return response.ViewResponse{
		Code: code,
		Result: model.SignUpResponse{
			UserResponse: mapUserToResponse(user),
			JWT: model.JWT{
				Token:        user.JWT.Token,
				RefreshToken: user.JWT.RefreshToken,
			},
		},
	}
}

// mapPrivateUserToResponse - маппинг пользователя в модель ответ
func mapUserToResponse(user entity.User) model.UserResponse {
	var (
		updatedDate *string
	)

	if user.UpdatedDate != nil {
		updateTime := user.UpdatedDate.Format(config.IsoTimeLayout)
		updatedDate = &updateTime
	}

	return model.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Surname:     user.Surname,
		Email:       user.Email,
		CreatedDate: user.CreatedDate.Format(config.IsoTimeLayout),
		UpdatedDate: updatedDate,
		Role:        string(user.Role),
	}
}

// MapToUserResponse - маппинг пользователя в модель ответ
func MapToUserResponse(code int, user entity.User) response.ViewResponse {
	return response.ViewResponse{
		Code:   code,
		Result: mapUserToResponse(user),
	}
}

// MapToPrivateUserResponse - маппинг приватного пользователя в модель ответ
func MapToPrivateUserResponse(code int, user entity.User) response.ViewResponse {
	return response.ViewResponse{
		Code:   code,
		Result: mapUserToResponse(user),
	}
}

// MapToUsersResponse - маппинг пользователей в модель ответ
func MapToUsersResponse(code int, users []entity.User) response.ViewResponse {
	result := make([]model.UserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, mapUserToResponse(user))
	}
	return response.ViewResponse{
		Code:   code,
		Result: result,
	}
}

// MapToJWTResponse - маппинг токена в модель с jwt
func MapToJWTResponse(code int, token, refreshToken string) response.ViewResponse {
	return response.ViewResponse{
		Code: code,
		Result: model.JWT{
			Token:        token,
			RefreshToken: refreshToken,
		},
	}
}
