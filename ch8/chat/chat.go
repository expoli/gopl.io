// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

//!+broadcaster
type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages 客户端数据输入 channel
)

func broadcaster() {
	/*
		他的内部变量clients会记录当前建立连接的客户端集合。
		其记录的内容是每一个客户端的消息发出channel的“资格”信息。
	*/
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			/*
				broadcaster也会监听全局的消息channel，
				所有的客户端都会向这个channel中发送消息。

				当broadcaster接收到什么消息时，
				就会将其广播至所有连接到服务端的客户端。
			*/
			for cli := range clients {
				cli <- msg
			}

			/*
				broadcaster监听来自全局的entering和leaving的channel来获知客户端的到来和离开事件。
				当其接收到其中的一个事件时，会更新clients集合，
				当该事件是离开行为时，它会关闭客户端的消息发送channel。
			*/
		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	/*
		handleConn为每一个客户端创建了一个clientWriter的goroutine，
		用来接收向客户端发送消息的channel中的广播消息，
		并将它们写入到客户端的网络连接。

		客户端的读取循环会在broadcaster接收到leaving通知并关闭了channel后终止。
	*/
	go clientWriter(conn, ch)
	/*
	   远程连接信息
	*/
	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	// 广播信息
	messages <- who + " has arrived"
	// 通过entering channel来通知客户端的到来。
	// 更新 clients map 数组
	entering <- ch
	// 广播消息
	input := bufio.NewScanner(conn)
	//  阻塞循环等待输入
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()
	// 连接断开
	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		/*
			新连接接入
		*/
		go handleConn(conn)
	}
}

//!-main
