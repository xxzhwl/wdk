// Package apollo 包描述
// Author: wanlizhan
// Date: 2023/6/11
package apollo

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/xxzhwl/wdk/dict"
	"github.com/xxzhwl/wdk/uconfig"
	"github.com/xxzhwl/wdk/ulog"
)

type ApolloClient struct {
	appId     string
	cluster   string
	ip        string
	nameSpace string
	secret    string
}

type apolloClientApply func(a *ApolloClient)

func NewDefaultApolloClient() *ApolloClient {
	apolloConf, err := uconfig.LocalMWithErr("Apollo")
	if err != nil {
		ulog.ErrorF("NewDefaultApolloClient", "初始化阿波罗配置中心失败，Err:%s", err.Error())
	}
	app, err := uconfig.LocalMWithErr("App")
	if err != nil {
		ulog.ErrorF("NewDefaultApolloClient", "初始化阿波罗配置中心失败，Err:%s", err.Error())
	}
	return NewApolloClient(WithAppId(dict.S(app, "Name")), WithCluster(dict.S(app, "Env")), WithIp(dict.S(apolloConf,
		"Ip")), WithNameSpace(dict.S(apolloConf, "NameSpace")), WithSecret(dict.S(apolloConf, "Secret")))
}

func NewApolloClient(f ...apolloClientApply) *ApolloClient {
	a := ApolloClient{}
	for _, apply := range f {
		apply(&a)
	}
	return &a
}

func WithAppId(appId string) apolloClientApply {
	return func(a *ApolloClient) {
		a.appId = appId
	}
}

func WithCluster(cluster string) apolloClientApply {
	return func(a *ApolloClient) {
		a.cluster = cluster
	}
}

func WithIp(ip string) apolloClientApply {
	return func(a *ApolloClient) {
		a.ip = ip
	}
}

func WithNameSpace(nameSpace string) apolloClientApply {
	return func(a *ApolloClient) {
		a.nameSpace = nameSpace
	}
}

func WithSecret(secret string) apolloClientApply {
	return func(a *ApolloClient) {
		a.secret = secret
	}
}

func (a *ApolloClient) GetRemoteConfigs() (map[string]any, error) {
	return nil, nil
}

func (a *ApolloClient) GetConfig(key string) (any, error) {
	c := &config.AppConfig{
		AppID:          a.appId,
		Cluster:        a.cluster,
		IP:             a.ip,
		NamespaceName:  a.nameSpace,
		IsBackupConfig: true,
		Secret:         a.secret,
	}
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		ulog.ErrorF("ApolloClient", "StartWithConfigErr:%s", err.Error())
		return nil, err
	}

	cache := client.GetConfigCache(c.NamespaceName)
	if cache == nil {
		ulog.Error("ApolloClient", "GetConfigCacheFailed,CacheIsNil")
		return nil, nil
	}
	return cache.Get(key)
}
