// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"errors"
)

var (
	_ error = _ErrorString("")
)

type _ErrorString string

func (e _ErrorString) Error() string { return string(e) }

type _ErrorStruct struct {
	XCode   int             `json:"Code"`
	XError  _ErrorString    `json:"Error"`
	XCaller []CallerInfo    `json:"Caller,omitempty"`
	XWraped []*_ErrorStruct `json:"Wraped,omitempty"`
}

func newErrorStruct(err error) *_ErrorStruct {
	if err == nil {
		return nil
	}

	x, ok := err.(Error)
	if !ok {
		return &_ErrorStruct{
			XError: _ErrorString(x.Error()),
		}
	}
	p := &_ErrorStruct{
		XCode:   x.Code(),
		XError:  _ErrorString(x.Error()),
		XCaller: x.Caller(),
	}
	for _, it := range x.Wraped() {
		if it, ok := it.(Error); ok {
			p.XWraped = append(p.XWraped, &_ErrorStruct{
				XCode:   it.Code(),
				XError:  _ErrorString(it.Error()),
				XCaller: it.Caller(),
			})
		} else {
			p.XWraped = append(p.XWraped, &_ErrorStruct{
				XError: _ErrorString(it.Error()),
			})
		}
	}
	return p
}

func newErrorStructFromJson(json string) (p *_ErrorStruct, err error) {
	if json == "" || reEmpty.MatchString(json) {
		return nil, nil
	}
	p = new(_ErrorStruct)
	if err = jsonDecode([]byte(json), p); err != nil {
		return nil, err
	}
	if p.XError == "" && p.XCode == 0 {
		return nil, nil
	}
	return
}

func (p *_ErrorStruct) ToError() *_Error {
	if p == nil {
		return nil
	}

	if p.XError == "" && p.XCode == 0 {
		return &_Error{}
	}

	errx := &_Error{
		XCode:   p.XCode,
		XError:  errors.New(string(p.XError)),
		XCaller: p.XCaller,
	}
	for i := len(p.XWraped) - 1; i >= 0; i-- {
		if p.XWraped[i].XError == "" && p.XWraped[i].XCode == 0 {
			continue
		}
		if len(errx.XWraped) == 0 {
			if p.XWraped[i].XCode == 0 && len(p.XWraped[i].XCaller) == 0 {
				errx.XWraped = []error{errors.New(string(p.XWraped[i].XError))}
				continue
			}
		}
		errx.XWraped = append(
			[]error{&_Error{
				XCode:   p.XWraped[i].XCode,
				XError:  errors.New(string(p.XWraped[i].XError)),
				XCaller: p.XWraped[i].XCaller,
				XWraped: errx.XWraped,
			}},
			errx.XWraped...,
		)
	}
	return errx
}

func (p *_ErrorStruct) ToErrorInterface() Error {
	if x := p.ToError(); x != nil {
		return x
	}
	return nil
}

func (p *_ErrorStruct) ToStdError() error {
	if p.XError == "" && p.XCode == 0 {
		return nil
	}
	return p.ToErrorInterface()
}
