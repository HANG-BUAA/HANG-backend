package dto

import "reflect"

type PaginationInfo struct {
	TotalRecords int `json:"total_records"` // 总记录数
	PageSize     int `json:"page_size"`     // 每页条数
	NextCursor   any `json:"next_cursor,omitempty"`
}

func BuildPaginationInfo(total, pageSize int, nextCursor any) *PaginationInfo {
	// nextCursor 判空
	if isZeroValue(nextCursor) {
		nextCursor = nil
	}
	return &PaginationInfo{
		TotalRecords: total,
		PageSize:     pageSize,
		NextCursor:   nextCursor,
	}
}

func isZeroValue(nextCursor any) bool {
	if nextCursor == nil {
		return true
	}
	v := reflect.ValueOf(nextCursor)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.String:
		return v.String() == ""
	default:
		return v.IsZero()
	}
}
