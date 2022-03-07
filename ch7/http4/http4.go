// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 195.

// Http4 is an e-commerce server that registers the /list and /price
// endpoint by calling http.HandleFunc.
package main

import (
	"fmt"
	"log"
	"net/http"
)

//!+main

func main() {
	db := database{"shoes": 50, "socks": 5}
	/*
		从上面的代码很容易看出应该怎么构建一个程序：
		由两个不同的web服务器监听不同的端口，并且定义不同的URL将它们指派到不同的handler。
		我们只要构建另外一个ServeMux并且再调用一次ListenAndServe（可能并行的）。
		但是在大多数程序中，一个web服务器就足够了。
		此外，在一个应用程序的多个文件中定义HTTP handler也是非常典型的，
		如果它们必须全部都显式地注册到这个应用的ServeMux实例上会比较麻烦。

		所以为了方便，net/http包提供了
		一个全局的ServeMux实例 DefaultServerMux 和包级别的 http.Handle 和 http.HandleFunc 函数。
		现在，为了使用DefaultServeMux作为服务器的主handler，
		我们不需要将它传给ListenAndServe函数；nil值就可以工作。
	*/
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//!-main

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}
