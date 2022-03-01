// 练习 1.3： 做实验测量潜在低效的版本和使用了strings.Join的版本的运行时间差异。（1.6节讲解了部分time包，11.4节展示了如何写标准测试程序，以得到系统性的性能评测。）
package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func echo1(args []string) {
	s, sep := "", ""
	for _, arg := range args {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

func echo2(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func main() {
	now := time.Now()
	echo1(os.Args)
	fmt.Printf("echo1: %d\n", time.Since(now))

	now = time.Now()
	echo2(os.Args)
	fmt.Printf("echo2: %d\n", time.Since(now))
}
