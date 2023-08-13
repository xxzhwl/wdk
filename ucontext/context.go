// Package ucontext 包描述
// Author: wanlizhan
// Date: 2023/7/16
package ucontext

import (
	"github.com/google/uuid"
	"github.com/xxzhwl/wdk/project"
	"github.com/xxzhwl/wdk/system"
	"github.com/xxzhwl/wdk/ustr"
	"sync"
)

var locker sync.RWMutex

var contextMap map[string]*Context

// Context 协程上下文信息
type Context struct {
	TraceId   string
	RequestId string
	LocalId   string
}

func init() {
	contextMap = make(map[string]*Context)
}

// BuildContext 构建上下文信息
// 如果该协程已经存在上下文,则直接返回
// 不存在则添加一个上下文映射到该协程Id
func BuildContext() *Context {
	goId := system.GetGoRoutineId()
	c := &Context{
		LocalId: uuid.New().String(),
	}
	locker.Lock()
	defer locker.Unlock()
	if _, ok := contextMap[goId]; !ok {
		contextMap[goId] = c
	} else {
		c = contextMap[goId]
	}
	return c
}

// GetContext 获取某个携程的上下文，前提是要知道该协程的Id
func GetContext(goId string) *Context {
	locker.RLock()
	defer locker.RUnlock()
	c := contextMap[goId]
	return c
}

// SetContext 为某个协程设置上下文
func SetContext(goId string, c *Context) {
	locker.Lock()
	defer locker.Unlock()
	contextMap[goId] = c
}

// ReSetContext 为当前协程重置上下文为c
// 通常在启动一个新的协程以后，为了追踪日志需要将父协程的上下文复制给它
func ReSetContext(c *Context) {
	goId := system.GetGoRoutineId()
	locker.Lock()
	defer locker.Unlock()
	contextMap[goId] = c
}

// RemoveContext 移除当前协程的上下文
// 协程运行结束以后要及时删掉上下文，避免内存泄漏
func RemoveContext() {
	goId := system.GetGoRoutineId()
	locker.Lock()
	defer locker.Unlock()
	delete(contextMap, goId)
}

// GetCurrentContext 获取当前协程的上下文
func GetCurrentContext() *Context {
	goId := system.GetGoRoutineId()
	locker.RLock()
	defer locker.RUnlock()
	if v, ok := contextMap[goId]; ok {
		return v
	}
	return BuildContext()
}

// NewTraceId 生成traceId
func NewTraceId() string {
	return "x-traceId-" + project.GetProjectName() + ustr.ABCCode(8)
}
