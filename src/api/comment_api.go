package api

import (
	"HANG-backend/src/global"
	"HANG-backend/src/service"
	"HANG-backend/src/service/dto"
	"github.com/gin-gonic/gin"
)

type CommentApi struct {
	BaseApi
	Service *service.CommentService
}

func NewCommentApi() CommentApi {
	return CommentApi{
		BaseApi: NewBaseApi(),
		Service: service.NewCommentService(),
	}
}

// 创建评论
func (m CommentApi) Create(c *gin.Context) {
	userID, _ := c.Get("id")
	var commentCreateRequestDTO dto.CommentCreateRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &commentCreateRequestDTO}).GetError(); err != nil {
		return
	}
	commentCreateRequestDTO.UserID = userID.(uint)

	commentCreateResponseDTO, err := m.Service.Create(&commentCreateRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_COMMENT_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *commentCreateResponseDTO,
	})
}
