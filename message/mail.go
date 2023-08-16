// Package message 包描述
// Author: wanlizhan
// Date: 2023/6/11
package message

import (
	mail "github.com/xhit/go-simple-mail/v2"
	"time"
)

// MailClient 邮件客户端
type MailClient struct {
	Conf   MailClientConf
	client mail.SMTPClient
}

// MailClientConf 邮件客户端配置
type MailClientConf struct {
	Server      string
	Port        int
	UserName    string
	Password    string
	sendTimeOut time.Duration
	connTimeOut time.Duration
}

// NewMailClient 根据配置新建客户端
func NewMailClient(conf MailClientConf) MailClient {
	return MailClient{
		Conf: conf,
	}
}

// SetSendTimeOut 设置邮件发送超时时间
func (m *MailClient) SetSendTimeOut(out time.Duration) *MailClient {
	m.Conf.sendTimeOut = out
	return m
}

// SetConnTimeOut 设置邮件服务器连接超时时间
func (m *MailClient) SetConnTimeOut(out time.Duration) *MailClient {
	m.Conf.connTimeOut = out
	return m
}

// GetClientInstance 获取邮件服务实例
func (m *MailClient) GetClientInstance() mail.SMTPClient {
	return m.client
}

// SendMailArg 发邮件的一般常用参数
type SendMailArg struct {
	From             string
	To               string
	Cc               string
	BodyBytes        []byte
	BodyType         mail.ContentType
	Subject          string
	AttachList       []AttachInfo
	Base64AttachList []Base64AttachInfo
	FileAttachList   []string
}

// AttachInfo 附件信息
type AttachInfo struct {
	Name string
	Data []byte
}

// Base64AttachInfo base64类型的附件
type Base64AttachInfo struct {
	Name string
	Data string
}

// SendMail 发送邮件
func (m *MailClient) SendMail(arg SendMailArg) error {
	client := mail.NewSMTPClient()
	client.Host = m.Conf.Server
	client.Port = m.Conf.Port
	client.Username = m.Conf.UserName
	client.Password = m.Conf.Password
	client.SendTimeout = m.Conf.sendTimeOut
	client.ConnectTimeout = m.Conf.connTimeOut

	dial, err := client.Connect()
	if err != nil {
		return err
	}
	dial.KeepAlive = false
	dial.SendTimeout = m.Conf.sendTimeOut

	msg := mail.NewMSG().SetSubject(arg.Subject).SetFrom(arg.From).
		AddTo(arg.To).AddCc(arg.Cc).SetBodyData(arg.BodyType, arg.BodyBytes)

	for i, _ := range arg.AttachList {
		msg = msg.Attach(&mail.File{
			Name: arg.AttachList[i].Name,
			Data: arg.AttachList[i].Data,
		})
	}

	for i, _ := range arg.FileAttachList {
		msg = msg.Attach(&mail.File{
			FilePath: arg.FileAttachList[i],
		})
	}

	for i, _ := range arg.Base64AttachList {
		msg = msg.Attach(&mail.File{
			Name:    arg.AttachList[i].Name,
			B64Data: arg.Base64AttachList[i].Data,
		})
	}

	return msg.Send(dial)
}
