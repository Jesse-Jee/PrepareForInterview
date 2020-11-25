package Practice

import (
	"fmt"
	"time"
)

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

// 解决办法： ch容量设为1
