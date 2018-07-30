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

// defaultLogger is the only global variable in logger package.
// It contains the default logger.
var defaultLogger = newDefaultLogger()

// newDefaultLogger initializes the default logger for logger package.
func newDefaultLogger() *Logger {
	logger := NewLogger()
	logger.AddBackend("default", Backend{
		Filter:     NewFilterPassAll(),
		Serializer: NewSerializerText(),
		Writer:     NewWriterStderr(),
	})
	return logger
}

// SetDefault sets the default logger.
func SetDefault(logger *Logger) {
	defaultLogger = logger
}

// SetThreshold defines default Logger's filter level.
// Only entries with level equal or less than threshold will be logged.
func SetThreshold(level Level) error {
	return defaultLogger.SetThreshold(level)
}

// Threshold returns current default Logger's filter level.
func Threshold() Level {
	return defaultLogger.Threshold()
}

// AddBackend adds or replaces a backend with given name in default logger.
func AddBackend(name string, b Backend) {
	defaultLogger.AddBackend(name, b)
}

// RemoveBackend removes a backend with given name from default logger.
func RemoveBackend(name string) error {
	return defaultLogger.RemoveBackend(name)
}

// RemoveAllBackends clears all backends from default logger.
func RemoveAllBackends() {
	defaultLogger.RemoveAllBackends()
}

// Log builds log message and logs entry to default logger.
func Log(level Level, args ...interface{}) {
	defaultLogger.Log(level, args...)
}

// Logf builds formatted log message and logs entry to default logger.
func Logf(level Level, format string, args ...interface{}) {
	defaultLogger.Logf(level, format, args...)
}

// Emergency logs emergency level message to default logger.
func Emergency(args ...interface{}) {
	defaultLogger.Emergency(args...)
}

// Alert logs alert level message to default logger.
func Alert(args ...interface{}) {
	defaultLogger.Alert(args...)
}

// Critical logs critical level message to default logger.
func Critical(args ...interface{}) {
	defaultLogger.Critical(args...)
}

// Error logs error level message to default logger.
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Warning logs warning level message to default logger.
func Warning(args ...interface{}) {
	defaultLogger.Warning(args...)
}

// Notice logs notice level message to default logger.
func Notice(args ...interface{}) {
	defaultLogger.Notice(args...)
}

// Info logs info level message to default logger.
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Debug logs debug level message to default logger.
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Emergencyf logs emergency level formatted message to default logger.
func Emergencyf(format string, args ...interface{}) {
	defaultLogger.Emergencyf(format, args...)
}

// Alertf logs alert level formatted message to default logger.
func Alertf(format string, args ...interface{}) {
	defaultLogger.Alertf(format, args...)
}

// Criticalf logs critical level formatted message to default logger.
func Criticalf(format string, args ...interface{}) {
	defaultLogger.Criticalf(format, args...)
}

// Errorf logs error level formatted message to default logger.
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Warningf logs warning level formatted message to default logger.
func Warningf(format string, args ...interface{}) {
	defaultLogger.Warningf(format, args...)
}

// Noticef logs notice level formatted message to default logger.
func Noticef(format string, args ...interface{}) {
	defaultLogger.Noticef(format, args...)
}

// Infof logs info level formatted message to default logger.
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Debugf logs debug level formatted message to default logger.
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// WithProperty creates a log message with a single property in default logger.
func WithProperty(key string, value interface{}) *Entry {
	return defaultLogger.WithProperty(key, value)
}

// WithProperties creates a log message with multiple properties in default logger.
func WithProperties(props Properties) *Entry {
	return defaultLogger.WithProperties(props)
}

// WithError creates a log message with an error property in default logger.
func WithError(err error) *Entry {
	return defaultLogger.WithError(err)
}
