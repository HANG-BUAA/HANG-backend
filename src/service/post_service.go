package service

import (
	"HANG-backend/src/custom_error"
	"HANG-backend/src/dao"
	"HANG-backend/src/global"
	"HANG-backend/src/service/dto"
	"errors"
)

var postService *PostService

type PostService struct {
	BaseService
	Dao *dao.PostDao
}

func NewPostService() *PostService {
	if postService == nil {
		postService = &PostService{
			Dao: dao.NewPostDao(),
		}
	}
	return postService
}

// Create 创建帖子
func (m *PostService) Create(postCreateDTO *dto.PostCreateRequestDTO) (res *dto.PostCreateResponseDTO, err error) {
	userID := postCreateDTO.UserID
	title := postCreateDTO.Title
	content := postCreateDTO.Content
	isAnonymous := postCreateDTO.IsAnonymous

	// todo 可能还要判断用户是否被禁言——中间件实现

	post, err := m.Dao.CreatePost(userID, title, content, *isAnonymous)
	if err != nil {
		return
	}
	tmp, err := m.Dao.ConvertPostModelToOverviewDTO(post, userID)
	if err != nil {
		return
	}
	res = (*dto.PostCreateResponseDTO)(tmp)
	return
}

// Like 用户喜欢帖子
func (m *PostService) Like(postLikeRequestDTO *dto.PostLikeRequestDTO) (err error) {
	userID := postLikeRequestDTO.UserID
	postID := postLikeRequestDTO.PostID
	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.Like(userID, postID)
		if err == nil {
			// 喜欢成功
			return
		}

		if errors.Is(err, &custom_error.OptimisticLockError{}) {
			// 并发版本冲突，重试
			continue
		}
		return
	}
	return custom_error.NewOptimisticLockError()
}

// Collect 收藏帖子
func (m *PostService) Collect(postCollectRequestDTO *dto.PostCollectRequestDTO) (err error) {
	userID := postCollectRequestDTO.UserID
	postID := postCollectRequestDTO.PostID
	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.Collect(userID, postID)
		if err == nil {
			// 收藏成功
			return
		}

		if errors.Is(err, &custom_error.OptimisticLockError{}) {
			// 并发版本冲突，重试
			continue
		}
		return
	}
	return custom_error.NewOptimisticLockError()
}

func (m *PostService) List(postListRequestDTO *dto.PostListRequestDTO) (res *dto.PostListResponseDTO, err error) {
	page := postListRequestDTO.Page
	pageSize := postListRequestDTO.PageSize
	userID := postListRequestDTO.UserID

	// todo 如果要做个性化推荐的话，后面这里要考虑把 user_id 传入，在 List 服务里使用
	posts, total, err := m.Dao.List(page, pageSize)
	if err != nil {
		return
	}
	overviews, err := m.Dao.ConvertPostModelsToOverviewDTOs(posts, userID)
	if err != nil {
		return
	}
	res = &dto.PostListResponseDTO{
		Pagination: *dto.BuildPaginationInfo(total, page, pageSize),
		Posts:      overviews,
	}
	return
}
