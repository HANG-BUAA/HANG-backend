package utils

import (
	"HANG-backend/src/global"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
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
