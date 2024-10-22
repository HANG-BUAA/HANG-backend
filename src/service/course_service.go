package service

import (
	"HANG-backend/src/dao"
)

var courseService *CourseService

type CourseService struct {
	BaseService
	Dao *dao.CourseDao
}

func NewCourseService() *CourseService {
	if courseService == nil {
		courseService = &CourseService{
			Dao: dao.NewCourseDao(),
		}
	}
	return courseService
}

//func (m *CourseService) CreateCourse(requestDTO *dto.AdminCourseCreateRequestDTO) (res *dto.AdminCourseCreateResponseDTO, err error) {
//	id := requestDTO.ID
//	name := requestDTO.Name
//	credits := requestDTO.Credits
//	campus := requestDTO.Campus
//	//tags := requestDTO.Tags  // todo add tags
//
//}
