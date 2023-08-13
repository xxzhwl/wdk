// Package server 包描述
// Author: wanlizhan
// Date: 2023/7/22
package server

import (
	"fmt"
	"github.com/xxzhwl/wdk"
	"github.com/xxzhwl/wdk/ucontext"
	"github.com/xxzhwl/wdk/ulog"
	"sync"
	"sync/atomic"
	"time"
)

var asyncCounter int32

func incrCount() {
	atomic.AddInt32(&asyncCounter, 1)
}

func decrCount() {
	atomic.AddInt32(&asyncCounter, -1)
}

func CounterZero() bool {
	return LoadCounter() == 0
}

func LoadCounter() int32 {
	return atomic.LoadInt32(&asyncCounter)
}

// Go 启动协程干活儿
func Go(fn func()) {
	_go(nil, fn)
}

// GoWait 启动协程并等待
func GoWait(wg *sync.WaitGroup, fn func()) {
	_go(wg, fn)
}

func _go(wg *sync.WaitGroup, fn func()) {
	incrCount()
	go func(f func(), ctx *ucontext.Context) {
		defer func() {
			if wg != nil {
				wg.Done()
			}
		}()
		defer decrCount()
		defer wdk.CatchPanic()
		ucontext.ReSetContext(ctx)
		defer ucontext.RemoveContext()
		f()
	}(fn, ucontext.GetCurrentContext())
}

func _goWithArg[T any](wg *sync.WaitGroup, fn func(T), arg T) {
	incrCount()
	go func() {

	}()
}

// GracefulExit 优雅推出
func GracefulExit(timeout time.Duration) error {
	ulog.InfoF("GracefulExit", "等待时间%.1f秒", timeout.Seconds())
	startTime := time.Now()
	for time.Now().Before(startTime.Add(timeout)) {
		time.Sleep(200 * time.Millisecond)
		if CounterZero() {
			ulog.Info("GracefulExit", "优雅退出.")
			return nil
		}
	}
	return fmt.Errorf("服务运行中，等待优雅 :%s 退出超时", timeout.String())
}
