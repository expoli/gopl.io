// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 276.

// Package memo provides a concurrency-safe memoization a function of
// a function.  Requests for different keys proceed in parallel.
// Concurrent requests for the same key block until the first completes.
// This implementation uses a Mutex.
package memo

import "sync"

// Func is the type of the function to memoize.
type Func func(string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

//!+
/*
条目里面拥有一个 result 结构体
*/
type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]*entry
}

func (memo *Memo) Get(key string) (value interface{}, err error) {
	/*
		获取互斥锁来保护共享变量cache map，
		查询map中是否存在指定条目，
		如果没有找到那么分配空间插入一个新条目，释放互斥锁。
	*/
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		// This is the first request for this key.
		// This goroutine becomes responsible for computing
		// the value and broadcasting the ready condition.
		// 无缓存的 channel 用于同步状态
		//  申请一个 token
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.mu.Unlock()

		e.res.value, e.res.err = memo.f(key)

		close(e.ready) // broadcast ready condition
	} else {
		// This is a repeat request for this key.
		memo.mu.Unlock()

		/*
			如果存在条目的话且其值没有写入完成（也就是有其它的goroutine在调用f这个慢函数）时，
			goroutine必须等待值ready之后才能读到条目的结果。

			而想知道是否ready的话，可以直接从ready channel中读取，
			由于这个读取操作在channel关闭之前一直是阻塞。
		*/
		<-e.ready // wait for ready condition
	}
	/*
		条目中的e.res.value和e.res.err变量是在多个goroutine之间共享的。
		创建条目的goroutine同时也会设置条目的值，
		其它goroutine在收到"ready"的广播消息之后立刻会去读取条目的值。

		尽管会被多个goroutine同时访问，但却并不需要互斥锁。
		ready channel的关闭一定会发生在其它goroutine接收到广播事件之前，
		因此第一个goroutine对这些变量的写操作是一定发生在这些读操作之前的。
		不会发生数据竞争。

		上面这样Memo的实现使用了一个互斥量来保护多个goroutine调用Get时的共享map变量。
	*/
	return e.res.value, e.res.err
}

//!-
