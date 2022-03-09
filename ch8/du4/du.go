// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 251.

// The du4 command computes the disk usage of the files in a directory.
package main

// The du4 variant includes cancellation:
// it terminates quickly when the user hits return.

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//!+1
var done = make(chan struct{})

/*
定义了一个工具函数，cancelled，这个函数在被调用的时候会轮询退出状态。
*/
func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

//!-1

func main() {
	// Determine the initial directories.
	roots := os.Args[1:]
	if len(roots) == 0 {
		roots = []string{"."}
	}

	//!+2
	// Cancel traversal when input is detected.
	// 取消协程，新建一个取消协程，当接收到键盘输入的时候，取消程序的运行
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		// 每当有输入被读到（比如用户按了回车键），
		// 这个goroutine就会把取消消息通过关闭done的channel广播出去。
		close(done)
	}()
	//!-2

	// Traverse each root of the file tree in parallel.
	fileSizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	// Print the results periodically.
	tick := time.Tick(500 * time.Millisecond)
	var nfiles, nbytes int64
loop:
	//!+3
	for {
		select {
		case <-done:
			// Drain fileSizes to allow existing goroutines to finish.
			// 清空 fileSize channel 保证协程能够正常退出
			for range fileSizes {
				// Do nothing.
			}
			return
		case size, ok := <-fileSizes:
			// ...
			//!-3
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes) // final totals
	/*
	   这里有一个方便的窍门我们可以一用：
	   取代掉直接从主函数返回，我们调用一个panic，然后runtime会把每一个goroutine的栈dump下来。
	   如果main goroutine是唯一一个剩下的goroutine的话，他会清理掉自己的一切资源。
	   但是如果还有其它的goroutine没有退出，他们可能没办法被正确地取消掉，
	   也有可能被取消但是取消操作会很花时间；所以这里的一个调研还是很有必要的。
	   我们用panic来获取到足够的信息来验证我们上面的判断，看看最终到底是什么样的情况。
	*/
	//panic("du exit")
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
//!+4
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	/*
		保证在处理文件夹的过程中，依旧能够感知到取消动作
		walkDir这个goroutine一启动就会轮询取消状态，
		如果取消状态被设置的话会直接返回，并且不做额外的事情。

		这样我们将所有在取消事件之后创建的goroutine改变为无操作。

		可以避免在取消事件发生时还去创建goroutine。
	*/
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		// ...
		//!-4
		if entry.IsDir() {
			// 新文件夹需要新的协程去处理
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
		//!+4
	}
}

//!-4
// 并发量限制 token 方法，channel 的缓存量为 20
var sema = make(chan struct{}, 20) // concurrency-limiting counting semaphore

// dirents returns the entries of directory dir.
//!+5
func dirents(dir string) []os.FileInfo {
	select {
	// 请求 token ，当缓存 channel 满时，此函数会堵塞，限制了并发数量
	// 此限制同时适用于 walkDir 函数
	case sema <- struct{}{}: // acquire token
	// 取消动作
	case <-done:
		return nil // cancelled
	}
	// 函数返回时释放 token
	defer func() { <-sema }() // release token

	// ...read directory...
	//!-5

	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	defer f.Close()

	// 读取文件夹下的所有文件的信息
	entries, err := f.Readdir(0) // 0 => no limit; read all entries
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		// Don't return: Readdir may return partial results.
	}
	return entries
}
