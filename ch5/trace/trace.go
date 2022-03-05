// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 146.

// The trace program uses defer to add entry/exit diagnostics to a function.
package main

import (
	"log"
	"time"
)

//!+main
func bigSlowOperation() {
	// 程序退出的时候执行
	// bigSlowOperation被调时，trace会返回一个函数值，
	// 该函数值会在bigSlowOperation退出时被调用。
	// 需要注意一点：
	// 		不要忘记defer语句后的圆括号，
	// 		否则本该在进入时执行的操作会在退出时执行，
	// 		而本该在退出时执行的，永远不会被执行。
	defer trace("bigSlowOperation")() // don't forget the extra parentheses
	// ...lots of work...
	time.Sleep(10 * time.Second) // simulate slow operation by sleeping
}

func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", msg)
	return func() { log.Printf("exit %s (%s)", msg, time.Since(start)) }
}

//!-main

func main() {
	bigSlowOperation()
}

/*
!+output
$ go build gopl.io/ch5/trace
$ ./trace
2015/11/18 09:53:26 enter bigSlowOperation
2015/11/18 09:53:36 exit bigSlowOperation (10.000589217s)
!-output
*/
