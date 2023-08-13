// Package wdk 包描述
// Author: wanlizhan
// Date: 2023/7/15
package wdk

import (
	"fmt"
	"github.com/xxzhwl/wdk/rfunc"
	"testing"
)

type Response struct {
	Data any
	Code int64
	Err  error
	Msg  string
}

type CUser struct {
}

type RegisterArg struct {
	Name string
	Age  int64
}

func (c *CUser) Register(arg RegisterArg) Response {
	return Response{
		Data: map[string]any{"Name": arg.Name},
		Code: 0,
		Err:  nil,
		Msg:  "",
	}
}

func init() {
	rfunc.Register("cuser", &CUser{})
}

func DoUser(key, method string, arg any) {
	fromStruct, err := CallMethodFromStruct(rfunc.GetAction(key), method, arg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(fromStruct[0])

}

func TestDoUser(t *testing.T) {
	DoUser("cuser", "Register", RegisterArg{
		Name: "wanlizhan",
		Age:  256,
	})
}

func (c *CUser) GetName() string {
	return "123"
}

func Add(a, b int64) int64 {
	return a + b
}

func TestReflectAdd(t *testing.T) {
	obj := GetTypeObj(Add)
	in := obj.NumIn()
	for i := 0; i < in; i++ {
		fmt.Println(obj.In(i))
	}
}

func TestReflect(t *testing.T) {
	valueObj := GetValueObj(&CUser{})

	mth := valueObj.MethodByName("Register")
	in := mth.Type().NumIn()
	for i := 0; i < in; i++ {
		fmt.Println(mth.Type().In(i))
	}
}

func TestCallMethodFromStruct(t *testing.T) {
	fromStruct, err := CallMethodFromStruct(&CUser{}, "Register", RegisterArg{
		Name: "wanliz",
		Age:  1234,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(fromStruct) == 0 {
		fmt.Println("Success")
	}
	fmt.Println(fromStruct[0])
}
