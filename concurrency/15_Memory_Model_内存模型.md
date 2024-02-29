# Go内存模型

## 重排和可见性

### 由于指令重排，代码并不一定会按照你写的顺序执行

- 程序在运行的时候，两个操作的顺序可能不会得到保证

### happens-before

- **在一个 goroutine 内部，程序的执行顺序和它们的代码指定的顺序是一样的**，即使编译器或者 CPU 重排了读写顺序，从行为上来看，也和代码指定的顺序一样。
- 但是，对于另一个 goroutine 来说，重排却会产生非常大的影响。因为 Go 只保证 goroutine 内部重排对读写的顺序没有影响
- **如果两个 action（read 或者 write）有明确的 happens-before 关系，你就可以确定它们之间的执行顺序（或者是行为表现上的顺序）。**
- 如果要保证对“变量 v 的读操作 r”能够观察到一个对“变量 v 的写操作 w”，并且 r 只能观察到 w 对变量 v 的写，没有其它对 v 的写操作，也就是说，我们要保证 r 绝对能观察到 w 操作的结果，那么就需要同时满足两个条件：
  - w happens before r；
  - 其它对 v 的写操作（w2、w3、w4, ......） 要么 happens before w，要么 happens after r，绝对不会和 w、r 同时发生，或者是在它们之间发生。
- 另外
  - 在 Go 语言中，对变量进行零值的初始化就是一个写操作。
  - 如果对超过机器 word（64bit、32bit 或者其它）大小的值进行读写，那么，就可以看作是对拆成 word 大小的几个读写无序进行。
  - Go 并不提供直接的 CPU 屏障（CPU fence）来提示编译器或者 CPU 保证顺序性，而是使用不同架构的内存屏障指令来实现统一的并发原语。

### Go语言中保证happens-before关系

- ### Init 函数

  - 应用程序的初始化是在单一的 goroutine 执行的。如果包 p 导入了包 q，那么，q 的 init 函数的执行一定 happens before p 的任何初始化代码。
  - 这里有一个特殊情况需要你记住：main 函数一定在导入的包的 init 函数之后执行。

- ### goroutine

  - 启动 goroutine 的 go 语句的执行，一定 happens before 此 goroutine 内的代码执行。

- Channel

  - 第 1 条规则是，往 Channel 中的发送操作，happens before 从该 Channel 接收相应数据的动作完成之前，即第 n 个 send 一定 happens before 第 n 个 receive 的完成。
  - 第 2 条规则是，close 一个 Channel 的调用，肯定 happens before 从关闭的 Channel 中读取出一个零值。
  - 第 3 条规则是，对于 unbuffered 的 Channel，也就是容量是 0 的 Channel，从此 Channel 中读取数据的调用一定 happens before 往此 Channel 发送数据的调用完成。
  - 第 4 条规则是，如果 Channel 的容量是 m（m>0），那么，第 n 个 receive 一定 happens before 第 n+m 个 send 的完成。

- Mutex/RWMutex

  - 第 n 次的 m.Unlock 一定 happens before 第 n+1 m.Lock 方法的返回；
  - 对于读写锁 RWMutex m，如果它的第 n 个 m.Lock 方法的调用已返回，那么它的第 n 个 m.Unlock 的方法调用一定 happens before 任何一个 m.RLock 方法调用的返回，只要这些 m.RLock 方法调用 happens after 第 n 次 m.Lock 的调用的返回。这就可以保证，只有释放了持有的写锁，那些等待的读请求才能请求到读锁。
  - 对于读写锁 RWMutex m，如果它的第 n 个 m.RLock 方法的调用已返回，那么它的第 k （k<=n）个成功的 m.RUnlock 方法的返回一定 happens before 任意的 m.RUnlockLock 方法调用，只要这些 m.Lock 方法调用 happens after 第 n 次 m.RLock。

- WaitGroup

  - Wait 方法等到计数值归零之后才返回。

- Once

  - 对于 once.Do(f) 调用，f 函数的那个单次调用一定 happens before 任何 once.Do(f) 调用的返回。

- atomic

  - 可以保证使用 atomic 的 Load/Store 的变量之间的顺序性。

