package service

import (
	"HANG-backend/src/dao"
	"HANG-backend/src/model"
	"HANG-backend/src/permission"
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

func (m *UserService) Login(userLoginRequestDTO *dto.UserLoginRequestDTO) (res *dto.UserLoginResponseDTO, err error) {
	var token string    // 生成的 token
	var user model.User // 查找的用户对象

	// 检查用户是否存在
	if userLoginRequestDTO.Username != "" {
		user, err = m.Dao.GetUserByName(userLoginRequestDTO.Username)
	} else if userLoginRequestDTO.StudentID != "" {
		user, err = m.Dao.GetUserByStudentID(userLoginRequestDTO.StudentID)
	} else {
		err = errors.New("either 'username' or 'student_id' must be provided")
	}
	if err != nil {
		return
	}

	// 检查密码匹配性
	if !utils.CompareHashAndPassword(user.Password, userLoginRequestDTO.Password) {
		err = errors.New("invalid password")
		return
	}

	token, err = utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		err = errors.New("failed to generate token")
	}

	res = &dto.UserLoginResponseDTO{
		ID:        user.ID,
		Token:     token,
		StudentID: user.StudentID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
	return
}

func (m *UserService) Register(userRegisterRequestDTO *dto.UserRegisterRequestDTO) (res *dto.UserRegisterResponseDTO, err error) {
	// 检查学号是否已经存在（被注册过）
	if m.Dao.CheckStudentIDExist(userRegisterRequestDTO.StudentID) {
		err = errors.New("the student_id is Already registered")
		return
	}

	// 检查验证码是否正确
	studentID := userRegisterRequestDTO.StudentID
	//containedCode, tmpErr := utils.GetRedis(studentID + "_verification")
	//if tmpErr != nil {
	//	err = errors.New("verification code expired")
	//	return
	//}
	//if containedCode != userRegisterRequestDTO.VerificationCode {
	//	err = errors.New("verification code expired")
	//	return
	//}
	user, err := m.Dao.AddUser(studentID, userRegisterRequestDTO.Password, permission.User)
	if err != nil {
		return
	}

	res = &dto.UserRegisterResponseDTO{
		ID:        user.ID,
		Username:  user.Username,
		StudentID: user.StudentID,
		Role:      user.Role,
	}
	return
}

func (m *UserService) SendEmail(userSendEmailRequestDTO *dto.UserSendEmailRequestDTO) error {
	studentID := userSendEmailRequestDTO.StudentID
	code, err := utils.SendEmail(studentID)
	if err != nil {
		return err
	}

	// 把 code 存到 redis里
	return utils.SetRedis(studentID+"_verification", code, time.Duration(viper.GetInt("smtp.expiration"))*time.Minute)
}

func (m *UserService) UpdateAvatar(userUpdateAvatarRequestDTO *dto.UserUpdateAvatarRequestDTO) (res *dto.UserUpdateAvatarResponseDTO, err error) {
	id := userUpdateAvatarRequestDTO.ID
	url := userUpdateAvatarRequestDTO.Url

	err = m.Dao.UpdateUser(id, map[string]interface{}{
		"avatar": url,
	})
	if err != nil {
		return
	}

	res = &dto.UserUpdateAvatarResponseDTO{
		Url: url,
	}
	return
}

func (m *UserService) AdminList(requestDTO *dto.AdminUserListRequestDTO) (responseDTO *dto.AdminUserListResponseDTO, err error) {
	id := requestDTO.ID
	studentID := requestDTO.StudentID
	username := requestDTO.Username
	role := requestDTO.Role
	page := requestDTO.Page
	pageSize := requestDTO.PageSize
	users, err := m.Dao.AdminList(id, studentID, username, role, page, pageSize)
	if err != nil {
		return
	}
	userOverviews := make([]dto.UserOverviewDTO, 0)
	for i := range users {
		userOverviews = append(userOverviews, dto.UserOverviewDTO{
			ID:        users[i].ID,
			StudentID: users[i].StudentID,
			Username:  users[i].Username,
			Role:      users[i].Role,
			CreatedAt: users[i].CreatedAt,
		})
	}
	responseDTO = &dto.AdminUserListResponseDTO{
		Users: userOverviews,
		Pagination: dto.PaginationInfo{
			TotalRecords: 777777,
			PageSize:     pageSize,
		},
	}
	return
}
