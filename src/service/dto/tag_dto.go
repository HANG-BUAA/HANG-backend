package dto

type AdminTagCreateRequestDTO struct {
	Type int    `json:"type" form:"type" binding:"required" required_err:"type is Required"`
	Name string `json:"name" form:"name" binding:"required" required_err:"name is Required"`
}

type AdminTagCreateResponseDTO struct {
	ID   uint   `json:"id"`
	Type int    `json:"type"`
	Name string `json:"name"`
}
