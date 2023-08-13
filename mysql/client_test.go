// Package mysql 包描述
// Author: wanlizhan
// Date: 2023/6/10
package mysql

import (
	"fmt"
	"github.com/duke-git/lancet/v2/maputil"
	"testing"
)

type User struct {
	NickName string `json:"nickname,omitempty"`
}

func TestName(t *testing.T) {
	clien, err := NewDynamicInstance("localhost", "3306", "root", "123456", "withu", "", false)
	if err != nil {
		t.Fatal(err)
	}
	query, err := clien.Query("select * from user")
	if err != nil {
		t.Fatal(err)
	}
	u1 := User{}
	maputil.MapTo(query[1], &u1)
	fmt.Println(u1)
}
