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
	"log/syslog"
)

// WriterSyslog writes to syslog using standard log/syslog package.
// It implements Writer interface.
type WriterSyslog struct {
	syslogClient *syslog.Writer
}

// NewWriterSyslog creates a new WriterSyslog object connecting it to the log daemon.
// Connection uses specified network and raddr address
func NewWriterSyslog(network, raddr string, facility syslog.Priority, tag string) *WriterSyslog {
	w, err := syslog.Dial(network, raddr, facility, tag)
	if err != nil {
		panic(err)
	}
	return &WriterSyslog{
		syslogClient: w,
	}
}

// Write writes to syslog. It implements Writer interface in WriterSyslog.
func (w *WriterSyslog) Write(level Level, p []byte) (int, error) {
	msg := string(p)
	switch level {
	case EmergLevel:
		return 0, w.syslogClient.Emerg(msg)
	case AlertLevel:
		return 0, w.syslogClient.Alert(msg)
	case CritLevel:
		return 0, w.syslogClient.Crit(msg)
	case ErrLevel:
		return 0, w.syslogClient.Err(msg)
	case WarningLevel:
		return 0, w.syslogClient.Warning(msg)
	case NoticeLevel:
		return 0, w.syslogClient.Notice(msg)
	case InfoLevel:
		return 0, w.syslogClient.Info(msg)
	case DebugLevel:
		return 0, w.syslogClient.Debug(msg)
	default:
		return 0, ErrInvalidLogLevel
	}
}
