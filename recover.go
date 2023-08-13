// Package wdk 包描述
// Author: wanlizhan
// Date: 2023/6/10
package wdk

import (
	"fmt"
	"github.com/xxzhwl/wdk/message"
	"github.com/xxzhwl/wdk/system"
	"github.com/xxzhwl/wdk/ulog"
)

// CatchPanic 捕获panic
func CatchPanic() {
	if r := recover(); r != nil {
		PanicLog(r)
	}
}

// PanicLog 对panic信息打日志处理
func PanicLog(r any) string {
	if r == nil {
		return ""
	}
	errMsg := fmt.Sprintf("encouter panic:%v\n", r)
	stack := system.GetStackFramesString(3, 0)
	content := errMsg + stack

	ulog.Error("CatchPanic", content)
	if err := message.SendAlarmMessage("CatchPanic", content); err != nil {
		ulog.Error("SendAlarmMessage", err.Error())
	}

	return content
}
