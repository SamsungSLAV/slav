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

// WriterStderr is a simple writer printing logs to standard error output (StdErr).
// It synchronizes writes with mutex.
// It implements Writer interface.
type WriterStderr struct {
	mutex sync.Locker
}

// NewWriterStderr creates a new WriterStderr object.
func NewWriterStderr() *WriterStderr {
	return &WriterStderr{
		mutex: new(sync.Mutex),
	}
}

// Write writes to StdErr. It implements Writer interface in WriterStderr.
func (w *WriterStderr) Write(_ Level, p []byte) (int, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return os.Stderr.Write(append(p, '\n'))
}
