// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 287.

//!+main

// The jpeg command reads a PNG image from the standard input
// and writes it as a JPEG image to the standard output.
package main

import (
	"fmt"
	"image"
	"image/jpeg"
	/*
		要注意image/png包的匿名导入语句。如果没有这一行语句，
		程序依然可以编译和运行，但是它将不能正确识别和解码PNG格式的图像：

		每个格式驱动列表的每个入口指定了四件事情：
		格式的名称；一个用于描述这种图像数据开头部分模式的字符串，用于解码器检测识别；
		一个Decode函数用于完成解码图像工作；
		一个DecodeConfig函数用于解码图像的大小和颜色空间的信息。

		每个驱动入口是通过调用image.RegisterFormat函数注册，
		一般是在每个格式包的init初始化函数中调用，例如image/png包是这样注册的：

		package png // image/png

		func Decode(r io.Reader) (image.Image, error)
		func DecodeConfig(r io.Reader) (image.Config, error)

		func init() {
			const pngHeader = "\x89PNG\r\n\x1a\n"
			image.RegisterFormat("png", pngHeader, Decode, DecodeConfig)
		}
	*/
	_ "image/png" // register PNG decoder
	"io"
	"os"
)

func main() {
	if err := toJPEG(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
		os.Exit(1)
	}
}

func toJPEG(in io.Reader, out io.Writer) error {
	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Input format =", kind)
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}

//!-main

/*
//!+with
$ go build gopl.io/ch3/mandelbrot
$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
Input format = png
//!-with

//!+without
$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
jpeg: image: unknown format
//!-without
*/
