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
	"os"
	"sync"
)

// WriterFile is a simple wrapper for os.File opened in append mode.
// It implements Writer interface.
type WriterFile struct {
	file  *os.File
	mutex sync.Locker
}

// NewWriterFile opens given file in write mode for appending
// and returns a new WriterFile object wrapping that file.
// The file permissions can be set using perm parameter.
// It panics when opening a file is not possible.
func NewWriterFile(filePath string, perm os.FileMode) *WriterFile {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		panic(err)
	}
	return &WriterFile{
		file:  f,
		mutex: new(sync.Mutex),
	}
}

// Write appends to file. It implements Writer interface in WriterFile.
func (w *WriterFile) Write(_ Level, p []byte) (int, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.file.Write(append(p, '\n'))
}
