package api

import (
	"HANG-backend/src/global"
	"HANG-backend/src/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"reflect"
)

type BaseApi struct {
	Ctx    *gin.Context
	Errors error
	Logger *zap.SugaredLogger
}

func NewBaseApi() BaseApi {
	return BaseApi{
		Logger: global.Logger,
	}
}

type BuildRequestOption struct {
	Ctx               *gin.Context
	DTO               any
	BindParamsFromUri bool
}

func (m *BaseApi) BuildRequest(option BuildRequestOption) *BaseApi {
	var errResult error

	// 绑定请求上下文
	m.Ctx = option.Ctx

	// 绑定请求数据
	if option.DTO != nil {
		// 每次只处理一种类型的绑定，如果两处都有，需要两次调用绑定
		if option.BindParamsFromUri {
			errResult = m.Ctx.ShouldBindUri(option.DTO)
		} else {
			errResult = m.Ctx.ShouldBind(option.DTO)
		}

		if errResult != nil {
			errResult = m.ParseValidateErrors(errResult, option.DTO)
			m.AdError(errResult)
			m.Fail(ResponseJson{
				Msg: m.GetError().Error(),
			})
		}
	}
	return m
}

func (m *BaseApi) AdError(errNew error) {
	m.Errors = utils.AppendError(m.Errors, errNew)
}

func (m *BaseApi) GetError() error {
	return m.Errors
}

func (m *BaseApi) ParseValidateErrors(errs error, target any) error {
	var errResult error

	var errValidation validator.ValidationErrors
	ok := errors.As(errs, &errValidation)
	if !ok {
		return errs
	}

	// 通过反射获取指针指向元素的类型对象
	fields := reflect.TypeOf(target).Elem()
	for _, fieldErr := range errValidation {
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

func (m *BaseApi) Fail(resp ResponseJson) {
	Fail(m.Ctx, resp)
}

func (m *BaseApi) OK(resp ResponseJson) {
	OK(m.Ctx, resp)
}

func (m *BaseApi) ServerFail(resp ResponseJson) {
	ServerFail(m.Ctx, resp)
}
