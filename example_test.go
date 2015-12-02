// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"fmt"

	"github.com/chai2010/errors"
)

func Example() {
	err0 := errors.New("err0")
	err1 := errors.Wrap(err0, "err1")
	err2 := func() error { return errors.Wrap(err1, "err2") }()
	err3 := func() error { return errors.Wrap(err2, "err3") }()

	fmt.Println(err0)
	fmt.Println(err1)
	fmt.Println(err2)
	fmt.Println(err3)
	// Output:
	// err0
	// err1 -> {err0}
	// err2 -> {err1 -> {err0}}
	// err3 -> {err2 -> {err1 -> {err0}}}
}

func Example_caller() {
	err := errors.New("error message")

	fmt.Println(err)
	for i, x := range err.(errors.Error).Caller() {
		fmt.Printf("caller:%d: %s\n", i, x.FuncName)
	}
	// Output:
	// error message
	// caller:0: github.com/chai2010/errors_test.Example_caller
	// caller:1: testing.runExample
	// caller:2: testing.RunExamples
	// caller:3: testing.(*M).Run
	// caller:4: main.main
}

func Example_wraped() {
	err0 := errors.New("err0")
	err1 := errors.Wrap(err0, "err1")
	err2 := func() error { return errors.Wrap(err1, "err2") }()
	err3 := func() error { return errors.Wrap(err2, "err3") }()

	fmt.Println(err3)
	for j, x := range err3.(errors.Error).Caller() {
		fmt.Printf("caller:%d: %s\n", j, x.FuncName)
	}
	for i, err := range err3.(errors.Error).Wraped() {
		fmt.Printf("wraped:%d: %v\n", i, err)
		for j, x := range err.(errors.Error).Caller() {
			fmt.Printf("    caller:%d: %s\n", j, x.FuncName)
		}
	}
	// Output:
	// err3 -> {err2 -> {err1 -> {err0}}}
	// caller:0: github.com/chai2010/errors_test.Example_wraped.func2
	// caller:1: github.com/chai2010/errors_test.Example_wraped
	// caller:2: testing.runExample
	// caller:3: testing.RunExamples
	// caller:4: testing.(*M).Run
	// caller:5: main.main
	// wraped:0: err2 -> {err1 -> {err0}}
	//     caller:0: github.com/chai2010/errors_test.Example_wraped.func1
	//     caller:1: github.com/chai2010/errors_test.Example_wraped
	//     caller:2: testing.runExample
	//     caller:3: testing.RunExamples
	//     caller:4: testing.(*M).Run
	//     caller:5: main.main
	// wraped:1: err1 -> {err0}
	//     caller:0: github.com/chai2010/errors_test.Example_wraped
	//     caller:1: testing.runExample
	//     caller:2: testing.RunExamples
	//     caller:3: testing.(*M).Run
	//     caller:4: main.main
	// wraped:2: err0
	//     caller:0: github.com/chai2010/errors_test.Example_wraped
	//     caller:1: testing.runExample
	//     caller:2: testing.RunExamples
	//     caller:3: testing.(*M).Run
	//     caller:4: main.main
}
