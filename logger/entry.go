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
	"fmt"
	"time"
)

// Entry defines a single log message entity.
type Entry struct {
	// Logger points to instance that manages this Entry.
	Logger *Logger
	// Level defines importance of log message.
	Level Level
	// Message contains actual log message.
	Message string
	// Properties hold key-value pairs of log message properties.
	Properties Properties
	// Timestamp stores point in time of log message creation.
	Timestamp time.Time
	// CallContext stores the source code context of log creation.
	CallContext *CallContext
	// depth is a call depth of a stack frame of a caller
	depth int
}

// IncDepth increases depth of an Entry for call stack frame calculation.
func (e *Entry) IncDepth(dep int) *Entry {
	e.depth += dep
	return e
}

// Log builds log message and logs entry.
func (e *Entry) Log(level Level, args ...interface{}) {
	e.IncDepth(1).process(level, fmt.Sprint(args...))
}

// Logf builds formatted log message and logs entry.
func (e *Entry) Logf(level Level, format string, args ...interface{}) {
	e.IncDepth(1).process(level, fmt.Sprintf(format, args...))
}

// process verifies if log level is above threshold and logs entry.
// It acquires timestamp, and source code context.
func (e *Entry) process(level Level, msg string) {
	if !e.Logger.PassThreshold(level) {
		return
	}
	e.Level = level
	e.Message = msg
	e.Timestamp = time.Now()
	e.CallContext = getCallContext(e.depth + 1)

	e.Logger.process(e)
}

// Emergency logs emergency level message.
func (e *Entry) Emergency(args ...interface{}) {
	e.IncDepth(1).Log(EmergLevel, args...)
}

// Alert logs alert level message.
func (e *Entry) Alert(args ...interface{}) {
	e.IncDepth(1).Log(AlertLevel, args...)
}

// Critical logs critical level message.
func (e *Entry) Critical(args ...interface{}) {
	e.IncDepth(1).Log(CritLevel, args...)
}

// Error logs error level message.
func (e *Entry) Error(args ...interface{}) {
	e.IncDepth(1).Log(ErrLevel, args...)
}

// Warning logs warning level message.
func (e *Entry) Warning(args ...interface{}) {
	e.IncDepth(1).Log(WarningLevel, args...)
}

// Notice logs notice level message.
func (e *Entry) Notice(args ...interface{}) {
	e.IncDepth(1).Log(NoticeLevel, args...)
}

// Info logs info level message.
func (e *Entry) Info(args ...interface{}) {
	e.IncDepth(1).Log(InfoLevel, args...)
}

// Debug logs debug level message.
func (e *Entry) Debug(args ...interface{}) {
	e.IncDepth(1).Log(DebugLevel, args...)
}

// Emergencyf logs emergency level formatted message.
func (e *Entry) Emergencyf(format string, args ...interface{}) {
	e.IncDepth(1).Logf(EmergLevel, format, args...)
}

// Alertf logs alert level formatted message.
func (e *Entry) Alertf(format string, args ...interface{}) {
	e.IncDepth(1).Logf(AlertLevel, format, args...)
}

// Criticalf logs critical level formatted message.
func (e *Entry) Criticalf(format string, args ...interface{}) {
	e.IncDepth(1).Logf(CritLevel, format, args...)
}

// Errorf logs error level formatted message.
func (e *Entry) Errorf(format string, args ...interface{}) {
	e.IncDepth(1).Logf(ErrLevel, format, args...)
}

// Warningf logs warning level formatted message.
func (e *Entry) Warningf(format string, args ...interface{}) {
	e.IncDepth(1).Logf(WarningLevel, format, args...)
}

// Noticef logs notice level formatted message.
func (e *Entry) Noticef(format string, args ...interface{}) {
	e.IncDepth(1).Logf(NoticeLevel, format, args...)
}

// Infof logs info level formatted message.
func (e *Entry) Infof(format string, args ...interface{}) {
	e.IncDepth(1).Logf(InfoLevel, format, args...)
}

// Debugf logs debug level formatted message.
func (e *Entry) Debugf(format string, args ...interface{}) {
	e.IncDepth(1).Logf(DebugLevel, format, args...)
}

// WithProperty adds a single property to the log message.
func (e *Entry) WithProperty(key string, value interface{}) *Entry {
	return e.WithProperties(Properties{key: value})
}

// WithProperties adds properties to the log message.
func (e *Entry) WithProperties(props Properties) *Entry {
	if e.Properties == nil {
		e.Properties = make(Properties)
	}
	for k, v := range props {
		e.Properties[k] = v
	}
	return e
}

// WithError adds error property to the log message.
func (e *Entry) WithError(err error) *Entry {
	return e.WithProperty(ErrorProperty, err.Error())
}
