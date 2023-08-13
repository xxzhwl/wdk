// Package queue 包描述
// Author: wanlizhan
// Date: 2023/6/10
package queue

import (
	"errors"
	"time"
)

// DelayQueue 延迟队列
type DelayQueue[T any] struct {
	bufSize  int64
	interval time.Duration
	delay    time.Duration
	queue    chan DelayQueueData[T]
}

// DelayQueueData 实现这个接口
type DelayQueueData[T any] struct {
	Ti   time.Time
	Data T
}

// NewDelayQueue 获取一个DelayQueue
func NewDelayQueue[T any](size int64, interval, delay time.Duration) (DelayQueue[T], error) {
	if interval <= 0 || delay <= 0 {
		return DelayQueue[T]{}, errors.New("invalid delay queue option")
	}
	if size <= 0 {
		size = 1e4
	}
	return DelayQueue[T]{
		bufSize:  size,
		interval: interval,
		delay:    delay,
		queue:    make(chan DelayQueueData[T], size),
	}, nil
}

func (d *DelayQueue[T]) Push(data DelayQueueData[T]) {
	for time.Since(data.Ti) <= d.delay {
		time.Sleep(d.interval)
	}
	d.queue <- data
}

func (d *DelayQueue[T]) Pop() (item DelayQueueData[T]) {
	return <-d.queue
}
