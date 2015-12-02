// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errors implements functions to manipulate errors.
package errors // import "github.com/chai2010/errors"

import (
	"errors"
	"fmt"
	"strings"
)

type Error interface {
	Caller() []CallerInfo
	Wraped() []error
	error
}

type errorInfo struct {
	_Caller []CallerInfo
	_Wraped []error
	_Error  error
}

type CallerInfo struct {
	FuncName string
	FileName string
	FileLine int
}

func New(msg string) error {
	return &errorInfo{
		_Caller: Caller(2),
		_Error:  errors.New(msg),
	}
}

func Newf(format string, args ...interface{}) error {
	return &errorInfo{
		_Caller: Caller(2),
		_Error:  fmt.Errorf(format, args...),
	}
}

func Wrap(err error, msg string) error {
	p := &errorInfo{
		_Caller: Caller(2),
		_Wraped: []error{err},
		_Error:  fmt.Errorf("%s -> {%v}", msg, err),
	}
	if e, ok := err.(Error); ok {
		p._Wraped = append(p._Wraped, e.Wraped()...)
	}
	return p
}

func Wrapf(err error, format string, args ...interface{}) error {
	p := &errorInfo{
		_Caller: Caller(2),
		_Wraped: []error{err},
		_Error:  fmt.Errorf("%s -> {%v}", fmt.Sprintf(format, args...), err),
	}
	if e, ok := err.(Error); ok {
		p._Wraped = append(p._Wraped, e.Wraped()...)
	}
	return p
}

func Caller(skip int) (infos []CallerInfo) {
	for ; ; skip++ {
		name, file, line, ok := callerInfo(skip + 1)
		if !ok {
			return
		}
		if strings.HasPrefix(name, "runtime.") {
			return
		}
		infos = append(infos, CallerInfo{
			FuncName: name,
			FileName: file,
			FileLine: line,
		})
	}
	return
}

func (p *errorInfo) Caller() []CallerInfo {
	return p._Caller
}

func (p *errorInfo) Wraped() []error {
	return p._Wraped
}

func (p *errorInfo) Error() string {
	return p._Error.Error()
}
