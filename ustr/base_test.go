// Package ustr 包描述
// Author: wanlizhan
// Date: 2023/7/2
package ustr

import (
	"testing"
)

func TestSnakeToBigCamel(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"test1", args{str: "sys_api_down"}, "SysApiDown"},
		{"test2", args{str: "user"}, "User"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnakeToBigCamel(tt.args.str); got != tt.want {
				t.Errorf("SnakeToBigCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnakeToSmallCamel(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test3", args{str: "sys_api_down"}, "sysApiDown"},
		{"test4", args{str: "user"}, "user"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnakeToSmallCamel(tt.args.str); got != tt.want {
				t.Errorf("SnakeToSmallCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}
