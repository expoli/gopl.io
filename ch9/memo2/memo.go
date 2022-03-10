// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 275.

// Package memo provides a concurrency-safe memoization a function of
// type Func.  Concurrent requests are serialized by a Mutex.
package memo

import "sync"

// Func is the type of the function to memoize.
type Func func(string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

//!+

type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache // 缓存同步锁
	cache map[string]result
}

// Get is concurrency-safe.
/*
不幸的是对于Memo的这一点改变使我们完全丧失了并发的性能优点。
每次对f的调用期间都会持有锁，Get将本来可以并行运行的I/O操作串行化了。

我们本章的目的是完成一个无锁缓存，而不是现在这样的将所有请求串行化的函数的缓存。
*/
func (memo *Memo) Get(key string) (value interface{}, err error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	memo.mu.Unlock()
	return res.value, res.err
}

//!-
