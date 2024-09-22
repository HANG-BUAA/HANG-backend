package dao

import (
    "HANG-backend/src/model"
    "errors"
)

var commentDao *CommentDao

type CommentDao struct {
    BaseDao
}

func NewCommentDao() *CommentDao {
    if commentDao == nil {
        commentDao = &CommentDao{
            NewBaseDao(),
        }
    }
    return commentDao
}

// CreateComment 创建评论
func (m *CommentDao) CreateComment(iUserID uint, iPostID uint, iReplyTo uint, iContent string, iIsAnonymous bool) (*model.Comment, error) {
    // 检测用户是否存在
    var iUser model.User
    err := m.Orm.Where("id = ?", iUserID).First(&iUser).Error
    if err != nil {
        return &model.Comment{}, err
    }

    // 检测帖子是否存在
    var iPost model.Post
    err = m.Orm.Where("id = ?", iPostID).First(&iPost).Error
    if err != nil {
        return &model.Comment{}, err
    }

    // 检测回复的评论是否存在
    var iComment model.Comment
    if iReplyTo != 0 {
        err = m.Orm.Where("id = ?", iReplyTo).First(&iComment).Error
        if err != nil {
            return &model.Comment{}, err
        }

        // 检测回复的评论是否属于该帖子
        if iComment.PostID != iPost.ID {
            return &model.Comment{}, errors.New("comment_id does not match the post")
        }
    }

    iComment = model.Comment{
        PostID:      iPostID,
        ReplyTo:     iReplyTo,
        UserID:      iUserID,
        Content:     iContent,
        IsAnonymous: iIsAnonymous,
    }

    if err := m.Orm.Create(&iComment).Error; err != nil {
        return &model.Comment{}, err
    }
    return &iComment, nil
}

//func (m *CommentDao) List(postID uint, page int)  {
//
//}
