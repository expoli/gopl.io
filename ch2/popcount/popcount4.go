// Package popcount 练习 2.5： 表达式x&(x-1)用于将x的最低的一个非零的bit位清零。使用这个算法重写PopCount函数，然后比较性能。
package popcount

func PopCount4(x uint64) int {
	sum := 0
	for x != 0 {
		temp := x & (x - 1)
		if temp != x {
			sum += 1
		}
		x = temp
	}
	return sum
}
