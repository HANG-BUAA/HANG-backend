package service

import (
	"HANG-backend/src/dao"
	"HANG-backend/src/service/dto"
	"errors"
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
	tmp, err := m.Dao.ConvertCourseModelToOverviewDTO(course)
	if err != nil {
		return nil, err
	}
	res = (*dto.AdminCourseCreateResponseDTO)(tmp)
	return
}

func (m *CourseService) CreateReview(requestDTO *dto.CreateCourseReviewRequestDTO) (res *dto.CreateCourseReviewResponseDTO, err error) {
	user := requestDTO.User
	courseID := requestDTO.CourseID
	Content := requestDTO.Content
	score := requestDTO.Score

	// 检查课程是否存在
	_, err = m.Dao.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}

	review, err := m.Dao.CreateReview(courseID, user, score, Content)
	if err != nil {
		return
	}
	res = &dto.CreateCourseReviewResponseDTO{
		ID:        review.ID,
		CourseID:  review.CourseID,
		Content:   review.Content,
		Score:     review.Score,
		IsSelf:    true,
		LikeNum:   0,
		HasLiked:  false,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
		DeletedAt: review.DeletedAt,
	}
	return
}

func (m *CourseService) LikeReview(requestDTO *dto.LikeCourseReviewRequestDTO) (err error) {
	user := requestDTO.User
	courseReview := requestDTO.CourseReview

	// 判断用户是否已经喜欢评论
	if m.Dao.CheckReviewLiked(user, courseReview) {
		return errors.New("liked review")
	}

	err = m.Dao.LikeReview(user, courseReview)
	return
}
