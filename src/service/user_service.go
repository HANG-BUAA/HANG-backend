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

func (m *UserService) Login(iUserLoginRequestDTO *dto.UserLoginRequestDTO) (res *dto.UserLoginResponseDTO, err error) {
	//var errResult error
	//var token string
	//
	//iUser, err := m.Dao.GetUserByName(iUserLoginRequestDTO.Username)
	//
	//// 用户名或密码不正确
	//if err != nil || !utils.CompareHashAndPassword(iUser.Password, iUserLoginRequestDTO.Password) {
	//	errResult = errors.New("invalid UserName Or Password")
	//} else {
	//	// 登录成功，生成 token
	//	token, err = utils.GenerateToken(iUser.ID, iUser.UserName)
	//	if err != nil {
	//		errResult = errors.New("generate Token Error")
	//	}
	//	iUserLoginRequestDTO.Token = token
	//	iUserLoginRequestDTO.StudentID = iUser.StudentID
	//	iUserLoginRequestDTO.Password = ""
	//	iUserLoginRequestDTO.ID = iUser.ID
	//}
	//return errResult
	var token string     // 生成的 token
	var iUser model.User // 查找的用户对象

	// 检查用户是否存在
	if iUserLoginRequestDTO.Username != "" {
		iUser, err = m.Dao.GetUserByName(iUserLoginRequestDTO.Username)
	} else if iUserLoginRequestDTO.StudentID != "" {
		iUser, err = m.Dao.GetUserByStudentID(iUserLoginRequestDTO.StudentID)
	} else {
		err = errors.New("Either 'username' or 'student_id' must be provided")
	}
	if err != nil {
		return
	}

	// 检查密码匹配性
	if !utils.CompareHashAndPassword(iUser.Password, iUserLoginRequestDTO.Password) {
		err = errors.New("invalid password")
		return
	}

	token, err = utils.GenerateToken(iUser.ID, iUser.UserName)
	if err != nil {
		err = errors.New("failed to generate token")
	}

	res = &dto.UserLoginResponseDTO{
		Token:     token,
		StudentID: iUser.StudentID,
		Username:  iUser.UserName,
	}
	res.ID = iUser.ID
	res.CreatedAt = iUser.CreatedAt
	res.UpdatedAt = iUser.UpdatedAt
	res.DeletedAt = iUser.DeletedAt
	return
}

func (m *UserService) Register(iUserRegisterRequestDTO *dto.UserRegisterRequestDTO) (res *dto.UserRegisterResponseDTO, err error) {
	//if m.Dao.CheckUserExit(iUserRegisterRequestDTO.Username) {
	//	return errors.New("username Exists")
	//}
	//
	//// 检查验证码是否正确
	//studentID := iUserRegisterRequestDTO.Username
	//containedCode, err := global.RedisClient.Get(studentID + "_verification")
	//if err != nil {
	//	return err
	//}
	//if containedCode != iUserRegisterRequestDTO.VerificationCode {
	//	return errors.New("verification Code Error")
	//}
	//
	//return m.Dao.AddUser(iUserRegisterRequestDTO)
	// 检查学号是否已经存在（被注册过）
	if m.Dao.CheckStudentIDExist(iUserRegisterRequestDTO.StudentID) {
		err = errors.New("the student_id is Already registered")
		return
	}

	// 检查验证码是否正确
	studentID := iUserRegisterRequestDTO.StudentID
	containedCode, tmpErr := global.RedisClient.Get(studentID + "_verification")
	if tmpErr != nil {
		err = errors.New("verification code expired")
		return
	}
	if containedCode != iUserRegisterRequestDTO.VerificationCode {
		err = errors.New("verification code expired")
		return
	}
	iUser, err := m.Dao.AddUser(studentID, iUserRegisterRequestDTO.Password)
	if err != nil {
		return
	}

	res = &dto.UserRegisterResponseDTO{
		ID:        iUser.ID,
		Username:  iUser.UserName,
		StudentID: iUser.StudentID,
	}
	return
}

func (m *UserService) SendEmail(iUserSendEmailRequestDTO *dto.UserSendEmailRequestDTO) error {
	studentID := iUserSendEmailRequestDTO.StudentID
	code, err := utils.SendEmail(studentID)
	if err != nil {
		return err
	}

	// 把 code 存到 redis里
	return global.RedisClient.Set(studentID+"_verification", code, time.Duration(viper.GetInt("smtp.expiration"))*time.Minute)
}
