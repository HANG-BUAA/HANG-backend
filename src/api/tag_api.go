package api

import (
	"HANG-backend/src/service"
	"HANG-backend/src/service/dto"
	"github.com/gin-gonic/gin"
)

type TagApi struct {
	BaseApi
	Service *service.TagService
}

func NewTagApi() TagApi {
	return TagApi{
		BaseApi: NewBaseApi(),
		Service: service.NewTagService(),
	}
}

func (m TagApi) AdminCreateTag(c *gin.Context) {
	var requestDTO dto.AdminTagCreateRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &requestDTO}).GetError(); err != nil {
		return
	}
	responseDTO, err := m.Service.AdminCreate(&requestDTO)
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
