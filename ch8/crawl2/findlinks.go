// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
/*
我们可以用一个有容量限制的buffered channel来控制并发，
这类似于操作系统里的计数信号量概念。

从概念上讲，channel里的n个空槽代表n个可以处理内容的token（通行证），
从channel里接收一个值会释放其中的一个token，并且生成一个新的空槽位。
这样保证了在没有接收介入时最多有n个发送操作。

由于channel里的元素类型并不重要，我们用一个零值的struct{}来作为其元素。
*/
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	// 申请一个资源
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	// 处理完毕之后，释放一个资源
	<-tokens // release the token

	if err != nil {
		log.Print(err)
	}
	return list
}

//!-sema

//!+
func main() {
	worklist := make(chan []string)
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	// 为了使这个程序能够终止，我们需要在worklist为空或者没有crawl的goroutine在运行时退出主循环。
	n++
	/*
		防止多个输入URL同时进入导致死锁
	*/
	go func() { worklist <- os.Args[1:] }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				/*
					每一次我们发现有元素需要被发送到worklist时，我们都会对n进行++操作，
					在向worklist中发送初始的命令行参数之前，我们也进行过一次++操作。

					这里的操作++是在每启动一个crawler的goroutine之前。主
					循环会在n减为0时终止，这时候说明没活可干了。
				*/
				n++
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}

//!-
