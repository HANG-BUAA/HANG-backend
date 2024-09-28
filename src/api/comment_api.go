package api

import (
	"HANG-backend/src/custom_error"
	"HANG-backend/src/global"
	"HANG-backend/src/service"
	"HANG-backend/src/service/dto"
	"errors"
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
	userID := c.MustGet("id")
	var commentCreateRequestDTO dto.CommentCreateRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &commentCreateRequestDTO}).GetError(); err != nil {
		return
	}
	commentCreateRequestDTO.UserID = userID.(uint)

	commentCreateResponseDTO, err := m.Service.Create(&commentCreateRequestDTO)
	if err != nil {
		code := global.ERR_CODE_COMMENT_FAILED
		if errors.Is(err, &custom_error.OptimisticLockError{}) {
			code = global.OptimisticLockMaxRetries
		}
		m.Fail(ResponseJson{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *commentCreateResponseDTO,
	})
}

func (m CommentApi) Like(c *gin.Context) {
	userID := c.MustGet("id")
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

func (m CommentApi) List(c *gin.Context) {
	userID := c.MustGet("id")
	page := c.MustGet("page").(int)
	pageSize := c.MustGet("page_size").(int)
	var commentListRequestDTO dto.CommentListRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &commentListRequestDTO}).GetError(); err != nil {
		return
	}
	commentListRequestDTO.UserID = userID.(uint)
	commentListRequestDTO.Page = page
	commentListRequestDTO.PageSize = pageSize

	comments, err := m.Service.List(&commentListRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_COMMENT_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *comments,
	})
}
