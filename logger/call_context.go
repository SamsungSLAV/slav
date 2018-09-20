/*
 *  Copyright (c) 2018 Samsung Electronics Co., Ltd All Rights Reserved
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License
 */

package logger

import (
	"runtime"
	"strings"
)

// CallContext defines log creation source code context.
type CallContext struct {
	Path     string `json:"path"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Package  string `json:"package"`
	Type     string `json:"type,omitempty"`
	Function string `json:"function"`
}

func getCallContext(depth int) *CallContext {
	pc, file, line, ok := runtime.Caller(depth + 1)
	if !ok {
		return nil
	}

	ret := new(CallContext)

	fi := strings.LastIndex(file, "/")
	if fi != -1 {
		ret.Path = file[:fi+1]
		ret.File = file[fi+1:]
	}

	ret.Line = line

	function := runtime.FuncForPC(pc).Name()

	si := strings.LastIndex(function, "/")
	// If there is no '/' move to index 0, otherwise move after '/'.
	si++

	fields := strings.Split(function[si:], ".")
	ret.Package = function[:si+len(fields[0])]
	if len(fields) == 3 {
		ret.Type = fields[1]
	}
	ret.Function = fields[len(fields)-1]

	return ret
}
