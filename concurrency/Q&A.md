1. WaitGroup

   64位原子操作支持操作32为系统内存？WaitGroup Add使用的都是64位原子操作？

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

2. 使用channel 如果停止子goroutine？

   使用stopChannel + select case？

   ```go
   func channel2_1() {
   	var chans [4]chan struct{}
   	stopCh := make(chan struct{})
   	for i, _ := range chans {
   		chans[i] = make(chan struct{})
   	}
   
   	for i := 0; i < 4; i++ {
   		go func(i int) {
   			for {
   				select {
   				case <- chans[i]:
   					fmt.Println(i, time.Now())
   					select {
   					case <- time.After(1 * time.Second):
   						chans[(i+1)%4] <- struct{}{}
   					case <- stopCh:
   						fmt.Println(i, "terminated")
   						return
   					}
   				case <- stopCh: // 如果在第一个case无法接收stopChannel
   					fmt.Println(i, "terminated")
   					return
   				}
   			}
   		}(i)
   	}
   
   	chans[0] <- struct{}{}
   	select {
   	case <-time.After(5 * time.Second):
   		close(stopCh)
   	}
   	time.Sleep(1 * time.Second)
   	for i, _ := range chans {
   		close(chans[i])
   	}
   }
   ```

3. Or Done模式的优势是什么，和2中的方法是主协程通知子携程&子携程通知主携程区别？

4. 扇入模式应用场景： 单一reciver场景？

5. 扇出模式 可以应用在观察者模式（别的方法实现观察者模式&订阅模式）

6. Mutex

   如果一直有新goroutine到来且自旋的执行去执行Lock()操作，此时一个goroutine执行UnLcok()时，由于有自旋的goroutine将Woken位置为1，它是否就不会执行唤醒操作了？而不唤醒goroutine，也就无法计算它等待的时间差，怎么能进入饥饿模式呢？

7. GC调优 控制内存分配速度，限制goroutine数量从而提高赋值器对·cpu的利用率？如何提高

8. GC：三色标级法相较于标记清除的优点：1. 不需要扫描整个堆栈heap，而是遍历对象，根据引用情况遍历被引用的对象；2. 三色标记法可以和其他goroutine并发执行？

   

