// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"fmt"

	"github.com/chai2010/errors"
)

var (
	e0 = errors.New("err:init0")
	e1 error
	e2 error
)

func init() {
	e1 = errors.New("err:init1")
}

func init() {
	e2 = errors.New("err:init2")
}

func main() {
	var e3 = errors.New("err:main3")
	var e4 = func() error {
		return errors.New("err:main4")
	}()
	var e5 = errors.Wrap(e1, "err:main5")
	var e6 = errors.Wrap(e5, "err:main6")
	var e7 = errors.Wrap(e6, "err:main7")

	fmt.Println(e0, e0.(errors.Error).Caller())
	fmt.Println(e1, e1.(errors.Error).Caller())
	fmt.Println(e2, e2.(errors.Error).Caller())
	fmt.Println(e3, e3.(errors.Error).Caller())
	fmt.Println(e4, e4.(errors.Error).Caller())
	fmt.Println(e5, e5.(errors.Error).Caller())
	fmt.Println(e6, e6.(errors.Error).Caller())

	for i, e := range e7.(errors.Error).Wraped() {
		fmt.Printf("err7: wraped(%d): %v\n", i, e)
	}

	fmt.Println("err7:", e7.(fmt.Stringer).String())
}

/*
Output:
err:init0 [{main.init hello.go 16}]
err:init1 [{main.init.1 hello.go 22} {main.init hello.go 51}]
err:init2 [{main.init.2 hello.go 26} {main.init hello.go 51}]
err:main3 [{main.main hello.go 30}]
err:main4 [{main.main.func1 hello.go 32} {main.main hello.go 33}]
err:main5 -> {err:init1} [{main.main hello.go 34}]
err:main6 -> {err:main5 -> {err:init1}} [{main.main hello.go 35}]
err7: wraped(0): err:main6 -> {err:main5 -> {err:init1}}
err7: wraped(1): err:main5 -> {err:init1}
err7: wraped(2): err:init1
err7: {
	"Code": 0,
	"Error": "err:main7 -> {err:main6 -> {err:main5 -> {err:init1}}}",
	"Caller": [
		{
			"FuncName": "main.main",
			"FileName": "hello.go",
			"FileLine": 36
		}
	],
	"Wraped": [
		{
			"Code": 0,
			"Error": "err:main6 -> {err:main5 -> {err:init1}}",
			"Caller": [
				{
					"FuncName": "main.main",
					"FileName": "hello.go",
					"FileLine": 35
				}
			]
		},
		{
			"Code": 0,
			"Error": "err:main5 -> {err:init1}",
			"Caller": [
				{
					"FuncName": "main.main",
					"FileName": "hello.go",
					"FileLine": 34
				}
			]
		},
		{
			"Code": 0,
			"Error": "err:init1",
			"Caller": [
				{
					"FuncName": "main.init.1",
					"FileName": "hello.go",
					"FileLine": 22
				},
				{
					"FuncName": "main.init",
					"FileName": "hello.go",
					"FileLine": 51
				}
			]
		}
	]
}
*/
