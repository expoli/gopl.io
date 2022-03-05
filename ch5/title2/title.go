// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 145.

// Title2 prints the title of an HTML document specified by a URL.
// It uses defer to simplify closing the response body stream.
package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

//!+
func title(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	// 你只需要在调用普通函数或方法前加上关键字defer，就完成了defer所需要的语法。
	// 当执行到该条语句时，函数和参数表达式得到计算，但直到包含该defer语句的函数执行完毕时，
	// defer后的函数才会被执行，不论包含defer语句的函数是通过return正常结束，还是由于panic导致的异常结束。
	// 你可以在一个函数中执行多条defer语句，它们的执行顺序与声明顺序相反。

	// defer语句经常被用于处理成对的操作，如打开、关闭、连接、断开连接、加锁、释放锁。
	// 通过defer机制，不论函数逻辑多复杂，都能保证在任何执行路径下，资源被释放。
	// 释放资源的defer应该直接跟在请求资源的语句后。
	// 在下面的代码中，一条defer语句替代了之前的所有resp.Body.Close
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		return fmt.Errorf("%s has type %s, not text/html", url, ct)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	// ...print doc's title element...
	//!-
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" &&
			n.FirstChild != nil {
			fmt.Println(n.FirstChild.Data)
		}
	}
	forEachNode(doc, visitNode, nil)
	//!+

	return nil
}

//!-

func main() {
	for _, arg := range os.Args[1:] {
		if err := title(arg); err != nil {
			fmt.Fprintf(os.Stderr, "title: %v\n", err)
		}
	}
}
