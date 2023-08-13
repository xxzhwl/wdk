// Package ugorm db连接配置
// Author: wanlizhan
// Date: 2023/3/9 22:58
package ugorm

import (
	"fmt"
	"github.com/xxzhwl/wdk/uconfig"
	"github.com/xxzhwl/wdk/ulog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const maxConnNum = 2000
const maxLifeTime = 20 * time.Second
const maxIdleNum = 2

type GormMysql struct {
	Db *gorm.DB
}

func newMysql(cfg MysqlConfig) (*GormMysql, error) {
	glog := ulog.GormLogger{}
	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{Logger: glog.LogMode(logger.Info)})
	if err != nil {
		return nil, err
	}
	dbInstance, err := db.DB()
	if err != nil {
		return nil, err
	}
	if cfg.MaxConnNum == 0 {
		cfg.MaxConnNum = maxConnNum
	}
	if cfg.MaxLifeTime <= maxLifeTime {
		cfg.MaxLifeTime = maxLifeTime
	}
	if cfg.MaxIdleNum == 0 {
		cfg.MaxIdleNum = maxIdleNum
	}
	dbInstance.SetMaxOpenConns(cfg.MaxConnNum)
	dbInstance.SetMaxIdleConns(cfg.MaxIdleNum)
	dbInstance.SetConnMaxLifetime(cfg.MaxLifeTime)
	return &GormMysql{db}, nil
}

// NewMysqlDefault 获取默认mysql连接
func NewMysqlDefault() (*GormMysql, error) {
	return NewMysqlBySchema("Default")
}

// NewMysqlBySchema 获取一个Db连接实例
func NewMysqlBySchema(schema string) (*GormMysql, error) {
	conf := MysqlConfig{}
	err := uconfig.LoadConfToStruct("Mysql."+schema+".Master", &conf)
	if err != nil {
		return nil, fmt.Errorf("schema:%s Not Found", schema)
	}
	return newMysql(conf)
}

// NewMysqlWithConfig 获取一个Db连接实例
func NewMysqlWithConfig(cfg MysqlConfig) (*GormMysql, error) {
	return newMysql(cfg)
}
