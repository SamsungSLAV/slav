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

// Package logger provides logging mechanism for SLAV projects.
package logger

import (
	"sync"
)

const (
	// DefaultThreshold is the default level of each newly created Logger.
	DefaultThreshold = InfoLevel
)

// Logger defines type for a single logger instance.
type Logger struct {
	// threshold defines filter level for entries.
	// Only entries with level equal or less than threshold will be logged.
	// The default threshold is set to InfoLevel.
	threshold Level

	// mutex protects Logger structure from concurrent access.
	mutex *sync.Mutex
}

// NewLogger creates a new Logger instance with default configuration.
// Default level threshold is set to InfoLevel.
func NewLogger() *Logger {
	return &Logger{
		threshold: DefaultThreshold,
		mutex:     new(sync.Mutex),
	}
}

// SetThreshold defines Logger's filter level.
// Only entries with level equal or less than threshold will be logged.
func (l *Logger) SetThreshold(level Level) error {
	if !level.IsValid() {
		return ErrInvalidLogLevel
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.threshold = level
	return nil
}

// Threshold returns current Logger's filter level.
func (l *Logger) Threshold() Level {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.threshold
}

// PassThreshold verifies if message with given level passes threshold and should be logged.
func (l *Logger) PassThreshold(level Level) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return (level <= l.threshold)
}
