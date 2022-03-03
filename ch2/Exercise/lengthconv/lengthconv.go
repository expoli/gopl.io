// Package lengthconv performs Feet and Meters conversions.
package lengthconv

import "fmt"

type Feet float64
type Meters float64

const (
	FtBase float64 = 3.2808
)

/*
英里的输出格式化处理
感觉像是对强制类型转换时的格式化处理
类似于Java中的toString方法
*/
func (f Feet) String() string { return fmt.Sprintf("%gft", f) }

/*
米的输出格式化处理
感觉像是对强制类型转换时的格式化处理
类似于Java中的toString方法
*/
func (m Meters) String() string { return fmt.Sprintf("%gm", m) }
