// Copyright 2018 Ryan Dahl <ry@tinyclouds.org>
// All rights reserved. MIT License.
package deno

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

func logDebug(format string, v ...interface{}) {
	// Unless the debug flag is specified, discard logs.
	if *flagDebug {
		fmt.Printf(format+"\n", v...)
	}
}

// 给定文件或目录是否存在。
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	// 如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
	if os.IsNotExist(err) {
		return false
	}
	panic(err)
}
// 断言
func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}
// 检测url是不是为真
func isRemote(filename string) bool {
	u, err := url.Parse(filename)
	check(err)
	return u.IsAbs()
}
// 检查错误
func check(e error) {
	if e != nil {
		panic(e)
	}
}
// os.Stderr 标准输出错误
func exitOnError(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
// 异步
func async(cb func()) {
	wg.Add(1)
	go func() {
		cb()
		wg.Done()
	}()
}

// 通配符
const wildcard = "[WILDCARD]"

// Matches the pattern string against the text string. The pattern can
// contain "[WILDCARD]" substrings which will match one or more characters.
// Returns true if matched.
func patternMatch(pattern string, text string) bool {
	// Empty pattern only match empty text.
	if len(pattern) == 0 {
		return len(text) == 0
	}

	if pattern == wildcard {
		return true
	}
	// 把字符串转换成数组
	parts := strings.Split(pattern, wildcard)

	if len(parts) == 1 {
		return pattern == text
	}
	// 判断字符串s是否以parts[0] 开头
	if strings.HasPrefix(text, parts[0]) {
	text = text[len(parts[0]):]
	} else {
		return false
	}

	for i := 1; i < len(parts); i++ {
		// If the last part is empty, we match.
		if i == len(parts)-1 {
			if parts[i] == "" || parts[i] == "\n" {
				return true
			}
		}
		// 判断子字符串在父字符串中第一次出现的位置
		index := strings.Index(text, parts[i])
		if index < 0 {
			return false
		}
		text = text[index+len(parts[i]):]
	}

	return len(text) == 0
}
