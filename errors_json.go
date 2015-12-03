// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"errors"
)

type errorString string

func (e errorString) Error() string { return string(e) }

type errorInfoJson struct {
	XCaller []CallerInfo     `json:"Caller,omitempty"`
	XWraped []*errorInfoJson `json:"Wraped,omitempty"`
	XError  errorString      `json:"Error"`
	XCode   int              `json:"Code"`
}

func newFrom(err error) *errorInfoJson {
	if err == nil {
		return nil
	}

	x, ok := err.(Error)
	if !ok {
		return &errorInfoJson{
			XError: errorString(x.Error()),
		}
	}

	getwraped := func(x Error) (wraped []*errorInfoJson) {
		for _, it := range x.Wraped() {
			if v := newFrom(it); v != nil {
				wraped = append(wraped, v) // call newFrom
			}
		}
		return
	}
	return &errorInfoJson{
		XCaller: x.Caller(),
		XWraped: getwraped(x),
		XError:  errorString(x.Error()),
		XCode:   x.Code(),
	}
}

func toErrorInfo(x *errorInfoJson) *errorInfo {
	if x == nil {
		return nil
	}
	if x.XError == "" && x.XCode == 0 {
		return &errorInfo{}
	}

	getwraped := func(x *errorInfoJson) (wraped []error) {
		for _, it := range x.XWraped {
			if v := toErrorInfo(it); v != nil {
				wraped = append(wraped, v) // call toErrorInfo
			}
		}
		return
	}
	return &errorInfo{
		XCaller: x.XCaller,
		XWraped: getwraped(x),
		XError:  errors.New(string(x.XError)),
		XCode:   x.XCode,
	}
}

func newFromJson(json string) (p *errorInfoJson, err error) {
	if json == "" || reEmpty.MatchString(json) {
		return nil, nil
	}
	p = new(errorInfoJson)
	if err = jsonDecode([]byte(json), p); err != nil {
		return nil, err
	}
	if p.XError == "" && p.XCode == 0 {
		return nil, nil
	}
	return
}

func (x *errorInfoJson) ToErrorInfo() *errorInfo {
	return toErrorInfo(x)
}

func (p *errorInfoJson) ToError() error {
	if p.XError == "" && p.XCode == 0 {
		return nil
	}
	return p.ToErrorInfo()
}
