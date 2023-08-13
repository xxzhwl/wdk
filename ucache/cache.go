// Package ucache 包描述
// Author: wanlizhan
// Date: 2023/6/18
package ucache

import (
	"github.com/xxzhwl/wdk/cvt"
	"sync"
	"time"
)

var m map[string]cacheData
var defaultExpire time.Duration
var defaultScan time.Duration
var enableExpire bool

type cacheApply func(c *Cache)

type cacheData struct {
	value     any
	clearTime time.Time
	canClear  bool
}

// Cache 缓存
type Cache struct {
	duration time.Duration

	scanDuration time.Duration

	data map[string]cacheData

	locker sync.RWMutex

	enableGlobalExpire bool
}

func init() {
	m = make(map[string]cacheData)
	defaultExpire = 30 * time.Second
	defaultScan = 5 * time.Second
	enableExpire = true
}

func NewCache(applyFunc ...cacheApply) *Cache {
	c := Cache{data: m, duration: defaultExpire, enableGlobalExpire: enableExpire, scanDuration: defaultScan}
	for _, fun := range applyFunc {
		fun(&c)
	}
	if c.enableGlobalExpire {
		go func() {
			c.clear()
		}()
	}
	return &c
}

// DisableGlobalExpire 关闭全局扫描过期
func DisableGlobalExpire() cacheApply {
	return func(c *Cache) {
		c.enableGlobalExpire = false
	}
}

// SetDuration 设置全局过期时间
func SetDuration(expire time.Duration) cacheApply {
	return func(c *Cache) {
		c.duration = expire
	}
}

// SetScanDuration 设置全局扫描时间
func SetScanDuration(expire time.Duration) cacheApply {
	return func(c *Cache) {
		c.scanDuration = expire
	}
}

func (c *Cache) Set(key string, v any, expire time.Duration) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if expire == 0 {
		expire = defaultExpire
	}
	c.data[key] = cacheData{
		value:     v,
		clearTime: time.Now().Add(expire),
		canClear:  true,
	}
}

// SetNoExpire 不过期的数据
func (c *Cache) SetNoExpire(key string, v any) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.data[key] = cacheData{
		value: v,
	}
}

func (c *Cache) S(key string) string {
	c.locker.RLock()
	res := c.data[key].value
	c.locker.RUnlock()
	return cvt.S(res)
}

func (c *Cache) I(key string) int64 {
	c.locker.RLock()
	res := c.data[key].value
	c.locker.RUnlock()
	return cvt.I(res)
}

func (c *Cache) M(key string) map[string]any {
	c.locker.RLock()
	res := c.data[key].value
	c.locker.RUnlock()
	return cvt.M(res)
}

func (c *Cache) clear() {
	ticker := time.NewTicker(c.scanDuration)
	for range ticker.C {
		for k, cd := range c.data {
			if time.Now().After(cd.clearTime) && cd.canClear {
				delete(c.data, k)
			}
		}
	}
}
