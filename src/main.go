package main

import (
	"HANG-backend/src/cmd"
)

// @title 小航书后端
// @version 0.0.1
// @description 后端api接口文档
func main() {
	defer cmd.Clean()
	cmd.Start()
}
