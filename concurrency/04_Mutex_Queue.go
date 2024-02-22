package concurrency

import (
	"fmt"
	"sync"
	"time"
)

type SliceQueue struct {
	data []interface{}
	mu   sync.Mutex
}

func NewSliceQueue(n int) (q *SliceQueue) {
	return &SliceQueue{data: make([]interface{}, 0, n)}
}

// Enqueue 把值放在队尾
func (q *SliceQueue) Enqueue(v interface{}) {
	q.mu.Lock()
	q.data = append(q.data, v)
	q.mu.Unlock()
}

// Dequeue 移去队头并返回
func (q *SliceQueue) Dequeue() interface{} {
	q.mu.Lock()
	if len(q.data) == 0 {
		q.mu.Unlock()
		return nil
	}
	v := q.data[0]
	q.data = q.data[1:]
	q.mu.Unlock()
	return v
}

func tryMutexQueue() {
	q := NewSliceQueue(50)
	var wg sync.WaitGroup
	wg.Add(500)
	for i := 0; i < 1000; i++ { // 启动1000个goroutine
		go func(i int) {
			q.Enqueue(i)
		}(i)
	}

	for i := 0; i < 500; i++ {
		go func(wg *sync.WaitGroup) {
			var cnt int
			for cnt < 50 {
				v := q.Dequeue()
				if v != nil {
					fmt.Println(v)
				} else {
					cnt++
				}
				time.Sleep(10 * time.Millisecond)
			}
			wg.Done()
		}(&wg)
	}

	wg.Wait()
}
