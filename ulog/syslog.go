// Package ulog 包描述
// Author: wanlizhan
// Date: 2023/7/16
package ulog

import (
	"github.com/bytedance/sonic"
	"github.com/xxzhwl/wdk/project"
	"github.com/xxzhwl/wdk/system"
	"github.com/xxzhwl/wdk/ucontext"
	"github.com/xxzhwl/wdk/utime"
)

func SyslogInfo(title, msg string) {
	logId := GetLogId()
	context := ucontext.GetCurrentContext()
	data := LogData{
		Env:        project.GetRunTime(),
		Message:    msg,
		Title:      title,
		LogId:      logId,
		LogLevel:   "INFO",
		LogTime:    utime.DateTimeFormat("2006-01-02 15:04:05.000"),
		TraceId:    context.TraceId,
		ReqId:      context.RequestId,
		LocalId:    context.LocalId,
		GoId:       system.GetGoRoutineId(),
		Stack:      system.GetStackFramesString(2, 0),
		SystemName: project.GetProjectName(),
	}
	if remoteLogger != nil {
		marshal, _ := sonic.Marshal(data)
		remoteLogger.Info("syslog", string(marshal))
	}
}

func SyslogWarn(title, msg string) {
	logId := GetLogId()
	context := ucontext.GetCurrentContext()
	data := LogData{
		Env:        project.GetRunTime(),
		Message:    msg,
		Title:      title,
		LogId:      logId,
		LogLevel:   "WARN",
		LogTime:    utime.DateTimeFormat("2006-01-02 15:04:05.000"),
		TraceId:    context.TraceId,
		ReqId:      context.RequestId,
		LocalId:    context.LocalId,
		GoId:       system.GetGoRoutineId(),
		Stack:      system.GetStackFramesString(2, 0),
		SystemName: project.GetProjectName(),
	}

	if remoteLogger != nil {
		marshal, _ := sonic.Marshal(data)
		remoteLogger.Warn("syslog", string(marshal))
	}
}

func SyslogError(title, msg string) {
	logId := GetLogId()
	context := ucontext.GetCurrentContext()
	data := LogData{
		Env:        project.GetRunTime(),
		Message:    msg,
		Title:      title,
		LogId:      logId,
		LogLevel:   "ERROR",
		LogTime:    utime.DateTimeFormat("2006-01-02 15:04:05.000"),
		TraceId:    context.TraceId,
		ReqId:      context.RequestId,
		LocalId:    context.LocalId,
		GoId:       system.GetGoRoutineId(),
		Stack:      system.GetStackFramesString(2, 0),
		SystemName: project.GetProjectName(),
	}

	if remoteLogger != nil {
		marshal, _ := sonic.Marshal(data)
		remoteLogger.Error("syslog", string(marshal))
	}

}
