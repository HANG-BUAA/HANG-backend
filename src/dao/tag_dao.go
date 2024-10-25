package dao

import "HANG-backend/src/model"

var tagDao *TagDao

type TagDao struct {
	BaseDao
}

func NewTagDao() *TagDao {
	if tagDao == nil {
		tagDao = &TagDao{
			NewBaseDao(),
		}
	}
	return tagDao
}

func (m *TagDao) Create(tagType int, name string) (*model.Tag, error) {
	tag := model.Tag{
		Name: name,
		Type: tagType,
	}
	if err := m.Orm.Create(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (m *TagDao) GetTagByID(id uint) (*model.Tag, error) {
	tag := &model.Tag{}
	if err := m.Orm.Where("id = ?", id).First(&tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (m *TagDao) GetTagByName(name string) (*model.Tag, error) {
	tag := &model.Tag{}
	if err := m.Orm.Where("name = ?", name).First(&tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (m *TagDao) ListByType(tagType int) ([]model.Tag, error) {
	var tags []model.Tag
	if err := m.Orm.Where("type = ?", tagType).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
