package dto

type PaginationInfo struct {
	TotalRecords int `json:"total_records"` // 总记录数
	PageSize     int `json:"page_size"`     // 每页条数
	NextCursor   any `json:"next_cursor"`
}

func BuildPaginationInfo(total, pageSize int, nextCursor any) *PaginationInfo {
	return &PaginationInfo{
		TotalRecords: total,
		PageSize:     pageSize,
		NextCursor:   nextCursor,
	}
}
