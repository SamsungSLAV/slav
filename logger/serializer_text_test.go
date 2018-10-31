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
	"errors"
	"io"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	T "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

// FailingWriter is a dummy io.Writer implementation that returns a given error after writting
// given amount of bytes.
type FailingWriter struct {
	failAfter int
	fail      error
}

func NewFailingWriter(failAfter int, fail error) io.Writer {
	return &FailingWriter{
		failAfter: failAfter,
		fail:      fail,
	}
}
func (w *FailingWriter) Write(p []byte) (n int, err error) {
	l := len(p)
	if l < w.failAfter {
		w.failAfter -= l
		return l, nil
	}
	return l - w.failAfter, w.fail
}

var _ = Describe("SerializerText", func() {
	const (
		empty     = ""
		nospecial = "Alice_has_got_7_cats."
		special   = "#$?"
		mixed     = "What?"
	)

	var (
		s         *SerializerText
		before    time.Time
		testError error
	)

	BeforeEach(func() {
		before = time.Now()
		s = NewSerializerText()
		testError = errors.New("TestError")
	})
	Describe("NewSerializerText", func() {
		It("should create a new object with default configuration", func() {
			Expect(s).NotTo(BeNil())
			Expect(s.TimestampMode).To(Equal(DefaultTimestampMode))
			Expect(s.TimeFormat).To(Equal(DefaultSerializerTextTimeFormat))
			Expect(s.UseColors).To(BeTrue())
			Expect(s.QuoteMode).To(Equal(DefaultQuoteMode))
			Expect(s.CallContextMode).To(Equal(DefaultCallContextMode))
			Expect(s.baseTime).To(BeTemporally(">", before))
			Expect(s.baseTime).To(BeTemporally("<", time.Now()))
			Expect(s.ascii).To(HaveLen(asciiTableSize))
		})
	})
	Describe("setDefaultsOnInvalid", func() {
		It("should set empty object to correct values", func() {
			before = time.Now()
			s = &SerializerText{}
			s.setDefaultsOnInvalid()

			Expect(s).NotTo(BeNil())
			Expect(s.TimeFormat).To(Equal(DefaultSerializerTextTimeFormat))
			Expect(s.baseTime).To(BeTemporally(">", before))
			Expect(s.baseTime).To(BeTemporally("<", time.Now()))
			Expect(s.ascii).To(HaveLen(asciiTableSize))
		})
		It("should set timestamp mode to default value if out of range", func() {
			s = &SerializerText{}
			s.TimestampMode = 0xBB
			s.setDefaultsOnInvalid()

			Expect(s.TimestampMode).To(Equal(DefaultTimestampMode))
		})
		It("should set time formatter mode to default value if empty", func() {
			s = &SerializerText{}
			s.TimeFormat = ""
			s.setDefaultsOnInvalid()

			Expect(s.TimeFormat).To(Equal(DefaultSerializerTextTimeFormat))
		})
		It("should set quote mode to default value if out of range", func() {
			s = &SerializerText{}
			s.QuoteMode = 0xBB
			s.setDefaultsOnInvalid()

			Expect(s.QuoteMode).To(Equal(DefaultQuoteMode))
		})
		It("should set call context mode to default value if out of range", func() {
			s = &SerializerText{}
			s.CallContextMode = 0xBB
			s.setDefaultsOnInvalid()

			Expect(s.CallContextMode).To(Equal(DefaultCallContextMode))
		})
		It("should set base time to now if not set before", func() {
			s = &SerializerText{}
			s.baseTime = time.Time{}
			before = time.Now()
			s.setDefaultsOnInvalid()

			Expect(s.baseTime).To(BeTemporally(">", before))
			Expect(s.baseTime).To(BeTemporally("<", time.Now()))
		})
		T.DescribeTable("should not change timestamp mode if valid",
			func(mode TimestampMode) {
				s = &SerializerText{}
				s.TimestampMode = mode
				s.setDefaultsOnInvalid()

				Expect(s.TimestampMode).To(Equal(mode))
			},
			T.Entry("TimestampModeNone", TimestampModeNone),
			T.Entry("TimestampModeDiff", TimestampModeDiff),
			T.Entry("TimestampModeFull", TimestampModeFull),
		)
		It("should not change time format if set", func() {
			s = &SerializerText{}
			s.TimeFormat = time.ANSIC
			s.setDefaultsOnInvalid()

			Expect(s.TimeFormat).To(Equal(time.ANSIC))
		})
		T.DescribeTable("should not change quote mode if valid",
			func(mode QuoteMode) {
				s = &SerializerText{}
				s.QuoteMode = mode
				s.setDefaultsOnInvalid()

				Expect(s.QuoteMode).To(Equal(mode))
			},
			T.Entry("QuoteModeNone", QuoteModeNone),
			T.Entry("QuoteModeSpecial", QuoteModeSpecial),
			T.Entry("QuoteModeSpecialAndEmpty", QuoteModeSpecialAndEmpty),
			T.Entry("QuoteModeAll", QuoteModeAll),
		)
		T.DescribeTable("should not change call context mode if valid",
			func(mode CallContextMode) {
				s = &SerializerText{}
				s.CallContextMode = mode
				s.setDefaultsOnInvalid()

				Expect(s.CallContextMode).To(Equal(mode))
			},
			T.Entry("CallContextModeNone", CallContextModeNone),
			T.Entry("CallContextModeCompact", CallContextModeCompact),
			T.Entry("CallContextModeFunction", CallContextModeFunction),
			T.Entry("CallContextModeFile", CallContextModeFile),
			T.Entry("CallContextModePackage", CallContextModePackage),
		)
		It("should not change base time if set", func() {
			magicTime := time.Unix(1234567890, 0)
			s = &SerializerText{}
			s.baseTime = magicTime
			s.setDefaultsOnInvalid()

			Expect(s.baseTime).To(BeTemporally("==", magicTime))
		})
	})
	Describe("initASCIITable", func() {
		const good = "abcdefghijklmnopqrstuvwxyz" +
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
			"0123456789" +
			"-._/@^+"
		It("should init table properly", func() {
			// Clear the table.
			for i := 0; i < asciiTableSize; i++ {
				s.ascii[i] = false
			}

			s.initASCIITable()
			for i := 0; i < asciiTableSize; i++ {
				expected := strings.ContainsAny(good, string([]byte{byte(i)}))
				Expect(s.ascii[i]).To(Equal(expected), "checking %d", i)
			}
		})
	})
	Describe("quotingFormat", func() {
		const badQuoteMode = QuoteMode(0xBB)

		T.DescribeTable("should return proper formatting string",
			func(mode QuoteMode, msg string, expected string) {
				s.QuoteMode = mode
				Expect(s.quotingFormat(msg)).To(Equal(expected))
			},
			T.Entry("None/empty", QuoteModeNone, empty, "%s"),
			T.Entry("None/nospecial", QuoteModeNone, nospecial, "%s"),
			T.Entry("None/special", QuoteModeNone, special, "%s"),
			T.Entry("None/mixed", QuoteModeNone, mixed, "%s"),

			T.Entry("SpecialAndEmpty/empty", QuoteModeSpecialAndEmpty, empty, "%q"),
			T.Entry("SpecialAndEmpty/nospecial", QuoteModeSpecialAndEmpty, nospecial, "%s"),
			T.Entry("SpecialAndEmpty/special", QuoteModeSpecialAndEmpty, special, "%q"),
			T.Entry("SpecialAndEmpty/mixed", QuoteModeSpecialAndEmpty, mixed, "%q"),

			T.Entry("Special/empty", QuoteModeSpecial, empty, "%s"),
			T.Entry("Special/nospecial", QuoteModeSpecial, nospecial, "%s"),
			T.Entry("Special/special", QuoteModeSpecial, special, "%q"),
			T.Entry("Special/mixed", QuoteModeSpecial, mixed, "%q"),

			T.Entry("All/empty", QuoteModeAll, empty, "%q"),
			T.Entry("All/nospecial", QuoteModeAll, nospecial, "%q"),
			T.Entry("All/special", QuoteModeAll, special, "%q"),
			T.Entry("All/mixed", QuoteModeAll, mixed, "%q"),

			T.Entry("BAD/empty", badQuoteMode, empty, "%s"),
			T.Entry("BAD/nospecial", badQuoteMode, nospecial, "%s"),
			T.Entry("BAD/special", badQuoteMode, special, "%s"),
			T.Entry("BAD/mixed", badQuoteMode, mixed, "%s"),
		)
	})
	Describe("serialization helpers", func() {
		var buf *bytes.Buffer
		BeforeEach(func() {
			buf = &bytes.Buffer{}
		})
		Describe("appendTimestamp", func() {
			It("should do nothing when time stamp mode is set to None", func() {
				s.TimestampMode = TimestampModeNone
				stamp := time.Now()
				err := s.appendTimestamp(buf, &stamp)

				Expect(err).NotTo(HaveOccurred())
				Expect(buf.Len()).To(BeZero())
			})
			T.DescribeTable("should serialize proper diff when time stamp is set to Diff",
				func(nanoseconds int, expected string) {
					s.TimestampMode = TimestampModeDiff
					stamp := s.baseTime.Add(time.Duration(nanoseconds) * time.Nanosecond)
					err := s.appendTimestamp(buf, &stamp)

					Expect(err).NotTo(HaveOccurred())
					Expect(buf.String()).To(Equal(expected))
				},
				T.Entry("Second", 1000*1000*1000, "[1.000000] "),
				T.Entry("Milisecond", 1000*1000, "[0.001000] "),
				T.Entry("Microsecond", 1000, "[0.000001] "),
				T.Entry("Nanosecond", 1, "[0.000000] "),
				T.Entry("Custom", 1234567890, "[1.234568] "),
			)
			T.DescribeTable("should serialize proper time when time stamp is set to Full",
				func(days int, expected string) {
					s.TimestampMode = TimestampModeFull
					stamp := time.Unix(1234567890, 0).AddDate(0, 0, days)
					err := s.appendTimestamp(buf, &stamp)

					Expect(err).NotTo(HaveOccurred())
					Expect(buf.String()).To(Equal(expected))
				},
				T.Entry("Day", 1, "[2009-02-15T00:31:30+01:00] "),
				T.Entry("Week", 7, "[2009-02-21T00:31:30+01:00] "),
				T.Entry("Month", 30, "[2009-03-16T00:31:30+01:00] "),
				T.Entry("Year", 365, "[2010-02-14T00:31:30+01:00] "),
			)
			T.DescribeTable("should return error if writing fails",
				func(mode TimestampMode) {
					w := NewFailingWriter(0, testError)
					s.TimestampMode = mode
					stamp := time.Now()
					err := s.appendTimestamp(w, &stamp)
					Expect(err).To(Equal(testError))
				},
				T.Entry("Diff", TimestampModeDiff),
				T.Entry("Full", TimestampModeFull),
			)
		})
		Describe("appendLevel", func() {
			T.DescribeTable("should serialize proper level without colors",
				func(level Level, expected string) {
					s.UseColors = false
					err := s.appendLevel(buf, level)

					Expect(err).NotTo(HaveOccurred())
					Expect(buf.String()).To(Equal(expected))
				},
				T.Entry("EmergLevel", EmergLevel, "[EME] "),
				T.Entry("AlertLevel", AlertLevel, "[ALE] "),
				T.Entry("CritLevel", CritLevel, "[CRI] "),
				T.Entry("ErrLevel", ErrLevel, "[ERR] "),
				T.Entry("WarningLevel", WarningLevel, "[WAR] "),
				T.Entry("NoticeLevel", NoticeLevel, "[NOT] "),
				T.Entry("InfoLevel", InfoLevel, "[INF] "),
				T.Entry("DebugLevel", DebugLevel, "[DEB] "),
			)
			T.DescribeTable("should serialize proper level with colors",
				func(level Level, expected string) {
					s.UseColors = true
					err := s.appendLevel(buf, level)

					Expect(err).NotTo(HaveOccurred())
					Expect(buf.String()).To(Equal(expected))
				},
				T.Entry("EmergLevel", EmergLevel, "["+red+bold+invert+"EME"+off+"] "),
				T.Entry("AlertLevel", AlertLevel, "["+cyan+bold+invert+"ALE"+off+"] "),
				T.Entry("CritLevel", CritLevel, "["+magenta+bold+invert+"CRI"+off+"] "),
				T.Entry("ErrLevel", ErrLevel, "["+red+bold+"ERR"+off+"] "),
				T.Entry("WarningLevel", WarningLevel, "["+yellow+bold+"WAR"+off+"] "),
				T.Entry("NoticeLevel", NoticeLevel, "["+blue+bold+"NOT"+off+"] "),
				T.Entry("InfoLevel", InfoLevel, "["+green+bold+"INF"+off+"] "),
				T.Entry("DebugLevel", DebugLevel, "["+bold+"DEB"+off+"] "),
			)
			T.DescribeTable("should return error if writing fails",
				func(useColors bool) {
					w := NewFailingWriter(0, testError)
					s.UseColors = useColors
					err := s.appendLevel(w, ErrLevel)
					Expect(err).To(Equal(testError))
				},
				T.Entry("WithColors", true),
				T.Entry("WithoutColors", false),
			)
		})
		Describe("appendCallContext", func() {
			var withType, withoutType CallContext
			BeforeEach(func() {
				withType = CallContext{
					Path:     "p",
					File:     "f",
					Line:     98765,
					Package:  "a",
					Type:     "t",
					Function: "u",
				}
				withoutType = withType
				withoutType.Type = ""
			})

			T.DescribeTable("should serialize context without colors",
				func(mode CallContextMode, ctx *CallContext, expected string) {
					s.UseColors = false
					s.CallContextMode = mode
					err := s.appendCallContext(buf, ctx)

					Expect(err).NotTo(HaveOccurred())
					Expect(buf.String()).To(Equal(expected))
				},
				T.Entry("None nil", CallContextModeNone, nil, ""),
				T.Entry("None withType", CallContextModeNone, &withType, ""),
				T.Entry("None withoutType", CallContextModeNone, &withoutType, ""),
				T.Entry("Compact nil", CallContextModeCompact, nil, ""),
				T.Entry("Compact withType", CallContextModeCompact, &withType, "[f:98765] "),
				T.Entry("Compact withoutType", CallContextModeCompact, &withoutType, "[f:98765] "),
				T.Entry("Function nil", CallContextModeFunction, nil, ""),
				T.Entry("Function withType", CallContextModeFunction, &withType, "[f:t.u:98765] "),
				T.Entry("Function withoutType", CallContextModeFunction, &withoutType,
					"[f:u:98765] "),
				T.Entry("File nil", CallContextModeFile, nil, ""),
				T.Entry("File withType", CallContextModeFile, &withType, "[pf:98765] "),
				T.Entry("File withoutType", CallContextModeFile, &withoutType, "[pf:98765] "),
				T.Entry("Package nil", CallContextModePackage, nil, ""),
				T.Entry("Package withType", CallContextModePackage, &withType, "[a:t.u:98765] "),
				T.Entry("Package withoutType", CallContextModePackage, &withoutType,
					"[a:u:98765] "),
			)
			T.DescribeTable("should serialize context with colors",
				func(mode CallContextMode, ctx *CallContext, expected string) {
					s.UseColors = true
					s.CallContextMode = mode
					err := s.appendCallContext(buf, ctx)

					Expect(err).NotTo(HaveOccurred())
					Expect(buf.String()).To(Equal(expected))
				},
				T.Entry("None nil", CallContextModeNone, nil, ""),
				T.Entry("None withType", CallContextModeNone, &withType, ""),
				T.Entry("None withoutType", CallContextModeNone, &withoutType, ""),
				T.Entry("Compact nil", CallContextModeCompact, nil, ""),
				T.Entry("Compact withType", CallContextModeCompact, &withType,
					"["+blue+bold+"f"+off+":"+yellow+"98765"+off+"] "),
				T.Entry("Compact withoutType", CallContextModeCompact, &withoutType,
					"["+blue+bold+"f"+off+":"+yellow+"98765"+off+"] "),
				T.Entry("Function nil", CallContextModeFunction, nil, ""),
				T.Entry("Function withType", CallContextModeFunction, &withType,
					"["+blue+bold+"f"+off+":"+green+"t.u"+off+":"+yellow+"98765"+off+"] "),
				T.Entry("Function withoutType", CallContextModeFunction, &withoutType,
					"["+blue+bold+"f"+off+":"+green+"u"+off+":"+yellow+"98765"+off+"] "),
				T.Entry("File nil", CallContextModeFile, nil, ""),
				T.Entry("File withType", CallContextModeFile, &withType,
					"["+blue+bold+"pf"+off+":"+yellow+"98765"+off+"] "),
				T.Entry("File withoutType", CallContextModeFile, &withoutType,
					"["+blue+bold+"pf"+off+":"+yellow+"98765"+off+"] "),
				T.Entry("Package nil", CallContextModePackage, nil, ""),
				T.Entry("Package withType", CallContextModePackage, &withType,
					"["+blue+bold+"a"+off+":"+green+"t.u"+off+":"+yellow+"98765"+off+"] "),
				T.Entry("Package withoutType", CallContextModePackage, &withoutType,
					"["+blue+bold+"a"+off+":"+green+"u"+off+":"+yellow+"98765"+off+"] "),
			)
			T.DescribeTable("should return error if writing fails",
				func(mode CallContextMode, ctx *CallContext) {
					s.UseColors = false
					s.CallContextMode = mode
					w := NewFailingWriter(0, testError)
					err := s.appendCallContext(w, ctx)
					Expect(err).To(Equal(testError))
				},
				T.Entry("Compact withType", CallContextModeCompact, &withType),
				T.Entry("Compact withoutType", CallContextModeCompact, &withoutType),
				T.Entry("Function withType", CallContextModeFunction, &withType),
				T.Entry("Function withoutType", CallContextModeFunction, &withoutType),
				T.Entry("File withType", CallContextModeFile, &withType),
				T.Entry("File withoutType", CallContextModeFile, &withoutType),
				T.Entry("Package withType", CallContextModePackage, &withType),
				T.Entry("Package withoutType", CallContextModePackage, &withoutType),
			)
		})
		Describe("appendMessage", func() {
			T.DescribeTable("should serialize properly quoted message",
				func(msg string, expected string) {
					err := s.appendMessage(buf, msg)
					Expect(err).NotTo(HaveOccurred())
					Expect(buf.String()).To(Equal(expected))
				},

				T.Entry("empty", empty, `"" `),
				T.Entry("nospecial", nospecial, nospecial+" "),
				T.Entry("special", special, `"`+special+`" `),
				T.Entry("mixed", mixed, `"`+mixed+`" `),
			)
			It("should return error if writing fails", func() {
				w := NewFailingWriter(0, testError)
				err := s.appendMessage(w, special)
				Expect(err).To(Equal(testError))
			})
		})
		Describe("appendProperties", func() {
			It("should do nothing when properties are empty", func() {
				p := Properties{}
				err := s.appendProperties(buf, p)

				Expect(err).NotTo(HaveOccurred())
				Expect(buf.Len()).To(BeZero())
			})
			It("should serialize all key-value pairs without color", func() {
				s.UseColors = false
				p := Properties{
					"name": "Alice",
					"hash": "#$%%@",
					"male": true,
					"age":  37,
					"skills": Properties{
						"coding": 7,
					},
					"issues": "",
				}
				err := s.appendProperties(buf, p)

				Expect(err).NotTo(HaveOccurred())
				expected := `{age:37;hash:"#$%%@";issues:"";male:true;name:Alice;skills:` +
					`"map[coding:7]";}`
				Expect(buf.String()).To(Equal(expected))
			})
			It("should serialize all key-value pairs with bolded key names", func() {
				s.UseColors = true
				p := Properties{
					"name": "Alice",
					"hash": "#$%%@",
					"male": true,
					"age":  37,
					"skills": Properties{
						"coding": 7,
					},
					"issues": "",
				}
				err := s.appendProperties(buf, p)

				Expect(err).NotTo(HaveOccurred())
				expected := "{" + propkey + "age" + off + ":37;" + propkey + "hash" + off +
					`:"#$%%@";` + propkey + "issues" + off + `:"";` + propkey + "male" + off +
					":true;" + propkey + "name" + off + ":Alice;" + propkey + "skills" + off +
					`:"map[coding:7]";}`
				Expect(buf.String()).To(Equal(expected))
			})
			T.DescribeTable("should return error if writing fails",
				func(bytes int) {
					w := NewFailingWriter(bytes, testError)
					p := Properties{
						"name": "Alice",
					}
					s.UseColors = false
					err := s.appendProperties(w, p)
					Expect(err).To(Equal(testError))
				},
				//                                          0        1
				//                                          1234567890123
				//                                          *  *     *  *
				// expected serialized properties string:	{name:Alice;}
				T.Entry("OpeningBracket", 1),
				T.Entry("Key", 4),
				T.Entry("Value", 10),
				T.Entry("ClosingBracket", 13),
			)
		})
	})
	Describe("serialize", func() {
		T.DescribeTable("should return error if writing fails",
			func(bytes int) {
				w := NewFailingWriter(bytes, testError)
				s.UseColors = false
				s.baseTime = time.Unix(1234567890, 0)
				entry := &Entry{
					Level:     ErrLevel,
					Message:   "message",
					Timestamp: s.baseTime.Add(time.Duration(11) * time.Minute),
					CallContext: &CallContext{
						Path:     "p",
						File:     "f",
						Line:     98765,
						Package:  "a",
						Type:     "t",
						Function: "u",
					},
					Properties: Properties{
						"name": "Alice",
						"hash": "#$%%@",
						"male": true,
					},
				}
				err := s.serialize(entry, w)
				Expect(err).To(Equal(testError))
			},
			// expected serialized properties string:
			// [660.000000] [ERR] [f:98765] message {hash:"#$%%@";male:true;name:Alice;}
			// 0        1         2         3         4         5         6         7
			// 1234567890123456789012345678901234567890123456789012345678901234567890123
			//       *        *       *        *                    *
			T.Entry("Timestamp", 7),
			T.Entry("Level", 16),
			T.Entry("Context", 24),
			T.Entry("Message", 33),
			T.Entry("Properties", 54),
		)
		It("should fail if Entry is invalid", func() {
			w := NewFailingWriter(0, testError)
			err := s.serialize(nil, w)
			Expect(err).To(Equal(ErrInvalidEntry))
		})
	})
	Describe("Serialize", func() {
		It("should serialize message with all elements", func() {
			s.baseTime = time.Unix(1234567890, 0)
			entry := &Entry{
				Level:     ErrLevel,
				Message:   "message",
				Timestamp: s.baseTime.Add(time.Duration(11) * time.Minute),
				CallContext: &CallContext{
					Path:     "p",
					File:     "f",
					Line:     98765,
					Package:  "a",
					Type:     "t",
					Function: "u",
				},
				Properties: Properties{
					"name": "Alice",
					"hash": "#$%%@",
					"male": true,
				},
			}
			byt, err := s.Serialize(entry)
			Expect(err).NotTo(HaveOccurred())
			expected := "[660.000000] [" + red + bold + "ERR" + off + "] [" + blue + bold + "f" +
				off + ":" + yellow + "98765" + off + "] message {" + propkey + "hash" + off +
				`:"#$%%@";` + propkey + "male" + off + ":true;" + propkey + "name" + off +
				":Alice;}"
			Expect(string(byt)).To(Equal(expected))
		})
		It("should fail if Entry is invalid", func() {
			byt, err := s.Serialize(nil)
			Expect(err).To(Equal(ErrInvalidEntry))
			Expect(byt).To(BeNil())
		})
	})
})
