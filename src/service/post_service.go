package service

import (
	"HANG-backend/src/dao"
	"HANG-backend/src/service/dto"
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

// CreatePost 创建帖子
func (m *PostService) CreatePost(postCreateDTO *dto.PostCreateRequestDTO) (res *dto.PostCreateResponseDTO, err error) {
	userID := postCreateDTO.UserID
	title := postCreateDTO.Title
	content := postCreateDTO.Content
	isAnonymous := postCreateDTO.IsAnonymous

	// todo 可能还要判断用户是否被禁言

	post, err := m.Dao.CreatePost(userID, title, content, *isAnonymous)
	if err != nil {
		return
	}
	res = &dto.PostCreateResponseDTO{
		ID:          post.ID,
		UserID:      post.UserID,
		Title:       post.Title,
		Content:     post.Content,
		IsAnonymous: post.IsAnonymous,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		DeletedAt:   post.DeletedAt,
	}
	return
}

// Like 用户喜欢帖子
func (m *PostService) Like(postLikeRequestDTO *dto.PostLikeRequestDTO) (err error) {
	userID := postLikeRequestDTO.UserID
	postID := postLikeRequestDTO.PostID
	err = m.Dao.LikePost(userID, postID)
	return
}

// Collect 收藏帖子
func (m *PostService) Collect(postCollectRequestDTO *dto.PostCollectRequestDTO) (err error) {
	userID := postCollectRequestDTO.UserID
	postID := postCollectRequestDTO.PostID
	err = m.Dao.CollectPost(userID, postID)
	return
}