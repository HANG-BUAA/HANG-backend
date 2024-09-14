package test

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "HANG <18646154381@163.com>"
	e.To = []string{"fanxingkun@buaa.edu.cn"}
	e.Subject = "Your Verification Code"
	e.Text = []byte("Your verification code is: " + "22371426")

	err := e.Send("smtp.163.com:25", smtp.PlainAuth("", "18646154381@163.com", "MSTLOQJPDZCLVMOZ", "smtp.163.com"))
	if err != nil {
		t.Error(err)
	}
}
