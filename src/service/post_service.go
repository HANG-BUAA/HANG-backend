package service

import (
	"HANG-backend/src/dao"
	"HANG-backend/src/global"
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

// List 查询帖子列表
func (m *PostService) List(postListRequestDTO *dto.PostListTRequestDTO) (res *dto.PostListTResponseDTO, err error) {
	userID := postListRequestDTO.UserID
	page := postListRequestDTO.Page

	// 目前的设定是该接口无法查询被删掉的帖子
	posts, total, err := m.Dao.ListPostOverviews(page, userID)
	if err != nil {
		return
	}

	// 过滤匿名信息
	for i := range posts {
		if posts[i].IsAnonymous {
			posts[i].UserAvatar = ""
			posts[i].UserID = 0
			posts[i].UserName = "匿名用户"
		}
	}

	res = &dto.PostListTResponseDTO{
		Posts: posts,
		Pagination: dto.PaginationInfo{
			TotalRecords: int(total),
			CurrentPage:  page,
			PageSize:     global.PageSize,
			TotalPages:   (int(total) + global.PageSize - 1) / global.PageSize,
		},
	}
	return
}
