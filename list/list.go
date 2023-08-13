// Package list 包描述
// Author: wanlizhan
// Date: 2023/6/9
package list

// UniqueList list去重
func UniqueList[T comparable](s []T) (res []T) {
	seen := make(map[T]int64, len(s))
	for _, itm := range s {
		if _, ok := seen[itm]; !ok {
			seen[itm] = 1
			res = append(res, itm)
		}
	}
	return
}

// InList 判断元素needle是否在list中
func InList[T comparable](needle T, list []T) bool {
	for _, item := range list {
		if item == needle {
			return true
		}
	}
	return false
}

// Intersection 取交集
func Intersection[T comparable](l1, l2 []T) []T {
	res := make([]T, 0)
	for _, t := range l1 {
		if InList(t, l2) {
			res = append(res, t)
		}
	}
	return res
}

// UnionList 取并集
func UnionList[T comparable](l1, l2 []T) []T {
	for _, t := range l1 {
		if !InList(t, l2) {
			l2 = append(l2, t)
			continue
		}
	}
	return l2
}

// DifferenceList 返回L1和L2的差集，也即只在l1中存在,l2中不存在的
func DifferenceList[T comparable](l1, l2 []T) []T {
	res := make([]T, 0)
	for _, t := range l1 {
		if !InList(t, l2) {
			res = append(res, t)
		}
	}
	return res
}
