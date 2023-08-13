// Package local
// Author: wanlizhan
// Date: 2023-02-19 09:37:07
package local

import (
	"fmt"
	"github.com/xxzhwl/wdk/project"
	"github.com/xxzhwl/wdk/system"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strconv"
	"time"
)

type LocalULog struct {
	*zap.Logger
}

// NewLocalLogger 注入自定义本地存储Logger
func NewLocalLogger(conf StoreConfig) *LocalULog {
	l := LocalULog{}
	core := zapcore.NewCore(initCore(), initWriter(conf), zapcore.DebugLevel)
	l.Logger = zap.New(core, zap.AddStacktrace(zap.ErrorLevel))
	return &l
}

func initCore() zapcore.Encoder {
	conf := zap.NewProductionEncoderConfig()
	conf.TimeKey = "time"
	conf.EncodeLevel = zapcore.CapitalLevelEncoder
	conf.CallerKey = "caller"
	conf.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(caller.FullPath())
	}
	conf.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Local().Format(time.DateTime))
	}
	return zapcore.NewJSONEncoder(conf)
}

// StoreConfig 本地日志存储配置
type StoreConfig struct {
	FilePath    string
	FileNamePre string
	MaxSize     int
	MaxAge      int
	MaxBackups  int
	LocalTime   bool
	Compress    bool
}

func initWriter(conf StoreConfig) zapcore.WriteSyncer {
	dir := project.GetRootDir()
	sept := string(os.PathSeparator)
	logger := &lumberjack.Logger{
		Filename: dir + sept + conf.FilePath + sept + conf.FileNamePre + "-" +
			time.Now().Format(time.DateOnly) + ".txt",
		MaxSize:    conf.MaxSize,
		MaxAge:     conf.MaxAge,
		MaxBackups: conf.MaxBackups,
		LocalTime:  conf.LocalTime,
		Compress:   conf.Compress,
	}

	sync := zapcore.AddSync(logger)
	return zapcore.NewMultiWriteSyncer(sync)
}

func (u LocalULog) Info(title string, msg string) {
	u.Logger.Info(msg, zap.String("Title", title), zap.String("GOID", system.GetGoRoutineId()), zap.String("PID", strconv.Itoa(os.Getpid())))
}

func (u LocalULog) InfoF(title, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	u.Logger.Info(msg, zap.String("Title", title), zap.String("GOID",
		system.GetGoRoutineId()), zap.String("PID", strconv.Itoa(os.Getpid())))
}

func (u LocalULog) Error(title, msg string) {
	u.Logger.Error(msg, zap.String("Title", title), zap.String("GOID", system.GetGoRoutineId()), zap.String("PID",
		strconv.Itoa(os.Getpid())))
}

func (u LocalULog) ErrorF(title, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	u.Logger.Error(msg, zap.String("Title", title), zap.String("GOID",
		system.GetGoRoutineId()), zap.String("PID", strconv.Itoa(os.Getpid())))
}

func (u LocalULog) Panic(title, msg string) {
	u.Logger.Panic(msg, zap.String("Title", title), zap.String("GOID", system.GetGoRoutineId()), zap.String("PID",
		strconv.Itoa(os.Getpid())))
}

func (u LocalULog) PanicF(title, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	u.Logger.Panic(msg, zap.String("Title", title), zap.String("GOID",
		system.GetGoRoutineId()), zap.String("PID", strconv.Itoa(os.Getpid())))
}

func (u LocalULog) Warn(title string, msg string) {
	u.Logger.Warn(msg, zap.String("Title", title), zap.String("GOID", system.GetGoRoutineId()), zap.String("PID",
		strconv.Itoa(os.Getpid())))
}

func (u LocalULog) WarnF(title, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	u.Logger.Warn(msg, zap.String("Title", title), zap.String("GOID",
		system.GetGoRoutineId()), zap.String("PID", strconv.Itoa(os.Getpid())))
}
