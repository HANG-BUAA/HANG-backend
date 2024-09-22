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
func (m *PostService) CreatePost(iPostCreateDTO *dto.PostCreateRequestDTO) (res *dto.PostCreateResponseDTO, err error) {
	iUserID := iPostCreateDTO.UserID
	iTitle := iPostCreateDTO.Title
	iContent := iPostCreateDTO.Content
	iIsAnonymous := iPostCreateDTO.IsAnonymous

	// todo 可能还要判断用户是否被禁言

	iPost, err := m.Dao.CreatePost(iUserID, iTitle, iContent, *iIsAnonymous)
	if err != nil {
		return
	}
	res = &dto.PostCreateResponseDTO{
		ID:          iPost.ID,
		UserID:      iPost.UserID,
		Title:       iPost.Title,
		Content:     iPost.Content,
		IsAnonymous: iPost.IsAnonymous,
		CreatedAt:   iPost.CreatedAt,
		UpdatedAt:   iPost.UpdatedAt,
		DeletedAt:   iPost.DeletedAt,
	}
	return
}

// Like 用户喜欢帖子
func (m *PostService) Like(iPostLikeRequestDTO *dto.PostLikeRequestDTO) (err error) {
	iUserID := iPostLikeRequestDTO.UserID
	iPostID := iPostLikeRequestDTO.PostID
	err = m.Dao.LikePost(iUserID, iPostID)
	return
}

// Collect 收藏帖子
func (m *PostService) Collect(iPostCollectRequestDTO *dto.PostCollectRequestDTO) (err error) {
	iUserID := iPostCollectRequestDTO.UserID
	iPostID := iPostCollectRequestDTO.PostID
	err = m.Dao.CollectPost(iUserID, iPostID)
	return
}

// List 查询帖子列表
func (m *PostService) List(iPostListRequestDTO *dto.PostListTRequestDTO) (res dto.PostListTResponseDTO, err error) {
	iUserID := iPostListRequestDTO.UserID
	iPage := iPostListRequestDTO.Page

	// todo 判断用户是否有查看被删掉的帖子的权限

	posts, total, err := m.Dao.ListPostOverviews(iPage, iUserID)
	if err != nil {
		return
	}

	// todo 查看用户是否有看匿名帖子的楼主信息的权限
	if true {
		for i := range posts {
			if posts[i].IsAnonymous {
				posts[i].UserAvatar = ""
				posts[i].UserID = 0
				posts[i].UserName = "匿名用户"
			}
		}
	}
	res = dto.PostListTResponseDTO{
		Posts: posts,
		Pagination: dto.PaginationInfo{
			TotalRecords: int(total),
			CurrentPage:  iPage,
			PageSize:     global.PageSize,
			TotalPages:   (int(total) + global.PageSize - 1) / global.PageSize,
		},
	}
	return
}
