// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 246.

// Countdown implements the countdown for a rocket launch.
package main

// NOTE: the ticker goroutine never terminates if the launch is aborted.
// This is a "goroutine leak".

import (
	"fmt"
	"os"
	"time"
)

//!+

func main() {
	// ...create abort channel...

	//!-

	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	//!+
	fmt.Println("Commencing countdown.  Press return to abort.")
	ticker := time.NewTicker(1 * time.Second)

	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-ticker.C: // receive from the ticker's channel
			// Do nothing.
		case <-abort:
			fmt.Println("Launch aborted!")
			/*
				time.Tick函数表现得好像它创建了一个在循环中调用time.Sleep的goroutine，
				每次被唤醒时发送一个事件。

				当countdown函数返回时，它会停止从tick中接收事件，但是ticker这个goroutine还依然存活，
				继续徒劳地尝试向channel中发送值，
				然而这时候已经没有其它的goroutine会从该channel中接收值了——这被称为goroutine泄露

				Tick函数挺方便，但是只有当程序整个生命周期都需要这个时间时我们使用它才比较合适。否则的话，我们应该使用下面的这种模式：

				ticker := time.NewTicker(1 * time.Second)
				<-ticker.C    // receive from the ticker's channel
				ticker.Stop() // cause the ticker's goroutine to terminate

			*/
			ticker.Stop() // cause the ticker's goroutine to terminate
			return
		}
	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}
