package http

import (
	"encoding/json"
	"github.com/GermanBogatov/auth-service/internal/common/apperror"
	"github.com/GermanBogatov/auth-service/internal/common/helpers"
	"github.com/GermanBogatov/auth-service/internal/common/response"
	"github.com/GermanBogatov/auth-service/internal/config"
	"github.com/GermanBogatov/auth-service/internal/handler/http/mapper"
	"github.com/GermanBogatov/auth-service/internal/handler/http/model"
	"github.com/GermanBogatov/auth-service/internal/handler/http/validator"
	"github.com/GermanBogatov/auth-service/pkg/logging"
	"github.com/pkg/errors"
	"net/http"
)

// SignUp - хэндлер регистрации
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var createUser model.SignUpRequest
	defer func() {
		err := r.Body.Close()
		if err != nil {
			logging.Error("error close request body")
		}
	}()

	if err := json.NewDecoder(r.Body).Decode(&createUser); err != nil {
		return apperror.BadRequestError(errors.Wrap(err, "json decode"))
	}

	err := validator.ValidateSignUpUser(createUser)
	if err != nil {
		return apperror.BadRequestError(errors.Wrap(err, "validate create user"))
	}

	user := mapper.MapToEntityUser(createUser)
	user.GenerateID()
	user.SetPasswordHash(helpers.GeneratePasswordHash(createUser.Password))
	user.GenerateCreatedDate()
	// todo когда админ появится условия предусмотреть
	user.AddRoleUser()

	token, refreshToken, err := h.jwtService.GenerateAccessAndRefreshTokens(ctx, user)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	user.SetJWT(token, refreshToken)

	err = h.userService.CreateUser(ctx, user)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return response.RespondSuccessCreate(w, mapper.MapToUserWithJWTResponse(http.StatusCreated, user))

}

// SignIn - хэндлер авторизации
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var signInUser model.SignInRequest
	defer func() {
		err := r.Body.Close()
		if err != nil {
			logging.Error("error close request body")
		}
	}()

	if err := json.NewDecoder(r.Body).Decode(&signInUser); err != nil {
		return apperror.BadRequestError(errors.Wrap(err, "json decode"))
	}

	err := validator.ValidateSignInUser(signInUser)
	if err != nil {
		return apperror.BadRequestError(errors.Wrap(err, "validate create user"))
	}

	passwordHash := helpers.GeneratePasswordHash(signInUser.Password)
	user, err := h.userService.GetUserByEmailAndPassword(ctx, signInUser.Email, passwordHash)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	token, refreshToken, err := h.jwtService.GenerateAccessAndRefreshTokens(ctx, user)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	user.SetJWT(token, refreshToken)

	return response.RespondSuccess(w, mapper.MapToUserWithJWTResponse(http.StatusOK, user))
}

// UpdateRefreshToken - хэндлер обновления jwt-токена
func (h *Handler) UpdateRefreshToken(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	refreshToken, err := helpers.GetStringFromPath(r, config.ParamID)
	if err != nil {
		return apperror.BadRequestError(errors.Wrap(err, "get refresh token from header"))
	}

	token, newRefreshToken, errToken := h.jwtService.UpdateRefreshToken(ctx, refreshToken)
	if errToken != nil {
		return apperror.InternalServerError(errToken)
	}

	return response.RespondSuccess(w, mapper.MapToJWTResponse(http.StatusOK, token, newRefreshToken))
}
