// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 113.

// Issuesreport prints a report of issues matching the search terms.
package main

import (
	"gopl.io/ch4/github"
	"log"
	"os"
	"text/template"
	"time"
)

//!+template
const templ = `{{.TotalCount}} issues:
{{range .Items}}----------------------------------------
Number: {{.Number}}
User:   {{.User.Login}}
Title:  {{.Title | printf "%.64s"}}
Age:    {{.CreatedAt | daysAgo}} days
{{end}}`

//!-template

//!+daysAgo
func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

//!-daysAgo

//!+exec
/*
生成模板的输出需要两个处理步骤。
	第一步是要分析模板并转为内部表示，然后基于指定的输入执行模板。分
	析模板部分一般只需要执行一次。
注意方法调用链的顺序：
	template.New先创建并返回一个模板；
	Funcs方法将daysAgo等自定义函数注册到模板中，并返回模板；
	最后调用Parse函数分析模板。
因为模板通常在编译时就测试好了，如果模板解析失败将是一个致命的错误。
	template.Must辅助函数可以简化这个致命错误的处理：
	它接受一个模板和一个error类型的参数，
	检测error是否为nil（如果不是nil则发出panic异常），
	然后返回传入的模板。
*/
var report = template.Must(template.New("issuelist").
	Funcs(template.FuncMap{"daysAgo": daysAgo}).
	Parse(templ))

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	/*
		一旦模板已经创建、注册了daysAgo函数、并通过分析和检测，
		我们就可以使用github.IssuesSearchResult作为输入源、
		os.Stdout作为输出源来执行模板：
	*/
	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}

//!-exec

func noMust() {
	//!+parse
	report, err := template.New("report").
		Funcs(template.FuncMap{"daysAgo": daysAgo}).
		Parse(templ)
	if err != nil {
		log.Fatal(err)
	}
	//!-parse
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}

/*
//!+output
$ go build gopl.io/ch4/issuesreport
$ ./issuesreport repo:golang/go is:open json decoder
13 issues:
----------------------------------------
Number: 5680
User:   eaigner
Title:  encoding/json: set key converter on en/decoder
Age:    750 days
----------------------------------------
Number: 6050
User:   gopherbot
Title:  encoding/json: provide tokenizer
Age:    695 days
----------------------------------------
...
//!-output
*/
