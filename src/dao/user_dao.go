package dao

import (
	"HANG-backend/src/model"
)

var userDao *UserDao

type UserDao struct {
	BaseDao
}

func NewUserDao() *UserDao {
	if userDao == nil {
		userDao = &UserDao{
			NewBaseDao(),
		}
	}
	return userDao
}

func (m *UserDao) GetUserByStudentID(iStudentID string) (model.User, error) {
	var iUser model.User
	err := m.Orm.Model(&iUser).Where("student_id = ?", iStudentID).Find(&iUser).Error
	return iUser, err
}

func (m *UserDao) GetUserByName(iUsername string) (model.User, error) {
	var iUser model.User
	err := m.Orm.Where("user_name = ?", iUsername).Find(&iUser).Error
	return iUser, err
}

func (m *UserDao) CheckStudentIDExist(iStudentID string) bool {
	var nTotal int64
	m.Orm.Model(&model.User{}).Where("student_id = ?", iStudentID).Count(&nTotal)
	return nTotal > 0
}

func (m *UserDao) AddUser(iStudentID, iPassword string) (model.User, error) {
	iUser := model.User{
		StudentID: iStudentID,
		UserName:  iStudentID,
		Password:  iPassword,
	}
	if err := m.Orm.Create(&iUser).Error; err != nil {
		return model.User{}, err
	}
	return iUser, nil
}

//func (m *UserDao) GetUserByNameAndPassword(stUsername, stPassword string) model.User {
//	var iUser model.User
//	m.Orm.Model(&iUser).Where("user_name = ? and password = ?", stUsername, stPassword).Find(&iUser)
//	return iUser
//}
//
//func (m *UserDao) CheckUserExit(stUserName string) bool {
//	var nTotal int64
//	m.Orm.Model(&model.User{}).Where("user_name = ?", stUserName).Count(&nTotal)
//	return nTotal > 0
//}
//
//func (m *UserDao) GetUserByName(stUsername string) (model.User, error) {
//	var iUser model.User
//	err := m.Orm.Where("user_name = ?", stUsername).Find(&iUser).Error
//	return iUser, err
//}
//
//func (m *UserDao) AddUser(iUserRegisterDTO *dto.UserRegisterDTO) error {
//	var iUser model.User
//	iUserRegisterDTO.ConvertToModel(&iUser)
//
//	err := m.Orm.Save(&iUser).Error
//	if err == nil {
//		iUserRegisterDTO.ID = iUser.ID
//		// 将验证码和密码清除，不返回
//		iUserRegisterDTO.Password = ""
//		iUserRegisterDTO.VerificationCode = ""
//	}
//	return err
//}
