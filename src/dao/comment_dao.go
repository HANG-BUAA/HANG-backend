package dao

import (
	"HANG-backend/src/custom_error"
	"HANG-backend/src/model"
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"sync"
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

func (m *CommentDao) ConvertCommentModelsToOverviewDTOs(comments []model.Comment, userID uint) ([]dto.CommentOverviewDTO, error) {
	res := make([]dto.CommentOverviewDTO, 0)
	for _, comment := range comments {
		tmp, err := m.ConvertCommentModelToOverviewDTO(&comment, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, *tmp)
	}
	return res, nil
}

func (m *CommentDao) ConvertCommentModelToOverviewDTO(comment *model.Comment, userID uint) (*dto.CommentOverviewDTO, error) {
	userName, userAvatar, err := m.getCommentUserNameAndAvatar(comment)
	if err != nil {
		return nil, err
	}

	// 计算回复的人的名字
	replyUserName, err := m.getReplyUserName(comment)
	if err != nil {
		return nil, err
	}

	var commentLike model.CommentLike

	return &dto.CommentOverviewDTO{
		ID:     comment.ID,
		PostID: comment.PostID,
		Author: dto.CommentAuthorDTO{
			UserID:     utils.IfThenElse(comment.IsAnonymous, uint(0), comment.UserID).(uint),
			UserName:   userName,
			UserAvatar: userAvatar,
		},
		ReplyCommentID:     comment.ReplyCommentID,
		ReplyRootCommentID: comment.ReplyRootCommentID,
		ReplyUserName:      replyUserName,
		Content:            comment.Content,
		LikeNum:            comment.LikeNum,
		HasLiked:           m.Orm.Where("comment_id = ? AND user_id = ?", comment.ID, userID).First(&commentLike).Error == nil,
		IsAnonymous:        comment.IsAnonymous,
		IsReplyAnonymous:   comment.IsReplyAnonymous,
		CreatedAt:          comment.CreatedAt,
		UpdatedAt:          comment.UpdatedAt,
		DeletedAt:          comment.DeletedAt,
	}, nil
}

func (m *CommentDao) CreateComment(userID uint, postID uint, replyCommentID uint, content string, isAnonymous bool) (*model.Comment, error) {
	// todo 这个校验放到 service 层
	// 检查用户、帖子、回复的评论是否存在，以及合法性
	var (
		user    model.User
		post    model.Post
		comment model.Comment
	)
	if err := m.Orm.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	if err := m.Orm.Where("id = ?", postID).First(&post).Error; err != nil {
		return nil, err
	}
	if replyCommentID != 0 {
		// 是二级评论，要检查回复的评论是否 match 当前的帖子
		if err := m.Orm.Where("id = ?", replyCommentID).First(&comment).Error; err != nil {
			return nil, err
		}
		if comment.PostID != postID {
			return nil, errors.New("comment doesn't match the post")
		}
	}
	var (
		isReplyAnonymous   bool
		replyRootCommentID uint
		replyUserName      string
		err                error
	)
	if replyCommentID == 0 {
		// 一级评论，其 root_id 由 Comment 的 AfterCreate 钩子生成
		replyUserName, err = m.getPostUserName(&post)
		if err != nil {
			return nil, err
		}
		isReplyAnonymous = post.IsAnonymous
	} else {
		// 二级评论
		replyRootCommentID = comment.ReplyRootCommentID
		replyUserName = comment.UserName // 这里不需要查的原因是：如果回复的对象是匿名的，这样就是正确的；否则该字段没有意义，随便记录一个都可以
		isReplyAnonymous = comment.IsAnonymous
	}

	// 分配“假名”——上锁
	lw := getMutex(postID)
	lw.mu.Lock()
	userName, err := m.allocateCommenterName(&user, &post, isAnonymous)
	if err != nil {
		releaseMutex(postID, lw)
		return nil, err
	}

	newComment := &model.Comment{
		PostID:             postID,
		UserID:             userID,
		UserName:           userName,
		ReplyCommentID:     replyCommentID,
		ReplyRootCommentID: replyRootCommentID,
		ReplyUserName:      replyUserName,
		Content:            content,
		IsAnonymous:        isAnonymous,
		IsReplyAnonymous:   isReplyAnonymous,
	}
	if err = m.Orm.Create(newComment).Error; err != nil {
		releaseMutex(postID, lw)
		return nil, err
	}
	releaseMutex(postID, lw)
	return newComment, nil
}

// 或者评论回复的评论/帖子的展示名
func (m *CommentDao) getReplyUserName(comment *model.Comment) (string, error) {
	if comment.ReplyCommentID == 0 {
		// 一级评论
		var post model.Post
		if err := m.Orm.Where("id = ?", comment.PostID).First(&post).Error; err != nil {
			return "", err
		}
		name, err := m.getPostUserName(&post)
		if err != nil {
			return "", err
		}
		return name, nil
	} else {
		// 二级评论
		var com model.Comment
		if err := m.Orm.Where("id = ?", comment.ReplyCommentID).First(&com).Error; err != nil {
			return "", err
		}
		name, err := m.getCommentUserName(&com)
		if err != nil {
			return "", err
		}
		return name, nil
	}
}

// 获取贴主的展示名
func (m *CommentDao) getPostUserName(post *model.Post) (string, error) {
	if post.IsAnonymous {
		return "洞主", nil
	}
	var user model.User
	if err := m.Orm.Where("id = ?", post.UserID).First(&user).Error; err != nil {
		return "", err
	}
	return user.Username, nil
}

// 获取评论作者的展示名
func (m *CommentDao) getCommentUserName(comment *model.Comment) (string, error) {
	if comment.IsAnonymous {
		return comment.UserName, nil
	}
	var user model.User
	if err := m.Orm.Where("id = ?", comment.UserID).First(&user).Error; err != nil {
		return "", err
	}
	return user.Username, nil
}

func (m *CommentDao) getCommentUserNameAndAvatar(comment *model.Comment) (string, string, error) {
	if comment.IsAnonymous {
		return comment.UserName, "匿名的头像的路径，还没想好", nil
	}
	var user model.User
	if err := m.Orm.Where("id = ?", comment.UserID).First(&user).Error; err != nil {
		return "", "", err
	}
	return user.Username, user.Avatar, nil
}

// 分配名字
func (m *CommentDao) allocateCommenterName(user *model.User, post *model.Post, isAnonymous bool) (string, error) {
	// 先判断是否是匿名贴，不是匿名的话这里直接返回用户名
	if !isAnonymous {
		return user.Username, nil
	}
	// 判断该用户是是洞主
	if user.ID == post.UserID {
		return "洞主", nil
	}

	// 判断该用户是否在该帖子下匿名评论过
	var comment model.Comment
	if err := m.Orm.Unscoped().Where("user_id = ? AND post_id = ? AND is_anonymous = ?", user.ID, post.ID, true).First(&comment).Error; err != nil {
		// 该用户没有匿名评论过
		var count int64
		err = m.Orm.Unscoped().Model(&model.Comment{}).Where("user_id = ? AND post_id = ? AND is_anonymous = ?", user.ID, post.ID, true).
			Distinct("user_id").Count(&count).Error
		if err != nil {
			return "", err
		}
		return "匿名用户" + strconv.Itoa(int(count)+1), nil
	} else {
		return comment.UserName, nil
	}
}

type LockWrapper struct {
	mu     sync.Mutex
	refCnt int // 引用计数，记录有多少人正在该帖子下评论
}

var mutexMap sync.Map

// 获取或创建一个新的锁，并增加引用计数
func getMutex(postID uint) *LockWrapper {
	lockInterface, ok := mutexMap.Load(postID)
	if !ok {
		// 如果没有找到锁，创建一个新的锁
		lockInterface, _ = mutexMap.LoadOrStore(postID, &LockWrapper{
			refCnt: 0, // 初始引用计数为0
		})
	}

	lw := lockInterface.(*LockWrapper)
	lw.refCnt++ // 增加引用计数
	return lw
}

// 释放锁，并在没有等待者时删除锁
func releaseMutex(postID uint, lw *LockWrapper) {
	lw.mu.Unlock() // 释放锁

	// 减少引用计数
	lw.refCnt--

	// 如果没有 Goroutine 在使用该锁，将其从 Map 中删除
	if lw.refCnt == 0 {
		mutexMap.Delete(postID)
		fmt.Printf("Lock for key %d has been deleted\n", postID)
	}
}

// Like 用户喜欢评论
func (m *CommentDao) Like(userID, commentID uint) error {
	// 检查用户和评论是否存在
	var (
		user    model.User
		comment model.Comment
	)
	if err := m.Orm.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user does not exist")
		}
		return err
	}
	if err := m.Orm.Where("id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("comment does not exist")
		}
		return err
	}

	// 查询用户是否已经喜欢了该帖子
	var commentLike model.CommentLike
	err := m.Orm.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&commentLike).Error
	if err == nil {
		return errors.New("liked comment")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 使用事务保证操作原子性
	return m.Orm.Transaction(func(tx *gorm.DB) error {
		newCommentLike := model.CommentLike{
			CommentID: commentID,
			UserID:    userID,
		}
		if err := tx.Create(&newCommentLike).Error; err != nil {
			return err
		}
		// 动态维护 Comment 表中的喜欢数字段，使用乐观锁防止在并发状态下数据不一致
		result := tx.Model(&model.Comment{}).Where("id = ? AND like_version = ?", commentID, comment.LikeVersion).Updates(map[string]interface{}{
			"like_num":     comment.LikeNum + 1,
			"like_version": comment.LikeVersion + 1,
		})
		if result.RowsAffected == 0 {
			return custom_error.NewOptimisticLockError()
		}
		return nil
	})
}

// ListFirstLevel 列出某个帖子下一级评论列表
func (m *CommentDao) ListFirstLevel(postID uint, page, pageSize int) ([]model.Comment, int, error) {
	// 计算总数
	var total int64
	if err := m.Orm.Model(&model.Comment{}).Where("post_id = ? AND reply_comment_id = ?", postID, 0).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var comments []model.Comment
	query := m.Orm.Model(&model.Comment{}).
		Where("post_id = ? AND reply_comment_id = ?", postID, 0).
		Limit(pageSize).
		Offset(offset).
		Order("id desc")
	if err := query.Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, int(total), nil
}

// ListSecondLevel 列出某个一级评论下二级评论列表
func (m *CommentDao) ListSecondLevel(commendID uint, page, pageSize int) ([]model.Comment, int, error) {
	// 计算总数
	var total int64
	if err := m.Orm.Model(&model.Comment{}).Where("reply_root_comment_id = ?", commendID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var comments []model.Comment
	query := m.Orm.Model(&model.Comment{}).
		Where("reply_root_comment_id = ?", commendID).
		Limit(pageSize).
		Offset(offset).
		Order("id desc")
	if err := query.Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, int(total), nil
}

// CheckCommentExist 检测帖子是否存在
func (m *CommentDao) CheckCommentExist(commentID uint) (*model.Comment, bool) {
	var comment model.Comment
	if err := m.Orm.Where("id = ?", commentID).First(&comment).Error; err != nil {
		return nil, false
	}
	return &comment, true
}

func (m *CommentDao) CheckPostExist(postID uint) (*model.Post, bool) {
	var post model.Post
	if err := m.Orm.Where("id = ?", postID).First(&post).Error; err != nil {
		return nil, false
	}
	return &post, true
}
