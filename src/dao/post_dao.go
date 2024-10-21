package dao

import (
	"HANG-backend/src/custom_error"
	"HANG-backend/src/model"
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"gorm.io/gorm"
	"time"
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

func (m *PostDao) ConvertPostModelsToOverviewDTOs(posts []model.Post, userID uint) ([]dto.PostOverviewDTO, error) {
	res := make([]dto.PostOverviewDTO, 0)
	for _, post := range posts {
		tmp, err := m.ConvertPostModelToOverviewDTO(&post, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, *tmp)
	}
	return res, nil
}

func (m *PostDao) ConvertPostModelToOverviewDTO(post *model.Post, userID uint) (*dto.PostOverviewDTO, error) {
	// 获取帖子作者信息
	userName, userAvatar, err := m.getPostUserNameAndAvatar(post)
	if err != nil {
		return nil, err
	}

	// 计算帖子的评论数
	var commentNum int64
	if err := m.Orm.Model(&model.Comment{}).Where("post_id = ?", post.ID).Count(&commentNum).Error; err != nil {
		return nil, err
	}

	var postLike model.PostLike
	var postCollect model.PostCollect

	return &dto.PostOverviewDTO{
		ID: post.ID,
		Author: dto.PostAuthorDTO{
			UserID:     utils.IfThenElse(post.IsAnonymous, uint(0), userID).(uint),
			UserName:   userName,
			UserAvatar: userAvatar,
		},
		Title:        post.Title,
		Content:      post.Content,
		IsAnonymous:  post.IsAnonymous,
		LikeNum:      post.LikeNum,
		CollectNum:   post.CollectNum,
		HasCollected: m.Orm.Where("post_id = ? AND user_id = ?", post.ID, userID).First(&postCollect).Error == nil,
		HasLiked:     m.Orm.Where("post_id = ? AND user_id = ?", post.ID, userID).First(&postLike).Error == nil,
		CommentNum:   int(commentNum),
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
		DeletedAt:    post.DeletedAt,
	}, nil
}

// CreatePost 创建帖子
func (m *PostDao) CreatePost(user *model.User, title string, content string, isAnonymous bool) (*model.Post, error) {
	post := model.Post{
		UserID:      user.ID,
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
func (m *PostDao) Like(user *model.User, post *model.Post) error {
	// 使用事务保证操作原子性
	return m.Orm.Transaction(func(tx *gorm.DB) error {
		newPostLike := model.PostLike{
			UserID: user.ID,
			PostID: post.ID,
		}
		if err := tx.Create(&newPostLike).Error; err != nil {
			return err
		}

		// 动态维护 Post 表中的喜欢数字段，使用乐观锁防止并发状态下数据不一致的情况
		result := tx.Model(&model.Post{}).Where("id = ? AND like_version = ?", post.ID, post.LikeVersion).Updates(map[string]interface{}{
			"like_num":     post.LikeNum + 1,
			"like_version": post.LikeVersion + 1,
		})
		if result.RowsAffected == 0 {
			return custom_error.NewOptimisticLockError()
		}
		return nil
	})
}

func (m *PostDao) Unlike(user *model.User, post *model.Post) error {
	return m.Orm.Transaction(func(tx *gorm.DB) error {
		tx.Where("user_id = ? AND post_id = ?", user.ID, post.ID).Delete(&model.PostLike{})

		result := tx.Model(&model.Post{}).Where("id = ? AND like_version = ?", post.ID, post.LikeVersion).Updates(map[string]interface{}{
			"like_num":     post.LikeNum - 1,
			"like_version": post.LikeVersion + 1,
		})
		if result.RowsAffected == 0 {
			return custom_error.NewOptimisticLockError()
		}
		return nil
	})
}

// Collect 用户收藏帖子
func (m *PostDao) Collect(user *model.User, post *model.Post) error {
	// 使用事务保证操作原子性
	return m.Orm.Transaction(func(tx *gorm.DB) error {
		newPostCollect := model.PostCollect{
			UserID: user.ID,
			PostID: post.ID,
		}
		if err := tx.Create(&newPostCollect).Error; err != nil {
			return err
		}

		// 动态维护 Post 表中的收藏数字段，使用乐观锁防止并发状态下数据不一致的情况
		result := tx.Model(&model.Post{}).Where("id = ? AND collect_version = ?", post.ID, post.CollectVersion).Updates(map[string]interface{}{
			"collect_num":     post.CollectNum + 1,
			"collect_version": post.CollectVersion + 1,
		})
		if result.RowsAffected == 0 {
			return custom_error.NewOptimisticLockError()
		}
		return nil
	})
}

func (m *PostDao) Uncollect(user *model.User, post *model.Post) error {
	return m.Orm.Transaction(func(tx *gorm.DB) error {
		tx.Where("user_id = ? AND post_id = ?", user.ID, post.ID).Delete(&model.PostCollect{})

		result := tx.Model(&model.Post{}).Where("id = ? AND collect_version = ?", post.ID, post.CollectVersion).Updates(map[string]interface{}{
			"collect_num":     post.CollectNum - 1,
			"collect_version": post.CollectVersion + 1,
		})
		if result.RowsAffected == 0 {
			return custom_error.NewOptimisticLockError()
		}
		return nil
	})
}

func (m *PostDao) CommonList(cursor uint, pageSize int) ([]model.Post, int, bool, error) {
	query := m.Orm.Model(&model.Post{})

	// 先计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, false, err
	}

	// 多查一条记录出来，判断是否到有下一页
	var posts []model.Post
	query = query.
		Limit(pageSize + 1).
		Order("id desc")
	if cursor != 0 {
		query = query.Where("id < ?", cursor)
	}
	if err := query.Find(&posts).Error; err != nil {
		return nil, 0, false, err
	}

	isEnd := len(posts) < pageSize+1
	return posts[:utils.IfThenElse(isEnd, len(posts), pageSize).(int)], int(total), isEnd, nil
}

// CheckLiked 判断用户是否已经喜欢该帖子
func (m *PostDao) CheckLiked(user *model.User, post *model.Post) bool {
	var postLike model.PostLike
	if err := m.Orm.Where("user_id = ? AND post_id = ?", user.ID, post.ID).First(&postLike).Error; err != nil {
		return false
	}
	return true
}

// GetListsByIDs 根据 id 列表查出记录
func (m *PostDao) GetListsByIDs(ids []uint) ([]model.Post, error) {
	var posts []model.Post
	if err := m.Orm.Where("id IN (?)", ids).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// CheckCollected 判断用户是否已经收藏该帖子
func (m *PostDao) CheckCollected(user *model.User, post *model.Post) bool {
	var postCollect model.PostCollect
	if err := m.Orm.Where("user_id = ? AND post_id = ?", user.ID, post.ID).First(&postCollect).Error; err != nil {
		return false
	}
	return true
}

func (m *PostDao) getPostUserNameAndAvatar(post *model.Post) (string, string, error) {
	if post.IsAnonymous {
		return "洞主", "匿名的头像，还没想好", nil
	}
	var user model.User
	if err := m.Orm.Where("id = ?", post.UserID).First(&user).Error; err != nil {
		return "", "", err
	}
	return user.Username, user.Avatar, nil
}

func (m *PostDao) GetCollections(user *model.User, cursor time.Time, pageSize int) ([]model.Post, int, bool, error) {
	// 基础查询：获取用户收藏的 Post，按 created_at 排序
	query := m.Orm.
		Model(&model.Post{}).
		Joins("JOIN post_collect ON post_collect.post_id = post.id").
		Where("post_collect.user_id = ?", user.ID).
		Order("post_collect.created_at desc")

	// 先计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, false, err
	}

	// 添加分页和 cursor 条件
	if !cursor.IsZero() {
		query = query.Where("post_collect.created_at < ?", cursor)
	}

	// 多查一条出来，为了判断当前是否到达了最后一页
	var posts []model.Post
	if err := query.Limit(pageSize + 1).Find(&posts).Error; err != nil {
		return nil, 0, false, err
	}

	isEnd := len(posts) < pageSize+1

	return posts[:utils.IfThenElse(isEnd, len(posts), pageSize).(int)], int(total), isEnd, nil
}

func (m *PostDao) GetCollectCursorByID(user *model.User, post *model.Post) (time.Time, error) {
	var postCollect model.PostCollect
	if err := m.Orm.Where("user_id = ? AND post_id = ?", user.ID, post.ID).First(&postCollect).Error; err != nil {
		return time.Time{}, err
	}
	return postCollect.CreatedAt, nil
}
