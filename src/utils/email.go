package utils

import (
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"math/rand"
	"net/smtp"
)

func generateVerificationCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func SendEmail(targetID string) (string, error) {
	code := generateVerificationCode()
	e := email.NewEmail()
	e.From = fmt.Sprintf("HANG <%s>", viper.GetString("smtp.user_addr"))
	e.To = []string{fmt.Sprintf("%s@buaa.edu.cn", targetID)}
	e.Subject = "Your Verification Code"
	e.Text = []byte("Your verification code is: " + code)

	// 发送验证码
	err := e.Send(
		viper.GetString("smtp.smtp_addr"),
		smtp.PlainAuth("", viper.GetString("smtp.user_addr"), viper.GetString("smtp.smtp_password"), viper.GetString("smtp.smtp_host")))
	if err != nil {
		return "", err
	}
	return code, nil
}
