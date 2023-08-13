// Package utils 包描述
// Author: wanlizhan
// Date: 2023/4/9
package uhttp

import (
	"github.com/bytedance/sonic"
	"github.com/xxzhwl/wdk/message"
	"testing"
	"time"
)

type Body struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func TestName2(t *testing.T) {
	client := NewHttpClient(3*time.Second, 30*time.Second, 3)
	marshal, _ := sonic.Marshal(map[string]any{"msg_type": "text", "content": map[string]string{"text": "测试发送告警消息"}})
	arg := PostArg{
		Title:  "测试",
		Url:    "",
		Body:   marshal,
		Header: nil,
		CallbackFunc: func(response []byte) error {
			message.SendText("发送成功")
			return nil
		},
	}
	client.Post(arg)
}
