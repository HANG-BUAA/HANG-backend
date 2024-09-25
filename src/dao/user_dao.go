package dao

import (
	"HANG-backend/src/model"
	"HANG-backend/src/permission"
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
	err := m.Orm.Where("username = ?", username).Find(&user).Error
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

func (m *UserDao) AddUser(studentID, password string, role permission.Role) (model.User, error) {
	user := model.User{
		StudentID: studentID,
		Username:  studentID,
		Password:  password,
		Role:      1,
	}
	if err := m.Orm.Create(&user).Error; err != nil {
		return model.User{}, err
	}
	// 初始化用户权限
	if err := permission.InitUserPermission(user.ID, role); err != nil {
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

func (m *UserDao) AdminList(id *uint, studentID *string, username *string, Role *uint, page int, pageSize int) ([]model.User, error) {
	var users []model.User
	query := m.Orm.Unscoped().Model(&model.User{})
	if id != nil {
		query = query.Where("id = ?", *id)
	}
	if studentID != nil {
		query = query.Where("student_id = ?", *studentID)
	}
	if username != nil {
		query = query.Where("username = ?", *username)
	}
	if Role != nil {
		query = query.Where("role = ?", *Role)
	}

	// 添加分页
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize).Find(&users)
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
