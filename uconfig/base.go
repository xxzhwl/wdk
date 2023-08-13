// Package uconfig 包描述
// Author: wanlizhan
// Date: 2023/6/11
package uconfig

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/xxzhwl/wdk"
	"github.com/xxzhwl/wdk/cvt"
	"github.com/xxzhwl/wdk/project"
	"github.com/xxzhwl/wdk/ucache"
	"github.com/xxzhwl/wdk/ucontext"
	"github.com/xxzhwl/wdk/ulog"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"time"
)

const (
	LocalConfigPath = "/config/"
	LocalConfigName = "config.yaml"

	ConfigCacheKey      = "ConfigCache"
	LocalConfigCacheKey = "LocalConfigCache"
)

var cache = ucache.NewCache()

type IRemoteConfig interface {
	GetRemoteConfigs() (map[string]any, error)

	GetConfig(key string) (any, error)
}
type ILocalConfig interface {
	GetLocalConfigs() (map[string]any, error)
}

var LocalConfigReader ILocalConfig

var RemoteConfigReader IRemoteConfig

var RefreshLogEnable bool

func init() {
	ucontext.BuildContext()
	res, err := getLocalConfigs()
	if err != nil {
		ulog.Error("InitLocalConfigs", err.Error())
	}
	cache.SetNoExpire(ConfigCacheKey, res)
}

func Refresh() {
	ctx := ucontext.BuildContext()
	go func(ctx *ucontext.Context) {
		defer func() {
			wdk.CatchPanic()
			ucontext.RemoveContext()
		}()
		ucontext.ReSetContext(ctx)
		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-ticker.C:
				_, err := getAllConfigs()
				if err != nil {
					ulog.Error("RefreshConfig", err.Error())
				} else if RefreshLogEnable {
					ulog.Info("RefreshConfig", "成功")
				}
			}
		}
	}(ctx)
}

// GetLocalConfigs 获取本地所有配置
func GetLocalConfigs() (map[string]any, error) {
	if LocalConfigReader != nil {
		return LocalConfigReader.GetLocalConfigs()
	}
	return getLocalConfigs()
}

func getLocalConfigs() (res map[string]any, err error) {
	filePath := project.GetRootDir() + LocalConfigPath + LocalConfigName
	file, err := os.ReadFile(filePath)
	if err != nil {
		ulog.Error("GetLocalConfigs", err.Error())
		return nil, err
	}
	if err = yaml.Unmarshal(file, &res); err != nil {
		ulog.Error("GetLocalConfigs", err.Error())
		return nil, err
	}
	cache.SetNoExpire(LocalConfigCacheKey, res)
	return
}

func GetRemoteConfigs() (map[string]any, error) {
	if RemoteConfigReader == nil {
		ulog.Error("GetRemoteConfigs", "IRemoteConfig interface not implemented")
		return nil, nil
	}
	return RemoteConfigReader.GetRemoteConfigs()
}

// GetAllConfigs 获取所有配置
func GetAllConfigs() (map[string]any, error) {
	return cache.M(ConfigCacheKey), nil
}

// getAllConfigs 获取所有配置
func getAllConfigs() (map[string]any, error) {
	res := cache.M(ConfigCacheKey)
	localRes, err := GetLocalConfigs()
	if err != nil {
		return nil, err
	}
	remoteRes, err := GetRemoteConfigs()
	if err != nil {
		return nil, err
	}
	if len(remoteRes) != 0 {
		res = remoteRes
	}
	for k, v := range localRes {
		res[k] = v
	}
	cache.SetNoExpire(ConfigCacheKey, res)
	return res, nil
}

func getConfig(key string) (any, error) {
	m := cache.M(ConfigCacheKey)
	if len(m) == 0 {
		m = make(map[string]any)
	}
	if v, ok := m[key]; ok {
		return v, nil
	} else {
		vt := getConfigFromMap(key, m)
		if vt != nil {
			return vt, nil
		}
	}
	if RemoteConfigReader != nil {
		config, err := RemoteConfigReader.GetConfig(key)
		if err != nil {
			return "", err
		}
		m[key] = config
		cache.Set(ConfigCacheKey, m, time.Minute)
		return config, nil
	}
	return nil, fmt.Errorf("未查询到该配置%s", key)
}

func getConfigFromMap(key string, m map[string]any) any {
	splitKeySlice := strings.SplitN(key, ".", 2)
	if vTemp, okTemp := m[splitKeySlice[0]]; okTemp {
		if len(splitKeySlice) == 1 {
			return vTemp
		}
		return getConfigFromMap(splitKeySlice[1], cvt.M(vTemp))
	}
	return nil
}

func getLocalConfig(key string) (any, error) {
	m := cache.M(LocalConfigCacheKey)
	if v, ok := m[key]; ok {
		return v, nil
	}
	res, err := getLocalConfigs()
	if err != nil {
		return nil, err
	}
	if v, ok := res[key]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("config:%s not found at local", key)
}

func LocalS(key string) string {
	config, err := getLocalConfig(key)
	if err != nil {
		return ""
	}
	return cvt.S(config)
}

func S(key string) string {
	config, err := getConfig(key)
	if err != nil {
		return ""
	}
	return cvt.S(config)
}

func I(key string) int64 {
	config, err := getConfig(key)
	if err != nil {
		return 0
	}
	return cvt.I(config)
}

func LocalI(key string) int64 {
	config, err := getLocalConfig(key)
	if err != nil {
		return 0
	}
	return cvt.I(config)
}

func M(key string) map[string]any {
	config, err := getConfig(key)
	if err != nil {
		return map[string]any{}
	}
	return cvt.M(config)
}

func LocalM(key string) map[string]any {
	config, err := getLocalConfig(key)
	if err != nil {
		return map[string]any{}
	}
	return cvt.M(config)
}

func MS(key string) map[string]string {
	config, err := getConfig(key)
	if err != nil {
		return map[string]string{}
	}
	return cvt.MS(config)
}

func LocalMS(key string) map[string]string {
	config, err := getLocalConfig(key)
	if err != nil {
		return map[string]string{}
	}
	return cvt.MS(config)
}

func LocalSWithErr(key string) (string, error) {
	config, err := getLocalConfig(key)
	if err != nil {
		return "", err
	}
	return cvt.ToString(config, "")
}

func SE(key string) (string, error) {
	config, err := getConfig(key)
	if err != nil {
		return "", err
	}
	return cvt.ToString(config, "")
}

func IE(key string) (int64, error) {
	config, err := getConfig(key)
	if err != nil {
		return 0, err
	}
	return cvt.ToInt(config, 0)
}

func LocalIWithErr(key string) (int64, error) {
	config, err := getLocalConfig(key)
	if err != nil {
		return 0, err
	}
	return cvt.ToInt(config, 0)
}

func ME(key string) (map[string]any, error) {
	config, err := getConfig(key)
	if err != nil {
		return map[string]any{}, err
	}
	return cvt.ToMap(config, map[string]any{})
}

func LocalMWithErr(key string) (map[string]any, error) {
	config, err := getLocalConfig(key)
	if err != nil {
		return map[string]any{}, err
	}
	return cvt.ToMap(config, map[string]any{})
}

func MSE(key string) (map[string]string, error) {
	config, err := getConfig(key)
	if err != nil {
		return map[string]string{}, err
	}
	return cvt.ToMapS(config, map[string]string{})
}

func LocalMSWithErr(key string) (map[string]string, error) {
	config, err := getLocalConfig(key)
	if err != nil {
		return map[string]string{}, err
	}
	return cvt.ToMapS(config, map[string]string{})
}

// LoadConfToStruct 将config转为结构体
func LoadConfToStruct(key string, value any) error {
	config, err := getConfig(key)
	if err != nil {
		return err
	}
	if err = sonic.Unmarshal([]byte(cvt.S(config)), &value); err != nil {
		return err
	}
	return nil
}

// StringListWithErr 获取一个json格式的字符串数组
func StringListWithErr(key string, defaultValue []string) ([]string, error) {
	config, err := getConfig(key)
	if err != nil {
		return defaultValue, err
	}
	var res []string
	if err = sonic.Unmarshal([]byte(config.(string)), &res); err != nil {
		return defaultValue, err
	}
	return res, nil
}

// StringList 获取一个json格式的字符串数组
func StringList(key string) []string {
	res, _ := StringListWithErr(key, []string{})
	return res
}
