package concurrency

import (
	"fmt"
	"sync"
)

type Counter struct {
	sync.Mutex
	Count int
}

func mutexCounter() {
	var c Counter
	c.Lock()
	defer c.Unlock()
	c.Count++
	foo(c) // 复制锁
}

// 这里Counter的参数是通过复制的方式传入的
func foo(c Counter) {
	c.Lock()
	defer c.Unlock()
	fmt.Println("in foo")
}

func mutexMultiLock() {
	l := &sync.Mutex{}
	l.Lock()
	l.Lock()
	foo1(l)
}


func foo1(l sync.Locker) {
	fmt.Println("in foo")
	l.Lock()
	bar(l)
	l.Unlock()
}


func bar(l sync.Locker) {
	l.Lock()
	fmt.Println("in bar")
	l.Unlock()
}
