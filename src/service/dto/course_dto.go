package dto

import (
	"HANG-backend/src/model"
	"gorm.io/gorm"
	"time"
)

type CourseReviewOverviewDTO struct {
	ID        uint           `json:"id"`
	CourseID  string         `json:"course_id"`
	Content   string         `json:"content"`
	Score     int            `json:"score"`
	IsSelf    bool           `json:"is_self"`
	LikeNum   int            `json:"like_num"`
	HasLiked  bool           `json:"has_liked"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type CourseOverviewDTO struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Credits      *float32       `json:"credits"`
	Campus       *int           `json:"campus"`
	ReviewNum    int            `json:"review_num"`
	MaterialNum  int            `json:"material_num"`
	AverageScore *float64       `json:"average_score"`
	Tags         []model.Tag    `json:"tags"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

type CourseMaterialOverviewDTO struct {
	ID       uint   `json:"id"`
	CourseID string `json:"course_id"`
	Author   struct {
		ID     uint   `json:"id"`
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	}
	Link        string         `json:"link"`
	Source      int            `json:"source"`
	Description string         `json:"description"`
	IsApproved  bool           `json:"is_approved"`
	IsOfficial  bool           `json:"is_official"`
	LikeNum     int            `json:"like_num"`
	HasLiked    bool           `json:"has_liked"`
	IsSelf      bool           `json:"is_self"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

type AdminCourseCreateRequestDTO struct {
	ID      string   `json:"id" form:"id" binding:"required" required_err:"id is Required"`
	Name    string   `json:"name" form:"name" binding:"required" required_err:"name is Required"`
	Credits *float32 `json:"credits" form:"credits"`
	Campus  *int     `json:"campus" form:"campus"`
	Tags    []uint   `json:"tags" form:"tags"`
}

type AdminCourseCreateResponseDTO CourseOverviewDTO

type CourseReviewCreateRequestDTO struct {
	User     *model.User
	CourseID string `json:"course_id" form:"course_id" binding:"required" required_err:"course_id is Required"`
	Content  string `json:"content" form:"content" binding:"required" required_err:"content is Required"`
	Score    int    `json:"score" form:"score" binding:"required" required_err:"score is Required"`
}

type CourseReviewCreateResponseDTO CourseReviewOverviewDTO

type CourseReviewLikeRequestDTO struct {
	User         *model.User
	CourseReview *model.CourseReview
}

type CourseReviewUnlikeRequestDTO struct {
	User         *model.User
	CourseReview *model.CourseReview
}

type CourseListRequestDTO struct {
	Cursor   string
	PageSize int
	Keyword  string `form:"keyword" json:"keyword"`
	Tags     []uint `form:"tags" json:"tags"`
}

type CourseListResponseDTO struct {
	Pagination PaginationInfo      `json:"pagination"`
	Courses    []CourseOverviewDTO `json:"courses"`
}

type CourseReviewListRequestDTO struct {
	Cursor   string
	PageSize int
	User     *model.User
	CourseID string  `json:"course_id" form:"course_id" binding:"required" required_err:"course_id is Required"`
	Query    *string `form:"query" json:"query"`
}

type CourseReviewListResponseDTO struct {
	Pagination PaginationInfo            `json:"pagination"`
	Reviews    []CourseReviewOverviewDTO `json:"reviews"`
}

type CourseRetrieveRequestDTO struct {
	Course *model.Course
}

type CourseRetrieveResponseDTO struct {
	Course CourseOverviewDTO `json:"course"`
}

type CourseMaterialCreateRequestDTO struct {
	User        *model.User
	IsOfficial  bool   `json:"is_official"`
	CourseID    string `json:"course_id" form:"course_id" binding:"required" required_err:"course_id is Required"`
	Link        string `json:"link" form:"link" binding:"required" required_err:"link is Required"`
	Source      int    `json:"source" form:"source" binding:"required" required_err:"source is Required"`
	Description string `json:"description" form:"description" binding:"required" required_err:"description is Required"`
}

type CourseMaterialCreateResponseDTO CourseMaterialOverviewDTO

type CourseMaterialLikeRequestDTO struct {
	User           *model.User
	CourseMaterial *model.CourseMaterial
}

type CourseMaterialUnlikeRequestDTO struct {
	User           *model.User
	CourseMaterial *model.CourseMaterial
}

type CourseMaterialListRequestDTO struct {
	Cursor     string
	PageSize   int
	User       *model.User
	CourseID   string `json:"course_id" form:"course_id" binding:"required" required_err:"course_id is Required"`
	IsOfficial *bool  `json:"is_official" form:"is_official" binding:"required" required_err:"is_official is Required"`
}

type CourseMaterialListResponseDTO struct {
	Pagination PaginationInfo              `json:"pagination"`
	Materials  []CourseMaterialOverviewDTO `json:"materials"`
}
