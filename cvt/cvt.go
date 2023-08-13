// Package cvt 包描述
// Author: wanlizhan
// Date: 2023/6/9
package cvt

import (
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ToBool 转换为bool类型
func ToBool(value any, defaultValue bool) (bool, error) {
	if value == nil {
		return false, nil
	}
	v := reflect.ValueOf(value)
	switch value.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return v.Uint() != 0, nil
	case int, int8, int16, int32, int64:
		return v.Int() != 0, nil
	case float32, float64:
		return v.Float() != 0, nil
	case bool:
		return v.Bool(), nil
	case string:
		parseBool, err := strconv.ParseBool(v.String())
		if err != nil {
			return defaultValue, fmt.Errorf("%v(%T)->bool Err:%s", value, value, err.Error())
		}
		return parseBool, nil
	default:
		return false, fmt.Errorf("%v(%T)->bool Err", value, value)
	}
}

func B(value any) bool {
	toBool, _ := ToBool(value, false)
	return toBool
}

// ToInt 转换value为int64类型，如果有误则返回defaultValue
func ToInt(value any, defaultValue int64) (int64, error) {
	if value == nil {
		return 0, nil
	}
	v := reflect.ValueOf(value)

	switch value.(type) {
	case int, int8, int16, int32, int64:
		return v.Int(), nil
	case uint, uint8, uint16, uint32, uint64:
		return int64(v.Uint()), nil
	case float32, float64:
		return int64(v.Float()), nil
	case string:
		result, err := strconv.ParseInt(v.String(), 0, 64)
		if err != nil {
			return defaultValue, fmt.Errorf("%v(%T) -> Int Err: %s ", value, value, err)
		}
		return result, nil
	case bool:
		if v.Bool() == false {
			return 0, nil
		}
		return 1, nil
	default:
		return defaultValue, fmt.Errorf("%v(%T) -> Int Err", value, value)
	}
}

// I 如果出错，则返回0值
func I(value any) int64 {
	toInt, _ := ToInt(value, 0)
	return toInt
}

// ToString 转换为string类型
func ToString(value any, defaultValue string) (string, error) {
	if value == nil {
		return defaultValue, nil
	}
	switch val := value.(type) {
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", value), nil
	case []byte:
		return string(val), nil
	default:
		b, err := sonic.Marshal(val)
		if err != nil {
			return defaultValue, fmt.Errorf("%v(%T)->String Err:%s", value, value, err.Error())
		}
		return string(b), nil
	}
}

// S 安全转换到String
func S(value any) string {
	toString, _ := ToString(value, "")
	return toString
}

// ToJsonS 将值转为json字符串
func ToJsonS(value any) (string, error) {
	result, err := sonic.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func M(value any) map[string]any {
	toMap, _ := ToMap(value, map[string]any{})
	return toMap
}

func MS(value any) map[string]string {
	toMapS, _ := ToMapS(value, map[string]string{})
	return toMapS
}

// ToMap 转换成map
func ToMap(value any, defaultValue map[string]any) (map[string]any, error) {
	res := make(map[string]any)
	if reflect.ValueOf(value).Kind() == reflect.String {
		err := sonic.Unmarshal([]byte(value.(string)), &res)
		if err != nil {
			return defaultValue, err
		}
	}
	marshal, err := sonic.Marshal(value)
	if err != nil {
		return defaultValue, err
	}
	err = sonic.Unmarshal(marshal, &res)
	if err != nil {
		return defaultValue, err
	}
	return res, nil
}

// ToMapS 转换成map string string
func ToMapS(value any, defaultValue map[string]string) (map[string]string, error) {
	res := make(map[string]string)
	marshal, err := sonic.Marshal(value)
	if err != nil {
		return defaultValue, err
	}
	err = sonic.Unmarshal(marshal, &res)
	if err != nil {
		return defaultValue, err
	}
	return res, nil
}

func MapToStruct[T any](m map[string]T, v any) error {
	marshal, err := sonic.Marshal(m)
	if err != nil {
		return err
	}
	if err := sonic.Unmarshal(marshal, &v); err != nil {
		return err
	}
	return nil
}

func ToDuration(a any) time.Duration {
	d, _ := ToDurationE(a)
	return d
}

func ToDurationE(a any) (d time.Duration, err error) {
	if a == nil {
		return 0, nil
	}

	var i any
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		i = a
	} else {
		v := reflect.ValueOf(a)
		for v.Kind() == reflect.Ptr && !v.IsNil() {
			v = v.Elem()
		}
		i = v.Interface()
	}

	switch s := i.(type) {
	case time.Duration:
		return s, nil
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		d = time.Duration(I(s))
		return
	case float32, float64:
		//d = time.Duration(F(s))
		return
	case string:
		if strings.ContainsAny(s, "nsuµmh") {
			d, err = time.ParseDuration(s)
		} else {
			d, err = time.ParseDuration(s + "ns")
		}
		return
	case json.Number:
		var v float64
		v, err = s.Float64()
		d = time.Duration(v)
		return
	default:
		err = fmt.Errorf("unable to cast %#v of type %T to Duration", i, i)
		return
	}
}
