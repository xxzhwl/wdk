// Package ustr 包描述
// Author: wanlizhan
// Date: 2023/7/2
package ustr

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math/rand"
	"strings"
	"time"
)

// SnakeToBigCamel 蛇形转大驼峰
func SnakeToBigCamel(str string) string {
	split := strings.Split(str, "_")
	res := ""

	caser := cases.Title(language.English)
	for i := range split {
		tempS := caser.String(split[i])
		res += tempS
	}
	return res
}

// SnakeToSmallCamel 蛇形转小驼峰
func SnakeToSmallCamel(str string) string {
	split := strings.Split(str, "_")
	res := ""
	caser := cases.Title(language.English)
	for i := range split {
		if i == 0 {
			res += split[i]
			continue
		}
		tempS := caser.String(split[i])
		res += tempS
	}
	return res
}

// CodeWithABC 带数字和英文验证码
func CodeWithABC(num int) (res string) {
	codeList := "0123456789qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		res += string(codeList[r.Intn(len(codeList))])
	}
	return
}

// NumCode 纯数字验证码
func NumCode(num int) (res string) {
	codeList := "0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		res += string(codeList[r.Intn(len(codeList))])
	}
	return
}

// ABCCode 纯英文验证码
func ABCCode(num int) (res string) {
	codeList := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		res += string(codeList[r.Intn(len(codeList))])
	}
	return
}

// Contains 不区分大小写的比较s1是否包含s2
func Contains(s1, s2 string) bool {
	return strings.Contains(strings.ToLower(s1), strings.ToLower(s2))
}

// CommonSplit 使用;或,分割字符串列表
func CommonSplit(s string) []string {
	res := make([]string, 0)
	sTemp := strings.Split(s, ";")
	for _, s2 := range sTemp {
		resTemp := strings.Split(s2, ",")
		res = append(res, resTemp...)
	}
	return res
}
