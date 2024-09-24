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
	userID, _ := c.Get("id")
	var commentCreateRequestDTO dto.CommentCreateRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &commentCreateRequestDTO}).GetError(); err != nil {
		return
	}
	commentCreateRequestDTO.UserID = userID.(uint)

	commentCreateResponseDTO, err := m.Service.CreateComment(&commentCreateRequestDTO)
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

// Like 喜欢评论
func (m CommentApi) Like(c *gin.Context) {
	userID, _ := c.Get("id")
	var commentLikeRequestDTO dto.CommentLikeRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &commentLikeRequestDTO, BindParamsFromUri: true}).GetError(); err != nil {
		return
	}
	commentLikeRequestDTO.UserID = userID.(uint)

	err := m.Service.Like(&commentLikeRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_COMMENT_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: gin.H{
			"status": "like success",
		},
	})
}

// List 查询评论列表
func (m CommentApi) List(c *gin.Context) {
	userID, _ := c.Get("id")
	var commentListRequestDTO dto.CommentListRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &commentListRequestDTO}).GetError(); err != nil {
		return
	}
	commentListRequestDTO.UserID = userID.(uint)

	commentListResponseDTO, err := m.Service.List(&commentListRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_COMMENT_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *commentListResponseDTO,
	})
}
