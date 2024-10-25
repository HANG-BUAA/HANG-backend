package dao

import (
	"HANG-backend/src/model"
	"gorm.io/gorm"
)

var courseDao *CourseDao

type CourseDao struct {
	BaseDao
}

func NewCourseDao() *CourseDao {
	if courseDao == nil {
		courseDao = &CourseDao{
			NewBaseDao(),
		}
	}
	return courseDao
}

func (m *CourseDao) Create(id, name string, credits *float32, campus *int, tags []model.Tag) (*model.Course, error) {
	course := model.Course{
		ID:      id,
		Name:    name,
		Credits: credits,
		Campus:  campus,
	}

	// 使用数据库事务
	err := m.Orm.Transaction(func(tx *gorm.DB) error {
		// 1. 创建课程
		if err2 := tx.Create(&course).Error; err2 != nil {
			return err2
		}

		// 2. 加标签
		if tags != nil {
			for _, tag := range tags {
				courseTag := model.CourseTag{
					CourseID: course.ID,
					TagID:    tag.ID,
				}
				if err2 := tx.Create(&courseTag).Error; err2 != nil {
					return err2
				}
			}
		}
		return nil
	})
	return &course, err
}

func (m *CourseDao) CheckTagsExist(tags []uint) ([]model.Tag, error) {
	var result []model.Tag
	for _, tagID := range tags {
		var tag model.Tag
		if err := m.Orm.First(&tag, tagID).Error; err != nil {
			return nil, err
		}
		result = append(result, tag)
	}
	return result, nil
}
