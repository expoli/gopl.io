// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 136.

// The toposort program prints the nodes of a DAG in topological order.
package main

import (
	"fmt"
	"sort"
)

//!+table
// prereqs maps computer science courses to their prerequisites.
// 考虑这样一个问题：给定一些计算机课程，每个课程都有前置课程，
// 只有完成了前置课程才可以开始当前课程的学习；我们的目标是选择出一组课程，
// 这组课程必须确保按顺序学习时，能全部被完成。每个课程的前置课程如下：
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

//!-table

//!+main
// 这类问题被称作拓扑排序。
// 从概念上说，前置条件可以构成有向图。图中的顶点表示课程，边表示课程间的依赖关系。
// 显然，图中应该无环，这也就是说从某点出发的边，最终不会回到该点。
// 下面的代码用深度优先搜索了整张图，获得了符合要求的课程序列。
func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func topoSort(m map[string][]string) []string {
	var order []string
	// 保存某些课程是否已经遍历完成
	seen := make(map[string]bool)
	// 深度优先搜索的递归体函数字面量
	// 当匿名函数需要被递归调用时，我们必须首先声明一个变量（在上面的例子中，我们首先声明了 visitAll），
	// 再将匿名函数赋值给这个变量。
	// 如果不分成两步，函数字面量无法与visitAll绑定，我们也无法递归调用该匿名函数。
	var visitAll func(items []string)

	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				// 递归完成，从后向前插入对应的课程
				order = append(order, item)
			}
		}
	}
	// 找出所有的课程名称
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	// 对课程进行排序
	// 哈希表prepreqs的value是遍历顺序固定的切片，
	// 而不再是遍历顺序随机的map，所以我们对prereqs的key值进行排序，
	// 保证每次运行toposort程序，都以相同的遍历顺序遍历prereqs。
	sort.Strings(keys)
	// 开始遍历所有的课程
	visitAll(keys)
	return order
}

//!-main
