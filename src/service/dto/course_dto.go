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
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Credits   *float32       `json:"credits"`
	Campus    *int           `json:"campus"`
	ReviewNum int            `json:"review_num"`
	Tags      []model.Tag    `json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type AdminCourseCreateRequestDTO struct {
	ID      string   `json:"id" form:"id" binding:"required" required_err:"id is Required"`
	Name    string   `json:"name" form:"name" binding:"required" required_err:"name is Required"`
	Credits *float32 `json:"credits" form:"credits"`
	Campus  *int     `json:"campus" form:"campus"`
	Tags    []uint   `json:"tags" form:"tags"`
}

type AdminCourseCreateResponseDTO CourseOverviewDTO

type CreateCourseReviewRequestDTO struct {
	User     *model.User
	CourseID string `json:"course_id" form:"course_id" binding:"required" required_err:"course_id is Required"`
	Content  string `json:"content" form:"content" binding:"required" required_err:"content is Required"`
	Score    int    `json:"score" form:"score" binding:"required" required_err:"score is Required"`
}

type CreateCourseReviewResponseDTO CourseReviewOverviewDTO

type LikeCourseReviewRequestDTO struct {
	User         *model.User
	CourseReview *model.CourseReview
}
