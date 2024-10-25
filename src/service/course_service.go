package service

import (
	"HANG-backend/src/dao"
	"HANG-backend/src/service/dto"
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

func (m *CourseService) CreateCourse(requestDTO *dto.AdminCourseCreateRequestDTO) (res *dto.AdminCourseCreateResponseDTO, err error) {
	id := requestDTO.ID
	name := requestDTO.Name
	credits := requestDTO.Credits
	campus := requestDTO.Campus
	tagIDs := requestDTO.Tags

	// 检查 tags 合法性
	tags, err := m.Dao.ListTagsByIDs(tagIDs)
	if err != nil {
		return nil, err
	}

	course, err := m.Dao.CreateCourse(id, name, credits, campus, tags)
	if err != nil {
		return
	}
	res = &dto.AdminCourseCreateResponseDTO{
		ID:        course.ID,
		Name:      course.Name,
		Credits:   course.Credits,
		Campus:    course.Campus,
		Tags:      tagIDs, // todo 返回格式会变化，此处可以直接返回 tags 详细信息
		CreatedAt: course.CreatedAt,
		UpdatedAt: course.UpdatedAt,
		DeletedAt: course.DeletedAt,
	}
	return
}

func (m *CourseService) CreateCourseReview(requestDTO *dto.CreateCourseReviewRequestDTO) (res *dto.CreateCourseReviewResponseDTO, err error) {
	user := requestDTO.User
	courseID := requestDTO.CourseID
	Content := requestDTO.Content
	score := requestDTO.Score

	// 检查课程是否存在
	_, err = m.Dao.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}

	review, err := m.Dao.CreateCourseReview(courseID, user, score, Content)
	if err != nil {
		return
	}
	res = &dto.CreateCourseReviewResponseDTO{
		ID:        review.ID,
		CourseID:  review.CourseID,
		Content:   review.Content,
		Score:     review.Score,
		IsSelf:    true,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
		DeletedAt: review.DeletedAt,
	}
	return
}
