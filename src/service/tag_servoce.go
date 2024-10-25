package service

import (
	"HANG-backend/src/dao"
	"HANG-backend/src/service/dto"
)

var tagService *TagService

type TagService struct {
	BaseService
	Dao *dao.TagDao
}

func NewTagService() *TagService {
	if tagService == nil {
		tagService = &TagService{
			Dao: dao.NewTagDao(),
		}
	}
	return tagService
}

func (m *TagService) AdminCreate(requestDTO *dto.AdminTagCreateRequestDTO) (res *dto.AdminTagCreateResponseDTO, err error) {
	name := requestDTO.Name
	tagType := requestDTO.Type

	tag, err := m.Dao.Create(tagType, name)
	if err != nil {
		return
	}
	res = &dto.AdminTagCreateResponseDTO{
		ID:   tag.ID,
		Type: tag.Type,
		Name: tag.Name,
	}
	return
}
