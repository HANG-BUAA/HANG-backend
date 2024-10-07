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

type PostApi struct {
	BaseApi
	Service *service.PostService
}

func NewPostApi() PostApi {
	return PostApi{
		BaseApi: NewBaseApi(),
		Service: service.NewPostService(),
	}
}

// Create 创建帖子
func (m PostApi) Create(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	var postCreateRequestDTO dto.PostCreateRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &postCreateRequestDTO}).GetError(); err != nil {
		return
	}
	postCreateRequestDTO.User = user

	postCreateResponseDTO, err := m.Service.Create(&postCreateRequestDTO)
	if err != nil {
		code := global.ERR_CODE_POST_FAILED
		if errors.Is(err, &custom_error.OptimisticLockError{}) {
			code = global.ERR_CODE_OPTISTIC_LOCK_RETRY_LIMIT
		}
		m.Fail(ResponseJson{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *postCreateResponseDTO,
	})
}

// Like 喜欢帖子
func (m PostApi) Like(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	post := c.MustGet("post").(*model.Post)
	var postLikeRequestDTO dto.PostLikeRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &postLikeRequestDTO}).GetError(); err != nil {
		return
	}
	postLikeRequestDTO.User = user
	postLikeRequestDTO.Post = post

	err := m.Service.Like(&postLikeRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_POST_FAILED,
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

// Collect 收藏帖子
func (m PostApi) Collect(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	post := c.MustGet("post").(*model.Post)
	var postCollectRequestDTO dto.PostCollectRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &postCollectRequestDTO, BindParamsFromUri: true}).GetError(); err != nil {
		return
	}
	postCollectRequestDTO.User = user
	postCollectRequestDTO.Post = post

	err := m.Service.Collect(&postCollectRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_POST_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: gin.H{
			"status": "collect success",
		},
	})
}

func (m PostApi) List(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	page := c.MustGet("page").(int)
	pageSize := c.MustGet("page_size").(int)
	var postListRequestDTO dto.PostListRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &postListRequestDTO}).GetError(); err != nil {
		return
	}
	postListRequestDTO.User = user
	postListRequestDTO.Page = page
	postListRequestDTO.PageSize = pageSize

	posts, err := m.Service.List(&postListRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_POST_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *posts,
	})
}
