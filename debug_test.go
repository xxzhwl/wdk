// Package wdk 包描述
// Author: wanlizhan
// Date: 2023/7/1
package wdk

import (
	"fmt"
	"github.com/xxzhwl/wdk/system"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

var lo sync.Mutex

var countMap map[string]int

func init() {
	countMap = make(map[string]int)
}
func TestGoId(t *testing.T) {
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		rand.Seed(time.Now().Unix())
		n := rand.Intn(10)
		duration := time.Duration(n)
		go func(d time.Duration) {
			defer wg.Done()
			Pln(d)
		}(duration)
	}
	wg.Wait()
	var have = false
	for _, i := range countMap {
		if i > 1 {
			have = true
		}
	}
	if have {
		fmt.Println("有协程被复用了")
	}
}

func Pln(d time.Duration) {
	goId := system.GetGoRoutineId()
	lo.Lock()
	if _, ok := countMap[goId]; ok {
		countMap[goId]++
	} else {
		countMap[goId] = 1
	}
	lo.Unlock()
	time.Sleep(d * time.Second)
}

func TestPrettyPrintln(t *testing.T) {
	type Person struct {
		Name    string
		Hobbies []string
	}

	p1 := Person{
		Name:    "xxxxx",
		Hobbies: []string{"game", "bool"},
	}

	PrettyPrintln(p1)
}
