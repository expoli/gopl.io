// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 191.

// Http1 is a rudimentary e-commerce server.
package main

import (
	"fmt"
	"log"
	"net/http"
)

//!+main

func main() {
	db := database{"shoes": 50, "socks": 5}
	log.Fatal(http.ListenAndServe("localhost:8000", db))
}

type dollars float32

/*
实现了 String 的接口，所以直接输出打印的时候可以有着正确的格式控制
*/
func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

/*
实现了 HTTPHandler 的接口
最前面的接受者就是 database
*/
func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

//!-main

/*
//!+handler
package http

type Handler interface {
	ServeHTTP(w ResponseWriter, r *Request)
}

func ListenAndServe(address string, h Handler) error
//!-handler
*/
