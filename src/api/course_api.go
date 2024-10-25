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

func (m CourseApi) CreateCourseReview(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	var requestDTO dto.CreateCourseReviewRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	requestDTO.User = user
	responseDTO, err := m.Service.CreateCourseReview(&requestDTO)
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
