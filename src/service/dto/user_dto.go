package dto

type UserLoginDTO struct {
	Username string `json:"username" binding:"required,first_is_a" message:"用户名错误" required_err:"用户名不可缺省"`
	Password string `json:"password" binding:"required"`
}
