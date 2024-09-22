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

func (m *UserDao) GetUserByID(id uint) (model.User, error) {
	var iUser model.User
	err := m.Orm.Where("id = ?", id).Find(&iUser).Error
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
		Role:      1,
	}
	if err := m.Orm.Create(&iUser).Error; err != nil {
		return model.User{}, err
	}
	return iUser, nil
}

func (m *UserDao) UpdateUser(userID uint, updatedFields map[string]interface{}) error {
	if err := m.Orm.Model(&model.User{}).Where("id = ?", userID).Updates(updatedFields).Error; err != nil {
		return err
	}
	return nil
}
