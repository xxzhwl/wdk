package cnacos

import "github.com/nacos-group/nacos-sdk-go/v2/common/constant"

type ClientConf struct {
	NameSpace string
	UserName  string
	Password  string
	TimeoutMs uint64
}

type NacosConf struct {
	ClientConfList []ClientConf
	ServerConfList []constant.ServerConfig
}
