// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"bytes"
	"encoding/json"
	"regexp"
	"runtime"
	"strings"
)

var (
	reEmpty   = regexp.MustCompile(`^\s*$`)
	reInit    = regexp.MustCompile(`init·?\d+$`) // main.init·1
	reClosure = regexp.MustCompile(`func·?\d+$`) // main.func·001
)

// name format:
// runtime.goexit
// runtime.main
// main.init
// main.main
// main.init·1 -> main.init
// main.func·001 -> main.func
// github.com/chai2010/errors.New
// ...
func callerInfo(skip int) (name, file string, line int, ok bool) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		name = "???"
		file = "???"
		line = 1
		return
	}

	name = runtime.FuncForPC(pc).Name()
	if reInit.MatchString(name) {
		name = reInit.ReplaceAllString(name, "init")
	}
	if reClosure.MatchString(name) {
		name = reClosure.ReplaceAllString(name, "func")
	}

	// Truncate file name at last file name separator.
	if idx := strings.LastIndex(file, "/"); idx >= 0 {
		file = file[idx+1:]
	} else if idx = strings.LastIndex(file, "\\"); idx >= 0 {
		file = file[idx+1:]
	}
	return
}

func jsonEncode(m interface{}) []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1) // <
	data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1) // >
	data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1) // &
	return data
}

func jsonEncodeIndent(m interface{}) []byte {
	data, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return nil
	}
	data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1) // <
	data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1) // >
	data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1) // &
	return data
}

func jsonEncodeString(m interface{}) string {
	return string(jsonEncodeIndent(m))
}

func jsonDecode(data []byte, m interface{}) error {
	return json.Unmarshal(data, m)
}
