package service

import (
	"HANG-backend/src/custom_error"
	"HANG-backend/src/dao"
	"HANG-backend/src/global"
	"HANG-backend/src/model"
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var postService *PostService

type PostService struct {
	BaseService
	Dao *dao.PostDao
}

func NewPostService() *PostService {
	if postService == nil {
		postService = &PostService{
			Dao: dao.NewPostDao(),
		}
	}
	return postService
}

// Create 创建帖子
func (m *PostService) Create(postCreateDTO *dto.PostCreateRequestDTO) (res *dto.PostCreateResponseDTO, err error) {
	user := postCreateDTO.User
	title := postCreateDTO.Title
	content := postCreateDTO.Content
	isAnonymous := postCreateDTO.IsAnonymous

	// todo 可能还要判断用户是否被禁言——中间件实现

	post, err := m.Dao.CreatePost(user, title, content, *isAnonymous)
	if err != nil {
		return
	}
	tmp, err := m.Dao.ConvertPostModelToOverviewDTO(post, user.ID)
	if err != nil {
		return
	}
	res = (*dto.PostCreateResponseDTO)(tmp)

	// 传输关键词数据
	go utils.SetKeyword(content)

	// 创建协程，异步地将数据传输到 rabbitmq中，进而同步到 es 里
	go func() {
		err := utils.PublishPostMessage(utils.PostMessage{
			ID:      post.ID,
			Title:   post.Title,
			Content: post.Content,
		})
		if err != nil {
			// todo 此处对于丢失的数据只写到了 logger 里，后期考虑更靠谱的方案
			global.Logger.Error(err)
		}
	}()
	return
}

// Like 用户喜欢帖子
func (m *PostService) Like(postLikeRequestDTO *dto.PostLikeRequestDTO) (err error) {
	user := postLikeRequestDTO.User
	post := postLikeRequestDTO.Post

	// 判断用户是否已经喜欢了该帖子
	if m.Dao.CheckLiked(user, post) {
		return errors.New("liked post")
	}

	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.Like(user, post)
		if err == nil {
			// 喜欢成功
			notification := model.Notification{
				Type:           "like post",
				OperatorID:     user.ID,
				OperatorName:   user.Username,
				OperatorAvatar: user.Avatar,
				NotifierID:     post.UserID,
				EntityID:       post.ID,
			}
			global.RDB.Create(&notification)
			return
		}

		if errors.Is(err, &custom_error.OptimisticLockError{}) {
			// 并发版本冲突，重试
			continue
		}
		return
	}
	return custom_error.NewOptimisticLockError()
}

func (m *PostService) Unlike(postUnlikeRequestDTO *dto.PostUnlikeRequestDTO) (err error) {
	user := postUnlikeRequestDTO.User
	post := postUnlikeRequestDTO.Post

	// 判断是否喜欢过该帖子
	if !m.Dao.CheckLiked(user, post) {
		return errors.New("unliked post")
	}

	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.Unlike(user, post)
		if err == nil {
			return
		}
		if errors.Is(err, &custom_error.OptimisticLockError{}) {
			continue
		}
	}
	return custom_error.NewOptimisticLockError()
}

// Collect 收藏帖子
func (m *PostService) Collect(postCollectRequestDTO *dto.PostCollectRequestDTO) (err error) {
	user := postCollectRequestDTO.User
	post := postCollectRequestDTO.Post

	// 判断用户是否已经收藏了该帖子
	if m.Dao.CheckCollected(user, post) {
		return errors.New("collected post")
	}

	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.Collect(user, post)
		if err == nil {
			// 收藏成功
			notification := model.Notification{
				Type:           "collect post",
				OperatorID:     user.ID,
				OperatorName:   user.Username,
				OperatorAvatar: user.Avatar,
				NotifierID:     post.UserID,
				EntityID:       post.ID,
			}
			global.RDB.Create(&notification)

			return
		}

		if errors.Is(err, &custom_error.OptimisticLockError{}) {
			// 并发版本冲突，重试
			continue
		}
		return
	}
	return custom_error.NewOptimisticLockError()
}

func (m *PostService) Uncollect(postUncollectRequestDTO *dto.PostUncollectRequestDTO) (err error) {
	user := postUncollectRequestDTO.User
	post := postUncollectRequestDTO.Post

	if !m.Dao.CheckCollected(user, post) {
		return errors.New("uncollected post")
	}

	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.Uncollect(user, post)
		if err == nil {
			return
		}
		if errors.Is(err, &custom_error.OptimisticLockError{}) {
			continue
		}
	}
	return custom_error.NewOptimisticLockError()
}

// CommonList 普通查询列表（不带搜索）
func (m *PostService) CommonList(postListRequestDTO *dto.PostListRequestDTO) (res *dto.PostListResponseDTO, err error) {
	pageSize := postListRequestDTO.PageSize
	user := postListRequestDTO.User

	cursor, err := strconv.ParseUint(postListRequestDTO.Cursor, 10, 32)
	if err != nil {
		cursor = 0
	}

	// todo 如果要做个性化推荐的话，后面这里要考虑把 user_id 传入，在 CommonList 服务里使用
	posts, total, isEnd, err := m.Dao.CommonList(uint(cursor), pageSize)
	if err != nil {
		return
	}
	if len(posts) == 0 {
		res = &dto.PostListResponseDTO{
			Posts: []dto.PostOverviewDTO{},
		}
		return
	}
	overviews, err := m.Dao.ConvertPostModelsToOverviewDTOs(posts, user.ID)
	if err != nil {
		return
	}

	var nextCursor any
	if isEnd {
		nextCursor = 0
	} else {
		nextCursor = posts[len(posts)-1].ID
	}

	res = &dto.PostListResponseDTO{
		Pagination: *dto.BuildPaginationInfo(total, len(overviews), nextCursor),
		Posts:      overviews,
	}
	return
}

func (m *PostService) SearchList(postListRequestDTO *dto.PostListRequestDTO) (res *dto.PostListResponseDTO, err error) {
	// todo baseURL 换成全局变量
	baseURL := fmt.Sprintf("http://%s:%s/post",
		viper.GetString("search_client.host"),
		viper.GetString("search_client.port"),
	)
	user := postListRequestDTO.User
	query := postListRequestDTO.Query
	cursor := postListRequestDTO.Cursor
	pageSize := postListRequestDTO.PageSize

	params := url.Values{}
	params.Add("query", query)
	params.Add("cursor", cursor)
	params.Add("page_size", strconv.Itoa(pageSize))
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 向搜索端发送请求
	resp, err := http.Get(fullURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// 响应结构
	var responseBody struct {
		Total int `json:"total"`
		Posts []struct {
			ID    uint    `json:"id"`
			Score float64 `json:"score"`
		} `json:"posts"`
		NextCursor string `json:"next_cursor"`
	}

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return
	}

	var ids []uint
	total := responseBody.Total
	for _, post := range responseBody.Posts {
		ids = append(ids, post.ID)
	}
	nextCursor := responseBody.NextCursor

	// 根据返回的 id 列表获取列表
	posts, err := m.Dao.GetListsByIDs(ids)
	if err != nil {
		return
	}

	// 结构体转换
	overviews, err := m.Dao.ConvertPostModelsToOverviewDTOs(posts, user.ID)
	if err != nil {
		return
	}
	res = &dto.PostListResponseDTO{
		Pagination: *dto.BuildPaginationInfo(total, len(overviews), nextCursor),
		Posts:      overviews,
	}
	return
}

func (m *PostService) CollectionList(postCollectionListRequestDTO *dto.PostCollectionListRequestDTO) (res *dto.PostCollectionListResponseDTO, err error) {
	pageSize := postCollectionListRequestDTO.PageSize
	user := postCollectionListRequestDTO.User
	tmp, err := utils.ParseTimeWithMultipleFormats(postCollectionListRequestDTO.Cursor)
	if err != nil {
		tmp = time.Time{}
	}
	cursor := tmp

	posts, total, isEnd, err := m.Dao.GetCollections(user, cursor, pageSize)
	if err != nil {
		return
	}
	overviews, err := m.Dao.ConvertPostModelsToOverviewDTOs(posts, user.ID)
	if err != nil {
		return
	}

	// 构建 nextCursor
	nextCursor := time.Time{}
	if !isEnd {
		nextCursor, err = m.Dao.GetCollectCursorByID(user, &posts[len(posts)-1])
		if err != nil {
			return
		}
	}
	res = &dto.PostCollectionListResponseDTO{
		Pagination: *dto.BuildPaginationInfo(total, len(overviews), nextCursor),
		Posts:      overviews,
	}
	return
}

func (m *PostService) Retrieve(postRetrieveRequestDTO *dto.PostRetrieveRequestDTO) (res *dto.PostRetrieveResponseDTO, err error) {
	user := postRetrieveRequestDTO.User
	post := postRetrieveRequestDTO.Post

	overview, err := m.Dao.ConvertPostModelToOverviewDTO(post, user.ID)
	if err != nil {
		return
	}
	res = &dto.PostRetrieveResponseDTO{
		Post: *overview,
	}
	return
}
