// Package ugorm 包描述
// Author: wanlizhan
// Date: 2023/7/2
package ugorm

import "time"

// MysqlConfig 数据源结构
type MysqlConfig struct {
	Dsn         string
	MaxConnNum  int
	MaxIdleNum  int
	MaxLifeTime time.Duration
}

// NewMysqlConfig 获取一个Db连接实例
func NewMysqlConfig(dsn string, maxConn, maxIdle int, maxLifeTime time.Duration) MysqlConfig {
	return MysqlConfig{Dsn: dsn, MaxConnNum: maxConn, MaxIdleNum: maxIdle, MaxLifeTime: maxLifeTime}
}
