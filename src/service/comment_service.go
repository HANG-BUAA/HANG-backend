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

func (m *CommentService) CreateComment(commentCreateRequestDTO *dto.CommentCreateRequestDTO) (res *dto.CommentCreateResponseDTO, err error) {
	userID := commentCreateRequestDTO.UserID
	postID := commentCreateRequestDTO.PostID
	replyTo := commentCreateRequestDTO.ReplyTo
	content := commentCreateRequestDTO.Content
	isAnonymous := commentCreateRequestDTO.IsAnonymous

	// todo 可能还要判断用户是否被禁言

	comment, err := m.Dao.CreateComment(userID, postID, *replyTo, content, *isAnonymous)
	if err != nil {
		return
	}
	res = &dto.CommentCreateResponseDTO{
		ID:          comment.ID,
		PostID:      comment.PostID,
		ReplyTo:     comment.ReplyTo,
		Content:     comment.Content,
		IsAnonymous: comment.IsAnonymous,
		CreatedAt:   comment.CreatedAt,
		UpdatedAt:   comment.UpdatedAt,
		DeletedAt:   comment.DeletedAt,
	}
	return
}

func (m *CommentService) Like(commentLikeRequestDTO *dto.CommentLikeRequestDTO) (err error) {
	userID := commentLikeRequestDTO.UserID
	commentID := commentLikeRequestDTO.CommentID
	err = m.Dao.LikeComment(userID, commentID)
	return
}

func (m *CommentService) List(commentListRequestDTO *dto.CommentListRequestDTO) (res *dto.CommentListResponseDTO, err error) {
	userID := commentListRequestDTO.UserID
	page := commentListRequestDTO.Page
	pageSize := commentListRequestDTO.PageSize
	postID := commentListRequestDTO.PostID

	// 评论是可以查询被删除的
	comments, total, err := m.Dao.ListCommentOverviews(postID, userID, page, pageSize)
	if err != nil {
		return
	}

	// 过滤信息
	for i := range comments {
		// 如果是匿名评论
		if comments[i].IsAnonymous {
			// todo 假名微服务，也可以考虑直接存到comment表里
			comments[i].Username = "匿名用户"
			comments[i].UserID = 0
			comments[i].UserAvatar = ""
		}

		// 如果已经被删除
		if comments[i].DeletedAt.Valid {
			comments[i].Content = "该评论已被删除"
		}

		// 回复的帖子/评论是否是匿名
		if comments[i].IsReplyToAnonymous {
			if comments[i].ReplyToID == 0 {
				comments[i].ReplyToName = "楼主"
			} else {
				// todo 假名微服务，也可以考虑直接存到comment表里
				comments[i].ReplyToName = "匿名用户（暂时没有分配假名）"
			}
		}
	}

	res = &dto.CommentListResponseDTO{
		Comments: comments,
		Pagination: dto.PaginationInfo{
			TotalRecords: int(total),
			CurrentPage:  page,
			PageSize:     pageSize,
			TotalPages:   (int(total) + pageSize - 1) / pageSize,
		},
	}
	return
}
