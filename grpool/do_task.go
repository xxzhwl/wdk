// Package grpool 包描述
// Author: wanlizhan
// Date: 2023/6/10
package grpool

import (
	"errors"
	"fmt"
	"github.com/xxzhwl/wdk"
	"github.com/xxzhwl/wdk/system"
	"sync"
)

// TaskHandler 对task做什么操作，需要定义
type TaskHandler[I, O any] func(task TaskRequest[I]) (outPut O, err error)

// ResponseHandler 对返回结果要做什么操作
type ResponseHandler[O any] func(response TaskResponse[O])

// Do 执行任务
func Do[I, O any](workerCount int, tasks []TaskRequest[I], taskHandler TaskHandler[I, O],
	responseHandler ResponseHandler[O]) (res []TaskResponse[O], err error) {
	var locker sync.RWMutex
	var taskNum = len(tasks)
	if taskNum == 0 {
		return nil, errors.New("协程池传入任务列表为空")
	}
	if taskHandler == nil {
		return nil, errors.New("任务处理器为nil")
	}
	var reqChan = make(chan TaskRequest[I], taskNum)
	var respChan = make(chan TaskResponse[O], taskNum)

	for _, task := range tasks {
		reqChan <- task
	}

	for i := 0; i < workerCount; i++ {
		go func(parentRoutineId string) {
			defer wdk.CatchPanic()

			for {
				locker.Lock()
				if len(reqChan) == 0 {
					locker.Unlock()
					break
				}
				taskData := <-reqChan
				locker.Unlock()

				func(task TaskRequest[I]) {
					defer func() {
						if r := recover(); r != nil {
							respChan <- TaskResponse[O]{Id: task.id, Err: fmt.Errorf("%v", r)}
						}
					}()
					outPut, err := taskHandler(task)
					respChan <- TaskResponse[O]{Id: task.id, OutputData: outPut, Err: err}
				}(taskData)
			}
		}(system.GetGoRoutineId())
	}

	if responseHandler != nil {
		for i := 0; i < taskNum; i++ {
			responseHandler(<-respChan)
		}
	}
	for i := 0; i < taskNum; i++ {
		res = append(res, <-respChan)
	}
	return
}
