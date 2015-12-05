// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"reflect"
	"testing"
)

func TestError(t *testing.T) {
	//
}

func TestError_callar(t *testing.T) {
	// TODO
}
func TestError_wraped(t *testing.T) {
	// TODO
}

func TestError_json(t *testing.T) {
	err0 := New("err0")
	err1 := Wrap(err0, "err1")
	err2 := func() error { return Wrap(err1, "err2") }()
	err3 := func() error { return Wrap(err2, "err3") }()

	errx := NewFromJson(string(jsonEncode(err3)))
	if !reflect.DeepEqual(errx, err3) {
		t.Logf("errx: %s\n", jsonEncodeString(errx))
		t.Logf("err3: %s\n", jsonEncodeString(err3))
		t.Fatal(errx, "!=", err3)
	}
}
