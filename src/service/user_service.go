package service

import (
	"HANG-backend/src/dao"
	"HANG-backend/src/global"
	"HANG-backend/src/model"
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"errors"
	"github.com/spf13/viper"
	"time"
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
		errResult = errors.New("invalid UserName Or Password")
	} else {
		// 登录成功，生成 token
		token, err = utils.GenerateToken(iUser.ID, iUser.UserName)
		if err != nil {
			errResult = errors.New("generate Token Error")
		}
	}
	return iUser, token, errResult
}

func (m *UserService) Register(iUserRegisterDTO *dto.UserRegisterDTO) error {
	if m.Dao.CheckUserExit(iUserRegisterDTO.Username) {
		return errors.New("username Exists")
	}

	// 检查验证码是否正确
	studentID := iUserRegisterDTO.Username
	containedCode, err := global.RedisClient.Get(studentID + "_verification")
	if err != nil {
		return err
	}
	if containedCode != iUserRegisterDTO.VerificationCode {
		return errors.New("verification Code Error")
	}

	return m.Dao.AddUser(iUserRegisterDTO)
}

func (m *UserService) SendEmail(iUserSendEmailDTO dto.UserSendEmailDTO) error {
	studentID := iUserSendEmailDTO.StudentID
	code, err := utils.SendEmail(studentID)
	if err != nil {
		return err
	}
	// 把 code 存到 redis里
	return global.RedisClient.Set(studentID+"_verification", code, time.Duration(viper.GetInt("smtp.expiration"))*time.Minute)
}
