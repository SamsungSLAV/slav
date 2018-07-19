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

// newEntry creates a new log entry.
func (l *Logger) newEntry() *Entry {
	return &Entry{
		Logger: l,
	}
}

// Log builds log message and logs entry.
func (l *Logger) Log(level Level, args ...interface{}) {
	l.newEntry().Log(level, args...)
}

// Logf builds formatted log message and logs entry.
func (l *Logger) Logf(level Level, format string, args ...interface{}) {
	l.newEntry().Logf(level, format, args...)
}

// Emergency logs emergency level message.
func (l *Logger) Emergency(args ...interface{}) {
	l.newEntry().Emergency(args...)
}

// Alert logs alert level message.
func (l *Logger) Alert(args ...interface{}) {
	l.newEntry().Alert(args...)
}

// Critical logs critical level message.
func (l *Logger) Critical(args ...interface{}) {
	l.newEntry().Critical(args...)
}

// Error logs error level message.
func (l *Logger) Error(args ...interface{}) {
	l.newEntry().Error(args...)
}

// Warning logs warning level message.
func (l *Logger) Warning(args ...interface{}) {
	l.newEntry().Warning(args...)
}

// Notice logs notice level message.
func (l *Logger) Notice(args ...interface{}) {
	l.newEntry().Notice(args...)
}

// Info logs info level message.
func (l *Logger) Info(args ...interface{}) {
	l.newEntry().Info(args...)
}

// Debug logs debug level message.
func (l *Logger) Debug(args ...interface{}) {
	l.newEntry().Debug(args...)
}

// Emergencyf logs emergency level formatted message.
func (l *Logger) Emergencyf(format string, args ...interface{}) {
	l.newEntry().Emergencyf(format, args...)
}

// Alertf logs alert level formatted message.
func (l *Logger) Alertf(format string, args ...interface{}) {
	l.newEntry().Alertf(format, args...)
}

// Criticalf logs critical level formatted message.
func (l *Logger) Criticalf(format string, args ...interface{}) {
	l.newEntry().Criticalf(format, args...)
}

// Errorf logs error level formatted message.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.newEntry().Errorf(format, args...)
}

// Warningf logs warning level formatted message.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.newEntry().Warningf(format, args...)
}

// Noticef logs notice level formatted message.
func (l *Logger) Noticef(format string, args ...interface{}) {
	l.newEntry().Noticef(format, args...)
}

// Infof logs info level formatted message.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.newEntry().Infof(format, args...)
}

// Debugf logs debug level formatted message.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.newEntry().Debugf(format, args...)
}
