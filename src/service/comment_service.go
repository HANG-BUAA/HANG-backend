package service

import (
	"HANG-backend/src/custom_error"
	"HANG-backend/src/dao"
	"HANG-backend/src/global"
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"errors"
	"strconv"
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
	user := commentCreateDTO.User
	replyCommentID := commentCreateDTO.ReplyCommentID
	content := commentCreateDTO.Content
	isAnonymous := commentCreateDTO.IsAnonymous

	// 判断 post 是否存在
	post, exist := m.Dao.CheckPostExist(commentCreateDTO.PostID)
	if !exist {
		return nil, errors.New("target post not found")
	}

	comment, err := m.Dao.CreateComment(user, post, *replyCommentID, content, *isAnonymous)
	if err != nil {
		return
	}

	tmp, err := m.Dao.ConvertCommentModelToOverviewDTO(comment, user.ID)
	if err != nil {
		return nil, err
	}
	res = (*dto.CommentCreateResponseDTO)(tmp)
	return
}

// Like 用户喜欢评论
func (m *CommentService) Like(commentLikeRequestDTO *dto.CommentLikeRequestDTO) (err error) {
	user := commentLikeRequestDTO.User
	comment := commentLikeRequestDTO.Comment

	// 判断是否已经喜欢该评论
	if m.Dao.CheckLiked(user, comment) {
		return errors.New("liked comment")
	}

	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.Like(user, comment)
		if err == nil {
			return
		}

		if errors.Is(err, &custom_error.OptimisticLockError{}) {
			continue
		}
	}
	return custom_error.NewOptimisticLockError()
}

// ListFirstLevel 列出某帖子下一级评论列表
func (m *CommentService) ListFirstLevel(commentListRequestDTO *dto.CommentListRequestDTO) (res *dto.CommentListResponseDTO, err error) {
	pageSize := commentListRequestDTO.PageSize
	user := commentListRequestDTO.User
	postID := commentListRequestDTO.PostID

	tmp, err := strconv.ParseUint(commentListRequestDTO.Cursor, 10, 32)
	if err != nil {
		tmp = 0
	}
	cursor := uint(tmp)

	// 校验帖子是否存在
	_, ok := m.Dao.CheckPostExist(postID)
	if !ok {
		return nil, errors.New("post not exist")
	}

	comments, total, err := m.Dao.ListFirstLevel(postID, cursor, pageSize)
	if err != nil {
		return
	}
	if cursor == 0 {
		cursor = uint(total)
	}
	overviews, err := m.Dao.ConvertCommentModelsToOverviewDTOs(comments, user.ID)
	if err != nil {
		return
	}

	nextCursor := utils.IfThenElse(int(cursor)-pageSize > 0, int(cursor)-pageSize, 0)
	res = &dto.CommentListResponseDTO{
		Pagination: *dto.BuildPaginationInfo(total, len(overviews), nextCursor),
		Comments:   overviews,
	}
	return
}

// ListSecondLevel 列出某一级评论下二级评论列表
func (m *CommentService) ListSecondLevel(commentListRequestDTO *dto.CommentListRequestDTO) (res *dto.CommentListResponseDTO, err error) {
	pageSize := commentListRequestDTO.PageSize
	userID := commentListRequestDTO.User.ID
	commentID := commentListRequestDTO.CommentID

	tmp, err := strconv.ParseUint(commentListRequestDTO.Cursor, 10, 32)
	if err != nil {
		tmp = 0
	}
	cursor := uint(tmp)

	// 校验一级评论是否存在
	_, ok := m.Dao.CheckCommentExist(commentID)
	if !ok {
		return nil, errors.New("comment not exist")
	}

	comments, total, err := m.Dao.ListSecondLevel(commentID, cursor, pageSize)
	if err != nil {
		return
	}

	if cursor == 0 {
		cursor = uint(total)
	}
	overviews, err := m.Dao.ConvertCommentModelsToOverviewDTOs(comments, userID)
	if err != nil {
		return
	}

	nextCursor := utils.IfThenElse(int(cursor)-pageSize > 0, int(cursor)-pageSize, 0)
	res = &dto.CommentListResponseDTO{
		Pagination: *dto.BuildPaginationInfo(total, len(overviews), nextCursor),
		Comments:   overviews,
	}
	return
}
