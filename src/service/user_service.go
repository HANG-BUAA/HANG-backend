package service

import (
	"HANG-backend/src/dao"
	"HANG-backend/src/model"
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"errors"
)

var userService *UserService

type UserService struct {
	BaseService
	Dao *dao.UserDao
}

func NewUserService() *UserService {
	if userService == nil {
		userService = &UserService{
			Dao: dao.NewUserDao(),
		}
	}
	return userService
}

func (m *UserService) Login(iUserDTO dto.UserLoginDTO) (model.User, string, error) {
	var errResult error
	var token string

	iUser, err := m.Dao.GetUserByName(iUserDTO.Username)

	// 用户名或密码不正确
	if err != nil || !utils.CompareHashAndPassword(iUser.Password, iUserDTO.Password) {
		errResult = errors.New("Invalid UserName Or Password")
	} else {
		// 登录成功，生成 token
		token, err = utils.GenerateToken(iUser.ID, iUser.UserName)
		if err != nil {
			errResult = errors.New("Generate Token Error")
		}
	}
	return iUser, token, errResult
}

func (m *UserService) AddUser(iUserAddDTO *dto.UserRegisterDTO) error {
	if m.Dao.CheckUserExit(iUserAddDTO.Username) {
		return errors.New("Username Exists")
	}
	return m.Dao.AddUser(iUserAddDTO)
}
