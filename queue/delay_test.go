// Package queue 包描述
// Author: wanlizhan
// Date: 2023/6/10
package queue

import (
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	q, err := NewDelayQueue[int](100, time.Second*5, time.Second*10)
	if err != nil {
		t.Fatal(err)
	}
	tk := time.NewTicker(time.Second * 2)
	t1 := time.Now()
	lq := make(chan int)
	select {
	case <-tk.C:
		q.Push(DelayQueueData[int]{Ti: time.Now(), Data: 1})
	}
	go func() {
		if time.Now().After(t1.Add(time.Second * 5)) {
			lq <- 1
		}
	}()
	select {
	case x := <-q.queue:
		fmt.Println(x.Ti.Format(time.DateTime), x.Data)
	case _ = <-lq:
		return
	}

}
