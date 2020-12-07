 
 # 加锁时机
 最晚加锁，最早释放。                         
 
 
 # sync.atomic
 
 ## copy-on-write
 写操作时复制全量老数据到一个新的对象中，携带上本次新写的数据，之后利用原子替换（atomic.Value）,更新调用者的变量。来完成无锁访问共享数据。                
 
 
 # sync.Mutex
 几种mutex锁的实现：
 - Barging。当锁被释放的时候，会唤醒第一个等待者，然后把锁给第一个等待者或者给第一个请求锁的人。                   
 - handoff。当锁释放时，锁会一直持有直到第一个等待者准备好获取锁。它降低了吞吐量，因为锁被持有，即使另一个goroutine准备获取它。                   
 
 - spinning。自旋在等待队列为空或者应用程序重度使用锁时效果不错。              
 
 
 # errGroup
 把依赖多个微服务RPC需要聚合数据的任务，分解为依赖和并行，依赖的意思是，需要上游a的数据才能访问下游b的数据进行组合。           
 但并行的意思是：分解为多个小任务并行执行，最终等全部执行完毕。            
 
 核心原理：利用sync.Waitgroup管理并行执行的goroutine。

- 并行工作流
- 错误处理或优雅降级
- context传播和取消
- 利用局部变量和闭包


# sync.Pool


# context
context是面向请求的，通常放在函数的第一个参数上。

context.WithValue内部基于valueCtx实现。                
```go
    // A valueCtx carries a key-value pair. It implements Value for that key and
    // delegates all other calls to the embedded Context.
    type valueCtx struct {
    	Context
    	key, val interface{}
    }
```
为了实现不断的withvalue，构建新的context，内部在查找key的时候，使用递归不断从当前，从父节点寻找匹配的key，
直到root context(Background和TODO value函数会返回nil)                   

## 怎么做到的数据安全无data-race？
因为每次调用WithValue，都会返回一个新的context。通过递归去找数值时，这里是只读的。没有人会去篡改它。              
因为每次去追加值的时候，都返回一个新的context对象，父context不会被改动。                 
                 
所以，context value应该是immutable的，每次重新赋值都是新的context（使用context.WithValue(ctx,oldValue)）。                    



