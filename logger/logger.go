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
	"fmt"
	"os"
	"sync"
	"sync/atomic"
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

	// backends contains all Backends used currently by the Logger.
	backends map[string]Backend
}

// NewLogger creates a new Logger instance with default configuration.
// Default level threshold is set to InfoLevel.
func NewLogger() *Logger {
	return &Logger{
		threshold: DefaultThreshold,
		mutex:     new(sync.Mutex),
		backends:  make(map[string]Backend),
	}
}

// SetThreshold defines Logger's filter level.
// Only entries with level equal or less than threshold will be logged.
func (l *Logger) SetThreshold(level Level) error {
	if !level.IsValid() {
		return ErrInvalidLogLevel
	}

	atomic.StoreUint32((*uint32)(&l.threshold), uint32(level))
	return nil
}

// Threshold returns current Logger's filter level.
func (l *Logger) Threshold() Level {
	return Level(atomic.LoadUint32((*uint32)(&l.threshold)))
}

// PassThreshold verifies if message with given level passes threshold and should be logged.
func (l *Logger) PassThreshold(level Level) bool {
	return (level <= l.Threshold())
}

// AddBackend adds or replaces a backend with given name.
func (l *Logger) AddBackend(name string, b Backend) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	b.Logger = l
	l.backends[name] = b
}

// RemoveBackend removes a backend with given name.
func (l *Logger) RemoveBackend(name string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	_, ok := l.backends[name]
	if !ok {
		return ErrInvalidBackendName
	}

	delete(l.backends, name)
	return nil
}

// RemoveAllBackends clears all backends.
func (l *Logger) RemoveAllBackends() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.backends = make(map[string]Backend)
}

// newEntry creates a new log entry.
func (l *Logger) newEntry() *Entry {
	return &Entry{
		Logger:     l,
		Properties: make(Properties),
	}
}

// process pass entry to backends.
func (l *Logger) process(entry *Entry) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for name, backend := range l.backends {
		err := backend.process(entry)
		if err != nil {
			// The error is printed to stderr. Potential fail of printing is ignored.
			_, _ = fmt.Fprintf(os.Stderr, "Error <%s> printing log message to <%s> backend.\n",
				err.Error(), name)
		}
	}
}

// Log builds log message and logs entry.
func (l *Logger) Log(level Level, args ...interface{}) {
	l.newEntry().IncDepth(1).Log(level, args...)
}

// Logf builds formatted log message and logs entry.
func (l *Logger) Logf(level Level, format string, args ...interface{}) {
	l.newEntry().IncDepth(1).Logf(level, format, args...)
}

// Emergency logs emergency level message.
func (l *Logger) Emergency(args ...interface{}) {
	l.newEntry().IncDepth(1).Emergency(args...)
}

// Alert logs alert level message.
func (l *Logger) Alert(args ...interface{}) {
	l.newEntry().IncDepth(1).Alert(args...)
}

// Critical logs critical level message.
func (l *Logger) Critical(args ...interface{}) {
	l.newEntry().IncDepth(1).Critical(args...)
}

// Error logs error level message.
func (l *Logger) Error(args ...interface{}) {
	l.newEntry().IncDepth(1).Error(args...)
}

// Warning logs warning level message.
func (l *Logger) Warning(args ...interface{}) {
	l.newEntry().IncDepth(1).Warning(args...)
}

// Notice logs notice level message.
func (l *Logger) Notice(args ...interface{}) {
	l.newEntry().IncDepth(1).Notice(args...)
}

// Info logs info level message.
func (l *Logger) Info(args ...interface{}) {
	l.newEntry().IncDepth(1).Info(args...)
}

// Debug logs debug level message.
func (l *Logger) Debug(args ...interface{}) {
	l.newEntry().IncDepth(1).Debug(args...)
}

// Emergencyf logs emergency level formatted message.
func (l *Logger) Emergencyf(format string, args ...interface{}) {
	l.newEntry().IncDepth(1).Emergencyf(format, args...)
}

// Alertf logs alert level formatted message.
func (l *Logger) Alertf(format string, args ...interface{}) {
	l.newEntry().IncDepth(1).Alertf(format, args...)
}

// Criticalf logs critical level formatted message.
func (l *Logger) Criticalf(format string, args ...interface{}) {
	l.newEntry().IncDepth(1).Criticalf(format, args...)
}

// Errorf logs error level formatted message.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.newEntry().IncDepth(1).Errorf(format, args...)
}

// Warningf logs warning level formatted message.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.newEntry().IncDepth(1).Warningf(format, args...)
}

// Noticef logs notice level formatted message.
func (l *Logger) Noticef(format string, args ...interface{}) {
	l.newEntry().IncDepth(1).Noticef(format, args...)
}

// Infof logs info level formatted message.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.newEntry().IncDepth(1).Infof(format, args...)
}

// Debugf logs debug level formatted message.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.newEntry().IncDepth(1).Debugf(format, args...)
}

// WithProperty creates a log message with a single property.
func (l *Logger) WithProperty(key string, value interface{}) *Entry {
	return l.newEntry().WithProperty(key, value)
}

// WithProperties creates a log message with multiple properties.
func (l *Logger) WithProperties(props Properties) *Entry {
	return l.newEntry().WithProperties(props)
}

// WithError creates a log message with an error property.
func (l *Logger) WithError(err error) *Entry {
	return l.newEntry().WithError(err)
}

// IncDepth increases depth of an Entry for call stack calculations.
func (l *Logger) IncDepth(dep int) *Entry {
	return l.newEntry().IncDepth(dep)
}
