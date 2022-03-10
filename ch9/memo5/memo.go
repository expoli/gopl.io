// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 278.

// Package memo provides a concurrency-safe non-blocking memoization
// of a function.  Requests for different keys proceed in parallel.
// Concurrent requests for the same key block until the first completes.
// This implementation uses a monitor goroutine.
package memo

//!+Func

// Func is the type of the function to memoize.
type Func func(key string) (interface{}, error)

// A result is the result of calling a Func.
type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

//!-Func

//!+get

// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	response chan<- result // the client wants a single result
}

/*
Memo 结构体，数据时 request 的 channel
*/
type Memo struct{ requests chan request }

// New returns a memoization of f.  Clients must subsequently call Close.
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

/*
上面的Get方法，会创建一个response channel，把它放进request结构中，
然后发送给monitor goroutine，然后马上又会接收它。
*/
func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

//!-get

//!+monitor
/*
cache变量被限制在了monitor goroutine `(*Memo).server中，下面会看到。
monitor会在循环中一直读取请求，直到request channel被Close方法关闭。

每一个请求都会去查询cache，如果没有找到条目的话，那么就会创建/插入一个新的条目。
*/
func (memo *Memo) server(f Func) {
	// 将 cache 限制在一个协程里面
	cache := make(map[string]*entry)
	// 遍历 requests 的 channel
	for req := range memo.requests {
		// 尝试获取对应的 cache
		e := cache[req.key]
		if e == nil {
			// This is the first request for this key.
			// 获取 ready 的 token
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			/*
				对call和deliver方法的调用必须让它们在自己的goroutine中进行
				以确保monitor goroutines不会因此而被阻塞住而没法处理新的请求。
			*/
			go e.call(f, req.key) // call f(key)
		}
		go e.deliver(req.response)
	}
}

/*
entry 的 call 方法
*/
func (e *entry) call(f Func, key string) {
	// Evaluate the function.
	e.res.value, e.res.err = f(key)
	// Broadcast the ready condition.
	// 归还 token
	close(e.ready)
}

/*
entry 的 deliver 方法，用于同步操作
1. 等待 ready 信号的完成
2. 将数据写入到 response channel
*/
func (e *entry) deliver(response chan<- result) {
	// Wait for the ready condition.
	<-e.ready
	// Send the result to the client.
	response <- e.res
}

//!-monitor
