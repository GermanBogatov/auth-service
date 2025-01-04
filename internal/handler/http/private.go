package http

import (
	"encoding/json"
	"fmt"
	"github.com/GermanBogatov/auth-service/internal/common/apperror"
	"github.com/GermanBogatov/auth-service/internal/common/helpers"
	"github.com/GermanBogatov/auth-service/internal/common/response"
	"github.com/GermanBogatov/auth-service/internal/config"
	"github.com/GermanBogatov/auth-service/internal/entity"
	"github.com/GermanBogatov/auth-service/internal/handler/http/mapper"
	"github.com/GermanBogatov/auth-service/internal/handler/http/model"
	"github.com/GermanBogatov/auth-service/internal/handler/http/validator"
	"github.com/GermanBogatov/auth-service/pkg/logging"
	"github.com/pkg/errors"
	"net/http"
)

func (h *Handler) PrivateUpdateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userID, err := helpers.GetUuidFromPath(r, config.ParamID)
	if err != nil {
		return apperror.BadRequestError(errors.Wrap(err, "get uuid from path"))
	}

	selfUserID := ctx.Value(config.ParamID).(string)
	role := ctx.Value(config.ParamRole).(string)

	// только админу можно редактировать пользователей любых
	if entity.RoleType(role) != entity.RoleAdmin || entity.RoleType(role) != entity.RoleSuperAdmin {
		return apperror.BadRequestError(fmt.Errorf("user [%s] does not have rights to update user [%s]", selfUserID, userID))
	}

	var userUpdate model.UserUpdatePrivate
	defer func() {
		errClose := r.Body.Close()
		if errClose != nil {
			logging.Error("error close request body")
		}
	}()

	if errDecode := json.NewDecoder(r.Body).Decode(&userUpdate); errDecode != nil {
		return apperror.BadRequestError(errors.Wrap(errDecode, "json decode"))
	}

	err = validator.ValidateUserUpdatePrivate(userUpdate)
	if err != nil {
		return apperror.BadRequestError(errors.Wrap(err, "validate user"))
	}

	user := mapper.MapToEntityUserUpdatePrivate(userUpdate)
	user.ID = userID.String()

	result, err := h.userService.UpdatePrivateUserByID(ctx, user)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return response.RespondSuccess(w, mapper.MapToPrivateUserResponse(http.StatusOK, result))
}
