// Package dict 包描述
// Author: wanlizhan
// Date: 2023/6/9
package dict

import (
	"fmt"
	"github.com/xxzhwl/wdk/cvt"
	"reflect"
	"strings"
)

// Mappam kv交换
func Mappam[K, V comparable](m map[K]V) map[V]K {
	res := make(map[V]K)
	for k, v := range m {
		res[v] = k
	}
	return res
}

// HaveKey map中是否有该key
func HaveKey[T comparable](m map[T]any, key T) bool {
	if len(m) == 0 {
		return false
	}
	if _, ok := m[key]; ok {
		return true
	}
	return false
}

// Keys 获取一个map中所有的key
func Keys[T comparable](m map[T]any) (res []T) {
	if len(m) == 0 {
		return
	}
	for k, _ := range m {
		res = append(res, k)
	}
	return
}

// Values 获取一个map中所有的Values
func Values[T comparable, V any](m map[T]V) (res []V) {
	if len(m) == 0 {
		return
	}
	for _, v := range m {
		res = append(res, v)
	}
	return
}

// ToSet 将map中的值转为set
func ToSet[K comparable, V any](m map[K]V) (res []V) {
	if len(m) == 0 {
		return
	}
	seen := make(map[K]int)
	for k, v := range m {
		if _, ok := seen[k]; !ok {
			seen[k] = 1
			res = append(res, v)
		}
	}
	return
}

// Merge 合并多个map
func Merge[K comparable, V any](maps ...map[K]V) (res map[K]V) {
	for _, m := range maps {
		for k, v := range m {
			res[k] = v
		}
	}
	return
}

// GetStringList map中获取key对应的value 一组slice
func GetStringList[K comparable](m map[K]any, key K, defaultValue []string) (res []string, err error) {
	if v, ok := m[key]; !ok {
		return defaultValue, nil
	} else {
		fv := reflect.ValueOf(v)
		switch fv.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < fv.Len(); i++ {
				res = append(res, cvt.S(fv.Index(i).Interface()))
			}
		}
	}
	return
}

// SL map中获取key对应的value 一组slice
func SL[K comparable](m map[K]any, key K) (res []string) {
	res, _ = GetStringList[K](m, key, []string{})
	return
}

// GetString 从map中获取键为key的string类型的value
func GetString(m map[string]any, key string, defaultValue string) (string, error) {
	if v, ok := m[key]; ok {
		return cvt.ToString(v, defaultValue)
	} else {
		splitKeySlice := strings.SplitN(key, ".", 2)
		if vTemp, okTemp := m[splitKeySlice[0]]; okTemp {
			if len(splitKeySlice) == 1 {
				return cvt.ToString(vTemp, defaultValue)
			}
			return GetString(cvt.M(vTemp), splitKeySlice[1], "")
		}
	}
	return defaultValue, fmt.Errorf("key %s in not exist in m", key)
}

// S 安全获取一个map中的string
func S(m map[string]any, key string) string {
	getString, _ := GetString(m, key, "")
	return getString
}

// GetInt64 从map中获取键为key的int64类型的value
func GetInt64(m map[string]any, key string, defaultValue int64) (int64, error) {
	if len(m) == 0 {
		return defaultValue, nil
	}
	if v, ok := m[key]; ok {
		return cvt.ToInt(v, defaultValue)
	} else {
		splitKeySlice := strings.SplitN(key, ".", 2)
		if vTemp, okTemp := m[splitKeySlice[0]]; okTemp {
			if len(splitKeySlice) == 1 {
				return cvt.ToInt(vTemp, defaultValue)
			}
			return GetInt64(cvt.M(vTemp), splitKeySlice[1], 0)
		}
	}
	return defaultValue, fmt.Errorf("key %s in not exist in m", key)
}

// I 安全获取一个map中的int64
func I(m map[string]any, key string) int64 {
	getInt64, _ := GetInt64(m, key, 0)
	return getInt64
}

func GetBool(m map[string]any, key string, defaultValue bool) (bool, error) {
	if len(m) == 0 {
		return defaultValue, nil
	}
	if v, ok := m[key]; ok {
		return cvt.ToBool(v, defaultValue)
	} else {
		splitKeySlice := strings.SplitN(key, ".", 2)
		if vTemp, okTemp := m[splitKeySlice[0]]; okTemp {
			if len(splitKeySlice) == 1 {
				return cvt.ToBool(vTemp, defaultValue)
			}
			return GetBool(cvt.M(vTemp), splitKeySlice[1], false)
		}
	}
	return defaultValue, fmt.Errorf("key %s in not exist in m", key)
}

// B 安全获取一个map中的bool
func B(m map[string]any, key string) bool {
	getBool, _ := GetBool(m, key, false)
	return getBool
}

func GetMap(m map[string]any, key string, defaultValue map[string]any) (map[string]any, error) {
	if len(m) == 0 {
		return defaultValue, nil
	}
	if v, ok := m[key]; ok {
		return cvt.ToMap(v, defaultValue)
	} else {
		splitKeySlice := strings.SplitN(key, ".", 2)
		if vTemp, okTemp := m[splitKeySlice[0]]; okTemp {
			if len(splitKeySlice) == 1 {
				return cvt.ToMap(vTemp, defaultValue)
			}
			return GetMap(cvt.M(vTemp), splitKeySlice[1], map[string]any{})
		}
	}
	return defaultValue, fmt.Errorf("key %s in not exist in m", key)
}

// M 安全获取一个map中的map
func M(m map[string]any, key string) map[string]any {
	getMap, _ := GetMap(m, key, map[string]any{})
	return getMap
}
