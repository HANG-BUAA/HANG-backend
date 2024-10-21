package api

import (
	"HANG-backend/src/custom_error"
	"HANG-backend/src/global"
	"HANG-backend/src/model"
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
	user := c.MustGet("user").(*model.User)
	var commentCreateRequestDTO dto.CommentCreateRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &commentCreateRequestDTO}).GetError(); err != nil {
		return
	}
	commentCreateRequestDTO.User = user

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
	user := c.MustGet("user").(*model.User)
	comment := c.MustGet("comment").(*model.Comment)
	var commentLikeRequestDTO dto.CommentLikeRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &commentLikeRequestDTO}).GetError(); err != nil {
		return
	}
	commentLikeRequestDTO.User = user
	commentLikeRequestDTO.Comment = comment

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

func (m CommentApi) Unlike(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	comment := c.MustGet("comment").(*model.Comment)
	var commentUnlikeRequestDTO dto.CommentUnlikeRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &commentUnlikeRequestDTO}).GetError(); err != nil {
		return
	}
	commentUnlikeRequestDTO.User = user
	commentUnlikeRequestDTO.Comment = comment

	err := m.Service.Unlike(&commentUnlikeRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_COMMENT_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: gin.H{
			"status": "unlike success",
		},
	})
}

func (m CommentApi) List(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	cursor := c.MustGet("cursor").(string)
	pageSize := c.MustGet("page_size").(int)
	var commentListRequestDTO dto.CommentListRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &commentListRequestDTO}).GetError(); err != nil {
		return
	}
	commentListRequestDTO.User = user
	commentListRequestDTO.Cursor = cursor
	commentListRequestDTO.PageSize = pageSize

	// 根据 level 选择服务
	var comments *dto.CommentListResponseDTO
	var err error
	level := commentListRequestDTO.Level
	if level == 1 {
		if commentListRequestDTO.PostID == 0 {
			m.Fail(ResponseJson{
				Code: global.ERR_CODE_COMMENT_FAILED,
				Msg:  "post_id is Required",
			})
			return
		}
		comments, err = m.Service.ListFirstLevel(&commentListRequestDTO)
	} else if level == 2 {
		if commentListRequestDTO.CommentID == 0 {
			m.Fail(ResponseJson{
				Code: global.ERR_CODE_COMMENT_FAILED,
				Msg:  "comment_id is Required",
			})
			return
		}
		comments, err = m.Service.ListSecondLevel(&commentListRequestDTO)
	} else {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_COMMENT_FAILED,
			Msg:  "invalid level",
		})
		return
	}

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
