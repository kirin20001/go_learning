package concurrency

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// 线程安全的计数器
type WGCounter struct {
	mu    sync.Mutex
	count uint64
}
// 对计数值加一
func (c *WGCounter) Incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}
// 获取当前的计数值
func (c *WGCounter) Count() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}
// sleep 1秒，然后计数值加1
func worker(c *WGCounter, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second)
	c.Incr()
}



func wgCounter() {
	var counter WGCounter

	var wg sync.WaitGroup
	wg.Add(10) // WaitGroup的值设置为10

	for i := 0; i < 10; i++ { // 启动10个goroutine执行加1任务
		go worker(&counter, &wg)
	}

	// 检查点，等待goroutine都完成任务
	wg.Wait()
	// 输出当前计数器的值
	fmt.Println(counter.Count())
}

func getWaitGroupCounter(wg *sync.WaitGroup) uint64 {
	v := atomic.LoadUint64((*uint64)(unsafe.Pointer(wg)))
	v = v >> 32
	return v
}

func getWGCounter() {
	var counter WGCounter
	var wg sync.WaitGroup
	wg.Add(10) // WaitGroup的值设置为10
	fmt.Println("wg number ", getWaitGroupCounter(&wg))
	for i := 0; i < 10; i++ { // 启动10个goroutine执行加1任务
		go worker(&counter, &wg)
	}
	fmt.Println("wg number ", getWaitGroupCounter(&wg))
	// 检查点，等待goroutine都完成任务
	wg.Wait()
	fmt.Println("wg number ", getWaitGroupCounter(&wg))
	// 输出当前计数器的值
	fmt.Println(counter.Count())
}
