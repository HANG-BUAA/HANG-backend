package global

// 用户相关错误码
const (
	ERR_CODE_LOGIN_FAILED = iota + 10010
	ERR_CODE_SEND_EMAIL
	ERR_CODE_ADD_USER
	ERR_CODE_INVALID_TOKEN
)

// 文件操作相关错误码
const (
	ERR_CODE_UPLOAD_MISSING_FIELD = iota + 20010
	ERR_CODE_UPLOAD_SERVER_FAILED
)

// 帖子相关错误码
const (
	ERR_CODE_POST_FAILED = iota + 30010
)

// 评论相关错误码
const (
	ERR_CODE_COMMENT_FAILED = iota + 40010
)
