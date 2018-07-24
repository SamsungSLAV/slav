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
}

// Log builds log message and logs entry.
func (e *Entry) Log(level Level, args ...interface{}) {
	e.process(level, fmt.Sprint(args...))
}

// Logf builds formatted log message and logs entry.
func (e *Entry) Logf(level Level, format string, args ...interface{}) {
	e.process(level, fmt.Sprintf(format, args...))
}

// process verifies if log level is above threshold and logs entry.
func (e *Entry) process(level Level, msg string) {
	if !e.Logger.PassThreshold(level) {
		return
	}
	e.Level = level
	e.Message = msg

	e.Logger.process(e)
}

// Emergency logs emergency level message.
func (e *Entry) Emergency(args ...interface{}) {
	e.Log(EmergLevel, args...)
}

// Alert logs alert level message.
func (e *Entry) Alert(args ...interface{}) {
	e.Log(AlertLevel, args...)
}

// Critical logs critical level message.
func (e *Entry) Critical(args ...interface{}) {
	e.Log(CritLevel, args...)
}

// Error logs error level message.
func (e *Entry) Error(args ...interface{}) {
	e.Log(ErrLevel, args...)
}

// Warning logs warning level message.
func (e *Entry) Warning(args ...interface{}) {
	e.Log(WarningLevel, args...)
}

// Notice logs notice level message.
func (e *Entry) Notice(args ...interface{}) {
	e.Log(NoticeLevel, args...)
}

// Info logs info level message.
func (e *Entry) Info(args ...interface{}) {
	e.Log(InfoLevel, args...)
}

// Debug logs debug level message.
func (e *Entry) Debug(args ...interface{}) {
	e.Log(DebugLevel, args...)
}

// Emergencyf logs emergency level formatted message.
func (e *Entry) Emergencyf(format string, args ...interface{}) {
	e.Logf(EmergLevel, format, args...)
}

// Alertf logs alert level formatted message.
func (e *Entry) Alertf(format string, args ...interface{}) {
	e.Logf(AlertLevel, format, args...)
}

// Criticalf logs critical level formatted message.
func (e *Entry) Criticalf(format string, args ...interface{}) {
	e.Logf(CritLevel, format, args...)
}

// Errorf logs error level formatted message.
func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Logf(ErrLevel, format, args...)
}

// Warningf logs warning level formatted message.
func (e *Entry) Warningf(format string, args ...interface{}) {
	e.Logf(WarningLevel, format, args...)
}

// Noticef logs notice level formatted message.
func (e *Entry) Noticef(format string, args ...interface{}) {
	e.Logf(NoticeLevel, format, args...)
}

// Infof logs info level formatted message.
func (e *Entry) Infof(format string, args ...interface{}) {
	e.Logf(InfoLevel, format, args...)
}

// Debugf logs debug level formatted message.
func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Logf(DebugLevel, format, args...)
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
