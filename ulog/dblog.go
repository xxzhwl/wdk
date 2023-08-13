// Package ugorm 包描述
// Author: wanlizhan
// Date: 2023/7/18
package ulog

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/xxzhwl/wdk/project"
	"github.com/xxzhwl/wdk/system"
	"github.com/xxzhwl/wdk/ucontext"
	"github.com/xxzhwl/wdk/ustr"
	"gorm.io/gorm/logger"
	"log"
	"strconv"
	"time"
)

type DbLogData struct {
	SqlStr     string
	SqlType    string
	LocalId    string
	ReqId      string
	TraceId    string
	Duration   string
	Env        string
	Affected   int64
	GoId       string
	Stack      string
	SystemName string
}

type GormLogger struct {
	LogLevel logger.LogLevel
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	temp := *g
	temp.LogLevel = level
	return &temp
}

func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	fmt.Println(msg, data)
	if g.LogLevel >= logger.Info {
		log.Println(msg)
		for _, datum := range data {
			log.Println(datum)
		}
	}
}

func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Warn {
		log.Println(msg)
		for _, datum := range data {
			log.Println(datum)
		}
	}
}

func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Error {
		log.Println(msg)
		for _, datum := range data {
			log.Println(datum)
		}
	}
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64),
	err error) {
	sql, affected := fc()
	duration := time.Since(begin)
	context := ucontext.GetCurrentContext()
	data := DbLogData{
		SqlStr:     sql,
		LocalId:    context.LocalId,
		ReqId:      context.RequestId,
		Env:        project.GetRunTime(),
		TraceId:    context.TraceId,
		Duration:   strconv.FormatInt(duration.Milliseconds(), 10) + "ms",
		Affected:   affected,
		GoId:       system.GetGoRoutineId(),
		Stack:      system.GetStackFramesString(2, 0),
		SystemName: project.GetProjectName(),
	}
	if ustr.Contains(sql, "insert") {
		data.SqlType = "Insert"
	} else if ustr.Contains(sql, "select") {
		data.SqlType = "Select"
	} else if ustr.Contains(sql, "update") {
		data.SqlType = "Update"
	} else if ustr.Contains(sql, "delete") {
		data.SqlType = "Delete"
	}
	marshal, _ := sonic.Marshal(data)
	Info("DbLog", string(marshal))

	if remoteLogger != nil {
		remoteLogger.Info("dblog", string(marshal))
	}
}
