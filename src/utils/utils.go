package utils

import (
	"HANG-backend/src/global"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func AppendError(existErr, newErr error) error {
	if existErr == nil {
		return newErr
	} else {
		return fmt.Errorf("%v, %w", existErr, newErr)
	}
}

func IfThenElse(condition bool, a any, b any) any {
	if condition {
		return a
	}
	return b
}

// UploadFile 上传文件到指定目录，返回最终目录
func UploadFile(file *multipart.FileHeader, dst string) (string, error) {
	folderPath := filepath.Join("images", dst)
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// 用 uuid 生成文件名
	ext := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + ext
	savePath := filepath.Join(folderPath, newFileName)

	f, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			global.Logger.Error("failed to close file: %v", cerr)
		}
	}()

	out, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			global.Logger.Error("failed to close file: %v", cerr)
		}
	}()

	if _, err := io.Copy(out, f); err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}
	return savePath, nil
}

// ParseTimeWithMultipleFormats 尝试使用多个格式解析时间字符串
func ParseTimeWithMultipleFormats(timeStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02T15:04:05.000-07:00", // 带三位毫秒和时区
		"2006-01-02T15:04:05.99-07:00",  // 带两位毫秒和时区
		"2006-01-02T15:04:05-07:00",     // 无毫秒带时区
		"2006-01-02T15:04:05Z07:00",     // 带时区（UTC 标准格式）
		"2006-01-02 15:04:05",           // 无时区和毫秒的格式
	}

	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.Parse(format, timeStr)
		if err == nil {
			return t, nil
		}
	}

	// 如果所有格式都无法解析，返回最后一个错误
	return time.Time{}, fmt.Errorf("unable to parse time: %v", err)
}
