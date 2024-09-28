package service

import (
	"HANG-backend/src/custom_error"
	"HANG-backend/src/dao"
	"HANG-backend/src/global"
	"HANG-backend/src/service/dto"
	"errors"
)

var commentService *CommentService

type CommentService struct {
	BaseService
	Dao *dao.CommentDao
}

func NewCommentService() *CommentService {
	if commentService == nil {
		commentService = &CommentService{
			Dao: dao.NewCommentDao(),
		}
	}
	return commentService
}

// Create 创建评论
func (m *CommentService) Create(commentCreateDTO *dto.CommentCreateRequestDTO) (res *dto.CommentCreateResponseDTO, err error) {
	userID := commentCreateDTO.UserID
	postID := commentCreateDTO.PostID
	replyCommentID := commentCreateDTO.ReplyCommentID
	content := commentCreateDTO.Content
	isAnonymous := commentCreateDTO.IsAnonymous

	// todo 判断用户是否被禁言——可以加中间件
	comment, err := m.Dao.CreateComment(userID, postID, *replyCommentID, content, *isAnonymous)
	if err != nil {
		return
	}

	tmp, err := m.Dao.ConvertCommentModelToOverviewDTO(comment, userID)
	if err != nil {
		return nil, err
	}
	res = (*dto.CommentCreateResponseDTO)(tmp)
	return
}

// Like 用户喜欢评论
func (m *CommentService) Like(commentLikeRequestDTO *dto.CommentLikeRequestDTO) (err error) {
	userID := commentLikeRequestDTO.UserID
	commentID := commentLikeRequestDTO.CommentID
	// todo 参数校验的工作应当放到 service 层
	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.Like(userID, commentID)
		if err == nil {
			return
		}

		if errors.Is(err, &custom_error.OptimisticLockError{}) {
			continue
		}
	}
	return custom_error.NewOptimisticLockError()
}

// ListFirstLevel 列出某帖子下一级评论列表
func (m *CommentService) ListFirstLevel(commentListRequestDTO *dto.CommentListRequestDTO) (res *dto.CommentListResponseDTO, err error) {
	page := commentListRequestDTO.Page
	pageSize := commentListRequestDTO.PageSize
	userID := commentListRequestDTO.UserID
	postID := commentListRequestDTO.PostID

	// 校验帖子是否存在
	_, ok := m.Dao.CheckPostExist(postID)
	if !ok {
		return nil, errors.New("post not exist")
	}

	comments, total, err := m.Dao.ListFirstLevel(postID, page, pageSize)
	if err != nil {
		return
	}
	overviews, err := m.Dao.ConvertCommentModelsToOverviewDTOs(comments, userID)
	if err != nil {
		return
	}

	res = &dto.CommentListResponseDTO{
		Pagination: *dto.BuildPaginationInfo(total, page, pageSize),
		Comments:   overviews,
	}
	return
}

// ListSecondLevel 列出某一级评论下二级评论列表
func (m *CommentService) ListSecondLevel(commentListRequestDTO *dto.CommentListRequestDTO) (res *dto.CommentListResponseDTO, err error) {
	page := commentListRequestDTO.Page
	pageSize := commentListRequestDTO.PageSize
	userID := commentListRequestDTO.UserID
	commentID := commentListRequestDTO.CommentID

	// 校验一级评论是否存在
	_, ok := m.Dao.CheckCommentExist(commentID)
	if !ok {
		return nil, errors.New("comment not exist")
	}

	comments, total, err := m.Dao.ListSecondLevel(commentID, page, pageSize)
	if err != nil {
		return
	}
	overviews, err := m.Dao.ConvertCommentModelsToOverviewDTOs(comments, userID)
	if err != nil {
		return
	}

	res = &dto.CommentListResponseDTO{
		Pagination: *dto.BuildPaginationInfo(total, page, pageSize),
		Comments:   overviews,
	}
	return
}
