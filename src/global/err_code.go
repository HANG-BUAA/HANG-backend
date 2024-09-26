package global

// 全局错误码
const (
	ERR_CODE_INVALID_TOKEN = iota + 10
	ERR_CODE_PERMISSION_DENIED
	ERR_CODE_OPTISTIC_LOCK_RETRY_LIMIT
)

// 用户相关错误码
const (
	ERR_CODE_LOGIN_FAILED = iota + 10010
	ERR_CODE_SEND_EMAIL
	ERR_CODE_ADD_USER
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
