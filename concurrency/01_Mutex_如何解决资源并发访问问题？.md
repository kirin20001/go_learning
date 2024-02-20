# 如何解决资源并发访问问题

## 并发问题
- 多个goroutine并发更新同一个资源（计数器，更新用户账户信息，秒杀系统，向同一个buffer中并发写入数据）
- 如果没有互斥系统，就会出现异常状况
- 解决方案：使用互斥锁（Mutex）
## 互斥锁的实现机制
- 临界区：并发编程中如果程序中一部分挥别并发访问或者需改，为了避免并发访问导致异常结果，这部分程序需要被保护起来，这部分被保护起来的程序叫做临界区。
- 多个线程同步访问临界区，就会造成访问或者操作错误，导致返回失败或者等待
- 解决方案：使用互斥锁，限定临界区只能同时由一个线程持有
- Go的标准库中提供了Mutex来实现互斥锁
- Mutex是使用最广泛的的同步原语（也可称并发原语）
- 同步原语（主要用于解决并发问题的一个基础数据结构）
  - 共享资源：并发的读写共享资源，会出现先数据竞争（data race）的问题，需要Mutex、RWMutex这样的并发原语来保护。
  - 任务编排：需要goroutine按照一定的规律执行，而goroutine之前有相互等待或者依赖的顺序关系，我们常常使用WaitGroup或者Channel来实现。
  - 消息传递：信息交流以及不同的groutine之间的线程安全的数据交流，常常使用Channel来实现。

## Mutex的基本使用方法
- Mutex实现的Locker接口
  ```go
  type Locker interface { 
    Lock() 
    Unlock()
  }
  ```
- 互斥锁 Mutex 就提供两个方法 Lock 和 Unlock
- 进入临界区之前调用 Lock 方法，退出临界区的时候调用 Unlock 方法
  ```go
  func(m *Mutex)Lock() 
  func(m *Mutex)Unlock()
  ```
- 当一个goroutine通过Lock方法获得了这个锁的拥有权后，其他请求锁的goroutine就会阻塞在Lock方法的调用上，直到锁被释放并且自己获取到了这个锁的拥有权。
- 并发场景：
  ```go
   import (
          "fmt"
          "sync"
      )
  
      func main() {
          var count = 0
          // 使用WaitGroup等待10个goroutine完成
          var wg sync.WaitGroup
          wg.Add(10)
          for i := 0; i < 10; i++ {
              go func() {
                  defer wg.Done()
                  // 对变量count执行10次加1
                  for j := 0; j < 100000; j++ {
                      count++
                  }
              }()
          }
          // 等待10个goroutine完成
          wg.Wait()
          fmt.Println(count)
      }
  ```
  - count ++ 不是原子操作，就可能有并发的问题
  - 解决方案：Go提供了一个检测并发访问共享资源是否又问题的工具-race detector，它可以帮助我们自动发现程序有没有data race的问题。
  - Go race detector 基于Google 的C/C++sanitizers技术实现，编译器通过探测所有的内存访问，在代码运行的时候，race detector 就能监控到对共享变量的非同步访问，出现race的时候，就会打印出警告信息。
### Mutex 嵌入到其它 struct 中使用
```
type Counter struct {
    mu    sync.Mutex
    Count uint64
}
```
- 如果嵌入的 struct 有多个字段，我们一般会把 Mutex 放在要控制的字段上面，然后使用空格把字段分隔开来。
- 在初始化嵌入的 struct 时，也不必初始化这个 Mutex 字段，不会因为没有初始化出现空指针或者是无法获取到锁的情况。
- 在初始化嵌入的 struct 时，也不必初始化这个 Mutex 字段，不会因为没有初始化出现空指针或者是无法获取到锁的情况。(?)

### 思考题
- 如果 Mutex 已经被一个 goroutine 获取了锁，其它等待中的 goroutine 们只能一直等待。那么，等这个锁释放后，等待中的 goroutine 中哪一个会优先获取 Mutex 呢？
- 当Mutex处于正常模式时，等待的goroutine是以FIFO排队的，
- 若锁释放，此时没有新goroutine与队头goroutine竞争，则队头goroutine获得锁。若有新goroutine竞争，大概率新goroutine获得锁
- 当队头goroutine竞争锁失败1ms后，它会将Mutex调整为饥饿模式，锁的所有权会直接从解锁goroutine移交给队头goroutine，此时新来的goroutine直接放入队尾
- 当一个goroutine获取锁后如果发现自己满足下列任意一个条件，则将锁从饥饿模式切换回正常模式：
	- 它是队列中最后一个
	- 它等待的时间少于1ms