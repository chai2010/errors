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

var (
	_ Error        = (*_Error)(nil)
	_ fmt.Stringer = (*_Error)(nil)
)

type Error interface {
	Caller() []CallerInfo
	Wraped() []error
	Code() int
	error

	private()
}

type _Error struct {
	XCode   int          `json:"Code"`
	XError  error        `json:"Error"`
	XCaller []CallerInfo `json:"Caller,omitempty"`
	XWraped []error      `json:"Wraped,omitempty"`
}

type CallerInfo struct {
	FuncName string
	FileName string
	FileLine int
}

func New(msg string) error {
	return &_Error{
		XCaller: Caller(2),
		XError:  errors.New(msg),
	}
}

func NewFrom(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(Error); ok {
		return e
	}
	return &_Error{
		XCaller: Caller(2),
		XError:  err,
	}
}

func Newf(format string, args ...interface{}) error {
	return &_Error{
		XCaller: Caller(2),
		XError:  fmt.Errorf(format, args...),
	}
}

func NewWithCode(code int, msg string) error {
	return &_Error{
		XCaller: Caller(2),
		XError:  errors.New(msg),
		XCode:   code,
	}
}

func NewWithCodef(code int, format string, args ...interface{}) error {
	return &_Error{
		XCaller: Caller(2),
		XError:  fmt.Errorf(format, args...),
		XCode:   code,
	}
}

func MustFromJson(json string) error {
	p, err := newErrorStructFromJson(json)
	if err != nil {
		panic(err)
	}
	return p.ToStdError()
}

func FromJson(json string) (Error, error) {
	p, err := newErrorStructFromJson(json)
	if err != nil {
		return nil, &_Error{
			XCaller: Caller(1), // skip == 1
			XWraped: []error{err},
			XError:  errors.New(fmt.Sprintf("errors.FromJson: jsonDecode failed: %v!", err)),
		}
	}

	return p.ToErrorInterface(), nil
}

func ToJson(err error) string {
	if p, ok := (err).(*_Error); ok {
		return p.String()
	}
	p := &_Error{XError: err}
	return p.String()
}

func Wrap(err error, msg string) error {
	p := &_Error{
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
	p := &_Error{
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
	p := &_Error{
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
	p := &_Error{
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

func (p *_Error) Caller() []CallerInfo {
	return p.XCaller
}

func (p *_Error) Wraped() []error {
	return p.XWraped
}

func (p *_Error) Error() string {
	return p.XError.Error()
}

func (p *_Error) Code() int {
	return p.XCode
}

func (p *_Error) String() string {
	return jsonEncodeString(p)
}

func (p *_Error) MarshalJSON() ([]byte, error) {
	return jsonEncode(newErrorStruct(p)), nil
}

func (p *_Error) UnmarshalJSON(data []byte) error {
	px, err := newErrorStructFromJson(string(data))
	if err != nil {
		return err
	}
	if px != nil {
		*p = *px.ToError()
	} else {
		*p = _Error{}
	}
	return nil
}

func (p *_Error) private() {
	panic("unreached!")
}
