package utils

import (
	"strconv"
	"strings"
)

// 解析float32，如果值为空，返回nil
func ParseFloat32(value string) *float32 {
	if value == "" {
		return nil
	}
	v, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil
	}
	floatValue := float32(v)
	return &floatValue
}

// 解析int，如果值为空，返回nil
func ParseInt(value string) *int {
	if value == "" {
		return nil
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &v
}

// 解析Tags（多个整数用逗号分隔）
func ParseTags(value string) []uint {
	if value == "" {
		return nil
	}
	tagStrings := strings.Split(value, ",")
	var tags []uint
	for _, tag := range tagStrings {
		tagInt, err := strconv.Atoi(tag)
		if err == nil {
			tags = append(tags, uint(tagInt))
		}
	}
	return tags
}
