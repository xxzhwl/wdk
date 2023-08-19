// Package ulog 包描述
// Author: wanlizhan
// Date: 2023/7/16
package ulog

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/xxzhwl/wdk/project"
	"github.com/xxzhwl/wdk/system"
	"github.com/xxzhwl/wdk/ucontext"
)

type ApiLogData struct {
	Env        string
	HttpMethod string
	Code       string
	GoId       string
	LogId      string
	LogLevel   string
	LogTime    string
	LocalId    string
	TraceId    string
	ReqId      string
	Stack      string
	SystemName string
	Path       string
	Request    string
	Response   string
	Method     string
	ClientIp   string
	StartTime  string
	EndTime    string
	Duration   string
}

func ApiOkLogStore(data ApiLogData) {
	logId := GetLogId()
	context := ucontext.GetCurrentContext()
	data.Env = project.GetRunTime()
	data.LogId = logId
	data.ReqId = context.RequestId
	data.TraceId = context.TraceId
	data.LocalId = context.LocalId
	data.GoId = system.GetGoRoutineId()

	logMsg := fmt.Sprintf("%v", data)
	marshal, err := sonic.Marshal(data)
	if err == nil {
		logMsg = string(marshal)
	}
	Info("ApiOkResponse", logMsg)
	if remoteLogger != nil {
		remoteLogger.Info("apilog", logMsg)
	}
}

func ApiFailLogStore(data ApiLogData) {
	logId := GetLogId()
	context := ucontext.GetCurrentContext()
	data.LogId = logId
	data.Env = project.GetRunTime()
	data.Stack = system.GetStackFramesString(2, 0)
	data.ReqId = context.RequestId
	data.TraceId = context.TraceId
	data.LocalId = context.LocalId
	data.GoId = system.GetGoRoutineId()

	logMsg := fmt.Sprintf("%v", data)
	marshal, err := sonic.Marshal(data)
	if err == nil {
		logMsg = string(marshal)
	}

	Error("ApiFailResponse", logMsg)
	if remoteLogger != nil {
		remoteLogger.Info("apilog", logMsg)
	}
}
