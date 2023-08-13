// Package defs 包描述
// Author: wanlizhan
// Date: 2023/6/10
package defs

// Integer 整型定义
type Integer interface {
	UInt | Int
}

// UInt 整型定义
type UInt interface {
	uint | uint8 | uint16 | uint32 | uint64
}

// Int 整型定义
type Int interface {
	int | int8 | int16 | int32 | int64
}
