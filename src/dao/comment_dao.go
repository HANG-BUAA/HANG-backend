package dao

import (
	"HANG-backend/src/model"
	"errors"
	"fmt"
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

func (m *CommentDao) CreateComment(userID uint, postID uint, replyCommentID uint, content string, isAnonymous bool) (*model.Comment, error) {
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
		// 一级评论，其 root_id 由 AfterCreate 钩子生成
		replyUserName, err = m.getPostUserName(&post)
		if err != nil {
			return nil, err
		}
		isReplyAnonymous = post.IsAnonymous
	} else {
		// 二级评论
		replyRootCommentID = comment.ReplyRootCommentID
		replyUserName = comment.UserName
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

func (m *CommentDao) GetCommentUserNameAndAvatar(comment *model.Comment) (string, string, error) {
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
