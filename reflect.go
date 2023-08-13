// Package wdk 包描述
// Author: wanlizhan
// Date: 2023/7/15
package wdk

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// GetTypeObj 反射获取一个type对象
func GetTypeObj(o any) reflect.Type {
	return reflect.TypeOf(o)
}

// GetValueObj 反射获取一个value对象
func GetValueObj(o any) reflect.Value {
	return reflect.ValueOf(o)
}

// CallMethodFromStruct 执行一个struct绑定的方法
func CallMethodFromStruct(o any, methodName string, arg any) ([]reflect.Value, error) {
	obj := GetValueObj(o)
	met := obj.MethodByName(methodName)
	if !met.IsValid() || met.IsNil() {
		return []reflect.Value{}, fmt.Errorf("方法无效")
	}
	in := make([]reflect.Value, 0)
	//校验参数
	if inTemp, err := checkParam(met, methodName, arg); err != nil {
		return nil, err
	} else {
		in = inTemp
	}

	return met.Call(in), nil
}

// CallMethodWithMulArg 复合参数
func CallMethodWithMulArg(o any, methodName string, args ...string) ([]reflect.Value, error) {
	obj := GetValueObj(o)
	met := obj.MethodByName(methodName)
	if !met.IsValid() || met.IsNil() {
		return []reflect.Value{}, fmt.Errorf("方法无效")
	}
	in := make([]reflect.Value, 0)
	//校验参数
	if inTemp, err := checkMulParam(met, methodName, args...); err != nil {
		return nil, err
	} else {
		in = inTemp
	}

	return met.Call(in), nil
}

func checkParam(met reflect.Value, methodName string, arg any) ([]reflect.Value, error) {
	var inNum = met.Type().NumIn()
	in := make([]reflect.Value, 0)
	if inNum == 0 {
		return in, nil
	}
	if inNum > 0 && arg == nil {
		content := fmt.Sprintf("%s方法需要参数%d个,分别是", methodName, inNum)
		for i := 0; i < inNum; i++ {
			content += met.Type().In(i).Name()
		}
		return nil, errors.New(content)
	}
	//校验具体参数信息
	metArgName := met.Type().In(0).Name()
	obj := GetValueObj(arg)
	argName := obj.Type().Name()
	if !strings.EqualFold(metArgName, argName) {
		return nil, errors.New(fmt.Sprintf("传入参数类型不符合%s类型", metArgName))
	}

	return append(in, obj), nil
}
func checkMulParam(met reflect.Value, methodName string, args ...string) ([]reflect.Value, error) {
	var inNum = met.Type().NumIn()
	in := make([]reflect.Value, 0)
	if inNum == 0 {
		return in, nil
	}
	if len(args) != inNum {
		content := fmt.Sprintf("%s方法需要参数%d个,分别是", methodName, inNum)
		for i := 0; i < inNum; i++ {
			content += met.Type().In(i).Name()
		}
		return nil, errors.New(content)
	}
	for _, arg := range args {
		in = append(in, GetValueObj(arg))
	}

	return in, nil
}
