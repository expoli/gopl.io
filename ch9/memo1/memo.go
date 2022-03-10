// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 272.

//!+

// Package memo provides a concurrency-unsafe
// memoization of a function of type Func.
package memo

// A Memo caches the results of calling a Func.
/*
Memo实例会记录需要缓存的函数f（类型为Func），
以及缓存内容（里面是一个string到result映射的map）。
*/
type Memo struct {
	f     Func
	cache map[string]result
}

// Func is the type of the function to memoize.
type Func func(key string) (interface{}, error)

/*
每一个result都是简单的函数返回的值对儿——一个值和一个错误值。
*/
type result struct {
	value interface{}
	err   error
}

/*
新建一个缓存，返回对应的地址信息
*/
func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// NOTE: not concurrency-safe!
/*
存在数据竞争
*/
func (memo *Memo) Get(key string) (interface{}, error) {
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	return res.value, res.err
}

//!-
