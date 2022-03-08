// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 240.

// Crawl1 crawls web links starting with the command-line arguments.
//
// This version quickly exhausts available file descriptors
// due to excessive concurrent calls to links.Extract.
//
// Also, it never terminates because the worklist is never closed.
/*
用bfs(广度优先)算法来抓取整个网站。
在本节中，我们会让这个爬虫并行化，这样每一个彼此独立的抓取命令可以并行进行IO，最大化利用网络资源。
*/
package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

//!+crawl
func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

//!-crawl

//!+main
func main() {
	worklist := make(chan []string)

	// Start with the command-line arguments.
	// 第一个要进行处理的链接
	// 将命令行参数传入worklist也是在一个另外的goroutine中进行的，
	// 这是为了避免channel两端的main goroutine与crawler goroutine都尝试向对方发送内容，
	// 却没有一端接收内容时发生死锁。(存在多个链接参数的时候，会出现这个问题）
	// 当然，这里我们也可以用buffered channel来解决问题，这里不再赘述。
	go func() { worklist <- os.Args[1:] }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				// 使用协程去并行获取所有的链接
				go func(link string) {
					// 将爬取到的新链接加入到自己的工作负载中去
					worklist <- crawl(link)
					// 注意这里的crawl所在的goroutine会将link作为一个显式的参数传入，来避免“循环变量快照”的问题（在5.6.1中有讲解）。
				}(link)
			}
		}
	}
}

//!-main

/*
//!+output
$ go build gopl.io/ch8/crawl1
$ ./crawl1 http://gopl.io/
http://gopl.io/
https://golang.org/help/

https://golang.org/doc/
https://golang.org/blog/
...
2015/07/15 18:22:12 Get ...: dial tcp: lookup blog.golang.org: no such host
2015/07/15 18:22:12 Get ...: dial tcp 23.21.222.120:443: socket:
                                                        too many open files
...
//!-output
*/
