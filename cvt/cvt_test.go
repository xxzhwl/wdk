// Package cvt 包描述
// Author: wanlizhan
// Date: 2023/6/9
package cvt

import (
	"fmt"
	"testing"
)

func TestToInt(t *testing.T) {
	type args struct {
		value        any
		defaultValue int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "BoolToInt", args: args{
			value:        true,
			defaultValue: 0,
		}},
		{
			name: "Float64ToInt",
			args: args{
				value:        3.14,
				defaultValue: 0,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "StringToInt",
			args: args{
				value:        "hhhh",
				defaultValue: 0,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt(tt.args.value, tt.args.defaultValue)
			if err != nil {
				t.Errorf("ToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}

func TestToBool(t *testing.T) {
	type args struct {
		value        any
		defaultValue bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "UIntToBool",
			args: args{
				value:        uint(0),
				defaultValue: false,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToBool(tt.args.value, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(uint(0) != 0)
			fmt.Println(got, tt.want)
		})
	}
}

func TestToInt1(t *testing.T) {
	type args struct {
		value        any
		defaultValue int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt(tt.args.value, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
