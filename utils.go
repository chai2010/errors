// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"regexp"
	"runtime"
	"strings"
)

var (
	reInit    = regexp.MustCompile(`init·\d+$`) // main.init·1
	reClosure = regexp.MustCompile(`func·\d+$`) // main.func·001
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
