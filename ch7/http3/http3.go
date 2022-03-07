// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 194.

// Http3 is an e-commerce server that registers the /list and /price
// endpoints by calling (*http.ServeMux).Handle.
package main

import (
	"fmt"
	"log"
	"net/http"
)

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

//!+main

func main() {
	db := database{"shoes": 50, "socks": 5}
	/*
		net/http包提供了一个请求多路器ServeMux来简化URL和handlers的联系。
		一个ServeMux将一批http.Handler聚集到一个单一的http.Handler中。
	*/
	mux := http.NewServeMux()
	/*
		分别绑定路由至相应的 Handler 处理函数
		第一个db.list是一个方法值
		也就是说db.list的调用会援引一个接收者是db的database.list方法。
		所以db.list是一个实现了handler类似行为的函数，但是因为它没有方法（理解：该方法没有它自己的方法），
		所以它不满足http.Handler接口并且不能直接传给mux.Handle。

		语句http.HandlerFunc(db.list)是一个转换而非一个函数调用，因为http.HandlerFunc是一个类型。

		ServeHTTP方法的行为是调用了它的函数本身。
		因此HandlerFunc是一个让函数值满足一个接口的适配器，
		这里函数和这个接口仅有的方法有相同的函数签名。
		实际上，这个技巧让一个单一的类型例如database以多种方式满足http.Handler接口：
		一种通过它的list方法，一种通过它的price方法等等。
	*/
	mux.Handle("/list", http.HandlerFunc(db.list))
	mux.Handle("/price", http.HandlerFunc(db.price))
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

//!-main

/*
//!+handlerfunc
package http

type HandlerFunc func(w ResponseWriter, r *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
//!-handlerfunc
*/
