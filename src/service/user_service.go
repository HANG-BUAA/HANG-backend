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
	var token string     // 生成的 token
	var iUser model.User // 查找的用户对象

	// 检查用户是否存在
	if iUserLoginRequestDTO.Username != "" {
		iUser, err = m.Dao.GetUserByName(iUserLoginRequestDTO.Username)
	} else if iUserLoginRequestDTO.StudentID != "" {
		iUser, err = m.Dao.GetUserByStudentID(iUserLoginRequestDTO.StudentID)
	} else {
		err = errors.New("either 'username' or 'student_id' must be provided")
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
		Role:      iUser.Role,
	}
	res.ID = iUser.ID
	res.CreatedAt = iUser.CreatedAt
	res.UpdatedAt = iUser.UpdatedAt
	res.DeletedAt = iUser.DeletedAt
	return
}

func (m *UserService) Register(iUserRegisterRequestDTO *dto.UserRegisterRequestDTO) (res *dto.UserRegisterResponseDTO, err error) {
	// 检查学号是否已经存在（被注册过）
	if m.Dao.CheckStudentIDExist(iUserRegisterRequestDTO.StudentID) {
		err = errors.New("the student_id is Already registered")
		return
	}

	// 检查验证码是否正确
	studentID := iUserRegisterRequestDTO.StudentID
	containedCode, tmpErr := utils.GetRedis(studentID + "_verification")
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
		Role:      iUser.Role,
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
	return utils.SetRedis(studentID+"_verification", code, time.Duration(viper.GetInt("smtp.expiration"))*time.Minute)
}

func (m *UserService) UpdateAvatar(iUserUpdateAvatarRequestDTO *dto.UserUpdateAvatarRequestDTO) (res *dto.UserUpdateAvatarResponseDTO, err error) {
	id := iUserUpdateAvatarRequestDTO.ID
	url := iUserUpdateAvatarRequestDTO.Url
	err = global.DB.Model(&model.User{}).Where("id = ?", id).Update("avatar", url).Error
	if err != nil {
		return
	}
	res = &dto.UserUpdateAvatarResponseDTO{
		Url: url,
	}
	return
}
