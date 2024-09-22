package service

import (
	"HANG-backend/src/dao"
	"HANG-backend/src/service/dto"
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

func (m *CommentService) CreateComment(iCommentCreateRequestDTO *dto.CommentCreateRequestDTO) (res *dto.CommentCreateResponseDTO, err error) {
	iUserID := iCommentCreateRequestDTO.UserID
	iPostID := iCommentCreateRequestDTO.PostID
	iReplyTo := iCommentCreateRequestDTO.ReplyTo
	iContent := iCommentCreateRequestDTO.Content
	iIsAnonymous := iCommentCreateRequestDTO.IsAnonymous

	// todo 可能还要判断用户是否被禁言

	iComment, err := m.Dao.CreateComment(iUserID, iPostID, *iReplyTo, iContent, *iIsAnonymous)
	if err != nil {
		return
	}
	res = &dto.CommentCreateResponseDTO{
		ID:          iComment.ID,
		PostID:      iComment.PostID,
		ReplyTo:     iComment.ReplyTo,
		Content:     iComment.Content,
		IsAnonymous: iComment.IsAnonymous,
		DisplayName: iComment.DisplayName,
		CreatedAt:   iComment.CreatedAt,
		UpdatedAt:   iComment.UpdatedAt,
		DeletedAt:   iComment.DeletedAt,
	}
	return
}
