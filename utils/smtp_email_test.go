package utils

import (
	"fmt"
	"testing"
)

func TestEmailServer_Send(t *testing.T) {
	es := NewEmailServer()
	fmt.Println(es)
	es.Message.SetHeader("From", es.FromEmail)
	es.Message.SetHeader("To", "xuthus5@qq.com")
	es.Message.SetHeader("Subject", "邮件测试功能!")
	es.Message.SetBody("text/html", "测试测试测试 - Hello <b>Bob</b> and <i>Cora</i>!")
	if err := es.Mail.DialAndSend(es.Message); err != nil {
		panic(err)
	}
}
