package main

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outline: %v\n", err)
		os.Exit(1)
	}
	outline(nil, doc)
}
func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		stack = append(stack, n.Data) // push tag
		fmt.Println(stack)
	}
	// 有一点值得注意：outline有入栈操作，但没有相对应的出栈操作。
	//当outline调用自身时，被调用者接收的是stack的拷贝。被调用者
	//对stack的元素追加操作，修改的是stack的拷贝，其可能会修改slice
	//底层的数组甚至是申请一块新的内存空间进行扩容；但这个过程并不会
	//修改调用方的stack。因此当函数返回时，调用方的stack与其调用自身
	//之前完全一致。
	/*
		大部分编程语言使用固定大小的函数调用栈，常见的大小从64KB到2MB不等。
		固定大小栈会限制递归的深度，当你用递归处理大量数据时，需要避免栈溢出；
		除此之外，还会导致安全性问题。与此相反，Go语言使用可变栈，栈的大小按
		需增加（初始时很小）。这使得我们使用递归时不必考虑溢出和安全问题。
	*/
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c)
	}
}
