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
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"time"
)

// TimestampMode defines possible time stamp logging modes.
type TimestampMode uint8

const (
	// TimestampModeNone - no time stamp used.
	TimestampModeNone TimestampMode = iota
	// TimestampModeDiff - seconds since creation of SerializerText.
	// Usualy it is equal to the binary start as logger and its serializer
	// are created mostly by packages' init functions.
	TimestampModeDiff
	// TimestampModeFull - date and time in UTC.
	TimestampModeFull
)

// QuoteMode defines possible quoting modes.
type QuoteMode uint8

const (
	// QuoteModeNone - no quoting is used.
	QuoteModeNone QuoteMode = iota
	// QuoteModeSpecial - values containing special characters are quoted.
	QuoteModeSpecial
	// QuoteModeSpecialAndEmpty - values containing special characters and empty are quoted.
	QuoteModeSpecialAndEmpty
	// QuoteModeAll - all values are quoted.
	QuoteModeAll
)

// CallContextMode defines possible modes of printing call source code context.
type CallContextMode uint8

const (
	// CallContextModeNone - no context is used.
	CallContextModeNone CallContextMode = iota
	// CallContextModeCompact - file name and line number are used.
	CallContextModeCompact
	// CallContextModeFunction - file name, line number and function name are used.
	CallContextModeFunction
	// CallContextModeFile - full file path and line number are used.
	CallContextModeFile
	// CallContextModePackage - package name, line number and function are used.
	CallContextModePackage
)

// Define default SerializerText properties.
const (
	// DefaultSerializerTextTimeFormat is the default date and time format.
	DefaultSerializerTextTimeFormat = time.RFC3339
	// DefaultTimestampMode is the default mode for logging time stamp.
	DefaultTimestampMode = TimestampModeDiff
	// DefaultQuoteMode is the default quoting mode.
	DefaultQuoteMode = QuoteModeSpecialAndEmpty
	// DefaultCallContextMode is the default context mode.
	DefaultCallContextMode = CallContextModeCompact
)

// Define termninal codes for colors.
const (
	red     = "\x1b[31m"
	green   = "\x1b[32m"
	yellow  = "\x1b[33m"
	blue    = "\x1b[34m"
	magenta = "\x1b[35m"
	cyan    = "\x1b[36m"

	bold   = "\x1b[1m"
	invert = "\x1b[7m"
	off    = "\x1b[0m"

	propkey  = cyan
	path     = blue + bold
	function = green
	line     = yellow
)

// levelColoring stores mapping between logger levels and their colors.
var levelColoring = map[Level]string{
	EmergLevel:   red + bold + invert,
	AlertLevel:   cyan + bold + invert,
	CritLevel:    magenta + bold + invert,
	ErrLevel:     red + bold,
	WarningLevel: yellow + bold,
	NoticeLevel:  blue + bold,
	InfoLevel:    green + bold,
	DebugLevel:   bold,
}

// asciiTableSize defines SerializerText.ascii tab size.
const asciiTableSize = 128

// SerializerText serializes entry to text format.
type SerializerText struct {
	// TimeFormat defines format for displaying date and time.
	// Used only when TimestampMode is set to TimestampModeFull.
	// See https://godoc.org/time#Time.Format description for details.
	TimeFormat string

	// TimestampMode defines mode for logging date and time.
	TimestampMode TimestampMode

	// QuoteMode defines which values are quoted.
	QuoteMode QuoteMode

	// CallContextMode defines way of serializing source code context.
	CallContextMode CallContextMode

	// UseColors set to true enables usage of colors.
	UseColors bool

	// ascii defines which characters are special and need quoting.
	ascii [asciiTableSize]bool

	// initASCII ensures that initialization of ascii is done only once.
	initASCII sync.Once

	// baseTime is the Timestamp from which elapsed time is calculated in TimestampModeDiff.
	baseTime time.Time
}

// NewSerializerText creates and returns a new default SerializerText with default values.
func NewSerializerText() *SerializerText {
	return &SerializerText{
		TimestampMode:   DefaultTimestampMode,
		TimeFormat:      DefaultSerializerTextTimeFormat,
		UseColors:       true,
		QuoteMode:       DefaultQuoteMode,
		CallContextMode: DefaultCallContextMode,
		baseTime:        time.Now(),
	}
}

// setDefaultsOnInvalid fixes invalid serializer's fields to default values.
func (s *SerializerText) setDefaultsOnInvalid() {
	if s.TimestampMode > TimestampModeFull {
		s.TimestampMode = DefaultTimestampMode
	}
	if len(s.TimeFormat) == 0 {
		s.TimeFormat = DefaultSerializerTextTimeFormat
	}
	if s.QuoteMode > QuoteModeAll {
		s.QuoteMode = DefaultQuoteMode
	}
	if s.CallContextMode > CallContextModePackage {
		s.CallContextMode = DefaultCallContextMode
	}
	if s.baseTime.IsZero() {
		s.baseTime = time.Now()
	}
}

// initASCIITable initializes ascii table.
func (s *SerializerText) initASCIITable() {
	for i := 'a'; i <= 'z'; i++ {
		s.ascii[i] = true
	}
	for i := 'A'; i <= 'Z'; i++ {
		s.ascii[i] = true
	}
	for i := '0'; i <= '9'; i++ {
		s.ascii[i] = true
	}
	for _, i := range []byte{'-', '.', '_', '/', '@', '^', '+'} {
		s.ascii[i] = true
	}
}

// quotingFormat returns formatting string for message depending on quoting settings.
func (s *SerializerText) quotingFormat(msg string) string {
	const (
		quoting    = "%q"
		notquoting = "%s"
	)

	switch s.QuoteMode {
	case QuoteModeSpecialAndEmpty:
		if len(msg) == 0 {
			return quoting
		}
		fallthrough
	case QuoteModeSpecial:
		s.initASCII.Do(s.initASCIITable)
		for _, i := range msg {
			if i >= asciiTableSize || !s.ascii[i] {
				return quoting
			}
		}
		return notquoting
	case QuoteModeAll:
		return quoting
	case QuoteModeNone:
	default:
	}
	return notquoting
}

// appendTimestamp to log message being created in buf.
func (s *SerializerText) appendTimestamp(buf io.Writer, t *time.Time) (err error) {
	switch s.TimestampMode {
	case TimestampModeDiff:
		const precision = 6
		_, err = fmt.Fprintf(buf, "[%.*f] ", precision, t.Sub(s.baseTime).Seconds())
	case TimestampModeFull:
		_, err = fmt.Fprintf(buf, "[%s] ", t.Format(s.TimeFormat))
	}
	return err
}

// appendLevel to log message being created in buf.
func (s *SerializerText) appendLevel(buf io.Writer, level Level) (err error) {
	var format string
	if s.UseColors {
		format = "[" + levelColoring[level] + "%.*s" + off + "] "
	} else {
		format = "[%.*s] "
	}
	const precision = 3
	_, err = fmt.Fprintf(buf, format, precision, strings.ToUpper(level.String()))
	return err
}

// appendCallContext to log message being created in buf.
func (s *SerializerText) appendCallContext(buf io.Writer, ctx *CallContext) (err error) {
	if ctx == nil || s.CallContextMode == CallContextModeNone {
		return nil
	}
	var fPath, fEnd, fFunc string
	if s.UseColors {
		fPath = "[" + path
		fEnd = off + ":" + line + "%d" + off + "] "
		fFunc = off + ":" + function
	} else {
		fPath = "["
		fEnd = ":%d] "
		fFunc = ":"
	}
	switch s.CallContextMode {
	case CallContextModeCompact:
		_, err = fmt.Fprintf(buf, fPath+"%s"+fEnd, ctx.File, ctx.Line)
	case CallContextModeFunction:
		if len(ctx.Type) > 0 {
			_, err = fmt.Fprintf(buf, fPath+"%s"+fFunc+"%s.%s"+fEnd, ctx.File, ctx.Type,
				ctx.Function, ctx.Line)
		} else {
			_, err = fmt.Fprintf(buf, fPath+"%s"+fFunc+"%s"+fEnd, ctx.File, ctx.Function, ctx.Line)
		}
	case CallContextModeFile:
		_, err = fmt.Fprintf(buf, fPath+"%s%s"+fEnd, ctx.Path, ctx.File, ctx.Line)
	case CallContextModePackage:
		if len(ctx.Type) > 0 {
			_, err = fmt.Fprintf(buf, fPath+"%s"+fFunc+"%s.%s"+fEnd, ctx.Package, ctx.Type,
				ctx.Function, ctx.Line)
		} else {
			_, err = fmt.Fprintf(buf, fPath+"%s"+fFunc+"%s"+fEnd, ctx.Package, ctx.Function,
				ctx.Line)
		}
	}
	return err
}

// appendMessage to log message being created in buf.
func (s *SerializerText) appendMessage(buf io.Writer, msg string) (err error) {
	format := s.quotingFormat(msg) + " "
	_, err = fmt.Fprintf(buf, format, msg)
	return err
}

// appendProperties to log message being created in buf.
func (s *SerializerText) appendProperties(buf io.Writer, properties Properties) (err error) {
	if len(properties) == 0 {
		return
	}
	_, err = fmt.Fprintf(buf, "{")
	if err != nil {
		return err
	}
	keys := make([]string, len(properties))
	i := 0
	for k := range properties {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := properties[k]
		format := s.quotingFormat(k)
		if s.UseColors {
			format = propkey + format + off
		}
		format = format + ":"
		_, err = fmt.Fprintf(buf, format, k)
		if err != nil {
			return err
		}

		value := fmt.Sprint(v)
		format = s.quotingFormat(value) + ";"
		_, err = fmt.Fprintf(buf, format, value)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(buf, "}")
	return err
}

// serialize writes parts of Entry to given writer.
func (s *SerializerText) serialize(entry *Entry, buf io.Writer) error {
	if entry == nil {
		return ErrInvalidEntry
	}
	err := s.appendTimestamp(buf, &entry.Timestamp)
	if err != nil {
		return err
	}
	err = s.appendLevel(buf, entry.Level)
	if err != nil {
		return err
	}
	err = s.appendCallContext(buf, entry.CallContext)
	if err != nil {
		return err
	}
	err = s.appendMessage(buf, entry.Message)
	if err != nil {
		return err
	}
	err = s.appendProperties(buf, entry.Properties)
	return err
}

// Serialize implements Serializer interface in SerializerText.
func (s *SerializerText) Serialize(entry *Entry) ([]byte, error) {
	s.setDefaultsOnInvalid()

	buf := &bytes.Buffer{}
	err := s.serialize(entry, buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
