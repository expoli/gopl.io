// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 218.

// Spinner displays an animation while computing the 45th Fibonacci number.
package main

import (
	"fmt"
	"time"
)

//!+
/*
下面的例子，main goroutine将计算菲波那契数列的第45个元素值。
由于计算函数使用低效的递归，所以会运行相当长时间，
在此期间我们想让用户看到一个可见的标识来表明程序依然在正常运行，
所以来做一个动画的小图标：
*/
func main() {
	go spinner(100 * time.Millisecond)
	const n = 45
	fibN := fib(n) // slow
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}

/*
输出一个动画的小图标，表示程序正在运行
*/
func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

/*
斐波纳契函数的递归求解
*/
func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}

//!-
