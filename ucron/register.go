// Package ucron 包描述
// Author: wanlizhan
// Date: 2023/7/23
package ucron

import (
	"github.com/robfig/cron/v3"
	"github.com/xxzhwl/wdk"
	"github.com/xxzhwl/wdk/ucontext"
)

type CronCenter struct {
	manager *cron.Cron
}

// NewCronCenter 构建一个新的cron-center
func NewCronCenter() *CronCenter {
	return &CronCenter{manager: cron.New(cron.WithSeconds())}
}

// Register 注册CRON任务
func (c *CronCenter) Register(spec string, ctx *ucontext.Context, fn func()) error {
	_, err := c.RegisterRetId(spec, ctx, fn)
	return err
}

// RegisterRetId 注册CRON任务并返回Id
func (c *CronCenter) RegisterRetId(spec string, ctx *ucontext.Context, fn func()) (int, error) {
	entryId, err := c.manager.AddFunc(spec, func() {
		defer wdk.CatchPanic()
		ucontext.ReSetContext(ctx)
		fn()
		ucontext.RemoveContext()
	})
	if err != nil {
		return 0, err
	}
	return int(entryId), err
}

// Run 运行任务
func (c *CronCenter) Run() {
	c.manager.Run()
}
