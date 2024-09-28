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

	res = &dto.CommentCreateResponseDTO{
		ID:     comment.ID,
		PostID: comment.PostID,
		Author: dto.CommentAuthorDTO{
			UserID:   comment.UserID,
			UserName: comment.UserName,
			UserAvatar: func() string {
				_, avatar, _ := m.Dao.GetCommentUserNameAndAvatar(comment)
				return avatar
			}(),
		},
		ReplyCommentID:     comment.ReplyCommentID,
		ReplyRootCommentID: comment.ReplyRootCommentID,
		ReplyUserName:      comment.ReplyUserName,
		Content:            comment.Content,
		LikeNum:            0, // 刚创建的评论，默认没有喜欢数
		IsLiked:            false,
		IsAnonymous:        comment.IsAnonymous,
		IsReplyAnonymous:   comment.IsReplyAnonymous,
		CreatedAt:          comment.CreatedAt,
		UpdatedAt:          comment.UpdatedAt,
		DeletedAt:          comment.DeletedAt,
	}
	return
}

// Like 用户喜欢评论
func (m *CommentService) Like(commentLikeRequestDTO *dto.CommentLikeRequestDTO) (err error) {
	userID := commentLikeRequestDTO.UserID
	commentID := commentLikeRequestDTO.CommentID
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
