// Package redis redis相关工具
// Author: wanlizhan
// Date: 2023/3/9 22:52
package redis

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/go-redis/redis"
	"github.com/xxzhwl/wdk/uconfig"
	"github.com/xxzhwl/wdk/ulog"
	"strconv"
	"time"
)

// RedisCli Redis客户端
type RedisCli struct {
	Client *redis.Client
}

// Config Redis配置
type Config struct {
	Addr     string
	Port     int
	Password string
	Db       int
}

// NewConfig 搞一个新配置
func NewConfig(addr, password string, port int) *Config {
	return &Config{
		Addr:     addr,
		Port:     port,
		Password: password,
	}
}

func (c *Config) WithDb(db int) *Config {
	c.Db = db
	return c
}

// NewRedis 指定配置获取Redis
func NewRedis(conf Config) RedisCli {
	options := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Addr, conf.Port),
		Password: conf.Password,
		DB:       conf.Db,
	}
	return RedisCli{redis.NewClient(&options)}
}

// NewDefaultRedis 获取默认Redis
func NewDefaultRedis() RedisCli {
	return NewRedisBySchema("Default")
}

// NewRedisBySchema 指定配置获取Redis
func NewRedisBySchema(schema string) RedisCli {
	config := Config{}
	err := uconfig.LoadConfToStruct("Redis."+schema+".Master", &config)
	if err != nil {
		ulog.Error("GetDefaultRedis", err.Error())
	}
	return NewRedis(config)
}

// I 拿到redis中命中key的int64值
func (r *RedisCli) I(key string) int64 {
	res, err := r.Client.Get(key).Result()
	if err != nil {
		return 0
	}
	ri, err := strconv.Atoi(res)
	if err != nil {
		return 0
	}
	return int64(ri)
}

// S 拿到redis中命中key的string值
func (r *RedisCli) S(key string) (x string) {
	r.GetObj(key, &x)
	return
}

// GetObj 拿到redis中命中key的值并赋给val
func (r *RedisCli) GetObj(key string, val any) {
	res, _ := r.Client.Get(key).Result()
	_ = sonic.Unmarshal([]byte(res), &val)
}

// Set 存储一个KV
func (r *RedisCli) Set(key string, val any) error {
	bts, err := sonic.Marshal(val)
	if err != nil {
		return err
	}
	if _, err = r.Client.Set(key, string(bts), 0).Result(); err != nil {
		return err
	}
	return nil
}

// SetEx 带过期时间的存储一个KV
func (r *RedisCli) SetEx(key string, val any, expiration time.Duration) error {
	bts, err := sonic.Marshal(val)
	if err != nil {
		return err
	}
	if _, err = r.Client.Set(key, string(bts), expiration).Result(); err != nil {
		return err
	}
	return nil
}

// Del 删除一个KV
func (r *RedisCli) Del(key string) error {
	if _, err := r.Client.Del(key).Result(); err != nil {
		return err
	}
	return nil
}

// GetExpireTime 获取过期时间
func (r *RedisCli) GetExpireTime(key string) (time.Duration, error) {
	return r.Client.TTL(key).Result()
}

// Expire 设置过期时间
func (r *RedisCli) Expire(key string, expiration time.Duration) error {
	if _, err := r.Client.Expire(key, expiration).Result(); err != nil {
		return err
	}
	return nil
}

// Incr 递增key的value
func (r *RedisCli) Incr(key string) error {
	if _, err := r.Client.Incr(key).Result(); err != nil {
		return err
	}
	return nil
}
