package utils

import (
	gomail "gopkg.in/gomail.v2"
)

type EmailServer struct {
	ServerHost string          //邮件STMP地址 host
	ServerPort int             //邮件STMP端口 port
	Username   string          //发送者昵称
	FromEmail  string          //发送者邮件地址
	Password   string          //发送者连接邮件服务密码
	Mail       *gomail.Dialer  //邮件连接句柄
	Message    *gomail.Message //消息
	To         []string        //接受者邮件列表
	Content    EmailContent    //邮件内容
}

type EmailContent struct {
	Subject string //邮件主题
	Content string //内容
}

// NewEmailServer 返回一个email服务对象
func NewEmailServer() *EmailServer {
	conf := GetConfig()
	var es = &EmailServer{
		ServerHost: conf.Email.ServerHost,
		ServerPort: conf.Email.ServerPort,
		Username:   conf.Email.Username,
		FromEmail:  conf.Email.UserEmail,
		Password:   conf.Email.Password,
	}
	es.Message = gomail.NewMessage()
	es.Mail = gomail.NewDialer(es.ServerHost, es.ServerPort, es.FromEmail, es.Password)
	return es
}

// Send 邮件发送
func (this *EmailServer) Send(ec *EmailContent) error {
	this.Message.SetHeader("From", this.FromEmail)
	this.Message.SetHeader("To", this.To...)
	this.Message.SetHeader("Subject", ec.Subject)
	this.Message.SetBody("text/html", ec.Content)
	return this.Mail.DialAndSend(this.Message)
}
