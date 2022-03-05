// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 139.

// Findlinks3 crawls the web, starting with the URLs on the command line.
package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

//!+breadthFirst
// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
// 广度优先搜索
func breadthFirst(f func(item string) []string, workList []string) {
	// 是否已经遍历过的控制数组
	seen := make(map[string]bool)
	for len(workList) > 0 {
		// 使用工作负载进行广度优先搜索
		items := workList
		// 重新创建切片（保存下一次的工作负载）
		workList = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				// append的参数“f(item)...”，会将f返回的一组元素一个个添加到worklist中。
				workList = append(workList, f(item)...)
			}
		}
	}
}

//!-breadthFirst

//!+crawl
// crawl函数会将URL输出，提取其中的新链接，并将这些新链接返回。
// 我们会将crawl作为参数传递给breadthFirst。
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
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	// 当所有发现的链接都已经被访问或电脑的内存耗尽时，程序运行结束。
	breadthFirst(crawl, os.Args[1:])
}

//!-main
