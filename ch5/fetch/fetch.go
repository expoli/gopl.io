// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 148.

// Fetch saves the contents of a URL into a local file.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

//!+
// Fetch downloads the URL and returns the
// name and length of the local file.
func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}
	n, err = io.Copy(f, resp.Body)
	// Close file, but prefer error from Copy, if any.
	// 关闭文件，但更喜欢复制中的错误（如果有）。
	// 上例中，通过os.Create打开文件进行写入，
	// 在关闭文件时，我们没有对f.close采用defer机制，因为这会产生一些微妙的错误。
	// 许多文件系统，尤其是NFS，写入文件时发生的错误会被延迟到文件关闭时反馈。
	// 如果没有检查文件关闭时的反馈信息，可能会导致数据丢失，而我们还误以为写入操作成功。
	// 如果io.Copy和f.close都失败了，我们倾向于将io.Copy的错误信息反馈给调用者，
	// 因为它先于f.close发生，更有可能接近问题的本质。
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return local, n, err
}

//!-

func main() {
	for _, url := range os.Args[1:] {
		local, n, err := fetch(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch %s: %v\n", url, err)
			continue
		}
		fmt.Fprintf(os.Stderr, "%s => %s (%d bytes).\n", url, local, n)
	}
}
