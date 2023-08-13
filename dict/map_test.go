// Package dict 包描述
// Author: wanlizhan
// Date: 2023/6/9
package dict

import (
	"fmt"
	"github.com/xxzhwl/wdk"
	"reflect"
	"testing"
)

func TestGetString(t *testing.T) {
	type args struct {
		m            map[string]any
		key          string
		defaultValue string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetString(tt.args.m, tt.args.key, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type Student struct {
	Name string
	Age  int
}

var m = map[string]any{
	"key1": "value1",
	"key2": 2,
	"key3": true,
	"key4": false,
	"key5": float64(3.14),
	"key6": float32(3.2),
	"key7": Student{
		Name: "hahhhh",
		Age:  24,
	},
}

func TestS(t *testing.T) {
	type args struct {
		m   map[string]any
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "GetS1",
			args: args{
				m:   m,
				key: "key1",
			},
			want: "value1",
		},
		{
			name: "GetS1",
			args: args{
				m:   m,
				key: "key2",
			},
			want: "2",
		},
		{
			name: "GetS1",
			args: args{
				m:   m,
				key: "key3",
			},
			want: "true",
		},
		{
			name: "GetS1",
			args: args{
				m:   m,
				key: "key4",
			},
			want: "false",
		},
		{
			name: "GetS1",
			args: args{
				m:   m,
				key: "key5",
			},
			want: "3.14",
		},
		{
			name: "GetS1",
			args: args{
				m:   m,
				key: "key6",
			},
			want: "3.2",
		},
		{
			name: "GetS1",
			args: args{
				m:   m,
				key: "key7",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(S(tt.args.m, tt.args.key), tt.want)
		})
	}
}

func TestGetInt64(t *testing.T) {
	type args struct {
		m            map[string]any
		key          string
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
			got, err := GetInt64(tt.args.m, tt.args.key, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetInt64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetString1(t *testing.T) {
	type args struct {
		m            map[string]any
		key          string
		defaultValue string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				m:            map[string]any{"mysql": map[string]any{"default": map[string]any{"schema": "nihao"}}},
				key:          "mysql.default",
				defaultValue: "",
			},
			want:    "{\"schema\":\"nihao\"}",
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				m:            map[string]any{"mysql": map[string]any{"default": "xxxxx"}},
				key:          "mysql.default",
				defaultValue: "",
			},
			want:    "xxxxx",
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				m:            map[string]any{"mysql.default": "xxxxx"},
				key:          "mysql.default",
				defaultValue: "",
			},
			want:    "xxxxx",
			wantErr: false,
		},
		{
			name: "test4",
			args: args{
				m:            map[string]any{"mysql": map[string]any{"default": map[string]any{"x1": "xxxxxxxxx"}}},
				key:          "mysql",
				defaultValue: "",
			},
			want:    "{\"default\":{\"x1\":\"xxxxxxxxx\"}}",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetString(tt.args.m, tt.args.key, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
			if got != tt.want {
				t.Errorf("GetString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStringList(t *testing.T) {
	type args[K comparable] struct {
		m            map[K]any
		key          K
		defaultValue []string
	}
	type testCase[K comparable] struct {
		name    string
		args    args[K]
		wantRes []string
		wantErr bool
	}
	var tests = []testCase[string /* TODO: Insert concrete types here */]{
		{
			name: "test1",
			args: args[string]{
				m:            map[string]any{"key1": []string{"12", "34", "56"}},
				key:          "key1",
				defaultValue: nil,
			},
			wantRes: []string{"12", "34", "56"},
			wantErr: false,
		},
		{
			name: "test2",
			args: args[string]{
				m:            map[string]any{"key1": []any{true, 1, "56", 200, 3.14}},
				key:          "key1",
				defaultValue: nil,
			},
			wantRes: []string{"true", "1", "56", "200", "3.14"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := GetStringList(tt.args.m, tt.args.key, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStringList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(gotRes, tt.wantRes)
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("GetStringList() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestI(t *testing.T) {
	type args struct {
		m   map[string]any
		key string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := I(tt.args.m, tt.args.key); got != tt.want {
				t.Errorf("I() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeys(t *testing.T) {
	type args[T comparable] struct {
		m map[T]any
	}
	type testCase[T comparable] struct {
		name    string
		args    args[T]
		wantRes []T
	}
	tests := []testCase[string /* TODO: Insert concrete types here */]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := Keys(tt.args.m); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Keys() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	type args[K comparable, V any] struct {
		maps []map[K]V
	}
	type testCase[K comparable, V any] struct {
		name    string
		args    args[K, V]
		wantRes map[K]V
	}
	tests := []testCase[string, any /* TODO: Insert concrete types here */]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := Merge(tt.args.maps...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Merge() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestS1(t *testing.T) {
	type args struct {
		m   map[string]any
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := S(tt.args.m, tt.args.key); got != tt.want {
				t.Errorf("S() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSet(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name    string
		args    args[K, V]
		wantRes []V
	}
	tests := []testCase[string, any /* TODO: Insert concrete types here */]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := ToSet(tt.args.m); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("ToSet() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestValues(t *testing.T) {
	type args[T comparable, V any] struct {
		m map[T]V
	}
	type testCase[T comparable, V any] struct {
		name    string
		args    args[T, V]
		wantRes []V
	}
	tests := []testCase[string, any /* TODO: Insert concrete types here */]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := Values(tt.args.m); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Values() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestM(t *testing.T) {
	type args struct {
		m   map[string]any
		key string
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				m:   map[string]any{"App": map[string]any{"Name": "Withu", "Host": "192.0.0.1"}},
				key: "App.Host",
			},
			want: map[string]any{"Name": "Withu", "Host": "192.0.0.1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := M(tt.args.m, tt.args.key)
			wdk.PrettyPrintln(got)
		})
	}
}
