package utils

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

var (
	STMP_HOST = ""
	STMP_PORT = 0
	STMP_USER = ""
	STMP_PASS = ""
)

const CONTENT_TEMPLATE = `
<h3>您正在重置密码</h3>
<p>您的验证码是：<strong>%s</strong></p>
<p>验证码 15分钟 有效，请勿泄露给他人</p>
`

// 发送重置密码 验证码邮箱
func SendResetCodeEmail(toEmail, code string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", STMP_USER)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", "密码重置验证码")
	msg.SetBody("text/html", fmt.Sprintf(CONTENT_TEMPLATE, code))

	dialer := gomail.NewDialer(STMP_HOST, STMP_PORT, STMP_USER, STMP_PASS)
	return dialer.DialAndSend(msg)
}
