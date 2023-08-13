// Package github.com/xxzhwl/wdk 包描述
// Author: wanlizhan
// Date: 2023/6/18
package ucache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache()
	cache.Set("name", "wanli", 10*time.Second)
	cache.Set("age", 25, 15*time.Second)
	cache.Set("hobby", map[string]any{"爱好": []string{"篮球", "足球"}}, 30*time.Second)
	time.Sleep(1 * time.Minute)
}
