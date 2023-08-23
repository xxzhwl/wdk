package cnacos

import (
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/xxzhwl/wdk/project"
	"github.com/xxzhwl/wdk/uconfig"
	"github.com/xxzhwl/wdk/ulog"
	"os"
)

// NacosClient 查询使用的NacosClient
type NacosClient struct {
	Cli []map[string]config_client.IConfigClient
}

// NewNacosByConf 根据配置文件获取一个nacosClient
func NewNacosByConf(conf NacosConf) NacosClient {
	var serverConfigs = conf.ServerConfList
	dir := project.GetRootDir()
	sept := string(os.PathSeparator)

	var clientMap []map[string]config_client.IConfigClient
	for _, clientConf := range conf.ClientConfList {
		clientConfig := constant.ClientConfig{
			NamespaceId:         clientConf.NameSpace,
			TimeoutMs:           clientConf.TimeoutMs,
			NotLoadCacheAtStart: true,
			LogDir:              dir + sept + "var/log/nacos",
			CacheDir:            dir + sept + "var/log/",
			LogLevel:            "debug",
		}
		client, err := clients.NewConfigClient(vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		})
		if err != nil {
			ulog.Error("NewNacosByConf", fmt.Sprintf("%s NewConfigClientErr:%s", clientConf.NameSpace, err))
			continue
		}
		clientMap = append(clientMap, map[string]config_client.IConfigClient{clientConf.NameSpace: client})
	}
	return NacosClient{Cli: clientMap}
}

// NewNacos 获取一个默认cli
func NewNacos() (NacosClient, error) {
	conf := NacosConf{}
	err := uconfig.LoadConfToStruct("NacosCenter", &conf)
	if err != nil {
		return NacosClient{}, err
	}
	return NewNacosByConf(conf), nil
}

// GetRemoteConfigs 获取远程所有配置
func (n NacosClient) GetRemoteConfigs() (map[string]any, error) {
	configs := map[string]any{}
	configList := []*model.ConfigPage{}
	for _, clientMap := range n.Cli {
		for space, client := range clientMap {
			configPage, err := client.SearchConfig(vo.SearchConfigParam{Search: "blur", PageNo: 0, PageSize: 200})
			if err != nil {
				ulog.Error("GetRemoteConfigs", fmt.Sprintf("[%s]SearchConfigErr%s", space, err))
				continue
			}
			configList = append(configList, configPage)
		}
	}
	for _, page := range configList {
		for _, item := range page.PageItems {
			configs[item.DataId] = item.Content
		}
	}

	return configs, nil
}

// GetConfig 获取配置
func (n NacosClient) GetConfig(key string) (any, error) {
	for _, clientMap := range n.Cli {
		for space, client := range clientMap {
			configs, err := client.SearchConfig(vo.SearchConfigParam{Search: "accurate", DataId: key, PageNo: 0, PageSize: 1})
			if err != nil {
				ulog.Error("GetConfig", fmt.Sprintf("[%s-%s]SearchConfigErr%s", space, key, err))
				continue
			}
			if len(configs.PageItems) >= 1 {
				return configs.PageItems[0].Content, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("RemoteConfig-[%s]-NotFound", key))
}
