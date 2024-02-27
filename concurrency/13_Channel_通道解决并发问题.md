# Channel

## Channel的发展

- Channel源自CSP模型

## Channel的应用场景

- 数据交流：当作并发的 buffer 或者 queue，解决生产者 - 消费者问题。多个 goroutine 可以并发当作生产者（Producer）和消费者（Consumer）。
- 数据传递：一个 goroutine 将数据交给另一个 goroutine，相当于把数据的拥有权 (引用) 托付出去。
- 信号通知：一个 goroutine 可以将信号 (closing、closed、data ready 等) 传递给另一个或者另一组 goroutine 。
- 任务编排：可以让一组 goroutine 按照一定的顺序并发或者串行的执行，这就是编排的功能。
- 锁：利用 Channel 也可以实现互斥锁的机制。

## Channel基本用法

- 发送数据 ch <- 2000

- 接收数据

  ```go
    x := <-ch // 把接收的一条数据赋值给变量x
    foo(<-ch) // 把接收的一个的数据作为参数传给函数
    <-ch // 丢弃接收的一条数据
  ```

- 其他操作

  - close 会把 chan 关闭掉，

  - cap 返回 chan 的容量，

  - len 返回 chan 中缓存的还未被取走的元素数量。

  - send 和 recv 都可以作为 select 语句的 case clause

    ```go
    func main() {
        var ch = make(chan int, 10)
        for i := 0; i < 10; i++ {
            select {
            case ch <- i:
            case v := <-ch:
                fmt.Println(v)
            }
        }
    }
    ```

  - chan 还可以应用于 for-range 语句中

## Channel 的实现原理

- qcount：代表 chan 中已经接收但还没被取走的元素的个数。内建函数 len 可以返回这个字段的值。
- dataqsiz：队列的大小。chan 使用一个循环队列来存放元素，循环队列很适合这种生产者 - 消费者的场景（我很好奇为什么这个字段省略 size 中的 e）。
- buf：存放元素的循环队列的 buffer。
- elemtype 和 elemsize：chan 中元素的类型和 size。因为 chan 一旦声明，它的元素类型是固定的，即普通类型或者指针类型，所以元素大小也是固定的。
- sendx：处理发送数据的指针在 buf 中的位置。一旦接收了新的数据，指针就会加上 elemsize，移向下一个位置。buf 的总大小是 elemsize 的整数倍，而且 buf 是一个循环列表。
- recvx：处理接收请求时的指针在 buf 中的位置。一旦取出数据，此指针会移动到下一个位置。
- recvq：chan 是多生产者多消费者的模式，如果消费者因为没有数据可读而被阻塞了，就会被加入到 recvq 队列中。s
- endq：如果生产者因为 buf 满了而阻塞，会被加入到 sendq 队列中。

## Channel常见错误

- 使用 Channel 最常见的错误是 panic 和 goroutine 泄漏。
- 首先，我们来总结下会 panic 的情况，总共有 3 种：
  - close 为 nil 的 chan；
  - send 已经 close 的 chan；
  - close 已经 close 的 chan。

## 总结

![](/Users/wyl/Desktop/go/src/go_learning/concurrency/img/13_Channel_1.webp)