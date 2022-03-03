//练习 2.4： 用移位算法重写PopCount函数，每次测试最右边的1bit，然后统计总数。比较和查表算法的性能差异。
package main

import (
	"fmt"
	"gopl.io/ch2/popcount"
	"time"
)

func main() {
	x := 1000000000
	now := time.Now()
	fmt.Printf("%d ", popcount.PopCount(uint64(x)))
	fmt.Printf("%s\n", time.Since(now))
	now = time.Now()
	fmt.Printf("%d ", popcount.PopCount2(uint64(x)))
	fmt.Printf("%s\n", time.Since(now))
	now = time.Now()
	fmt.Printf("%d ", popcount.PopCount3(uint64(x)))
	fmt.Printf("%s\n", time.Since(now))
	now = time.Now()
	fmt.Printf("%d ", popcount.PopCount4(uint64(x)))
	fmt.Printf("%s", time.Since(now))
}
