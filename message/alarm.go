// Package message 包描述
// Author: wanlizhan
// Date: 2023/6/11
package message

import (
	"github.com/xxzhwl/wdk/ulog"
)

// IAlarm 告警接口
type IAlarm interface {
	SendAlarmMail(title, content string) error
	SendAlarmMessage(title, content string) error
}

// AlarmClient 告警客户端
var AlarmClient IAlarm

// SendAlarmMessage  发送告警消息
func SendAlarmMessage(title, content string) error {
	if AlarmClient == nil {
		ulog.Error("Alarm", "IAlarm interface not implemented")
		return nil
	}
	return AlarmClient.SendAlarmMessage(title, content)
}
