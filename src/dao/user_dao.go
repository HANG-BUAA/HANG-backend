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

func (m *UserDao) GetUserByStudentID(studentID string) (model.User, error) {
	var user model.User
	err := m.Orm.Model(&user).Where("student_id = ?", studentID).Find(&user).Error
	return user, err
}

func (m *UserDao) GetUserByName(username string) (model.User, error) {
	var user model.User
	err := m.Orm.Where("user_name = ?", username).Find(&user).Error
	return user, err
}

func (m *UserDao) GetUserByID(id uint) (model.User, error) {
	var user model.User
	err := m.Orm.Where("id = ?", id).Find(&user).Error
	return user, err
}

func (m *UserDao) CheckStudentIDExist(studentID string) bool {
	var total int64
	m.Orm.Model(&model.User{}).Where("student_id = ?", studentID).Count(&total)
	return total > 0
}

func (m *UserDao) AddUser(studentID, password string) (model.User, error) {
	user := model.User{
		StudentID: studentID,
		UserName:  studentID,
		Password:  password,
		Role:      1,
	}
	if err := m.Orm.Create(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (m *UserDao) UpdateUser(userID uint, updatedFields map[string]interface{}) error {
	if err := m.Orm.Model(&model.User{}).Where("id = ?", userID).Updates(updatedFields).Error; err != nil {
		return err
	}
	return nil
}
