// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 244.

// Countdown implements the countdown for a rocket launch.
package main

import (
	"fmt"
	"os"
	"time"
)

//!+

func main() {
	// ...create abort channel...

	//!-

	//!+abort
	/*
		我们启动一个goroutine，
		这个goroutine会尝试从标准输入中读入一个单独的byte并且，
		如果成功了，会向名为abort的channel发送一个值。
	*/
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()
	//!-abort

	//!+
	fmt.Println("Commencing countdown.  Press return to abort.")
	/*
		。我们无法做到从每一个channel中接收信息，
		如果我们这么做的话，如果第一个channel中没有事件发过来那么程序就会立刻被阻塞，
		这样我们就无法收到第二个channel中发过来的事件。

		这时候我们需要多路复用（multiplex）这些操作了，
		为了能够多路复用，我们使用了select语句。
	*/
	select {
	// time.After函数会立即返回一个channel，并起一个新的goroutine在经过特定的时间后向该channel发送一个独立的值。
	case <-time.After(10 * time.Second):
		// Do nothing.
	case <-abort:
		fmt.Println("Launch aborted!")
		return
	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}
