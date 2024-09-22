package api

import (
	"HANG-backend/src/global"
	"HANG-backend/src/service"
	"HANG-backend/src/service/dto"
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
	iUserID, _ := c.Get("id")
	var iPostCreateRequestDTO dto.PostCreateRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iPostCreateRequestDTO}).GetError(); err != nil {
		return
	}
	iPostCreateRequestDTO.UserID = iUserID.(uint)

	iPostCreateResponseDTO, err := m.Service.CreatePost(&iPostCreateRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_POST_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *iPostCreateResponseDTO,
	})
}

// Like 喜欢帖子
func (m PostApi) Like(c *gin.Context) {
	iUserID, _ := c.Get("id")
	var iPostLikeRequestDTO dto.PostLikeRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iPostLikeRequestDTO, BindParamsFromUri: true}).GetError(); err != nil {
		return
	}
	iPostLikeRequestDTO.UserID = iUserID.(uint)

	err := m.Service.Like(&iPostLikeRequestDTO)
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
	iUserID, _ := c.Get("id")
	var iPostCollectRequestDTO dto.PostCollectRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iPostCollectRequestDTO, BindParamsFromUri: true}).GetError(); err != nil {
		return
	}
	iPostCollectRequestDTO.UserID = iUserID.(uint)

	err := m.Service.Collect(&iPostCollectRequestDTO)
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

// List 查询帖子列表
func (m PostApi) List(c *gin.Context) {
	iUserID, _ := c.Get("id")
	var iPostListRequestDTO dto.PostListTRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iPostListRequestDTO}).GetError(); err != nil {
		return
	}
	iPostListRequestDTO.UserID = iUserID.(uint)

	iPostListResponseDTO, err := m.Service.List(&iPostListRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_POST_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: iPostListResponseDTO,
	})
}
