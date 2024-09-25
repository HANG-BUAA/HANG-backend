package dao

import (
	"HANG-backend/src/model"
	"HANG-backend/src/service/dto"
	"errors"
	"gorm.io/gorm"
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
func (m *CommentDao) CreateComment(userID uint, postID uint, replyTo uint, content string, isAnonymous bool) (*model.Comment, error) {
	// 检测用户是否存在
	var user model.User
	err := m.Orm.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return &model.Comment{}, err
	}

	// 检测帖子是否存在
	var post model.Post
	err = m.Orm.Where("id = ?", postID).First(&post).Error
	if err != nil {
		return &model.Comment{}, err
	}

	// 检测回复的评论是否存在
	var comment model.Comment
	if replyTo != 0 {
		err = m.Orm.Where("id = ?", replyTo).First(&comment).Error
		if err != nil {
			return &model.Comment{}, err
		}

		// 检测回复的评论是否属于该帖子
		if comment.PostID != post.ID {
			return &model.Comment{}, errors.New("comment_id does not match the post")
		}
	}

	comment = model.Comment{
		PostID:      postID,
		ReplyTo:     replyTo,
		UserID:      userID,
		Username:    user.Username,
		Content:     content,
		IsAnonymous: isAnonymous,
	}

	if err := m.Orm.Create(&comment).Error; err != nil {
		return &model.Comment{}, err
	}
	return &comment, nil
}

func (m *CommentDao) LikeComment(userID uint, commentID uint) error {
	var commentLike model.CommentLike

	// 查询用户是否已经喜欢了该评论
	err := m.Orm.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&commentLike).Error
	if err == nil {
		// 用户已经喜欢过该评论
		return errors.New("liked comment")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	newCommentLike := model.CommentLike{
		UserID:    userID,
		CommentID: commentID,
	}
	if err := m.Orm.Create(&newCommentLike).Error; err != nil {
		return err
	}
	return nil
}

func (m *CommentDao) ListCommentOverviews(postID uint, userID uint, page int, pageSize int) (commentOverviews []dto.CommentOverviewDTO, total int64, err error) {
	var comments []model.Comment
	offset := (page - 1) * pageSize

	// todo 超出分页页数限制

	// 查询包括已经被删除的评论
	if err = m.Orm.Model(&model.Comment{}).
		Unscoped().
		Where("post_id = ?", postID).
		Offset(offset).
		Limit(pageSize).
		Order("id desc").
		Find(&comments).Error; err != nil {
		return
	}

	for _, comment := range comments {
		// 查询当前 postID 下评论总数
		m.Orm.Model(&model.Comment{}).Unscoped().Where("post_id = ?", postID).Count(&total)

		// 统计喜欢数
		var likeCount int64
		m.Orm.Model(&model.CommentLike{}).Where("comment_id = ?", comment.ID).Count(&likeCount)

		// 获取用户
		var user model.User
		if err = m.Orm.First(&user, comment.UserID).Error; err != nil {
			return
		}

		// 检查当前用户是否已经喜欢该评论
		var hasLiked bool
		if err = m.Orm.Model(&model.CommentLike{}).
			Where("user_id = ? AND comment_id = ?", userID, comment.ID).
			First(&model.CommentLike{}).Error; err == nil {
			hasLiked = true
		}

		var (
			replyToName        string
			isReplyToAnonymous bool
		)

		// 判断 ReplyToID 是否为 0，进而判断是直接回复帖子还是回复评论
		if comment.ReplyTo != 0 {
			var replyToComment model.Comment
			if err = m.Orm.Unscoped().First(&replyToComment, comment.ReplyTo).Error; err == nil {
				isReplyToAnonymous = replyToComment.IsAnonymous
				var replyToUser model.User
				if err = m.Orm.First(&replyToUser, replyToComment.UserID).Error; err == nil {
					replyToName = replyToUser.Username
				}
			}
		} else {
			var post model.Post
			if err = m.Orm.Unscoped().First(&post, comment.PostID).Error; err == nil {
				isReplyToAnonymous = post.IsAnonymous
				var postUser model.User
				if err = m.Orm.First(&postUser, post.UserID).Error; err == nil {
					replyToName = postUser.Username
				}
			}
		}

		commentOverviews = append(commentOverviews, dto.CommentOverviewDTO{
			ID:                 comment.ID,
			PostID:             comment.PostID,
			ReplyToID:          comment.ReplyTo,
			ReplyToName:        replyToName,
			IsReplyToAnonymous: isReplyToAnonymous,
			UserID:             comment.UserID,
			Username:           user.Username,
			UserAvatar:         user.Avatar,
			Content:            comment.Content,
			IsAnonymous:        comment.IsAnonymous,
			LikeNum:            int(likeCount),
			HasLiked:           hasLiked,
			CreatedAt:          comment.CreatedAt,
			UpdatedAt:          comment.UpdatedAt,
			DeletedAt:          comment.DeletedAt,
		})
	}

	return
}
