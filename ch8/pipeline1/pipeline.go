// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 228.

// Pipeline1 demonstrates an infinite 3-stage pipeline.
package main

import "fmt"

//!+
func main() {
	/*
		整数序列 channel
	*/
	naturals := make(chan int)
	/*
		平方后的整数序列 channel
	*/
	squares := make(chan int)

	// Counter
	// 整数发生器的协程
	go func() {
		for x := 0; ; x++ {
			naturals <- x
		}
	}()

	// Squarer
	// 平方计算协程
	go func() {
		for {
			x := <-naturals
			squares <- x * x
		}
	}()

	// Printer (in main goroutine)
	// 主协程，进行数据打印
	for {
		fmt.Println(<-squares)
	}
}

//!-
