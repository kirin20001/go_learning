1. WaitGroup

   64位原子操作如何操作32为系统内存

   state1总共需要96bit，state()函数获得的是计数值加waiter数所在64位的指针？

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

   