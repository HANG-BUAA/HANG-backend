package dao

import (
	"HANG-backend/src/global"
	"HANG-backend/src/model"
	"HANG-backend/src/service/dto"
	"errors"
	"gorm.io/gorm"
)

var postDao *PostDao

type PostDao struct {
	BaseDao
}

func NewPostDao() *PostDao {
	if postDao == nil {
		postDao = &PostDao{
			NewBaseDao(),
		}
	}
	return postDao
}

// CreatePost 创建帖子
func (m *PostDao) CreatePost(userID uint, title string, content string, isAnonymous bool) (*model.Post, error) {
	// 检查用户是否存在
	var user model.User
	err := m.Orm.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return &model.Post{}, err
	}

	// todo 检查用户是否被禁言

	post := model.Post{
		UserID:      userID,
		Title:       title,
		Content:     content,
		IsAnonymous: isAnonymous,
	}
	if err := m.Orm.Create(&post).Error; err != nil {
		return &model.Post{}, err
	}
	return &post, nil
}

// LikePost 用户喜欢某个帖子
func (m *PostDao) LikePost(userID uint, postID uint) error {
	var postLike model.PostLike

	// 查询用户是否已经喜欢了该帖子
	err := m.Orm.Where("user_id = ? AND post_id = ?", userID, postID).First(&postLike).Error
	if err == nil {
		// 用户已经喜欢了该帖子
		return errors.New("liked post")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	newPostLike := model.PostLike{
		UserID: userID,
		PostID: postID,
	}
	if err := m.Orm.Create(&newPostLike).Error; err != nil {
		return err
	}
	return nil
}

func (m *PostDao) CollectPost(userID uint, postID uint) error {
	var postCollect model.PostCollect

	// 查询用户是否已经收藏了该帖子
	err := m.Orm.Where("user_id = ? AND post_id = ?", userID, postID).First(&postCollect).Error
	if err == nil {
		return errors.New("collected post")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	newPostCollect := model.PostCollect{
		UserID: userID,
		PostID: postID,
	}
	if err := m.Orm.Create(&newPostCollect).Error; err != nil {
		return err
	}
	return nil
}

// ListPostOverviews 查看帖子列表
func (m *PostDao) ListPostOverviews(page int, userID uint) (posts []dto.PostOverviewDTO, total int64, err error) {
	offset := (page - 1) * global.PageSize
	limit := global.PageSize

	// 查询总的帖子数（包括或不包括软删除，视你的需求而定）
	err = m.Orm.Model(&model.Post{}).Count(&total).Error
	if err != nil {
		return
	}

	// todo 超出分页数量检测

	// 子查询：获取点赞数、收藏数、评论数
	subQueryLike := m.Orm.Model(&model.PostLike{}).Select("post_id, COUNT(*) as like_num").Group("post_id")
	subQueryCollect := m.Orm.Model(&model.PostCollect{}).Select("post_id, COUNT(*) as collect_num").Group("post_id")
	subQueryComment := m.Orm.Model(&model.Comment{}).Select("post_id, COUNT(*) as comment_num").Group("post_id")

	// 子查询：判断当前用户是否喜欢该帖子
	subQueryUserLike := m.Orm.Model(&model.PostLike{}).
		Select("post_id, 1 as has_liked").
		Where("user_id = ?", userID).
		Group("post_id")

	// 子查询：判断当前用户是否收藏该帖子
	subQueryUserCollect := m.Orm.Model(&model.PostCollect{}).
		Select("post_id, 1 as has_collected").
		Where("user_id = ?", userID).
		Group("post_id")

	// 查询分页数据
	err = m.Orm.
		Table("post"). // 使用单数表名
		Select("post.id, post.user_id, user.user_name, user.avatar as user_avatar, post.title, post.content, post.is_anonymous, "+
			"COALESCE(like_num, 0) as like_num, COALESCE(collect_num, 0) as collect_num, COALESCE(comment_num, 0) as comment_num, "+
			"COALESCE(has_liked, 0) as has_liked, COALESCE(has_collected, 0) as has_collected").
		Joins("LEFT JOIN user ON user.id = post.user_id").
		Joins("LEFT JOIN (?) as post_like ON post_like.post_id = post.id", subQueryLike).
		Joins("LEFT JOIN (?) as post_collect ON post_collect.post_id = post.id", subQueryCollect).
		Joins("LEFT JOIN (?) as comment ON comment.post_id = post.id", subQueryComment).
		Joins("LEFT JOIN (?) as user_like ON user_like.post_id = post.id", subQueryUserLike).
		Joins("LEFT JOIN (?) as user_collect ON user_collect.post_id = post.id", subQueryUserCollect).
		Order("post.id desc").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	return
}
