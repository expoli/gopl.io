// Package popcount /*用移位算法重写PopCount函数，每次测试最右边的1bit，然后统计总数。比较和查表算法的性能差异。
package popcount

/*
生成了和数字大小相关的对应的数值中的二进制中1的个数
*/
func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount2 returns the population count (number of set bits) of x.
func PopCount3(x uint64) int {
	sum := 0
	for i := 0; i < 64; i++ {
		temp := (x >> i) & 1
		if temp != 0 {
			sum += 1
		}
	}
	return sum
}
