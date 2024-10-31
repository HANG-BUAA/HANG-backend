package api

import (
	"HANG-backend/src/model"
	"HANG-backend/src/service"
	"HANG-backend/src/service/dto"
	"github.com/gin-gonic/gin"
)

type CourseApi struct {
	BaseApi
	Service *service.CourseService
}

func NewCourseApi() CourseApi {
	return CourseApi{
		BaseApi: NewBaseApi(),
		Service: service.NewCourseService(),
	}
}

func (m CourseApi) CreateCourse(c *gin.Context) {
	var requestDTO dto.AdminCourseCreateRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	responseDTO, err := m.Service.CreateCourse(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *responseDTO,
	})
}

func (m CourseApi) CreateReview(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	var requestDTO dto.CreateCourseReviewRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	requestDTO.User = user
	responseDTO, err := m.Service.CreateReview(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Data: *responseDTO,
	})
}

func (m CourseApi) LikeReview(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	review := c.MustGet("course_review").(*model.CourseReview)
	var requestDTO dto.LikeCourseReviewRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	requestDTO.User = user
	requestDTO.CourseReview = review

	err := m.Service.LikeReview(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: gin.H{
			"status": "like success",
		},
	})
}

func (m CourseApi) ListCourse(c *gin.Context) {
	cursor := c.MustGet("cursor").(string)
	pageSize := c.MustGet("page_size").(int)
	var requestDTO dto.CourseListRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	requestDTO.Cursor = cursor
	requestDTO.PageSize = pageSize

	courses, err := m.Service.ListCourse(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *courses,
	})
}

func (m CourseApi) ListReview(c *gin.Context) {
	cursor := c.MustGet("cursor").(string)
	pageSize := c.MustGet("page_size").(int)
	user := c.MustGet("user").(*model.User)
	var requestDTO dto.CourseReviewListRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	requestDTO.Cursor = cursor
	requestDTO.PageSize = pageSize
	requestDTO.User = user

	// todo 搜索服务
	var reviews *dto.CourseReviewListResponseDTO
	var err error
	if requestDTO.Query != nil {
		reviews, err = m.Service.SearchListReview(&requestDTO)
	} else {
		reviews, err = m.Service.CommonListReview(&requestDTO)
	}
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *reviews,
	})
}

func (m CourseApi) Retrieve(c *gin.Context) {
	course := c.MustGet("course").(*model.Course)
	var requestDTO dto.CourseRetrieveRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	requestDTO.Course = course

	courseOverview, err := m.Service.Retrieve(&requestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *courseOverview,
	})
}
