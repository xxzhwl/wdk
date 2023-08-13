// Package system 包描述
// Author: wanlizhan
// Date: 2023/6/10
package system

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// GetGoRoutineId 获取goroutineId
func GetGoRoutineId() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	fields := strings.Fields(string(buf[:n]))
	return fields[1]
}

// GetStackFrames 获取当前函数的调用栈 返回原始[]runtime.Frame
// skip为2 跳过本函数与callers
// maxSize 为最大返回的调用栈行数
func GetStackFrames(skip int, maxSize int) []runtime.Frame {
	skip += 2
	if maxSize == 0 {
		maxSize = 16
	}

	pc := make([]uintptr, maxSize)
	callers := runtime.Callers(skip, pc)
	if callers == 0 {
		return []runtime.Frame{}
	}

	pc = pc[:callers]
	frames := runtime.CallersFrames(pc)

	ret := make([]runtime.Frame, 0)
	for {
		frame, more := frames.Next()
		ret = append(ret, frame)
		if !more {
			break
		}
	}
	return ret
}

// GetStackFramesString 把栈信息整合为string
func GetStackFramesString(skip, maxSize int) (ret string) {
	frames := GetStackFrames(skip, maxSize)
	if len(frames) <= 0 {
		return ""
	}
	for _, frame := range frames {
		ret += frame.Function + "\n"
		ret += "\t" + frame.File + "\t" + strconv.Itoa(frame.Line) + "\t\n"
	}
	return
}

// GetCurrentFilePath 获取当前文件路径
func GetCurrentFilePath() string {
	frames := GetStackFrames(2, 0)
	if len(frames) == 0 {
		return ""
	}
	return frames[0].File
}

// GetCurrentFolderPath 获取当前文件路径
func GetCurrentFolderPath() string {
	path := GetCurrentFilePath()
	if len(path) <= 1 {
		return ""
	}
	return filepath.Dir(path)
}

// GetCurrentFileNameWithSuffix 获取当前文件名带后缀
func GetCurrentFileNameWithSuffix() string {
	path := GetCurrentFilePath()
	if len(path) <= 1 {
		return ""
	}
	return path[strings.LastIndex(path, "/")+1:]
}

// GetCurrentFileName 获取当前文件名且不带
func GetCurrentFileName() string {
	fileName := GetCurrentFileNameWithSuffix()
	return fileName[:strings.LastIndex(fileName, ".")]
}
