// Package wdk 包描述
// Author: wanlizhan
// Date: 2023/6/10
package wdk

import (
	"testing"
)

func TestName(t *testing.T) {
	t.Skip()
	func() {
		defer CatchPanic()
		panic("测试panic")
	}()

}
