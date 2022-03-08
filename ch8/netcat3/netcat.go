// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"io"
	"log"
	"net"
	"os"
)

//!+
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		/*
			当用户关闭了标准输入，
			主goroutine中的mustCopy函数调用将返回，
			然后调用conn.Close()关闭读和写方向的网络连接。

			关闭网络连接中的写方向的连接将导致server程序收到一个文件（end-of-file）结束的信号。
			关闭网络连接中读方向的连接将导致后台goroutine的io.Copy函数调用返回一个“read from closed connection”（“从关闭的连接读”）类似的错误，
			因此我们临时移除了错误日志语句；

			基于channels发送消息有两个重要方面。
			首先每个消息都有一个值，但是有时候通讯的事实和发生的时刻也同样重要。
			当我们更希望强调通讯发生的时刻时，我们将它称为消息事件。

			有些消息事件并不携带额外的信息，它仅仅是用作两个goroutine之间的同步，
			这时候我们可以用struct{}空结构体作为channels元素的类型，
			虽然也可以使用bool或int类型实现同样的功能，
			done <- 1语句也比done <- struct{}{}更短。
		*/
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	/*
		将 std 复制到连接中去
	*/
	mustCopy(conn, os.Stdin)
	conn.Close()
	/*
		它在主goroutine中（译注：就是执行main函数的goroutine）将标准输入复制到server，
		因此当客户端程序关闭标准输入时，后台goroutine可能依然在工作。

		我们需要让主goroutine等待后台goroutine完成工作后再退出，
		我们使用了一个channel来同步两个goroutine：
	*/
	<-done // wait for background goroutine to finish
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
