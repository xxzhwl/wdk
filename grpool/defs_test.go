// Package grpool 包描述
// Author: wanlizhan
// Date: 2023/6/10
package grpool

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewTaskRequests(t *testing.T) {
	type args[I any] struct {
		data []I
	}
	type testCase[I any] struct {
		name string
		args args[I]
		want []TaskRequest[I]
	}
	tests := []testCase[Student /* TODO: Insert concrete types here */]{
		// TODO: Add test cases.
		{
			name: "Test1",
			args: args[Student]{
				data: []Student{
					{Age: 25},
					{Age: 15},
					{Age: 35},
					{Age: 45},
					{Age: 55},
					{Age: 65},
					{Age: 75},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTaskRequests(tt.args.data)
			fmt.Println(got)
		})
	}
}

func NewStudentTasks() []TaskRequest[Student] {
	source := rand.NewSource(time.Now().Unix())
	r := rand.New(source)
	var data []Student
	for i := 0; i < 10000; i++ {
		data = append(data, Student{Age: r.Int63n(100)})
	}
	res := NewTaskRequests(data)
	return res
}

func TestName(t *testing.T) {
	tasks := NewStudentTasks()
	fmt.Println(tasks)
	res, err := Do[Student, Student](10, tasks, func(task TaskRequest[Student]) (outPut Student, err error) {
		return Student{task.inputData.Age - 100}, nil
	}, nil)
	if err != nil {
		return
	}
	fmt.Println(res)
}
