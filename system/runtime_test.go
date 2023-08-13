// Package inner 包描述
// Author: wanlizhan
// Date: 2023/6/10
package system

import (
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	ret := GetStackFramesString(0, 0)
	fmt.Println(ret)
	path := GetCurrentFilePath()
	name := GetCurrentFileName()
	suffix := GetCurrentFileNameWithSuffix()
	fmt.Println(path, name, suffix)
	fmt.Println("outer:" + GetGoRoutineId())

	go func() {
		fmt.Println("inner:" + GetGoRoutineId())
	}()
	time.Sleep(5 * time.Second)
}
