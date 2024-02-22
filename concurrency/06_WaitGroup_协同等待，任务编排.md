## WaitGroup 基本用法

- WaitGroup提供三个方法

  ```go
  func (wg *WaitGroup) Add(delta int)  
  func (wg *WaitGroup) Done()
  func (wg *WaitGroup) Wait()
  ```

  - Add，用来设置 WaitGroup 的计数值；
  - Done，用来将 WaitGroup 的计数值减 1，其实就是调用了 Add(-1)；
  - Wait，调用这个方法的 goroutine 会一直阻塞，直到 WaitGroup 的计数值变为 0。

## WaitGroup的实现
- noCopy 的辅助字段，主要就是辅助 vet 工具检查是否通过 copy 赋值这个 WaitGroup 实例。
- state1，一个具有复合意义的字段，包含 WaitGroup 的计数、阻塞在检查点的 waiter 数和信号量。
```go
type WaitGroup struct {
    // 避免复制使用的一个技巧，可以告诉vet工具违反了复制使用的规则
    noCopy noCopy
    // 64bit(8bytes)的值分成两段，高32bit是计数值，低32bit是waiter的计数
    // 另外32bit是用作信号量的
    // 因为64bit值的原子操作需要64bit对齐，但是32bit编译器不支持，所以数组中的元素在不同的架构中不一样，具体处理看下面的方法
    // 总之，会找到对齐的那64bit作为state，其余的32bit做信号量
    state1 [3]uint32
}


// 得到state的地址和信号量的地址
func (wg *WaitGroup) state() (statep *uint64, semap *uint32) {
    if uintptr(unsafe.Pointer(&wg.state1))%8 == 0 {
        // 如果地址是64bit对齐的，数组前两个元素做state，后一个元素做信号量
        return (*uint64)(unsafe.Pointer(&wg.state1)), &wg.state1[2]
    } else {
        // 如果地址是32bit对齐的，数组后两个元素用来做state，它可以用来做64bit的原子操作，第一个元素32bit用来做信号量
        return (*uint64)(unsafe.Pointer(&wg.state1[1])), &wg.state1[0]
    }
}
```

## Add & Done

- Add 方法主要操作的是 state 的计数部分。你可以为计数值增加一个 delta 值，内部通过原子操作把这个值加到计数值上。需要注意的是，这个 delta 也可以是个负数，相当于为计数值减去一个值，Done 方法内部其实就是通过 Add(-1) 实现的。

  ```go
  func (wg *WaitGroup) Add(delta int) {
      statep, semap := wg.state()
      // 高32bit是计数值v，所以把delta左移32，增加到计数上
      state := atomic.AddUint64(statep, uint64(delta)<<32)
      v := int32(state >> 32) // 当前计数值
      w := uint32(state) // waiter count
  
      if v > 0 || w == 0 {
          return
      }
  
      // 如果计数值v为0并且waiter的数量w不为0，那么state的值就是waiter的数量
      // 将waiter的数量设置为0，因为计数值v也是0,所以它们俩的组合*statep直接设置为0即可。此时需要并唤醒所有的waiter
      *statep = 0
      for ; w != 0; w-- {
          runtime_Semrelease(semap, false, 0)
      }
  }
  
  
  // Done方法实际就是计数器减1
  func (wg *WaitGroup) Done() {
      wg.Add(-1)
  }
  ```

  ## Wait

  - 不断检查 state 的值。如果其中的计数值变为了 0，那么说明所有的任务已完成，调用者不必再等待，直接返回。如果计数值大于 0，说明此时还有任务没完成，那么调用者就变成了等待者，需要加入 waiter 队列，并且阻塞住自己。

    ```go
    func (wg *WaitGroup) Wait() {
        statep, semap := wg.state()
        
        for {
            state := atomic.LoadUint64(statep)
            v := int32(state >> 32) // 当前计数值
            w := uint32(state) // waiter的数量
            if v == 0 {
                // 如果计数值为0, 调用这个方法的goroutine不必再等待，继续执行它后面的逻辑即可
                return
            }
            // 否则把waiter数量加1。期间可能有并发调用Wait的情况，所以最外层使用了一个for循环
            if atomic.CompareAndSwapUint64(statep, state, state+1) {
                // 阻塞休眠等待
                runtime_Semacquire(semap)
                // 被唤醒，不再阻塞，返回
                return
            }
        }
    }
    ```

    

## WaitGroup常见错误

- 计数器设置为负值
  - 调用 Add 的时候传递一个负数。如果你能保证当前的计数器加上这个负数后还是大于等于 0 的话，也没有问题，否则就会导致 panic。
  - 调用 Done 方法的次数过多，超过了 WaitGroup 的计数值。
  - 使用 WaitGroup 的正确姿势是，预先确定好 WaitGroup 的计数值，然后调用相同次数的 Done 完成相应的任务。比如，在 WaitGroup 变量声明之后，就立即设置它的计数值，或者在 goroutine 启动之前增加 1，然后在 goroutine 中调用 Done。如果你没有遵循这些规则，就很可能会导致 Done 方法调用的次数和计数值不一致，进而造成死锁（Done 调用次数比计数值少）或者 panic（Done 调用次数比计数值多）。
- 不期望的 Add 时机
  - 在使用 WaitGroup 的时候，你一定要遵循的原则就是，等所有的 Add 方法调用之后再调用 Wait，否则就可能导致 panic 或者不期望的结果。
- 前一个 Wait 还没结束就重用 WaitGroup
  - 因为 WaitGroup 是可以重用的。只要 WaitGroup 的计数值恢复到零值的状态，那么它就可以被看作是新创建的 WaitGroup，被重复使用。
  - 但是，如果我们在 WaitGroup 的计数值还没有恢复到零值的时候就重用，就会导致程序 panic。
  - 我们看一个例子，初始设置 WaitGroup 的计数值为 1，启动一个 goroutine 先调用 Done 方法，接着就调用 Add 方法，Add 方法有可能和主 goroutine 并发执行。

## noCopy：辅助 vet 检查

- 指示 vet 工具在做检查的时候，这个数据结构不能做值复制使用。

## 思考题

- 通常我们可以把 WaitGroup 的计数值，理解为等待要完成的 waiter 的数量。你可以试着扩展下 WaitGroup，来查询 WaitGroup 的当前的计数值吗？
