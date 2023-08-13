// Package message 包描述
// Author: wanlizhan
// Date: 2023/6/11
package message

import (
	"github.com/xxzhwl/wdk/ulog"
)

// IRobotMessage 机器人消息接口
type IRobotMessage interface {
	SendText(msg string) error
}

var RobotClient IRobotMessage

// SendText 机器人发送消息
func SendText(content string) error {
	if RobotClient == nil {
		ulog.Error("RobotMessage", "IRobotMessage interface not implemented")
		return nil
	}
	return RobotClient.SendText(content)
}
