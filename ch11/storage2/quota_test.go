// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

//!+test
package storage

import (
	"strings"
	"testing"
)

func TestCheckQuotaNotifiesUser(t *testing.T) {
	var notifiedUser, notifiedMsg string

	/*
		我们必须修改测试代码恢复notifyUsers原先的状态以便后续其他的测试没有影响，
		要确保所有的执行路径后都能恢复，包括测试失败或panic异常的情形。

		在这种情况下，我们建议使用defer语句来延后执行处理恢复的代码。

		这种处理模式可以用来暂时保存和恢复所有的全局变量，
		包括命令行标志参数、调试选项和优化参数；
		安装和移除导致生产代码产生一些调试信息的钩子函数；
		还有有些诱导生产代码进入某些重要状态的改变，
		比如超时、错误，甚至是一些刻意制造的并发行为等因素。

		以这种方式使用全局变量是安全的，
		因为go test命令并不会同时并发地执行多个测试。
	*/
	saved := notifyUser
	defer func() { notifyUser = saved }()

	/*
		伪邮件发送函数，覆盖 CheckQuota 中的同名函数
	*/
	notifyUser = func(user, msg string) {
		notifiedUser, notifiedMsg = user, msg
	}

	const user = "joe@example.org"
	usage[user] = 980000000 // simulate a 980MB-used condition

	CheckQuota(user)
	if notifiedUser == "" && notifiedMsg == "" {
		t.Fatalf("notifyUser not called")
	}
	if notifiedUser != user {
		t.Errorf("wrong user (%s) notified, want %s",
			notifiedUser, user)
	}
	const wantSubstring = "98% of your quota"
	if !strings.Contains(notifiedMsg, wantSubstring) {
		t.Errorf("unexpected notification message <<%s>>, "+
			"want substring %q", notifiedMsg, wantSubstring)
	}
}

//!-test

/*
//!+defer
func TestCheckQuotaNotifiesUser(t *testing.T) {
	// Save and restore original notifyUser.
	saved := notifyUser
	defer func() { notifyUser = saved }()

	// Install the test's fake notifyUser.
	var notifiedUser, notifiedMsg string
	notifyUser = func(user, msg string) {
		notifiedUser, notifiedMsg = user, msg
	}
	// ...rest of test...
}
//!-defer
*/
