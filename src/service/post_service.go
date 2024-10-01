package service

import (
	"HANG-backend/src/custom_error"
	"HANG-backend/src/dao"
	"HANG-backend/src/global"
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
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
	userID := postCreateDTO.UserID
	title := postCreateDTO.Title
	content := postCreateDTO.Content
	isAnonymous := postCreateDTO.IsAnonymous

	// todo 可能还要判断用户是否被禁言——中间件实现

	post, err := m.Dao.CreatePost(userID, title, content, *isAnonymous)
	if err != nil {
		return
	}
	tmp, err := m.Dao.ConvertPostModelToOverviewDTO(post, userID)
	if err != nil {
		return
	}
	res = (*dto.PostCreateResponseDTO)(tmp)

	// 创建协程，异步地将数据传输到 rabbitmq中，进而同步到 es 里
	go func() {
		err := utils.PublishPostMessage(utils.PostMessage{
			post.ID,
			post.Title,
			post.Content,
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
	userID := postLikeRequestDTO.UserID
	postID := postLikeRequestDTO.PostID
	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.Like(userID, postID)
		if err == nil {
			// 喜欢成功
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

// Collect 收藏帖子
func (m *PostService) Collect(postCollectRequestDTO *dto.PostCollectRequestDTO) (err error) {
	userID := postCollectRequestDTO.UserID
	postID := postCollectRequestDTO.PostID
	for retries := 0; retries < global.OptimisticLockMaxRetries; retries++ {
		err = m.Dao.Collect(userID, postID)
		if err == nil {
			// 收藏成功
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

func (m *PostService) List(postListRequestDTO *dto.PostListRequestDTO) (res *dto.PostListResponseDTO, err error) {
	page := postListRequestDTO.Page
	pageSize := postListRequestDTO.PageSize
	userID := postListRequestDTO.UserID
	query := postListRequestDTO.Query

	var ids []uint = nil
	if query != "" {
		ids, err = searchPostsByQuery(query)
		if err != nil {
			return nil, errors.New("search end error!")
		}
		if len(ids) == 0 {
			// 没有匹配的结果
			return &dto.PostListResponseDTO{}, nil
		}
	}

	// todo 如果要做个性化推荐的话，后面这里要考虑把 user_id 传入，在 ListFirstLevel 服务里使用
	posts, total, err := m.Dao.List(page, pageSize, ids)
	if err != nil {
		return
	}
	overviews, err := m.Dao.ConvertPostModelsToOverviewDTOs(posts, userID)
	if err != nil {
		return
	}
	res = &dto.PostListResponseDTO{
		Pagination: *dto.BuildPaginationInfo(total, page, pageSize),
		Posts:      overviews,
	}
	return
}

func searchPostsByQuery(query string) ([]uint, error) {
	baseURL := fmt.Sprintf("http://%s:%s/post",
		viper.GetString("search_client.host"),
		viper.GetString("search_client.port"),
	)
	params := url.Values{}
	params.Add("query", query)
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		return []uint{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []uint{}, err
	}

	// 解析响应数据
	var items []map[string]interface{}
	err = json.Unmarshal(body, &items)
	if err != nil {
		return []uint{}, err
	}

	// 提取所有id
	var ids []uint
	for _, item := range items {
		if id, ok := item["id"].(float64); ok { // JSON数字解析为float64
			ids = append(ids, uint(id))
		}
	}
	return ids, nil
}
