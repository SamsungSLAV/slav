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

// Level of log entries importance.
type Level uint32

// The log level's definitions are consistent with Unix syslog levels.
const (
	// EmergLevel is used when system is unusable.
	// It matches syslog's LOG_EMERG level.
	EmergLevel Level = iota
	// AlertLevel is used when action must be taken immediately.
	// It matches syslog's LOG_ALERT level.
	AlertLevel
	// CritLevel is used when critical conditions occur.
	// It matches syslog's LOG_CRIT level.
	CritLevel
	// ErrLevel is used when error conditions occur.
	// It matches syslog's LOG_ERR level.
	ErrLevel
	// WarningLevel is used when warning conditions occur.
	// It matches syslog's LOG_WARNING level.
	WarningLevel
	// NoticeLevel is used when normal, but significant, conditions occur.
	// It matches syslog's LOG_NOTICE level.
	NoticeLevel
	// InfoLevel is used for logging informational message.
	// It matches syslog's LOG_INFO level.
	InfoLevel
	// DebugLevel is used for logging debug-level message.
	// It matches syslog's LOG_DEBUG level.
	DebugLevel
)

// Log levels strings
const (
	// EmergLevelStr is string representation of EmergLevel
	EmergLevelStr = "emergency"
	// AlertLevelStr is string representation of AlertLevel
	AlertLevelStr = "alert"
	// CritLevelStr is string representation of CritLevel
	CritLevelStr = "critical"
	// ErrLevelStr is string representation of ErrLevel
	ErrLevelStr = "error"
	// WarningLevelStr is string representation of WarningLevel
	WarningLevelStr = "warning"
	// NoticeLevelStr is string representation of NoticeLevel
	NoticeLevelStr = "notice"
	// InfoLevelStr is string representation of InfoLevel
	InfoLevelStr = "info"
	// DebugLevelStr is string representation of DebugLevel
	DebugLevelStr = "debug"
	// UnknownLevel is string representation of unknown logging level
	UnknownLevelStr = "unknown"
)

// String converts Level to human readable string.
func (l Level) String() string {
	switch l {
	case EmergLevel:
		return EmergLevelStr
	case AlertLevel:
		return AlertLevelStr
	case CritLevel:
		return CritLevelStr
	case ErrLevel:
		return ErrLevelStr
	case WarningLevel:
		return WarningLevelStr
	case NoticeLevel:
		return NoticeLevelStr
	case InfoLevel:
		return InfoLevelStr
	case DebugLevel:
		return DebugLevelStr
	default:
		return UnknownLevelStr
	}
}

// StringToLevel converts string value to loggers' Level type.
// It is to be used when providing user with ability to specify
// logging level - e.g. setting log level via cli flag..
// If string is not matched, invalid level (DebugLevel+1) and
// ErrInvalidLogLevel is returned.
func StringToLevel(l string) (Level, error) {
	switch l {
	case EmergLevelStr:
		return EmergLevel, nil
	case AlertLevelStr:
		return AlertLevel, nil
	case CritLevelStr:
		return CritLevel, nil
	case ErrLevelStr:
		return ErrLevel, nil
	case WarningLevelStr:
		return WarningLevel, nil
	case NoticeLevelStr:
		return NoticeLevel, nil
	case InfoLevelStr:
		return InfoLevel, nil
	case DebugLevelStr:
		return DebugLevel, nil
	}
	return DebugLevel + 1, ErrInvalidLogLevel
}

// IsValid verifies if level has a valid value.
func (l Level) IsValid() bool {
	return (l <= DebugLevel)
}
