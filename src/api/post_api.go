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

func (m PostApi) Unlike(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	post := c.MustGet("post").(*model.Post)
	var postUnlikeRequestDTO dto.PostUnlikeRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &postUnlikeRequestDTO}).GetError(); err != nil {
		return
	}
	postUnlikeRequestDTO.User = user
	postUnlikeRequestDTO.Post = post

	err := m.Service.Unlike(&postUnlikeRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_POST_FAILED,
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

func (m PostApi) Uncollect(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	post := c.MustGet("post").(*model.Post)
	var postUncollectRequestDTO dto.PostUncollectRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &postUncollectRequestDTO}).GetError(); err != nil {
		return
	}
	postUncollectRequestDTO.User = user
	postUncollectRequestDTO.Post = post

	err := m.Service.Uncollect(&postUncollectRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_POST_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: gin.H{
			"status": "uncollect success",
		},
	})
}

func (m PostApi) List(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	cursor := c.MustGet("cursor").(string)
	pageSize := c.MustGet("page_size").(int)
	var postListRequestDTO dto.PostListRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &postListRequestDTO}).GetError(); err != nil {
		return
	}
	postListRequestDTO.User = user
	postListRequestDTO.Cursor = cursor
	postListRequestDTO.PageSize = pageSize

	var posts *dto.PostListResponseDTO
	var err error

	// 根据是否有 query 判断走哪个服务
	if postListRequestDTO.Query != "" {
		posts, err = m.Service.SearchList(&postListRequestDTO)
	} else {
		posts, err = m.Service.CommonList(&postListRequestDTO)
	}

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

func (m PostApi) CollectionList(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	cursor := c.MustGet("cursor").(string)
	pageSize := c.MustGet("page_size").(int)
	var postCollectionListRequestDTO dto.PostCollectionListRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &postCollectionListRequestDTO}).GetError(); err != nil {
		return
	}
	postCollectionListRequestDTO.User = user
	postCollectionListRequestDTO.Cursor = cursor
	postCollectionListRequestDTO.PageSize = pageSize

	posts, err := m.Service.CollectionList(&postCollectionListRequestDTO)
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

func (m PostApi) Retrieve(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	post := c.MustGet("post").(*model.Post)
	var postRetrieveRequestDTO dto.PostRetrieveRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &postRetrieveRequestDTO}).GetError(); err != nil {
		return
	}
	postRetrieveRequestDTO.User = user
	postRetrieveRequestDTO.Post = post

	postOverview, err := m.Service.Retrieve(&postRetrieveRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_POST_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *postOverview,
	})
}
