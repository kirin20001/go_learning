# 线程安全的Map

- Go原生的Map采用Hash Table数据结构，通过key的索引，能够快速的找到对应的值

## map的基本使用方法

- 内建的map中key类型的k必须是可比较的
- bool、整数、浮点数、复数、字符串、指针、Channel、接口都是可比较的，包含可比较元素的 struct 和数组，这俩也是可比较的，而 slice、map、函数值都是不可比较的。

## map常见错误

- 未初始化
- 并发读写
  - Go 内建的 map 对象不是线程（goroutine）安全的，并发读写的时候运行时会有检查，遇到并发问题就会导致 panic。

## 如何实现线程安全的map类型

### 加读写锁：扩展map，支持并发读写

``` go
type RWMap struct { // 一个读写锁保护的线程安全的map
    sync.RWMutex // 读写锁保护下面的map字段
    m map[int]int
}
// 新建一个RWMap
func NewRWMap(n int) *RWMap {
    return &RWMap{
        m: make(map[int]int, n),
    }
}
func (m *RWMap) Get(k int) (int, bool) { //从map中读取一个值
    m.RLock()
    defer m.RUnlock()
    v, existed := m.m[k] // 在锁的保护下从map中读取
    return v, existed
}

func (m *RWMap) Set(k int, v int) { // 设置一个键值对
    m.Lock()              // 锁保护
    defer m.Unlock()
    m.m[k] = v
}

func (m *RWMap) Delete(k int) { //删除一个键
    m.Lock()                   // 锁保护
    defer m.Unlock()
    delete(m.m, k)
}

func (m *RWMap) Len() int { // map的长度
    m.RLock()   // 锁保护
    defer m.RUnlock()
    return len(m.m)
}

func (m *RWMap) Each(f func(k, v int) bool) { // 遍历map
    m.RLock()             //遍历期间一直持有读锁
    defer m.RUnlock()

    for k, v := range m.m {
        if !f(k, v) {
            return
        }
    }
}

```

### 分片加锁：更高效的并发map

- 减少锁的粒度常用的方法就是分片（Shard），将一把锁分成几把锁，每个锁控制一个分片。
- 加锁和分片加锁这两种方案都比较常用，如果是追求更高的性能，显然是分片加锁更好，因为它可以降低锁的粒度，进而提高访问此 map 对象的吞吐。如果并发性能要求不是那么高的场景，简单加锁方式更简单。

## sync.Map 

- 在特定场景下的Map

  - 只会增长的缓存系统中，一个 key 只写入一次而被读很多次；
  - 多个 goroutine 为不相交的键集读、写和重写键值对。

- 优点

  - 空间换时间。通过冗余的两个数据结构（只读的 read 字段、可写的 dirty），来减少加锁对性能的影响。对只读字段（read）的操作不需要加锁。
  - 优先从 read 字段读取、更新、删除，因为对 read 字段的读取不需要锁。
  - 动态调整。miss 次数多了之后，将 dirty 数据提升为 read，避免总是从 dirty 中加锁读取。
  - double-checking。加锁之后先还要再检查 read 字段，确定真的不存在才操作 dirty 字段。
  - 延迟删除。删除一个键值只是打标记，只有在提升 dirty 字段为 read 字段的时候才清理删除的数据。

  ```go
  type Map struct {
      mu Mutex
      // 基本上你可以把它看成一个安全的只读的map
      // 它包含的元素其实也是通过原子操作更新的，但是已删除的entry就需要加锁操作了
      read atomic.Value // readOnly
  
      // 包含需要加锁才能访问的元素
      // 包括所有在read字段中但未被expunged（删除）的元素以及新加的元素
      dirty map[interface{}]*entry
  
      // 记录从read中读取miss的次数，一旦miss数和dirty长度一样了，就会把dirty提升为read，并把dirty置空
      misses int
  }
  
  type readOnly struct {
      m       map[interface{}]*entry
      amended bool // 当dirty中包含read没有的数据时为true，比如新增一条数据
  }
  
  // expunged是用来标识此项已经删掉的指针
  // 当map中的一个项目被删除了，只是把它的值标记为expunged，以后才有机会真正删除此项
  var expunged = unsafe.Pointer(new(interface{}))
  
  // entry代表一个值
  type entry struct {
      p unsafe.Pointer // *interface{}
  }
  ```

  

### store

- 可以看出，Store 既可以是新增元素，也可以是更新元素。如果运气好的话，更新的是已存在的未被删除的元素，直接更新即可，不会用到锁。

- 如果运气不好，需要更新（重用）删除的对象、更新还未提升的 dirty 中的对象，或者新增加元素的时候就会使用到了锁，这个时候，性能就会下降。

  ```go
  func (m *Map) Store(key, value interface{}) {
      read, _ := m.read.Load().(readOnly)
      // 如果read字段包含这个项，说明是更新，cas更新项目的值即可
      if e, ok := read.m[key]; ok && e.tryStore(&value) {
          return
      }
  
      // read中不存在，或者cas更新失败，就需要加锁访问dirty了
      m.mu.Lock()
      read, _ = m.read.Load().(readOnly)
      if e, ok := read.m[key]; ok { // 双检查，看看read是否已经存在了
          if e.unexpungeLocked() {
              // 此项目先前已经被删除了，通过将它的值设置为nil，标记为unexpunged
              m.dirty[key] = e
          }
          e.storeLocked(&value) // 更新
      } else if e, ok := m.dirty[key]; ok { // 如果dirty中有此项
          e.storeLocked(&value) // 直接更新
      } else { // 否则就是一个新的key
          if !read.amended { //如果dirty为nil
              // 需要创建dirty对象，并且标记read的amended为true,
              // 说明有元素它不包含而dirty包含
              m.dirtyLocked()
              m.read.Store(readOnly{m: read.m, amended: true})
          }
          m.dirty[key] = newEntry(value) //将新值增加到dirty对象中
      }
      m.mu.Unlock()
  }
  
  func (m *Map) dirtyLocked() {
      if m.dirty != nil { // 如果dirty字段已经存在，不需要创建了
          return
      }
  
      read, _ := m.read.Load().(readOnly) // 获取read字段
      m.dirty = make(map[interface{}]*entry, len(read.m))
      for k, e := range read.m { // 遍历read字段
          if !e.tryExpungeLocked() { // 把非punged的键值对复制到dirty中
              m.dirty[k] = e
          }
      }
  }
  ```

  ### Load

  - Load 方法用来读取一个 key 对应的值。它也是从 read 开始处理，一开始并不需要锁。

    ```go
    func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
        // 首先从read处理
        read, _ := m.read.Load().(readOnly)
        e, ok := read.m[key]
        if !ok && read.amended { // 如果不存在并且dirty不为nil(有新的元素)
            m.mu.Lock()
            // 双检查，看看read中现在是否存在此key
            read, _ = m.read.Load().(readOnly)
            e, ok = read.m[key]
            if !ok && read.amended {//依然不存在，并且dirty不为nil
                e, ok = m.dirty[key]// 从dirty中读取
                // 不管dirty中存不存在，miss数都加1
                m.missLocked()
            }
            m.mu.Unlock()
        }
        if !ok {
            return nil, false
        }
        return e.load() //返回读取的对象，e既可能是从read中获得的，也可能是从dirty中获得的
    }
    ```

    ### Delete

    - 同样地，Delete 方法是先从 read 操作开始，原因我们已经知道了，因为不需要锁。

      ``` go
      func (m *Map) LoadAndDelete(key interface{}) (value interface{}, loaded bool) {
          read, _ := m.read.Load().(readOnly)
          e, ok := read.m[key]
          if !ok && read.amended {
              m.mu.Lock()
              // 双检查
              read, _ = m.read.Load().(readOnly)
              e, ok = read.m[key]
              if !ok && read.amended {
                  e, ok = m.dirty[key]
                  // 这一行长坤在1.15中实现的时候忘记加上了，导致在特殊的场景下有些key总是没有被回收
                  delete(m.dirty, key)
                  // miss数加1
                  m.missLocked()
              }
              m.mu.Unlock()
          }
          if ok {
              return e.delete()
          }
          return nil, false
      }
      
      func (m *Map) Delete(key interface{}) {
          m.LoadAndDelete(key)
      }
      func (e *entry) delete() (value interface{}, ok bool) {
          for {
              p := atomic.LoadPointer(&e.p)
              if p == nil || p == expunged {
                  return nil, false
              }
              if atomic.CompareAndSwapPointer(&e.p, p, nil) {
                  return *(*interface{})(p), true
              }
          }
      }
      ```

      