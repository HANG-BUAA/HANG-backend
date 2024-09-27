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
                avatar, _ := m.Dao.GetCommentUserAvatar(comment)
                return avatar
            }(),
        },
        ReplyCommentID:     comment.ReplyCommentID,
        ReplyRootCommentID: comment.ReplyRootCommentID,
        ReplyUserName:      comment.ReplyUserName,
        Content:            comment.Content,
        IsAnonymous:        comment.IsAnonymous,
        IsReplyAnonymous:   comment.IsReplyAnonymous,
        CreatedAt:          comment.CreatedAt,
        UpdatedAt:          comment.UpdatedAt,
        DeletedAt:          comment.DeletedAt,
    }
    return
}
