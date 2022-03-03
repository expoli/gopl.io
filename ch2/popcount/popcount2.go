// Package popcount /*下面的代码定义了一个PopCount函数，用于返回一个数字中含二进制1bit的个数。
//它使用init初始化函数来生成辅助表格pc，pc表格用于处理每个8bit宽度的数字含二进制的1bit的bit个数，
//这样的话在处理64bit宽度的数字时就没有必要循环64次，只需要8次查表就可以了。（这并不是最快的统计1bit
//数目的算法，但是它可以方便演示init函数的用法，并且演示了如何预生成辅助表格，这是编程中常用的技术）。
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
func PopCount2(x uint64) int {
	sum := 0
	for i := 0; i < 8; i++ {
		sum += int(pc[byte(x>>(i*8))])
	}
	return sum
}
