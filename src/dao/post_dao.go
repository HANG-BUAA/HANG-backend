package dao

import (
	"HANG-backend/src/custom_error"
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

// Like 用户喜欢某个帖子
func (m *PostDao) Like(userID uint, postID uint) error {
	// 检查用户和帖子是否存在
	var (
		user model.User
		post model.Post
	)
	if err := m.Orm.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user does not exist")
		}
		return err
	}
	if err := m.Orm.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post does not exist")
		}
		return err
	}

	// 查询用户是否已经喜欢了该帖子
	var postLike model.PostLike
	err := m.Orm.Where("user_id = ? AND post_id = ?", userID, postID).First(&postLike).Error
	if err == nil {
		// 用户已经喜欢了该帖子
		return errors.New("liked post")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 使用事务保证操作原子性
	return m.Orm.Transaction(func(tx *gorm.DB) error {
		newPostLike := model.PostLike{
			UserID: userID,
			PostID: postID,
		}
		if err := tx.Create(&newPostLike).Error; err != nil {
			return err
		}

		// 动态维护 Post 表中的喜欢数字段，使用乐观锁防止并发状态下数据不一致的情况
		result := tx.Model(&model.Post{}).Where("id = ? AND like_version = ?", postID, post.LikeVersion).Updates(map[string]interface{}{
			"like_num":     post.LikeNum + 1,
			"like_version": post.LikeVersion + 1,
		})
		if result.RowsAffected == 0 {
			return custom_error.NewOptimisticLockError()
		}
		return nil
	})
}

// Collect 用户收藏帖子
func (m *PostDao) Collect(userID uint, postID uint) error {
	// 检查用户和帖子是否存在
	var (
		user model.User
		post model.Post
	)
	if err := m.Orm.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user does not exist")
		}
		return err
	}
	if err := m.Orm.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post does not exist")
		}
		return err
	}

	// 查询用户是否已经收藏了该帖子
	var postCollect model.PostCollect
	err := m.Orm.Where("user_id = ? AND post_id = ?", userID, postID).First(&postCollect).Error
	if err == nil {
		return errors.New("collected post")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 使用事务保证操作原子性
	return m.Orm.Transaction(func(tx *gorm.DB) error {
		newPostCollect := model.PostCollect{
			UserID: userID,
			PostID: postID,
		}
		if err := tx.Create(&newPostCollect).Error; err != nil {
			return err
		}

		// 动态维护 Post 表中的收藏数字段，使用乐观锁防止并发状态下数据不一致的情况
		result := tx.Model(&model.Post{}).Where("id = ? AND collect_version = ?", postID, post.CollectVersion).Updates(map[string]interface{}{
			"collect_num":     post.CollectNum + 1,
			"collect_version": post.CollectVersion + 1,
		})
		if result.RowsAffected == 0 {
			return custom_error.NewOptimisticLockError()
		}
		return nil
	})
}

func (m *PostDao) GetPostUserNameAndAvatar(post *model.Post) (string, string, error) {
	if post.IsAnonymous {
		return "洞主", "匿名的头像，还没想好", nil
	}
	var user model.User
	if err := m.Orm.Where("id = ?", post.UserID).First(&user).Error; err != nil {
		return "", "", err
	}
	return user.Username, user.Avatar, nil
}
