package dao

import (
	"HANG-backend/src/model"
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
