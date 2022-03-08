// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 229.

// Pipeline2 demonstrates a finite 3-stage pipeline.
package main

import "fmt"

//!+
func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// Counter
	go func() {
		for x := 0; x < 100; x++ {
			naturals <- x
		}
		/*
			数据发送完毕后，关闭 channel
		*/
		close(naturals)
	}()

	// Squarer
	go func() {
		/*
			使用 for range 方法迭代 channel 的时候，
			如果 channel 关闭之后，range 自动退出循环
		*/
		for x := range naturals {
			squares <- x * x
		}
		close(squares)
	}()

	// Printer (in main goroutine)
	/*
		其实你并不需要关闭每一个channel。
		只有当需要告诉接收者goroutine，所有的数据已经全部发送时才需要关闭channel。
		不管一个channel是否被关闭，当它没有被引用时将会被Go语言的垃圾自动回收器回收。

		（不要将关闭一个打开文件的操作和关闭一个channel操作混淆。
		对于每个打开的文件，都需要在不使用的时候调用对应的Close方法来关闭文件。）

		试图重复关闭一个channel将导致panic异常，
		试图关闭一个nil值的channel也将导致panic异常。

		关闭一个channels还会触发一个广播机制，
	*/
	for x := range squares {
		fmt.Println(x)
	}
}

//!-
