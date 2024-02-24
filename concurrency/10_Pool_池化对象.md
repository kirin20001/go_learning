# Sync.Pool池化

- Go自带垃圾回收
- 如果你想使用 Go 开发一个高性能的应用程序的话，就必须考虑垃圾回收给性能带来的影响
- 因为Go 的自动垃圾回收机制还是有一个 STW（stop-the-world，程序暂停）的时间，
- 而且，大量地创建在堆上的对象，也会影响垃圾回收标记的时间。
- 另外数据库连接、TCP 的长连接，这些连接在创建的时候是一个非常耗时的操作。

## sync.Pool

- sync.Pool 数据类型用来保存一组可独立访问的临时对象。请注意这里加粗的“临时”这两个字，它说明了 sync.Pool 这个数据类型的特点，也就是说，它池化的对象会在未来的某个时候被毫无预兆地移除掉。而且，如果没有别的对象引用这个被移除的对象的话，这个被移除的对象就会被垃圾回收掉。
- sync.Pool 本身就是线程安全的，多个 goroutine 可以并发地调用它的方法存取对象；
- sync.Pool 不可在使用之后再复制使用。

## sync.Pool使用方法

- 1.NewPool struct 包含一个 New 字段，这个字段的类型是函数 func() interface{}。当调用 Pool 的 Get 方法从池中获取元素，没有更多的空闲元素可返回时，就会调用这个 New 方法来创建新的元素。如果你没有设置 New 字段，没有更多的空闲元素可返回时，Get 方法将返回 nil，表明当前没有可用的元素。有趣的是，New 是可变的字段。这就意味着，你可以在程序运行的时候改变创建元素的方法。当然，很少有人会这么做，因为一般我们创建元素的逻辑都是一致的，要创建的也是同一类的元素，所以你在使用 Pool 的时候也没必要玩一些“花活”，在程序运行时更改 New 的值。
- 2.Get如果调用这个方法，就会从 Pool取走一个元素，这也就意味着，这个元素会从 Pool 中移除，返回给调用者。不过，除了返回值是正常实例化的元素，Get 方法的返回值还可能会是一个 nil（Pool.New 字段没有设置，又没有空闲元素可以返回），所以你在使用的时候，可能需要判断。
- 3.Put这个方法用于将一个元素返还给 Pool，Pool 会把这个元素保存到池中，并且可以复用。但如果 Put 一个 nil 值，Pool 就会忽略这个值。

## 实现原理

- 存在两个主要字段 victim和local
- victim字段中的对象在gc过程中会被清除
- gc结束后local的数据会被赋给victim，local同样被清空
- ![pool](/Users/wyl/Desktop/go/src/go_learning/concurrency/img/10_Pool.webp)

### Get

- 首先会从local的private中获取可用元素。
- 若没有获取到，从local的shared中获取一个
- 如果依然没哟获取到，通过getslow方法从其他poolLocal的shared队列中获取一个（getslow方法会遍历所有shared队列，然后查找victim）
- 如果均没有可用对象，使用New函数创建一个新对象

### Put

- 首先会设置local的private
- 如果private已经有值，则将元素push到local的shared队列中

