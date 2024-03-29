package concurrency

import (
	"sync"
	"time"
)

func RWMutexCounter() {
	var counter RWCounter
	for i := 0; i < 10; i++ { // 10个reader
		go func() {
			for {
				counter.Count() // 计数器读操作
				time.Sleep(time.Millisecond)
			}
		}()
	}

	for { // 一个writer
		counter.Incr() // 计数器写操作
		time.Sleep(time.Second)
	}
}
// 一个线程安全的计数器
type RWCounter struct {
	mu    sync.RWMutex
	count uint64
}

// 使用写锁保护
func (c *RWCounter) Incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

// 使用读锁保护
func (c *RWCounter) Count() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.count
}
