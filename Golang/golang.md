# 协程
线程是进程中的执行体，拥有一个执行入口，以及从进程虚拟空间中分配的栈，包括用户栈和内核栈，操作系统会记录线程控制信息，而线程获得CPU时间片后才可以执行。
CPU这里的栈指针，指令指针等寄存器都要切换到对应的线程。如果线程自己又创建几个执行体，再给执行体指定自己的执行入口，申请一些内存给它们用作执行栈，
那么线程就可以按需调度这些执行体了。为了实现这些执行体的切换，线程也需要记录它们的控制信息。包括ID，栈的位置， 执行入口地址，执行现场等，
线程可以选择一个执行体来执行，此时CPU中的指令指针就会指向这个执行体的执行入口，栈基和栈指针寄存器也会指向线程给它分配的执行栈，要切换执行体时，
需要先保存当前执行体的执行现场，然后切换到另一个执行体。通过同样的方式，可以恢复到之前的执行体，这样就可以从中断的地方继续执行。
这些由线程创建的执行体就是所谓的协程。
因为用户程序不能操作内核空间，所有只能给协程分配用户栈，操作系统对协程一无所知，所有协程又被称为"用户态线程"。

协程思想的关键在于控制流的主动让出和恢复。
每个协程有自己的执行栈可以保存执行现场，所以可以由用户程序按需创建协程，协程主动让出执行权时，保存执行现场，然后切换到其他协程，
协程恢复执行时，会根据执行保存的执行现场，恢复到中断前的的状态继续执行，这样就通过协程实现了既轻量又灵活的由用户态进行调度的多任务模型。


# 并发模型 GMP
## 由来
早期单进程操作系统，进程顺序执行。
单进程带来了两个问题：
- 计算机只能一个接一个处理
- 进程阻塞带来cpu时间浪费
于是，多线程/多进程走向舞台。CPU调度器，根据时间片调度。
多进程/多线程解决了阻塞问题，但CPU切换成本提高。
- 并发设计复杂
- CPU调度消耗高
- 内存占用高

虚拟内存划分用户空间和内核空间，用户空间由用户线程使用，内核空间由操作系统线程使用。
我们把使用用户空间的线程，叫为协程。

用协程调度器管理协程，CPU只关注内核线程。
由此创建出协程对线程
 N：1 阻塞问题
 1：1 代价昂贵
 M：N  只需要优化协程调度器即可。
 
 golang语言的，就是GMP调度模型。
 
## 协程改进
 GO语言对协程做了改进。
 - 内存：几KB   可大量创建
 - 灵活调度     可常切换


## go调度器的历史
### 早期GM模型
M从全局先获取到锁，获取到锁后获取G。释放锁后，把G归还给队列。
- 创建，销毁，调度G需要每个M获取锁，形成锁竞争。
- 转移G造成额外的系统负载。
- 频繁的线程阻塞和取消阻塞增加了系统开销。

##  GMP模型设计思想
### 模型简介
G goroutine 协程
M machine 线程
P processor proccessor处理器

P 的个数可通过GOMAXPROCS进行设置。
P上面有G的本地队列 localP
全局G队列存放等待运行的G
P本地队列数量限制，最多不超过256个
优先把新创建的G放到P的本地队列，如果满了，会放到全局队列中。

P列表 程序启动时创建，最多有GOMAXPROCS个（可配置）
M列表 当前操作系统分配到Go程序的内核线程数

#### P和M的数量问题
P的数量
- 可以通过GOMAXPROCS环境变量设置。
- 可以在程序中通过runtime.GOMAXPROCS()设置。
M的数量
- go语言本身设定的M最大量为10000个。
- runtime/debug/SetMaxThreads设置


### 设计策略
- 复用线程：
    - work stealing机制
    - hand off机制
- 利用并行
    - 可利用GOMAXPROCS限定P的个数，= CPU/2
- 抢占
    - 每个G最多与CPU结合10ms，新的G就去抢占这个CPU
- 全局G队列
    - 优先从其他队列偷G，如果偷不到的时候，会从全局队列中取。
    


#### work stealing 机制
当某个P队列中没有G时，会去其他队列中偷取G

#### hand off 机制
当某个与M1结合的G1阻塞时，会创建或唤醒一个M3，让M1继续执行G1，而剩下的P带着本地队列里的G2 与 M3结合。
当G1执行完成后 M1休眠或销毁。


### go func 经历了什么
首先，会创建一个G，然后加入当前执行的MP组合的本地队列。如果本地队列已满，则加入到全局队列中。
然后，M就会尝试去本地队列获取G执行。如果本地队列为空，则从其他本地队列偷取或从全局队列中获取。
然后去真正运行G中的函数，时间片结束返回给M，G被放回到本地队列中。
如果在运行G的过程中发生系统调用阻塞，则会尝试创建或唤醒新的M，新的M去接管被阻塞G的P及其本地队列。
如果阻塞了的G执行完了，原来的M要么加入到休眠的队列，要么被销毁。而G会放到原来的P本地队列或全局队列中。


### 调度器的生命周期
#### M0
启动进程后，编号为0的主线程。
保存在runtime.m0，不在堆上分配。
负责执行初始化操作和启动第一个G。
启动完第一个G后，M0就跟其他M一样了。

#### G0
每次启动一个M，都会第一个创建的Goroutine，就是G0。G0不是每个进程中唯一的，而是每个线程中唯一的。
G0仅用于调度其他G。
G0不指向任何可执行的函数。

每个M都会有一个自己的G0，在调度或系统调用时会使用G0的栈空间。
M0的G0会放到全局空间。

### 可视化GMP编程
使用trace工具。
1. 创建trace文件 
    ```go
       	f, err := os.Create("trace.out")
    ```
2. 启动trace
    ```go
       	err = trace.Start(f)
    ```
3. 停止trace
    ```go
       	trace.Stop()
    ```
4. go build并执行后，会得到trace.out文件。

通过 go tool trace 工具打开trace文件。

### 通过debug trace 查看GMP信息
GODEBUG=schedtrace=1000  ./可执行程序
1000 指的是1000毫秒打印一次。
```go
    SCHED 0ms: gomaxprocs=12 idleprocs=9 threads=5 spinningthreads=1 idlethreads=0 runqueue=0 [0 0 0 0 0 0 0 0 0 0 0 0]
    SCHED 1005ms: gomaxprocs=12 idleprocs=10 threads=5 spinningthreads=1 idlethreads=2 runqueue=0 [0 0 0 0 0 0 0 0 0 0 0 0]
    hello gmp
```
SCHED 代表调试信息
0ms 代表程序从启动到输出的时间。
gomaxprocs 代表P的数量，一般默认与CPU核心数一致。
idleprocs 处于idle状态的p的数量
threads 线程数量，包括M0，也包括当前调试的线程。
spinningthreads 处于自旋状态的线程。
idlethreads 处于idle状态的线程
runqueue 全局G队列等待运行的G的数量
0数组，各个P的本地G队列中的G数量


### G调度场景
1. G的创建，由G1创建出的G3，为了保证局部性，会优先放到本地的P队列中。
2. G1执行完毕后，切换到G0，M1优先从本地队列中获取G2
3. G2开辟过多G，导致本地队列满的情况。对本地队列进行分割，将本地队列队首的一半，打乱顺序，放到全局队列中。同时将新创的也加入到全局队列中。
   然后剩余本地队列中的G向前移动到队首。本地队列再次变成未满状态，这时再创建的G优先加入到本地队列中。
4. 唤醒正在休眠的M，当正在执行的G2要创建一个新的G3时，会尝试从M休眠队列中唤醒一个M，新唤醒的M会尝试和一个P2去绑定，P2先和G0结合进行调度，
   这时，P2本地队列中还没有G，这种状态就是自旋线程。自旋就是为了不断寻找G。
5. 自旋线程从全局队列获取G。自旋线程优先从全局队列获取G，满足公式n = min(len(GQ)/gomaxprocs + 1, len(GQ/2)),其中GQ表示全局队列总长度，
   如果全局队列中没有G，才会从其他MP队列中偷取G
6. 自旋线程偷取G。从要偷取的P队列中，将队列一分为2，取后半部分偷取过来。
7. 自旋线程的最大限制。自旋线程+执行线程总数 <= GOMAXPROCS
8. G发生调用阻塞。执行中的G发生系统调用阻塞时，会尝试唤醒休眠M队列中的M与P进行绑定，如果没有M，那么P就进入空闲P队列。
   为什么不找自旋线程呢？因为自旋线程是抢占G的不是抢占P的。
9. 当发生阻塞的G不阻塞了，M会记住之前与之绑定的P，优先去获取原配，而P如果已经与新的M绑定了，则抢占失败。M会再尝试从空闲P队列中获取P。
   如果空闲P队列中也没有P，那么M就会放弃G，将G放到全局G队列中，M则会被强制休眠，加入到休眠M队列中。 

### 结构体
```go
    type g struct {
    	// Stack parameters.
    	// stack describes the actual stack memory: [stack.lo, stack.hi).
    	// stackguard0 is the stack pointer compared in the Go stack growth prologue.
    	// It is stack.lo+StackGuard normally, but can be StackPreempt to trigger a preemption.
    	// stackguard1 is the stack pointer compared in the C stack growth prologue.
    	// It is stack.lo+StackGuard on g0 and gsignal stacks.
    	// It is ~0 on other goroutine stacks, to trigger a call to morestackc (and crash).
    	stack       stack   // offset known to runtime/cgo
    	stackguard0 uintptr // offset known to liblink
    	stackguard1 uintptr // offset known to liblink
    
    	_panic       *_panic // innermost panic - offset known to liblink
    	_defer       *_defer // innermost defer
    	m            *m      // current m; offset known to arm liblink
    	sched        gobuf
    	syscallsp    uintptr        // if status==Gsyscall, syscallsp = sched.sp to use during gc
    	syscallpc    uintptr        // if status==Gsyscall, syscallpc = sched.pc to use during gc
    	stktopsp     uintptr        // expected sp at top of stack, to check in traceback
    	param        unsafe.Pointer // passed parameter on wakeup
    	atomicstatus uint32
    	stackLock    uint32 // sigprof/scang lock; TODO: fold in to atomicstatus
    	goid         int64
    	schedlink    guintptr
    	waitsince    int64      // approx time when the g become blocked
    	waitreason   waitReason // if status==Gwaiting
    
    	preempt       bool // preemption signal, duplicates stackguard0 = stackpreempt
    	preemptStop   bool // transition to _Gpreempted on preemption; otherwise, just deschedule
    	preemptShrink bool // shrink stack at synchronous safe point
    
    	// asyncSafePoint is set if g is stopped at an asynchronous
    	// safe point. This means there are frames on the stack
    	// without precise pointer information.
    	asyncSafePoint bool
    
    	paniconfault bool // panic (instead of crash) on unexpected fault address
    	gcscandone   bool // g has scanned stack; protected by _Gscan bit in status
    	throwsplit   bool // must not split stack
    	// activeStackChans indicates that there are unlocked channels
    	// pointing into this goroutine's stack. If true, stack
    	// copying needs to acquire channel locks to protect these
    	// areas of the stack.
    	activeStackChans bool
    	// parkingOnChan indicates that the goroutine is about to
    	// park on a chansend or chanrecv. Used to signal an unsafe point
    	// for stack shrinking. It's a boolean value, but is updated atomically.
    	parkingOnChan uint8
    
    	raceignore     int8     // ignore race detection events
    	sysblocktraced bool     // StartTrace has emitted EvGoInSyscall about this goroutine
    	sysexitticks   int64    // cputicks when syscall has returned (for tracing)
    	traceseq       uint64   // trace event sequencer
    	tracelastp     puintptr // last P emitted an event for this goroutine
    	lockedm        muintptr
    	sig            uint32
    	writebuf       []byte
    	sigcode0       uintptr
    	sigcode1       uintptr
    	sigpc          uintptr
    	gopc           uintptr         // pc of go statement that created this goroutine
    	ancestors      *[]ancestorInfo // ancestor information goroutine(s) that created this goroutine (only used if debug.tracebackancestors)
    	startpc        uintptr         // pc of goroutine function
    	racectx        uintptr
    	waiting        *sudog         // sudog structures this g is waiting on (that have a valid elem ptr); in lock order
    	cgoCtxt        []uintptr      // cgo traceback context
    	labels         unsafe.Pointer // profiler labels
    	timer          *timer         // cached timer for time.Sleep
    	selectDone     uint32         // are we participating in a select and did someone win the race?
    
    	// Per-G GC state
    
    	// gcAssistBytes is this G's GC assist credit in terms of
    	// bytes allocated. If this is positive, then the G has credit
    	// to allocate gcAssistBytes bytes without assisting. If this
    	// is negative, then the G must correct this by performing
    	// scan work. We track this in bytes to make it fast to update
    	// and check for debt in the malloc hot path. The assist ratio
    	// determines how this corresponds to scan work debt.
    	gcAssistBytes int64
    }
    
    type m struct {
    	g0      *g     // goroutine with scheduling stack
    	morebuf gobuf  // gobuf arg to morestack
    	divmod  uint32 // div/mod denominator for arm - known to liblink
    
    	// Fields not known to debuggers.
    	procid        uint64       // for debuggers, but offset not hard-coded
    	gsignal       *g           // signal-handling g
    	goSigStack    gsignalStack // Go-allocated signal handling stack
    	sigmask       sigset       // storage for saved signal mask
    	tls           [6]uintptr   // thread-local storage (for x86 extern register)
    	mstartfn      func()
    	curg          *g       // current running goroutine
    	caughtsig     guintptr // goroutine running during fatal signal
    	p             puintptr // attached p for executing go code (nil if not executing go code)
    	nextp         puintptr
    	oldp          puintptr // the p that was attached before executing a syscall
    	id            int64
    	mallocing     int32
    	throwing      int32
    	preemptoff    string // if != "", keep curg running on this m
    	locks         int32
    	dying         int32
    	profilehz     int32
    	spinning      bool // m is out of work and is actively looking for work
    	blocked       bool // m is blocked on a note
    	newSigstack   bool // minit on C thread called sigaltstack
    	printlock     int8
    	incgo         bool   // m is executing a cgo call
    	freeWait      uint32 // if == 0, safe to free g0 and delete m (atomic)
    	fastrand      [2]uint32
    	needextram    bool
    	traceback     uint8
    	ncgocall      uint64      // number of cgo calls in total
    	ncgo          int32       // number of cgo calls currently in progress
    	cgoCallersUse uint32      // if non-zero, cgoCallers in use temporarily
    	cgoCallers    *cgoCallers // cgo traceback if crashing in cgo call
    	park          note
    	alllink       *m // on allm
    	schedlink     muintptr
    	lockedg       guintptr
    	createstack   [32]uintptr // stack that created this thread.
    	lockedExt     uint32      // tracking for external LockOSThread
    	lockedInt     uint32      // tracking for internal lockOSThread
    	nextwaitm     muintptr    // next m waiting for lock
    	waitunlockf   func(*g, unsafe.Pointer) bool
    	waitlock      unsafe.Pointer
    	waittraceev   byte
    	waittraceskip int
    	startingtrace bool
    	syscalltick   uint32
    	freelink      *m // on sched.freem
    
    	// these are here because they are too large to be on the stack
    	// of low-level NOSPLIT functions.
    	libcall   libcall
    	libcallpc uintptr // for cpu profiler
    	libcallsp uintptr
    	libcallg  guintptr
    	syscall   libcall // stores syscall parameters on windows
    
    	vdsoSP uintptr // SP for traceback while in VDSO call (0 if not in call)
    	vdsoPC uintptr // PC for traceback while in VDSO call
    
    	// preemptGen counts the number of completed preemption
    	// signals. This is used to detect when a preemption is
    	// requested, but fails. Accessed atomically.
    	preemptGen uint32
    
    	// Whether this is a pending preemption signal on this M.
    	// Accessed atomically.
    	signalPending uint32
    
    	dlogPerM
    
    	mOS
    
    	// Up to 10 locks held by this m, maintained by the lock ranking code.
    	locksHeldLen int
    	locksHeld    [10]heldLockInfo
    }
    
    type p struct {
    	id          int32
    	status      uint32 // one of pidle/prunning/...
    	link        puintptr
    	schedtick   uint32     // incremented on every scheduler call
    	syscalltick uint32     // incremented on every system call
    	sysmontick  sysmontick // last tick observed by sysmon
    	m           muintptr   // back-link to associated m (nil if idle)
    	mcache      *mcache
    	pcache      pageCache
    	raceprocctx uintptr
    
    	deferpool    [5][]*_defer // pool of available defer structs of different sizes (see panic.go)
    	deferpoolbuf [5][32]*_defer
    
    	// Cache of goroutine ids, amortizes accesses to runtime·sched.goidgen.
    	goidcache    uint64
    	goidcacheend uint64
    
    	// Queue of runnable goroutines. Accessed without lock.
    	runqhead uint32
    	runqtail uint32
    	runq     [256]guintptr
    	// runnext, if non-nil, is a runnable G that was ready'd by
    	// the current G and should be run next instead of what's in
    	// runq if there's time remaining in the running G's time
    	// slice. It will inherit the time left in the current time
    	// slice. If a set of goroutines is locked in a
    	// communicate-and-wait pattern, this schedules that set as a
    	// unit and eliminates the (potentially large) scheduling
    	// latency that otherwise arises from adding the ready'd
    	// goroutines to the end of the run queue.
    	runnext guintptr
    
    	// Available G's (status == Gdead)
    	gFree struct {
    		gList
    		n int32
    	}
    
    	sudogcache []*sudog
    	sudogbuf   [128]*sudog
    
    	// Cache of mspan objects from the heap.
    	mspancache struct {
    		// We need an explicit length here because this field is used
    		// in allocation codepaths where write barriers are not allowed,
    		// and eliminating the write barrier/keeping it eliminated from
    		// slice updates is tricky, moreso than just managing the length
    		// ourselves.
    		len int
    		buf [128]*mspan
    	}
    
    	tracebuf traceBufPtr
    
    	// traceSweep indicates the sweep events should be traced.
    	// This is used to defer the sweep start event until a span
    	// has actually been swept.
    	traceSweep bool
    	// traceSwept and traceReclaimed track the number of bytes
    	// swept and reclaimed by sweeping in the current sweep loop.
    	traceSwept, traceReclaimed uintptr
    
    	palloc persistentAlloc // per-P to avoid mutex
    
    	_ uint32 // Alignment for atomic fields below
    
    	// The when field of the first entry on the timer heap.
    	// This is updated using atomic functions.
    	// This is 0 if the timer heap is empty.
    	timer0When uint64
    
    	// Per-P GC state
    	gcAssistTime         int64    // Nanoseconds in assistAlloc
    	gcFractionalMarkTime int64    // Nanoseconds in fractional mark worker (atomic)
    	gcBgMarkWorker       guintptr // (atomic)
    	gcMarkWorkerMode     gcMarkWorkerMode
    
    	// gcMarkWorkerStartTime is the nanotime() at which this mark
    	// worker started.
    	gcMarkWorkerStartTime int64
    
    	// gcw is this P's GC work buffer cache. The work buffer is
    	// filled by write barriers, drained by mutator assists, and
    	// disposed on certain GC state transitions.
    	gcw gcWork
    
    	// wbBuf is this P's GC write barrier buffer.
    	//
    	// TODO: Consider caching this in the running G.
    	wbBuf wbBuf
    
    	runSafePointFn uint32 // if 1, run sched.safePointFn at next safe point
    
    	// Lock for timers. We normally access the timers while running
    	// on this P, but the scheduler can also do it from a different P.
    	timersLock mutex
    
    	// Actions to take at some time. This is used to implement the
    	// standard library's time package.
    	// Must hold timersLock to access.
    	timers []*timer
    
    	// Number of timers in P's heap.
    	// Modified using atomic instructions.
    	numTimers uint32
    
    	// Number of timerModifiedEarlier timers on P's heap.
    	// This should only be modified while holding timersLock,
    	// or while the timer status is in a transient state
    	// such as timerModifying.
    	adjustTimers uint32
    
    	// Number of timerDeleted timers in P's heap.
    	// Modified using atomic instructions.
    	deletedTimers uint32
    
    	// Race context used while executing timer functions.
    	timerRaceCtx uintptr
    
    	// preempt is set to indicate that this P should be enter the
    	// scheduler ASAP (regardless of what G is running on it).
    	preempt bool
    
    	pad cpu.CacheLinePad
    }
```

M里面存了两个比较重要的东西，一个是g0，一个是curg。

g0：会深度参与运行时的调度过程，比如goroutine的创建、内存分配等
curg：代表当前正在线程上执行的goroutine。

## 调度策略
- 第一步，为了保证调度的公平性，每个工作线程每进行61次调度就需要优先从全局运行队列中获取goroutine出来运行，因为如果只调度本地运行队列中的goroutine，
则全局运行队列中的goroutine有可能得不到运行。
- 第二步，从工作线程本地对列中找G
- 第三步，如果全局队列为空，用findrunnable从其他工作线程的运行队列中偷取goroutine。

### 从全局获取G
从全局运行队列中获取可运行的goroutine是通过globrunqget函数来完成的。
```go
    func globrunqget(_p_ *p, max int32) *g {
    	if sched.runqsize == 0 {
    		return nil
    	}
    	// 计算要从全局拿的数量
    	n := sched.runqsize/gomaxprocs + 1
    	if n > sched.runqsize {
    		n = sched.runqsize
    	}
    	// 不能超过最大数
    	if max > 0 && n > max {
    		n = max
    	}
    	// 最多只能取本地队列的一半
    	if n > int32(len(_p_.runq))/2 {
    		n = int32(len(_p_.runq)) / 2
    	}
    
    	sched.runqsize -= n
    
    	gp := sched.runq.pop()
    	n--
    	for ; n > 0; n-- {
    		gp1 := sched.runq.pop()
    		runqput(_p_, gp1, false)
    	}
    	return gp
    }
```
该函数的第一个参数是与当前工作线程绑定的p，第二个参数max表示最多可以从全局队列中拿多少个g到当前工作线程的本地运行队列中来。
globrunqget函数首先会根据全局运行队列中goroutine的数量，函数参数max以及_p_的本地队列的容量计算出到底应该拿多少个goroutine，
然后把第一个g结构体对象通过返回值的方式返回给调用函数，其它的则通过runqput函数放入当前工作线程的本地运行队列。这段代码值得一提的是，
计算应该从全局运行队列中拿走多少个goroutine时根据p的数量（gomaxprocs）做了负载均衡

### 从工作线程本地队列获取
```go
    func runqget(_p_ *p) (gp *g, inheritTime bool) {
    	// If there's a runnext, it's the next G to run.
    	for {
    		next := _p_.runnext
    		if next == 0 {
    			break
    		}
    		if _p_.runnext.cas(next, 0) {
    			return next.ptr(), true
    		}
    	}
    
    	for {
    		h := atomic.LoadAcq(&_p_.runqhead) // load-acquire, synchronize with other consumers
    		t := _p_.runqtail
    		if t == h {
    			return nil, false
    		}
    		gp := _p_.runq[h%uint32(len(_p_.runq))].ptr()
    		if atomic.CasRel(&_p_.runqhead, h, h+1) { // cas-release, commits consume
    			return gp, false
    		}
    	}
    }
```

这里首先需要注意的是不管是从runnext还是从循环队列中拿取goroutine都使用了cas操作，这里的cas操作是必需的，
因为可能有其他工作线程此时此刻也正在访问这两个成员，从这里偷取可运行的goroutine。

其次，代码中对runqhead的操作使用了atomic.LoadAcq和atomic.CasRel，它们分别提供了load-acquire和cas-release语义。

对于atomic.LoadAcq来说，其语义主要包含如下几条：

- 原子读取，也就是说不管代码运行在哪种平台，保证在读取过程中不会有其它线程对该变量进行写入；
- 位于atomic.LoadAcq之后的代码，对内存的读取和写入必须在atomic.LoadAcq读取完成后才能执行，编译器和CPU都不能打乱这个顺序；
- 当前线程执行atomic.LoadAcq时可以读取到其它线程最近一次通过atomic.CasRel对同一个变量写入的值，与此同时，位于atomic.LoadAcq之后的代码，
不管读取哪个内存地址中的值，都可以读取到其它线程中位于atomic.CasRel（对同一个变量操作）之前的代码最近一次对内存的写入。
对于atomic.CasRel来说，其语义主要包含如下几条：

- 原子的执行比较并交换的操作；
- 位于atomic.CasRel之前的代码，对内存的读取和写入必须在atomic.CasRel对内存的写入之前完成，编译器和CPU都不能打乱这个顺序；
- 线程执行atomic.CasRel完成后其它线程通过atomic.LoadAcq读取同一个变量可以读到最新的值，与此同时，位于atomic.CasRel之前的代码对内存写入的值，
可以被其它线程中位于atomic.LoadAcq（对同一个变量操作）之后的代码读取到。因为可能有多个线程会并发的修改和读取runqhead，
以及需要依靠runqhead的值来读取runq数组的元素，所以需要使用atomic.LoadAcq和atomic.CasRel来保证上述语义。

我们可能会问，为什么读取p的runqtail成员不需要使用atomic.LoadAcq或atomic.load？因为runqtail不会被其它线程修改，
只会被当前工作线程修改，此时没有人修改它，所以也就不需要使用原子相关的操作。

### sysmon
sysmon是我们的保洁阿姨，它是一个M，又叫监控线程，不需要P就可以独立运行，每20us~10ms会被唤醒一次出来打扫卫生，
主要工作就是回收垃圾、回收长时间系统调度阻塞的P、向长时间运行的G发出抢占调度等等。

# GC
## 1.3版本最开始采用标记清除法。
带来的问题：
- 需要长时间stoptheworld
- 标记需要扫描整个heap
- 清除数据会产生很多的heap碎片（不清楚对象是否关联）

## 1.5版本的三色标记法 
白色标记集合、灰色标记集合、黑色标记集合。
- 创建时，所有对象标记为白色，对象放入白色集合。
- 程序遍历根节点集合，只遍历一层，将第一层对象标记为灰色节点放到灰色集合中。
  再遍历一层灰色标记集合，将可达对象，从白色标记为灰色，之前标记为灰色的，标记为黑色。
- 重复上一步，直到灰色标记表中无任何对象。
- 收集所有白色对象

### 三色标记法问题
如果三色标记不使用STW的话，
- 已经标记为黑色的对象，重新引用到白色的对象；
- 且灰色对象与这个白色对象之间正好解除引用。

因为黑色的对象不会再被扫描，而白色对象就会等待被回收。
那么，就要使用STW。

## 强弱三色不变式
破坏掉三色标记出问题的两个条件即可。
### 强三色不变式
强制要求黑色对象不能引用白色对象。这样就破坏了条件1

### 弱三色不变式
黑色对象可以引用白色对象，但要求有灰色对象对它的引用。这样就破坏了条件2


**屏障机制就是为了满足强三色不变式或者弱三色不变式。** 


## 屏障
屏障就是在程序执行中，额外增加的判断机制。

### 插入写屏障
当对象被引用时触发的机制
在A对象引用B对象时，B对象被标记为灰色。
满足：强三色不变式

栈本身空间比较小，为了保证性能，插入写屏障不在栈上使用。
栈在清除白色对象前，启动STW，同时将所有对象置为白色，重新扫描一遍，再做清除。

#### 插入写屏障的不足
就是在最后时，需要STW来扫描栈，大约需要10-100ms


### 删除写屏障
当对象被删除时触发的机制
被删除的对象，如果自身为灰色或者白色，那么被标记为灰色。
满足：弱三色不变式 

#### 删除写屏障的不足
回收精度比较低
一个对象即使被删除了，最后一个指向它的指针依然可以活过一轮。在下一轮GC中被清理掉。


## 1.8版本后三色标记法+混合写屏障机制
1. GC开始，优先扫描栈，将栈上的所有可达对象全部扫描并标记为黑色（之后不再进行二次扫描，无需STW）。栈中不启用屏障。
2. GC期间，任何在栈上创建的新对象，都为黑色。
3. 被删除对象标记为灰色。
4. 被添加对象标记为灰色。

满足：变形的弱三色不变式

### 几个典型场景
#### 对象被一个堆对象删除引用，成为栈对象的下游。
前提： 堆对象4 -> 对象7
此时，栈对象1引用对象7，而堆对象4删除对对象7的引用。

经历什么过程呢？
GC开始，堆对象4被标记为灰色。
栈对象1增加对对象7的引用，同时堆对象4删除对对象7的引用。
此时栈中对象1依然为黑色，不启用屏障，而堆中因为启用着屏障，于是对象7被标记为灰色。

#### 对象被一个栈对象删除引用，成为另一个栈对象的下游
栈对象2引用栈对象3
新建的栈对象9引用了栈对象3，同时栈对象2删除对栈对象3的引用。

因为在GC过程中新创建的栈对象均为黑色，因此对象9就为黑色。

#### 对象被一个堆对象删除引用，成为另一个堆对象的下游
假设堆对象10当前颜色为黑色。
堆对象10添加下游引用堆对象7，触发屏障机制，堆对象7被标记为灰色。堆对象4删除堆对象7的引用，触发屏障机制，堆对象7被标记为灰色。

#### 对象从一个栈对象被删除引用，成为另一个堆对象的下游。
栈对象1删除对栈对象2的引用
堆对象3删除对堆对象4的引用，同时堆对象3添加堆栈对象2的引用。


栈对象1删除堆对象2的引用，不触发屏障机制。
堆对象3删除堆对象4的引用，触发屏障，堆对象4被标记为灰色。
堆对象3增加对栈对象2的引用，因为栈对象2已经是黑色了，无屏障操作。这样就保护了堆对象4及其下游引用。


## GC触发条件
1. 内存大小阈值,内存达到上次GC的2倍
2. 达到定时时间 2m interval 


# Go的内存模型
在并发环境中多goroutine读取相同变量时，变量的可见性条件。       
即：在什么情况下，goroutine在读取一个变量的值时，能够看到其他goroutine对这个变量的写的结果。     

程序运行的时候，两个操作的执行顺序不一定得到保证。由此引出一个重要概念：        
**Happens-before**
在一个goroutine内部，程序的执行顺序和它们的代码指定的顺序是一样的，即便编译器和CPU重排了读写顺序，从行为上来看，也和代码的指定顺序一样。      
    
**go只保证goroutine的内部重排对读写的顺序没有影响。**

如果要保证多个goroutine对一个共享变量的顺序，可以使用并发原语为读写建立happens-before关系，来保证顺序。


保证的happens-before：  
- init函数    
    - main函数一定是在导入包的init函数之后执行。
- 后面还有很多...

  

# 闭包用法

函数是头等对象，可以作为参数传递，可以做函数的返回值，也可以绑定到变量。Go语言称这样的参数、返回值或变量为function value。
function value本质上是一个指针，但并不直接指向函数指令入口，而是指向一个runtime.funcval结构体。这个结构体里只有一个地址，就是函数指令的入口地址。
## 闭包
- 必须要有在函数外部定义，在函数内部被引用的自由变量
- 脱离了形成闭包的上下文，闭包也能照常使用这些自由变量。
这个自由变量，称之为捕获变量。

```go
    package main
    
    import (
    	"fmt"
    )
    
    func create() func() int {
    	c := 2
    	return func() int {
    		return c
    	}
    }
    
    func main() {
    	f1 := create()
    	f2 := create()
    	fmt.Println(f1())
    	fmt.Println(f2())
    
    }
```
函数create被赋值给f1和f2两个变量。这种情况编译器会做出优化，让f1和f2共用一个funcval结构体。
闭包函数的指令在编译阶段完成，但因为每个闭包对象都要保存自己的捕获变量，所以要到执行阶段才创建对应的闭包对象。
到执行阶段，main函数栈帧有两个局部变量f1和f2，然后是返回值空间。到create函数栈帧这里，有一个局部变量c=2，create函数会在堆上分配一个funcval结构体，
fn指向闭包函数入口，除此之外，还有一个捕获列表，这里只捕获了一个变量c。然后这个结构体的起始地址就作为返回值写入返回值空间，所以，f1被赋值为addr2.

go语言中，通过funcval调用函数时，会把对应的funcval结构体地址存入特定的寄存器中，amd64平台用的是DX寄存器。
这样在闭包函数中，就可以通过寄存器取出funcval结构体的地址，然后，加上相应的偏移量来找到每个被捕获的变量。
所以，GO语言中，闭包就是有捕获列表的function value。


被闭包捕获的变量，要在外层函数与闭包函数中表现一致，好像他们在使用同一个变量。为此，go语言的编译器针对不同情况做了不同处理，
被捕获变量没有被任何修改的话，直接拷贝值到捕获列表就行了。
但如果除了初始化赋值外，还被修改过，那就要再细分了。

### 初始化赋值后被修改的局部变量场景
```go
    package main
    
    import (
    	"fmt"
    )
    
    func create() (fs [2]func()) {
    	for i := 0; i < 2; i++ {
    		fs[i] = func() {
    			fmt.Println(i)
    		}
    	}
    	return
    }
    
    func main() {
    	f1 := create()
    	for i := 0; i < len(f1); i++ {
    		f1[i]()
    	}
    
    }
```
闭包函数指令入口地址addrf,main函数栈帧中局部变量f1是一个长度为2的function value数组，返回值也是。
到create栈帧，由于被闭包捕获，局部变量i改为堆分配，在栈上只存一个地址&i。
第一次for循环，在堆上创建funcval结构体，捕获i的地址。这样，闭包函数就和外层函数操作同一个变量了。
addr0 
&i
fn=addrf

返回值第一个元素存储addr0
第一次for循环结束i自增1 i=1

第二次for循环，再次堆分配funcval捕获变量i地址。
addr1
&i
fn=addrf

返回值第二个元素存储addr1
第二次循环结束i自增1  i=2

达到退出循环条件，create函数结束，把返回值拷贝到局部变量f1。通过f1[0]调用函数时，把addr0存入寄存器，闭包函数通过寄存器的地址加上偏移找到捕获变量i地址。
通过f1[1]调度用函数时，把addr1存入寄存器，闭包函数找到addr1。被捕获的地址都指向i=2,所以每次都会打印2。

**闭包导致的局部变量堆分配，也是变量逃逸的一种场景。**

### 修改并被捕获的是参数场景
涉及到函数原型，就不能像局部变量那样处理了。
参数通过调用者栈帧传入，但是编译器会把栈上这个参数拷贝到堆上一份，然后外层函数和闭包函数都使用堆上分配的这个。


### 如果被捕获的是返回值场景
调用者栈帧上依然会分配返回值空间，不过，闭包的外层函数会在堆上也分配一个，外层函数和闭包函数都使用堆上分配的这个。但是在外层函数返回前，
需要把堆上的返回值拷贝到栈上的返回值空间。

# defer
defer在定义时，对外部变量的引用有两种方式。
- 作为函数参数： 传值，cache起来。
- 作为闭包： 在defer真正调用的时候，要根据上下文。


## defer执行顺序
函数返回前倒序执行。

### 倒序执行的实现
defer指令对应两部分的内容。
- deferproc 负责把要执行的函数信息保存起来，我们称之为defer注册，deferproc函数会返回0。
- 返回之前通过deferreturn执行

注册的defer函数。
正是先注册，后调用，实现了延迟执行的效果。

defer信息会注册到一个链表，当前执行的goroutine持有这个链表的头指针。存在runtime.g结构体中的_defer变量中，指向链表头。defer链表链起来的是
一个一个的_defer结构体,新注册的defer会添加到链表头，执行时也是从头开始，所以defer才会表现为倒序执行。

### _defer
```go
    type _defer struct {
    	siz     int32   // 参数和返回值占的字节数
    	started bool    // 是否已执行
    	heap    bool    // 是否是堆分配
    	sp      uintptr // 调用者栈指针
    	pc      uintptr // deferproc返回地址
    	fn     *funcval // 注册的函数
    	_panic *_panic  // 
    	link   *_defer  // 链到前一个注册的结构体
    }
```

### deferproc函数
func deferproc(siz int32, fn *funval)
siz: 函数参数+返回值占用空间大小
### deferreturn 函数


## defer 1.13 1.14 做了优化
1.12版本 defer通过 deferproc注册函数信息，_defer结构体分配在堆上。
1.13中，通过使用局部变量，将变量保存在栈上。再通过deferprocStack将栈上这个_defer结构体注册到defer链表中。
1.13版本主要的优化点：
**减少defer信息的堆分配**

1.14版本中，
在编译阶段插入代码，把defer函数的执行逻辑展开在所属函数内，从而免于创建_defer结构体，而且不需要注册到defer链表中。
但是在panic时会比较复杂，因为没有注册到defer链表中，需要采用栈扫描的方式来发现，于是_defer结构体又增加了几个字段


显式循环和隐式循环依然使用1.12版本处理方式。


## defer和return在一起的时候
拆解：
返回值 = xxx
defer
空return 


## defer recover使用
recover（）只有在defer的上下文中才有效，（且只有通过defer中用匿名函数调用才有效）直接调用的话只会返回nil


# 逃逸分析
在编译原理中，分析指针动态范围的方法称之为逃逸分析。
通常来说，当一个对象的指针被多个方法或线程引用时，我们称这个指针发生了逃逸。
**逃逸分析决定一个变量是分配在堆上还是栈上** 

- 堆：适合不可预知大小的内存分配，分配速度慢，会形成内存碎片。需要通过垃圾回收去释放内存。
- 栈：分配内存只需要 PUSH RELEASE
 
通过逃逸分析，可以尽量将不需要分配到堆上的变量直接分配到栈上。堆上变量少，减少分配堆内存的开销，也会减轻gc压力。

## 逃逸分析基本原则
如果一个函数返回对一个变量的引用，那么它就会发生逃逸。

## 逃逸分析是如何完成的。
编译器分析代码的特征和代码生命周期，GO中的变量只有在编译器可以证明函数返回后不会再被引用，才分配到栈上，其他情况都分配到堆上。

编译器会根据变量是否被外部引用来决定是否逃逸。
- 如果函数外部没有引用，分配到栈上。有一种情况除外，局部变量所需内存过大。
- 如有，堆上。

## 观察逃逸分析的命令
```go
    go build -gcflags '-m' xx.go
```

# mutex
用于解决资源并发访问问题。  
如：
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

count++ 由于并发访问，加的数就出问题。  
使用命令 go run -race xx.go 就可以检测data race问题。它是编译器通过探测所有内存访问，加入代码能监视对这些内存的访问，在代码运行时，
就可以监控到堆共享变量的非同步访问。

有Lock（）和Unlock（）两个方法。

在实现上，有：  
- lockfast： 直接获得锁
- lockslow： 经过一系列判断获得锁。

## mutex实现
mutex有两种模式：  
- 正常模式
- 饥饿模式


**在正常模式下**，waiter FIFO，被唤醒的waiter不是直接获得锁，而是和新来的进行竞争。（因为新来的在时间片内，不需要进行上下文的切换）  
如果没有竞争过新来的，被唤醒的waiter就会被插到等待队列的队首。如果waiter获取不到锁的时间超过了1ms，就会进入到饥饿模式。    

**饥饿模式下**，mutex的拥有者将直接把锁交给队首的waiter（这是为了防止老的goroutine一直获取不到锁苦苦等待），新来的goroutine不会尝试获取锁，
即使看起来锁没有被持有，也不会去抢，也不会自旋，而是加入到等待队列的队尾。

**什么时候切换回正常模式呢？**
- 当前持有mutex的waiter发现，自己已经是最后一个等待的waiter了，后面没人来了。
- 自己获取锁少于1ms了。

## 等待waiter队列最大数量是多少？
state 是int32类型 出去3位标记位，就是32-3 = 29,即2的29次方 - 1 约等于5亿个。一个goroutine差不多占2k,5亿个也不过占1T左右。



# RWMutex
reader/writer 互斥锁  
同一时刻，可以由任意数量的reader持有，或者被单个writer持有。

- Lock()、Unlock() 写操作时调用方法
- RLock()、RUnlock() 读操作时调用方法
- RLocker 返回一个读对象。

大量并发读，少量并发写的场景，可以使用RWMutex。  

## RWMutex实现
- 写优先
一个正在阻塞的Lock调用，会排除新的reader请求到锁。  
```go
type RWMutex struct {
  w           Mutex   // 互斥锁解决多个writer的竞争
  writerSem   uint32  // writer信号量
  readerSem   uint32  // reader信号量
  readerCount int32   // reader的数量
  readerWait  int32   // writer等待完成的reader的数量
}

const rwmutexMaxReaders = 1 << 30
```
两个信号量

#channel原理及实现
基于通信顺序进程模型（CSP）思想，设计而来。

## channel几种使用类型
- 数据交流
    - 当作并发的buffer或queue，解决生产者消费者问题。
- 数据传递
    - 一个goroutine把数据交给另一个goroutine
- 信号通知
    - 一个goroutine把信号传递给另一个goroutine
- 任务编排
    - 让一组goroutine按照一定顺序并发
- 锁
    - 利用channel实现互斥锁

## 基本用法

- 只接收  
    <-chan int 只能从chan接收int 
- 只发送    
    chan<- struct{}  只能发送struct
- 既接收又发送   
    chan string 可以发送、接收string

箭头指向chan，就表示可以往里塞数据。  
箭头原理chan，就表示往外吐数据。   

chan中的元素可以是任意类型。  
如：  
chan<- chan int   
chan<- <-chan int   
<-chan <-chan int   
chan (<-chan int)   

**未初始化的chan零值是nil。对nil的chan发送接收数据总会阻塞**  


## 实现原理
### 数据结构
```go
    type hchan struct {
    	qcount   uint           // total data in the queue 队列元素数量
    	dataqsiz uint           // size of the circular queue 队列大小
    	buf      unsafe.Pointer // points to an array of dataqsiz elements 队列指针
    	elemsize uint16         // chan 中元素大小
    	closed   uint32         // 是否已关闭    
    	elemtype *_type // element type chan中元素类型
    	sendx    uint   // send index send在buf中的索引
    	recvx    uint   // receive index recv在buf中的索引
    	recvq    waitq  // list of recv waiters receiver的等待队列
    	sendq    waitq  // list of send waiters send的发送队列
    
    	// lock protects all fields in hchan, as well as several
    	// fields in sudogs blocked on this channel.
    	//
    	// Do not change another G's status while holding this lock
    	// (in particular, do not ready a G), as this can deadlock
    	// with stack shrinking.
    	lock mutex  // 互斥锁，用于保护所有字段
    }
    
    type waitq struct {
    	first *sudog
    	last  *sudog
    }
```

### 初始化
编译器根据容量大小，选择调用makechan64还是makechan。  
```go
    func makechan64(t *chantype, size int64) *hchan {
    	if int64(int(size)) != size {
    		panic(plainError("makechan: size out of range"))
    	}
    
    	return makechan(t, int(size))
    }
    
    func makechan(t *chantype, size int) *hchan {
    	elem := t.elem
    
    	// compiler checks this but be safe.
    	if elem.size >= 1<<16 {
    		throw("makechan: invalid channel element type")
    	}
    	if hchanSize%maxAlign != 0 || elem.align > maxAlign {
    		throw("makechan: bad alignment")
    	}
    
    	mem, overflow := math.MulUintptr(elem.size, uintptr(size))
    	if overflow || mem > maxAlloc-hchanSize || size < 0 {
    		panic(plainError("makechan: size out of range"))
    	}
    
    	// Hchan does not contain pointers interesting for GC when elements stored in buf do not contain pointers.
    	// buf points into the same allocation, elemtype is persistent.
    	// SudoG's are referenced from their owning thread so they can't be collected.
    	// TODO(dvyukov,rlh): Rethink when collector can move allocated objects.
    	var c *hchan
    	switch {
    	case mem == 0:
    		// Queue or element size is zero. 如果chan的大小或者元素的size是0，不必创建buf.
    		c = (*hchan)(mallocgc(hchanSize, nil, true))
    		// Race detector uses this location for synchronization.
    		c.buf = c.raceaddr()
    	case elem.ptrdata == 0:
    		// Elements do not contain pointers. 如果元素类型不是指针，分配一块连续的内存给hchan和buf
    		// Allocate hchan and buf in one call.
    		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
    		c.buf = add(unsafe.Pointer(c), hchanSize)
    	default:
    		// Elements contain pointers. 元素包含指针，单独分配buf
    		c = new(hchan)
    		c.buf = mallocgc(mem, elem, true)
    	}
    
    	c.elemsize = uint16(elem.size)
    	c.elemtype = elem
    	c.dataqsiz = uint(size)
    	lockInit(&c.lock, lockRankHchan)
    
    	if debugChan {
    		print("makechan: chan=", c, "; elemsize=", elem.size, "; dataqsiz=", size, "\n")
    	}
    	return c
    }
```

### send
在发送数据给chan的时候，会把send语句转换为chansend1函数，chansend1调用chansend。        
先判断是否为nil，如果是调用gopark阻塞休眠。   
如果chan已满，但还没有close，直接返回。      
如果已经close，再发送数据报panic。      
如果等待队列中有等待的receiver，就把它从队列中弹出，然后直接把数据交给它，而不需要放到buf中，速度可以快一些。        
如果没有等待的receiver，就把数据放到buf中，放入后就返回。  
如果buf满了，发送者的goroutine就会加入到发送者等待队列，直到被唤醒。    

### recv
调用chanrecv1函数，要两个返回值，会调用chanrecv2。他们俩都会调用chanrecv函数。    
chan为nil时，调用者会被永远阻塞。    
如果chan已经被close，并且队列中没有缓存的元素，那么返回true,false. 
如果是unbuffer的chan，就把sender的数据复制给receiver，否则就从队列头部取一个值，并把sender的放到队列尾部。    
如果没有等待的sender，如果buf中有元素，就取个元素给receiver。     
如果buf中没有元素，当前receiver就会阻塞，直到从sender接收了数据，或者chan被close，才会返回。     

### close
通过close，可以把chan关闭。调用closechan函数。        

如果chan 为nil，close会panic。如果chan已经被close，再次close会panic。   
如果chan不为nil，也没有被close过，就把等待队列中的sender，recver从队列中全部移除并唤醒。


## 易出问题的坑
- 关闭nil的chan
- send已经close的chan
- close已经close的chan

![Image_text](https://raw.githubusercontent.com/jizengguang/PrepareForInterview/master/Picture/chan_method.png)



# 并发使用的手段
chan+go                 
waitgroup





# 内存泄露分析
goroutine被阻塞，无法被gc，会造成内存泄漏。
通常情况下，是在使用chan时，一个goroutine往无缓冲的chan中写入数据，而数据因为某些原因没有被接收。
如：      
```go
    func process(timeout time.Duration) bool {
    	ch := make(chan bool)
    
    	go func() {
    		time.Sleep(timeout + time.Second)
    		ch <- true
    		fmt.Println("exit goroutine")
    	}()
    
    	select {
    	case result := <-ch:
    		return result
    	case <-time.After(timeout):
    		return false
    	}
    }
// 如果超时先发生，第13行将被永远阻塞。造成goroutine泄漏。
// 因为unbuffer的chan必须 reader，writer同时准备好才行。
```

# context
## context用途
- 上下文信息传递，如：处理HTTP请求，在请求链路上传递信息     
- 控制子goroutine的运行
- 超时控制的方法调用
- 可以取消的方法调用

## 实现
```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```
4个方法：        
- Deadline()  返回这个context被取消的截止日期。如果没有设置，ok返回false。     
- Done()  返回一个chan对象，在context被取消时，这个chan会被close。        
    - 如果Done（）没有被close，Err返回nil，如果被close了，Err返回close的原因。
- Err()
- Value()   返回与context绑定的key的value。     

- context.Background()
- context.ToDo()
都是返回一个非nil的空context对象。无截止时间，不会被cancel，不会超时。一般用在主函数、初始化、测试等。。。       

一般性规则：
- 一般当参数时放第一个
- 不把nil当context参数值。
- 用来做临时上下文透传的，不要持久化和长久保存。
- key不要是字符串类型或者其他内建类型。
- 常使用struct{}当key的类型。

**几个特殊用途的方法：**
- WithValue 基于parent context 生成了新的context，保存了key-value，用于传递上下文.
```go        
type valueCtx struct {
    Context
    key, val interface{}
}
```     
优先从自己的key-val中找，如果没找到，就去parent里找。

- WithCancel 
    返回parent的副本，只是done channel是新建对象，它的类型是cancelCtx。     
    主要用于需要主动取消的长时间任务时。正常完成了，需要调用cancel()方法，切记。 
- WithTimeout
    与WithDeadline增加了超时时间的参数，超时时间+当前时间就是截止时间。        
- WithDeadline
    设置一个timerCtx。
    timerCtx被close主要有三个原因：    
    - 截止时间到了
    - cancel被调用了
    - parent的Done被close了。       


# slice和数组
数组array和切片slice都是集合类型，用来存储某一种类型的值。  
**区别：**  
- 数组的长度是固定的，切片是可变长的。             
- 数组是值类型，切片是引用类型。
- 数组的容量永远等于长度。

数组的长度声明的时候就需要给定。        

我们可以看成切片是对数组的简单封装。
数组可以叫成是切片的底层数组。而切片也可以看作是对数组的某个连续片段的引用。

## 切片扩容
不会改变原有切片，而是生成一个容量更大的切片，把原有元素和新元素一并拷贝到新切片中。      
**扩容原则：**
1. 一般情况下，新容量会是原容量的两倍。       
2. 如果原长度大于等于1024，将会以原容量的1.25倍为基准。新基准会不断与1.25相乘，直到结果不小于原长度与要追加的元素数量之和。
3. 如果一次性追加的元素过多，使新长度比原容量的2倍还要大。那么会以新长度为基准。

无需扩容时，append指向的是原底层数组的新切片，扩容时，append指向的新底层数组的新切片。
因此，切片的底层数组其实是不会被替换的。            


# interface原理及实现
类型元数据结构体_type，作为每个类型元数据的header。           
```go
    type _type struct {
    	size       uintptr 
    	ptrdata    uintptr // size of memory prefix holding all pointers
    	hash       uint32
    	tflag      tflag
    	align      uint8
    	fieldAlign uint8
    	kind       uint8
    	// function for comparing objects of this type
    	// (ptr to object A, ptr to object B) -> ==?
    	equal func(unsafe.Pointer, unsafe.Pointer) bool
    	// gcdata stores the GC type data for the garbage collector.
    	// If the KindGCProg bit is set in kind, gcdata is a GC program.
    	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
    	gcdata    *byte
    	str       nameOff
    	ptrToThis typeOff
    }
```
自定义类型
```go
    type uncommontype struct {
    	pkgpath nameOff  //类型所在包路径
    	mcount  uint16 // number of methods  该类型关联方法数量
    	xcount  uint16 // number of exported methods
    	moff    uint32 // offset from this uncommontype to [mcount]method 方法元数据数组的偏移量
    	_       uint32 // unused
    }
```

type Mytype = int32 // 这种是为int32起别名，对应的都是int32类型元数据。
type Mytype int32  // 这种是自定义类型，Mytype


**interface有空接口和非空接口两种**   
## 空接口 interface{}
空接口可以接收任意类型的数据，只需要记录这个数据在哪，是什么类型的就足够了。
runtime.eface
```go
     type eface struct {
        	_type *_type // 指向动态类型接口元数据
        	data  unsafe.Pointer // 指向接口的动态值
        }
```
var e interface{}       
空接口e在未赋值之前，_type和data都是nil      
当e被赋值时，data等于赋值的变量，_type指向该变量的类型元数据。
## 非空接口
就是有方法的接口
interface {
    A()
    B()
}
变量赋值给非空接口类型时，必须实现非空接口的所有方法。

```go
    type iface struct {
    	tab  *itab      // 接口方法列表和接口动态类型信息
    	data unsafe.Pointer // 指向接口的动态值
    }
    // layout of Itab known to compilers
    // allocated in non-garbage-collected memory
    // Needs to be in sync with
    // ../cmd/compile/internal/gc/reflect.go:/^func.dumptabs.
    type itab struct {
    	inter *interfacetype  // 指向interface的类型元数据
    	_type *_type    // 指向动态类型元数据
    	hash  uint32 // copy of _type.hash. Used for type switches. 类型hash值，用于快速判断类型是否相等时使用。
    	_     [4]byte
    	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter. 方法地址
    }
    type interfacetype struct {
    	typ     _type  
    	pkgpath name       //      
    	mhdr    []imethod  // 接口方法列表
    }
```
itab结构体是可复用的，go会把用到的itab缓存起来，以接口类型和动态类型组合为key，以itab结构体指针为value，构造一个hash表。
用于存储和查询itab缓存信息。

需要一个itab时会先到这个hash表中查找，用接口类型的hash值，与动态类型的hash值做异或运算。如果有，就拿来使用，没有的话就创建itab结构体，添加到
hash表中。


# sync.Once
用来执行且仅执行一次的动作。常用于单例对象初始化。
once只有一个方法：Do（f func()）
只有第一次调用f才会被执行。

```go
    // Once is an object that will perform exactly one action.
    type Once struct {
    	// done indicates whether the action has been performed.
    	// It is first in the struct because it is used in the hot path.
    	// The hot path is inlined at every call site.
    	// Placing done first allows more compact instructions on some architectures (amd64/x86),
    	// and fewer instructions (to calculate offset) on other architectures.
    	done uint32
    	m    Mutex
    }
```
使用了 done标识是否已经执行过。
使用互斥锁保证只有一个goroutine进行初始化。利用双检查机制，保证同时到来的多个goroutine看到的值是1。

## 易出错情况
- 死锁：
    once中套once，执行两次do
    
- 未初始化
    


# waitGroup用法
主要用于并发-等待问题。
- Add（delta int）设置计数值
- Done（）    计数值减一，实际上就是调用了Add(-1)
- Wait（）  阻塞等待，直到计数器为0


## waitgroup实现

- state1 用于记录状态的数组。
- nocopy

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

Add()，为计数器增加一个delta值。       
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

Wait(),不断检查state的值，如果变为了0，那么调用者不再等待，直接返回。如果大于0，调用者加入waiter队列并阻塞。       
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

### waitgroup常见错误
- 计数器设置为负值
- done调用过多，超过了计数器的值
- 调用add之前调用了wait
- 前一个wait还没结束，就重用waitgroup。


# make 和 new的区别
make和new都是内建函数。
new只接收类型参数，分配好内存后，返回指向该类型内存地址的指针。同时会把分配的内存置为零值。
make只用于slice，map，chan的内存创建，返回的类型就是他们本身。因为他们本身就是引用类型了。


# 实现一个线程安全的map
- map的key必须是可比较的，bool、整数、浮点数、复数、字符串、指针、Channel、接口都是可比较的，包含可比较元素的 struct 和数组。
slice，map，函数值都是不可比较的。

- map是无序的。
```go
    package main
    
    import (
    	"fmt"
    	"sync"
    )
    
    type RwMap struct {
    	sync.RWMutex
    	m map[int]int
    }
    
    func newRwMap(n int) *RwMap {
    	return &RwMap{m: make(map[int]int, n)}
    }
    
    func (rm *RwMap) Get(key int) (value int, exist bool) {
    	rm.RLock()
    	defer rm.RUnlock()
    	value, exist = rm.m[key]
    	return
    }
    
    func (rm *RwMap) Set(key, value int) {
    	rm.Lock()
    	defer rm.Unlock()
    	rm.m[key] = value
    }
    
    func (rm *RwMap) Delete(key int) {
    	rm.Lock()
    	defer rm.Unlock()
    	delete(rm.m, key)
    }
    
    func (rm *RwMap) Len() int {
    	rm.RLock()
    	defer rm.RUnlock()
    	return len(rm.m)
    }
    
    func (rm *RwMap) Each(f func(k, v int) bool) {
    	rm.RLock()
    	defer rm.RUnlock()
    	for k, v := range rm.m {
    		if !f(k, v) {
    			return
    		}
    	}
    }
    
    func main() {
    	m := newRwMap(5)
    	m.Set(1, 1)
    	m.Set(2, 2)
    	m.Set(3, 3)
    	m.Set(4, 4)
    	m.Set(5, 5)
    
    	fmt.Println(m.Len())
    	t := make(map[int]int,5)
    	m.Each(func(k, v int) bool {
    		t[k] =v
    		return true
    	})
    	fmt.Println(t)
    }

```

## map常见使用错误
- 未初始化
- 并发读写

## 官方sync.Map是线程安全的。
适用场景
- 只一次写，多次读的并发场景。

# copy
从源切片复制到目标切片。
一个特殊情况是，也可以从一个字符串中拷贝字节到字节切片中。

```go
    package main
    
    import "fmt"
    
    func main() {
    	s := "hello world"
    	b := make([]byte, 11)
    	copy(b, s)
    	fmt.Println(string(b))
    }

```

# go怎么做深拷贝
## 浅拷贝
创建一个新对象，这个对象有着原始对象属性值的一份精确拷贝，如果属性是基本类型，拷贝的就是基本类型的值。         
如果属性是引用类型，拷贝的就是内存地址。如果其中一个对象改变了这个地址，就会影响另一个对象。      

## 深拷贝          
将一个对象从内存中完整的拷贝出来一份，从堆内存中开辟一个新的区域存放对象，且修改新对象不会影响原对象。         

slice扩容时就是做的深拷贝。


# 协程栈空间大小
栈空间的演变：     
v1.0 ~ v1.1 — 最小栈内存空间为 4KB；         
v1.2 — 将最小栈内存提升到了 8KB；              
v1.3 — 使用连续栈替换之前版本的分段栈；             
v1.4 — 将最小栈内存降低到了 2KB；                  

给goroutine初始堆栈大小为2Kb，随着程序运行使用而增加。  
最大值32位设置为250Mb,64位最大为1Gb，由主goroutine G0设置。        


# golang如何知道或者检测死锁
- 自测时可以启动一个goroutine，运行pprof，登录监控界面，查看goroutine的调用栈来定位分析。  
- 使用 vet工具 go vet xx.go 检查


# 如何实现只开100个协程
sync.waitGroup
或者channel


# go mod命令
go mod init xxx项目 新建go.mod文件            
go list -m all 查看当前模块及其所有依赖项。       
go list -m -versions 包名 该包的可用版本             
go.sum 文件中包含了依赖项的特殊哈希值加密，来确保这些依赖项在将来下载时，与第一次下载的一致，确保项目依赖的模块不会被恶意或非预期的修改。            

go.mod文件中的 indirect标识表示一个依赖项不被该项目直接使用是其他模块的间接依赖。                
可以为每个不同的主要版本使用不同的模块路径 如：xx/xx/v1   xx/xx/v2             
使用 go mod tidy 清除未使用的依赖项。                   


- go mod init创建一个新模块，初始化go.mod描述它的文件。                   
- go build，go test以及其他软件包构建命令go.mod根据需要添加新的依赖项。                 
- go list -m all 打印当前模块的依赖关系。                       
- go get 更改所需的依赖版本（或添加新的依赖）。                    
- go mod tidy 删除未使用的依赖项。                    


# reflect
通过反射可以获得对象的类型和对象的值。             



# proof

# gdb


# 协程交叉打印数组

# 性能问题排查

# 压力测试如何实现
go test -bench ./        

