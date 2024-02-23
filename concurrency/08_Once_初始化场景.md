# Once 一个简约而不简单的并发原语

- Once 可以用来执行且仅仅执行一次动作，常常用于单例对象的初始化场景。

## Once的使用场景

- sync.Once 只暴露了一个方法 Do，你可以多次调用 Do 方法，但是只有第一次调用 Do 方法时 f 参数才会执行，这里的 f 是一个无参数无返回值的函数。

  ```go
  func (o *Once) Do(f func())
  ```

- Once 常常用来初始化单例资源，或者并发访问只需初始化一次的共享资源，或者在测试的时候初始化一次测试资源。

- 例：

  ```go
     // 值是3.0或者0.0的一个数据结构
     var threeOnce struct {
      sync.Once
      v *Float
    }
    
      // 返回此数据结构的值，如果还没有初始化为3.0，则初始化
    func three() *Float {
      threeOnce.Do(func() { // 使用Once初始化
        threeOnce.v = NewFloat(3.0)
      })
      return threeOnce.v
    }
  ```

- 它将 sync.Once 和 *Float 封装成一个对象，提供了只初始化一次的值 v。 你看它的 three 方法的实现，虽然每次都调用 threeOnce.Do 方法，但是参数只会被调用一次。

  ## Once的实现
- Once 实现要使用一个互斥锁，这样初始化的时候如果有并发的 goroutine，就会进入doSlow 方法。
- 互斥锁的机制保证只有一个 goroutine 进行初始化，同时利用双检查的机制（double-checking），再次判断 o.done 是否为 0，如果为 0，则是第一次执行，执行完毕后，就将 o.done 设置为 1，然后释放锁。
```go
type Once struct {
    done uint32
    m    Mutex
}

func (o *Once) Do(f func()) {
    if atomic.LoadUint32(&o.done) == 0 {
        o.doSlow(f)
    }
}


func (o *Once) doSlow(f func()) {
    o.m.Lock()
    defer o.m.Unlock()
    // 双检查
    if o.done == 0 {
        defer atomic.StoreUint32(&o.done, 1)
        f()
    }
}
```

## Once常见错误

- 死锁

  - 你已经知道了 Do 方法会执行一次 f，但是如果 f 中再次调用这个 Once 的 Do 方法的话，就会导致死锁的情况出现。这还不是无限递归的情况，而是的的确确的 Lock 的递归调用导致的死锁。

    ```go
    func main() {
        var once sync.Once
        once.Do(func() {
            once.Do(func() {
                fmt.Println("初始化")
            })
        })
    }
    ```

- 未初始化

  - 如果 f 方法执行的时候 panic，或者 f 执行初始化资源的时候失败了，这个时候，Once 还是会认为初次执行已经成功了，即使再次调用 Do 方法，也不会再次执行 f。