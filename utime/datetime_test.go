// Package utime 包描述
// Author: wanlizhan
// Date: 2023/6/10
package utime

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	str1 := "2023-05-10 22:15:30"
	t2, err := StrToTime(str1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(t2.Unix())
	fmt.Println(TimeStamp(), Date(), DateTime())

}
