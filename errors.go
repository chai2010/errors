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
	XError  error        `json:"Error"`
	XCode   int          `json:"Code"`
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
	p, err := newFromJson(json)
	if err != nil {
		return &errorInfo{
			XCaller: Caller(1), // skip == 1
			XWraped: []error{err},
			XError:  errors.New(fmt.Sprintf("errors.NewFromJson: jsonDecode failed: %v!", err)),
		}
	}
	return p.ToError()
}

func Wrap(err error, msg string) error {
	p := &errorInfo{
		XCaller: Caller(2),
		XWraped: []error{err},
		XError:  errors.New(fmt.Sprintf("%s -> {%v}", msg, err)),
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
		XError:  errors.New(fmt.Sprintf("%s -> {%v}", fmt.Sprintf(format, args...), err)),
	}
	if e, ok := err.(Error); ok {
		p.XWraped = append(p.XWraped, e.Wraped()...)
	}
	return p
}

func WrapWithCode(code int, err error, msg string) error {
	p := &errorInfo{
		XCaller: Caller(2),
		XWraped: []error{err},
		XError:  errors.New(fmt.Sprintf("%s -> {%v}", msg, err)),
		XCode:   code,
	}
	if e, ok := err.(Error); ok {
		p.XWraped = append(p.XWraped, e.Wraped()...)
	}
	return p
}

func WrapWithCodef(code int, err error, format string, args ...interface{}) error {
	p := &errorInfo{
		XCaller: Caller(2),
		XWraped: []error{err},
		XError:  errors.New(fmt.Sprintf("%s -> {%v}", fmt.Sprintf(format, args...), err)),
		XCode:   code,
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

func (p *errorInfo) MarshalJSON() ([]byte, error) {
	return jsonEncode(newFrom(p, nil)), nil
}

func (p *errorInfo) UnmarshalJSON(data []byte) error {
	px, err := newFromJson(string(data))
	if err != nil {
		return err
	}
	*p = *px.ToErrorInfo()
	return nil
}
