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
	Code() int
	error
}

type errorInfo struct {
	XCaller []CallerInfo `json:"Caller,omitempty"`
	XWraped []error      `json:"Wraped,omitempty"`
	XError  error        `json:"Error,omitempty"`
	XCode   int          `json:"Code,omitempty"`
}

type CallerInfo struct {
	FuncName string
	FileName string
	FileLine int
}

func New(msg string) error {
	return &errorInfo{
		XCaller: Caller(2),
		XError:  errors.New(msg),
	}
}

func Newf(format string, args ...interface{}) error {
	return &errorInfo{
		XCaller: Caller(2),
		XError:  fmt.Errorf(format, args...),
	}
}

func NewWithCode(code int, msg string) error {
	return &errorInfo{
		XCaller: Caller(2),
		XError:  errors.New(msg),
		XCode:   code,
	}
}

func NewWithCodef(code int, format string, args ...interface{}) error {
	return &errorInfo{
		XCaller: Caller(2),
		XError:  fmt.Errorf(format, args...),
		XCode:   code,
	}
}

func NewFromJson(json string) error {
	panic("TODO")
	if json == "" || reEmpty.MatchString(json) {
		return nil
	}
	var err errorInfo
	if e := jsonDecode([]byte(json), &err); e != nil {
		fmt.Println(e.Error())
		return &errorInfo{
			XCaller: Caller(1), // skip == 1
			XWraped: []error{e},
			XError:  fmt.Errorf("errors.NewFromJson: jsonDecode failed: %v!", e),
		}
	}
	if err.XError == nil {
		return nil
	}
	return &err
}

func Wrap(err error, msg string) error {
	p := &errorInfo{
		XCaller: Caller(2),
		XWraped: []error{err},
		XError:  fmt.Errorf("%s -> {%v}", msg, err),
	}
	if e, ok := err.(Error); ok {
		p.XWraped = append(p.XWraped, e.Wraped()...)
	}
	return p
}

func Wrapf(err error, format string, args ...interface{}) error {
	p := &errorInfo{
		XCaller: Caller(2),
		XWraped: []error{err},
		XError:  fmt.Errorf("%s -> {%v}", fmt.Sprintf(format, args...), err),
	}
	if e, ok := err.(Error); ok {
		p.XWraped = append(p.XWraped, e.Wraped()...)
	}
	return p
}

func Caller(skip int) []CallerInfo {
	var infos []CallerInfo
	for ; ; skip++ {
		name, file, line, ok := callerInfo(skip + 1)
		if !ok {
			return infos
		}
		if strings.HasPrefix(name, "runtime.") {
			return infos
		}
		infos = append(infos, CallerInfo{
			FuncName: name,
			FileName: file,
			FileLine: line,
		})
	}
	panic("unreached!")
}

func (p *errorInfo) Caller() []CallerInfo {
	return p.XCaller
}

func (p *errorInfo) Wraped() []error {
	return p.XWraped
}

func (p *errorInfo) Error() string {
	return p.XError.Error()
}

func (p *errorInfo) Code() int {
	return p.XCode
}

func (p *errorInfo) UnmarshalJSON([]byte) error {
	panic("TODO")
}
