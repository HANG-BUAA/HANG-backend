package dto

type PaginationInfo struct {
	TotalRecords int `json:"total_records"` // 总记录数
	CurrentPage  int `json:"current_page"`  // 当前页
	PageSize     int `json:"page_size"`     // 每页条数
	TotalPages   int `json:"total_pages"`   // 总页数
}

func BuildPaginationInfo(total, curPage, pageSize int) *PaginationInfo {
	return &PaginationInfo{
		TotalRecords: total,
		CurrentPage:  curPage,
		PageSize:     pageSize,
		TotalPages:   (total + pageSize - 1) / pageSize,
	}
}
