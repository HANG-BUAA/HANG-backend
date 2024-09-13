package dao

import (
	"HANG-backend/src/model"
	"HANG-backend/src/service/dto"
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

func (m *UserDao) GetUserByNameAndPassword(stUsername, stPassword string) model.User {
	var iUser model.User
	m.Orm.Model(&iUser).Where("user_name = ? and password = ?", stUsername, stPassword).Find(&iUser)
	return iUser
}

func (m *UserDao) CheckUserExit(stUserName string) bool {
	var nTotal int64
	m.Orm.Model(&model.User{}).Where("user_name = ?", stUserName).Count(&nTotal)
	return nTotal > 0
}

func (m *UserDao) GetUserByName(stUsername string) (model.User, error) {
	var iUser model.User
	err := m.Orm.Where("user_name = ?", stUsername).Find(&iUser).Error
	return iUser, err
}

func (m *UserDao) AddUser(iUserAddDTO *dto.UserRegisterDTO) error {
	var iUser model.User
	iUserAddDTO.ConvertToModel(&iUser)

	err := m.Orm.Save(&iUser).Error
	if err == nil {
		iUserAddDTO.ID = iUser.ID
		iUserAddDTO.Password = ""
	}
	return err
}
