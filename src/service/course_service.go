package service

import (
	"HANG-backend/src/custom_error"
	"HANG-backend/src/dao"
	"HANG-backend/src/global"
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"errors"
	"fmt"
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

	// todo 用转换函数
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

	go func() {
		err := utils.PublishCourseReviewMessage(utils.CourseReviewMessage{
			ID:      review.ID,
			Content: review.Content,
		})
		if err != nil {
			global.Logger.Error(err)
		}
	}()
	return
}

func (m *CourseService) LikeReview(requestDTO *dto.LikeCourseReviewRequestDTO) (err error) {
	user := requestDTO.User
	courseReview := requestDTO.CourseReview

	// 判断用户是否已经喜欢评论
	if m.Dao.CheckReviewLiked(user, courseReview) {
		return errors.New("liked review")
	}

	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.LikeReview(user, courseReview)
		if err == nil {
			return
		}
		if errors.Is(err, &custom_error.OptimisticLockError{}) {
			continue
		}
	}
	return custom_error.NewOptimisticLockError()
}

func (m *CourseService) ListCourse(requestDTO *dto.CourseListRequestDTO) (res *dto.CourseListResponseDTO, err error) {
	pageSize := requestDTO.PageSize
	keyword := requestDTO.Keyword
	tags := requestDTO.Tags
	cursor := requestDTO.Cursor

	courses, total, isEnd, err := m.Dao.ListCourse(cursor, pageSize, keyword, tags)
	if err != nil {
		return
	}
	if len(courses) == 0 {
		res = &dto.CourseListResponseDTO{
			Courses: []dto.CourseOverviewDTO{},
		}
		return
	}

	overviews, err := m.Dao.ConvertCourseModelsToOverviewDTOs(courses)
	if err != nil {
		return
	}
	nextCursor := utils.IfThenElse(isEnd, 0, courses[len(courses)-1].ID)
	res = &dto.CourseListResponseDTO{
		Pagination: *dto.BuildPaginationInfo(total, len(overviews), nextCursor),
		Courses:    overviews,
	}
	return
}

func (m *CourseService) CommonListReview(requestDTO *dto.CourseReviewListRequestDTO) (res *dto.CourseReviewListResponseDTO, err error) {
	pageSize := requestDTO.PageSize
	user := requestDTO.User
	courseID := requestDTO.CourseID

	// todo 检查课程是否存在

	var cursorLikeNum int
	var cursorID uint
	cursor := new(struct {
		LikeNum int
		ID      uint
	})
	_, err = fmt.Sscanf(requestDTO.Cursor, "%d %d", &cursorLikeNum, &cursorID)
	if err != nil {
		cursor = nil
	} else {
		cursor.LikeNum = cursorLikeNum
		cursor.ID = cursorID
	}

	reviews, total, isEnd, err := m.Dao.CommonListReview(cursor, pageSize, courseID)
	if err != nil {
		return
	}
	if len(reviews) == 0 {
		res = &dto.CourseReviewListResponseDTO{
			Reviews: []dto.CourseReviewOverviewDTO{},
		}
		return
	}

	overviews, err := m.Dao.ConvertReviewModelsToOverviewDTOs(reviews, user)
	if err != nil {
		return
	}
	nextCursor := utils.IfThenElse(isEnd, nil, fmt.Sprintf("%d %d", reviews[len(reviews)-1].LikeNum, reviews[len(reviews)-1].ID))

	res = &dto.CourseReviewListResponseDTO{
		Pagination: *dto.BuildPaginationInfo(total, len(overviews), nextCursor),
		Reviews:    overviews,
	}
	return
}

func (m *CourseService) Retrieve(requestDTO *dto.CourseRetrieveRequestDTO) (res *dto.CourseRetrieveResponseDTO, err error) {
	course := requestDTO.Course

	overview, err := m.Dao.ConvertCourseModelToOverviewDTO(course)
	if err != nil {
		return
	}
	res = &dto.CourseRetrieveResponseDTO{
		Course: *overview,
	}
	return
}
