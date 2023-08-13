// Package ulog 包描述
// Author: wanlizhan
// Date: 2023/6/11
package ulog

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/xxzhwl/wdk/project"
	"github.com/xxzhwl/wdk/system"
	"github.com/xxzhwl/wdk/ulog/local"
	"log"
	"os"
)

const (
	LogTypeInfo  = "INFO"
	LogTypeWarn  = "WARN"
	LogTypeErr   = "ERROR"
	LogTypePanic = "PANIC"
)

type LogData struct {
	Env        string
	Message    string
	Title      string
	LogId      string
	GoId       string
	LogType    string
	LogLevel   string
	LogTime    string
	LocalId    string
	TraceId    string
	ReqId      string
	Stack      string
	SystemName string
}

type ILogger interface {
	Info(title string, msg string)
	InfoF(title, template string, args ...any)
	Warn(title string, msg string)
	WarnF(title, template string, args ...any)
	Error(title string, msg string)
	ErrorF(title, template string, args ...any)
}

var localLogger ILogger

var remoteLogger ILogger

func init() {
	InjectLocalLogger(local.NewLocalLogger(local.StoreConfig{
		FilePath:    "var/log",
		FileNamePre: project.GetProjectName(),
		LocalTime:   true,
		Compress:    true,
	}))
}

// InjectLocalLogger 注入logger
func InjectLocalLogger(local ILogger) {
	localLogger = local
}

// InjectRemoteLogger 注入logger
func InjectRemoteLogger(remote ILogger) {
	remoteLogger = remote
}

// GetLogId  获取日志id
func GetLogId() string {
	name := os.Getenv("AppName")
	goId := system.GetGoRoutineId()
	u := uuid.New()
	return fmt.Sprintf("%s%s%s", name, goId, u.String())
}

func Info(title, msg string) {
	if localLogger != nil {
		localLogger.Info(title, msg)
	}
	goId := system.GetGoRoutineId()
	log.Printf("[%s][PID=%d][GOID=%s][%s]%s\n", LogTypeInfo, os.Getpid(), goId, title, msg)
	SyslogInfo(title, msg)
}

func InfoF(title, template string, args ...any) {
	if localLogger != nil {
		localLogger.InfoF(title, template, args...)
	}
	msg := fmt.Sprintf(template, args...)
	goId := system.GetGoRoutineId()
	log.Printf("[%s][PID=%d][GOID=%s][%s]%s\n", LogTypeInfo, os.Getpid(), goId, title, msg)
	SyslogInfo(title, msg)
}

func Warn(title, msg string) {
	if localLogger != nil {
		localLogger.Warn(title, msg)
	}
	goId := system.GetGoRoutineId()
	log.Printf("[%s][PID=%d][GOID=%s][%s]%s\n", LogTypeWarn, os.Getpid(), goId, title, msg)
	SyslogWarn(title, msg)
}

func WarnF(title, template string, args ...any) {
	if localLogger != nil {
		localLogger.WarnF(title, template, args...)
	}
	msg := fmt.Sprintf(template, args...)
	goId := system.GetGoRoutineId()
	log.Printf("[%s][PID=%d][GOID=%s][%s]%s\n", LogTypeWarn, os.Getpid(), goId, title, msg)
	SyslogWarn(title, msg)
}

func Error(title, msg string) {
	if localLogger != nil {
		localLogger.Error(title, msg)
	}
	goId := system.GetGoRoutineId()
	log.Printf("[%s][PID=%d][GOID=%s][%s]%s\n", LogTypeErr, os.Getpid(), goId, title, msg)
	SyslogError(title, msg)
}

func ErrorF(title, template string, args ...any) {
	if localLogger != nil {
		localLogger.ErrorF(title, template, args...)
	}
	msg := fmt.Sprintf(template, args...)
	goId := system.GetGoRoutineId()
	log.Printf("[%s][PID=%d][GOID=%s][%s]%s\n", LogTypeErr, os.Getpid(), goId, title, msg)
	SyslogError(title, msg)
}
