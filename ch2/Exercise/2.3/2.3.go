//练习 2.3： 重写PopCount函数，用一个循环代替单一的表达式。比较两个版本的性能。（11.4节将展示如何系统地比较两个不同实现的性能。）
package main

import (
	"fmt"
	"gopl.io/ch2/popcount"
	"time"
)

func main() {
	x := 1000000000
	now := time.Now()
	popcount.PopCount(uint64(x))
	fmt.Printf("%d s\n", int(time.Since(now)))
	now = time.Now()
	popcount.PopCount2(uint64(x))
	fmt.Printf("%d s", int(time.Since(now)))
}
