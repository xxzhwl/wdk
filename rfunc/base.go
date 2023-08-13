// Package rfunc 包描述
// Author: wanlizhan
// Date: 2023/7/16
package rfunc

import (
	"sync"
)

var locker sync.RWMutex

var router map[string]any

func init() {
	router = make(map[string]any)
}

// Register 注册对象路由
func Register(key string, objPtr any) {
	locker.Lock()
	defer locker.Unlock()
	router[key] = objPtr
}

func GetAction(key string) any {
	return router[key]
}
