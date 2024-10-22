package dto

type AdminCourseCreateRequestDTO struct {
	ID      string   `json:"id" form:"id" binding:"required" required_err:"id is Required"`
	Name    string   `json:"name" form:"name" binding:"required" required_err:"name is Required"`
	Credits *float32 `json:"credits" form:"credits"`
	Campus  *int     `json:"campus" form:"campus"`
	Tags    []int    `json:"tags" form:"tags"`
}

type AdminCourseCreateResponseDTO struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Credits   *float32 `json:"credits"`
	Campus    *int     `json:"campus"`
	Tags      []int    `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	DeletedAt string   `json:"deleted_at"`
}
