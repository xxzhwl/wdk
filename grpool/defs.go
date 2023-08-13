// Package grpool 包描述
// Author: wanlizhan
// Date: 2023/6/10
package grpool

import (
	"strconv"
)

// TaskRequest 请求任务
type TaskRequest[I any] struct {
	id        string
	inputData I
}

// TaskResponse 任务结果
type TaskResponse[O any] struct {
	Id         string
	OutputData O
	Err        error
}

// NewTaskRequests 构建一组请求任务
func NewTaskRequests[I any](data []I) []TaskRequest[I] {
	res := make([]TaskRequest[I], len(data))
	for i, itm := range data {
		res[i] = TaskRequest[I]{strconv.Itoa(i), itm}
	}

	return res
}
