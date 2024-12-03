package service

import (
	"HANG-backend/src/custom_error"
	"HANG-backend/src/dao"
	"HANG-backend/src/global"
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

func (m *CourseService) CreateReview(requestDTO *dto.CourseReviewCreateRequestDTO) (res *dto.CourseReviewCreateResponseDTO, err error) {
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

	overview, err := m.Dao.ConvertReviewModelToOverviewDTO(review, user)
	if err != nil {
		return
	}
	res = (*dto.CourseReviewCreateResponseDTO)(overview)

	go func() {
		err := utils.PublishCourseReviewMessage(utils.CourseReviewMessage{
			ID:       review.ID,
			CourseID: review.CourseID,
			Content:  review.Content,
		})
		if err != nil {
			global.Logger.Error(err)
		}
	}()
	return
}

func (m *CourseService) LikeReview(requestDTO *dto.CourseReviewLikeRequestDTO) (err error) {
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

func (m *CourseService) UnlikeReview(requestDTO *dto.CourseReviewUnlikeRequestDTO) (err error) {
	user := requestDTO.User
	courseReview := requestDTO.CourseReview
	if !m.Dao.CheckReviewLiked(user, courseReview) {
		return errors.New("unliked review")
	}

	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.UnlikeReview(user, courseReview)
		if err == nil {
			return
		}
		if errors.Is(err, &custom_error.OptimisticLockError{}) {
			continue
		}
	}
	return custom_error.NewOptimisticLockError()
}

func (m *CourseService) UnlikeMaterial(requestDTO *dto.CourseMaterialUnlikeRequestDTO) (err error) {
	user := requestDTO.User
	courseMaterial := requestDTO.CourseMaterial

	if !m.Dao.CheckMaterialLiked(user, courseMaterial) {
		return errors.New("unliked material")
	}

	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.UnlikeMaterial(user, courseMaterial)
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

	var nextCursor any
	if isEnd {
		nextCursor = 0
	} else {
		nextCursor = courses[len(courses)-1].ID
	}
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

	// 检查课程是否存在
	_, err = m.Dao.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}

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

	var nextCursor any
	if isEnd {
		nextCursor = nil
	} else {
		nextCursor = fmt.Sprintf("%d %d", reviews[len(reviews)-1].LikeNum, reviews[len(reviews)-1].ID)
	}

	res = &dto.CourseReviewListResponseDTO{
		Pagination: *dto.BuildPaginationInfo(total, len(overviews), nextCursor),
		Reviews:    overviews,
	}
	return
}

func (m *CourseService) SearchListReview(requestDTO *dto.CourseReviewListRequestDTO) (res *dto.CourseReviewListResponseDTO, err error) {
	// todo baseURL 换成全局变量
	baseURL := fmt.Sprintf("http://%s:%s/course_review",
		viper.GetString("search_client.host"),
		viper.GetString("search_client.port"),
	)
	user := requestDTO.User
	courseID := requestDTO.CourseID
	query := requestDTO.Query
	cursor := requestDTO.Cursor
	pageSize := requestDTO.PageSize

	params := url.Values{}
	params.Add("query", *query)
	params.Add("course_id", courseID)
	params.Add("page_size", strconv.Itoa(pageSize))
	params.Add("cursor", cursor)
	fullUrl := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 发送请求
	resp, err := http.Get(fullUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// 响应结构
	var responseBody struct {
		Total         int `json:"total"`
		CourseReviews []struct {
			ID    uint    `json:"id"`
			Score float64 `json:"score"`
		} `json:"course_reviews"`
		NextCursor string `json:"next_cursor"`
	}

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return
	}

	var ids []uint
	total := responseBody.Total
	for _, courseReview := range responseBody.CourseReviews {
		ids = append(ids, courseReview.ID)
	}
	nextCursor := responseBody.NextCursor

	// 获取 review 列表
	reviews, err := m.Dao.GetReviewsByIDs(ids)
	if err != nil {
		return
	}

	overviews, err := m.Dao.ConvertReviewModelsToOverviewDTOs(reviews, user)
	if err != nil {
		return
	}
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

func (m *CourseService) CreateMaterial(requestDTO *dto.CourseMaterialCreateRequestDTO) (res *dto.CourseMaterialCreateResponseDTO, err error) {
	user := requestDTO.User
	courseID := requestDTO.CourseID
	link := requestDTO.Link
	source := requestDTO.Source
	description := requestDTO.Description

	// 检查课程是否存在
	_, err = m.Dao.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}

	material, err := m.Dao.CreateMaterial(user, courseID, link, description, source, false, false)
	if err != nil {
		return nil, err
	}

	overview, err := m.Dao.ConvertMaterialModelToOverviewDTO(material, user)
	if err != nil {
		return
	}
	res = (*dto.CourseMaterialCreateResponseDTO)(overview)
	return
}

func (m *CourseService) LikeMaterial(requestDTO *dto.CourseMaterialLikeRequestDTO) (err error) {
	user := requestDTO.User
	courseMaterial := requestDTO.CourseMaterial

	// 判断用户是否已经喜欢
	if m.Dao.CheckMaterialLiked(user, courseMaterial) {
		return errors.New("liked material")
	}

	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.LikeMaterial(user, courseMaterial)
		if err == nil {
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
	}
	return custom_error.NewOptimisticLockError()
}

func (m *CourseService) ListMaterial(requestDTO *dto.CourseMaterialListRequestDTO) (res *dto.CourseMaterialListResponseDTO, err error) {
	pageSize := requestDTO.PageSize
	user := requestDTO.User
	courseID := requestDTO.CourseID
	isOfficial := requestDTO.IsOfficial

	// 检查课程是否存在
	_, err = m.Dao.GetCourseByID(courseID)
	if err != nil {
		return
	}

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

	materials, total, isEnd, err := m.Dao.ListMaterial(cursor, pageSize, courseID, *isOfficial)
	if err != nil {
		return
	}

	overviews, err := m.Dao.ConvertMaterialModelsToOverviews(materials, user)
	if err != nil {
		return
	}

	var nextCursor any
	if isEnd {
		nextCursor = nil
	} else {
		nextCursor = fmt.Sprintf("%d %d", materials[len(materials)-1].LikeNum, materials[len(materials)-1].ID)
	}

	res = &dto.CourseMaterialListResponseDTO{
		Pagination: *dto.BuildPaginationInfo(total, len(materials), nextCursor),
		Materials:  overviews,
	}
	return
}
