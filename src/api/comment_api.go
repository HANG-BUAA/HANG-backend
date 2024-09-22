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

// Create 创建评论
func (m CommentApi) Create(c *gin.Context) {
	iUserID, _ := c.Get("id")
	var iCommentCreateRequestDTO dto.CommentCreateRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iCommentCreateRequestDTO}).GetError(); err != nil {
		return
	}
	iCommentCreateRequestDTO.UserID = iUserID.(uint)

	iCommentCreateResponseDTO, err := m.Service.CreateComment(&iCommentCreateRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_COMMENT_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *iCommentCreateResponseDTO,
	})
}
