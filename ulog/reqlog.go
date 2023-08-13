// Package ulog 包描述
// Author: wanlizhan
// Date: 2023/7/19
package ulog

import (
	"github.com/bytedance/sonic"
	"github.com/xxzhwl/wdk/project"
	"github.com/xxzhwl/wdk/system"
	"github.com/xxzhwl/wdk/utime"
)

type ReqLogData struct {
	Env        string
	HttpMethod string
	ErrMsg     string
	Success    bool
	Code       string
	Status     string
	GoId       string
	LogId      string
	Title      string
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
	StartTime  string
	EndTime    string
	Duration   string
}

func ReqLog(data ReqLogData) {
	logId := GetLogId()
	data.Env = project.GetRunTime()
	data.LogId = logId
	data.GoId = system.GetGoRoutineId()
	data.LogTime = utime.DateTime()
	data.SystemName = project.GetProjectName()
	marshal, _ := sonic.Marshal(data)
	Info("ReqLogData", string(marshal))
	if remoteLogger != nil {
		remoteLogger.Info("reqlog", string(marshal))
	}
}
