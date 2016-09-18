// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"reflect"
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	err0 := New("err0 message")
	tAssert(t, "err0 message" == err0.Error())
	tAssert(t, len(err0.(Error).Wraped()) == 0)

	err1 := Wrap(err0, "err1 message")
	tAssert(t, strings.HasPrefix(err1.Error(), "err1 message"))
	tAssert(t, strings.HasSuffix(err1.Error(), "{"+err0.Error()+"}"))
	tAssert(t, strings.Contains(err1.Error(), err0.Error()))
	tAssert(t, len(err1.(Error).Wraped()) == 1)
	tAssert(t, err1.(Error).Wraped()[0] == err0)
}

func TestError_callar(t *testing.T) {
	err := New("error message")
	callers := err.(Error).Caller()
	tAssert(t, len(callers) == 2)
	tAssert(t, callers[0].FuncName == "github.com/chai2010/errors.TestError_callar")
	tAssert(t, callers[len(callers)-2].FuncName == "github.com/chai2010/errors.TestError_callar")
	tAssert(t, callers[len(callers)-1].FuncName == "testing.tRunner")
}

func TestError_wraped(t *testing.T) {
	err0 := New("err0")
	err1 := Wrap(err0, "err1")
	err2 := func() error { return Wrap(err1, "err2") }()
	err3 := func() error { return Wrap(err2, "err3") }()

	tAssert(t, len(err0.(Error).Wraped()) == 0)
	tAssert(t, len(err1.(Error).Wraped()) == 1)
	tAssert(t, len(err2.(Error).Wraped()) == 2)
	tAssert(t, len(err3.(Error).Wraped()) == 3)

	tAssert(t, err1.(Error).Wraped()[0] == err0)

	tAssert(t, err2.(Error).Wraped()[0] == err1)
	tAssert(t, err2.(Error).Wraped()[1] == err0)

	tAssert(t, err3.(Error).Wraped()[0] == err2)
	tAssert(t, err3.(Error).Wraped()[1] == err1)
	tAssert(t, err3.(Error).Wraped()[2] == err0)
}

func TestError_json(t *testing.T) {
	err0 := New("err0")
	err1 := Wrap(err0, "err1")
	err2 := func() error { return Wrap(err1, "err2") }()
	err3 := func() error { return Wrap(err2, "err3") }()

	errx := MustFromJson(string(jsonEncode(err3)))
	if !reflect.DeepEqual(errx, err3) {
		t.Logf("errx: %s\n", jsonEncodeString(errx))
		t.Logf("err3: %s\n", jsonEncodeString(err3))
		t.Fatal(errx, "!=", err3)
	}
}

func TestCaller(t *testing.T) {
	skip0Caller := Caller(0)
	tAssert(t, len(skip0Caller) >= 2)
	tAssert(t, skip0Caller[0].FuncName == "github.com/chai2010/errors.Caller")
	tAssert(t, skip0Caller[1].FuncName == "github.com/chai2010/errors.TestCaller")

	skip1Caller := Caller(1)
	tAssert(t, len(skip1Caller) >= 2)
	tAssert(t, skip1Caller[0].FuncName == "github.com/chai2010/errors.TestCaller")
	tAssert(t, skip1Caller[1].FuncName == "testing.tRunner")
}
