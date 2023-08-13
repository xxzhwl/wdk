// Package grpool 包描述
// Author: wanlizhan
// Date: 2023/6/10
package grpool

import (
	"reflect"
	"testing"
)

type Student struct {
	Age int64
}

func TestDo(t *testing.T) {
	type args[I any, O any] struct {
		workerCount     int
		tasks           []TaskRequest[I]
		taskHandler     TaskHandler[I, O]
		responseHandler ResponseHandler[O]
	}
	type testCase[I any, O any] struct {
		name    string
		args    args[I, O]
		wantRes []TaskResponse[O]
		wantErr bool
	}
	tests := []testCase[Student, Student /* TODO: Insert concrete types here */]{
		// TODO: Add test cases.
		{
			name: "Test1",
			args: args[Student, Student]{
				workerCount:     5,
				tasks:           nil,
				taskHandler:     nil,
				responseHandler: nil,
			},
			wantRes: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := Do(tt.args.workerCount, tt.args.tasks, tt.args.taskHandler, tt.args.responseHandler)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Do() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
