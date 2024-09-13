package api

import (
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type UserApi struct{}

func NewUserApi() UserApi {
	return UserApi{}
}

func (m UserApi) Login(ctx *gin.Context) {
	var iUserLoginDTO dto.UserLoginDTO
	errs := ctx.ShouldBind(&iUserLoginDTO)
	if errs != nil {
		Fail(ctx, ResponseJson{
			Msg: parseValidateErrors(errs.(validator.ValidationErrors), &iUserLoginDTO).Error(),
		})
		return
	}
	OK(ctx, ResponseJson{
		Data: iUserLoginDTO,
	})
}

func parseValidateErrors(errs validator.ValidationErrors, target any) error {
	var errResult error

	// 通过反射获取指针指向元素的类型对象
	fields := reflect.TypeOf(target).Elem()
	for _, fieldErr := range errs {
		field, _ := fields.FieldByName(fieldErr.Field())
		errMessageTag := fmt.Sprintf("%s_err", fieldErr.Tag())
		errMessage := field.Tag.Get(errMessageTag)
		if errMessage == "" {
			errMessage = field.Tag.Get("message")
		}

		if errMessage == "" {
			errMessage = fmt.Sprintf("%s: %s Error", fieldErr.Field(), fieldErr.Tag())
		}

		errResult = utils.AppendError(errResult, errors.New(errMessage))
	}

	return errResult
}
