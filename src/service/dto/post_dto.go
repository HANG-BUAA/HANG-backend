package dto

import (
	"HANG-backend/src/model"
	"gorm.io/gorm"
	"time"
)

type PostAuthorDTO struct {
	UserID     uint   `json:"user_id,omitempty"`
	UserName   string `json:"user_name"`
	UserAvatar string `json:"user_avatar"`
}

type PostOverviewDTO struct {
	ID           uint           `json:"id"`
	Author       PostAuthorDTO  `json:"author"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	IsAnonymous  bool           `json:"is_anonymous"`
	CollectNum   int            `json:"collect_num"`
	LikeNum      int            `json:"like_num"`
	CommentNum   int            `json:"comment_num"`
	HasLiked     bool           `json:"has_liked"`
	HasCollected bool           `json:"has_collected"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

type PostCreateRequestDTO struct {
	User        *model.User
	Title       string `json:"title" form:"title" binding:"required" required_err:"title is Required"`
	Content     string `json:"content" form:"content" binding:"required" required_err:"content is Required"`
	IsAnonymous *bool  `json:"is_anonymous" form:"is_anonymous" binding:"required" required_err:"is_anonymous is Required"`
}

type PostCreateResponseDTO PostOverviewDTO

type PostLikeRequestDTO struct {
	Post *model.Post
	User *model.User
}

type PostCollectRequestDTO struct {
	Post *model.Post
	User *model.User
}

type PostUnlikeRequestDTO struct {
	Post *model.Post
	User *model.User
}

type PostUncollectRequestDTO struct {
	Post *model.Post
	User *model.User
}

type PostListRequestDTO struct {
	Cursor   string
	PageSize int
	User     *model.User
	Query    string `form:"query"`
}

type PostListResponseDTO struct {
	Pagination PaginationInfo    `json:"pagination"`
	Posts      []PostOverviewDTO `json:"posts"`
}

type PostCollectionListRequestDTO struct {
	Cursor   string
	PageSize int
	User     *model.User
}

type PostCollectionListResponseDTO struct {
	Pagination PaginationInfo    `json:"pagination"`
	Posts      []PostOverviewDTO `json:"posts"`
}

type PostRetrieveRequestDTO struct {
	User *model.User
	Post *model.Post
}

type PostRetrieveResponseDTO struct {
	Post PostOverviewDTO `json:"post"`
}
