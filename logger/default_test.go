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
	"errors"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	T "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Default", func() {
	const (
		backendName        = string("backendName")
		anotherBackendName = string("anotherBackendName")
		testMessage        = string("Test Message")
		anotherTestMessage = string("Another Test Message")
		format             = string("%s >>> %s")
		expectedMessage    = string(testMessage + " >>> " + anotherTestMessage)
	)

	It("defaultLogger should be already initialized", func() {
		Expect(defaultLogger).NotTo(BeNil())
		Expect(defaultLogger.mutex).NotTo(BeNil())
		Expect(defaultLogger.threshold).To(Equal(DefaultThreshold))
		defaultLogger.mutex.Lock()
		defer defaultLogger.mutex.Unlock()
		Expect(defaultLogger.backends).To(HaveKey("default"))
	})
	It("defaultLogger should log to stderr in default text format", func() {
		stderr := withStderrMocked(func() {
			defaultLogger.WithProperties(Properties{
				"name": "Alice",
				"hash": "#$%%@",
				"male": true,
				"age":  37,
				"skills": Properties{
					"coding": 7,
				},
				"issues": "",
			}).Emergencyf("test <%s>:%d log", "string", 67)
		})
		expected := "[" + red + bold + invert + "EME" + off + `] "test <string>:67 log" {` +
			propkey + "age" + off + ":37;" + propkey + "hash" + off + `:"#$%%@";` + propkey +
			"issues" + off + `:"";` + propkey + "male" + off + ":true;" + propkey + "name" +
			off + ":Alice;" + propkey + "skills" + off + `:"map[coding:7]";}`
		// ContainSubstring used instead of Equal, because the timestamps are hard to be compared.
		Expect(stderr).To(ContainSubstring(expected))
	})
	Describe("with default logger substituted", func() {
		var (
			ctrl    *gomock.Controller
			mf, amf *MockFilter
			ms, ams *MockSerializer
			mw, amw *MockWriter
			mb, amb Backend

			L *Logger
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mf = NewMockFilter(ctrl)
			ms = NewMockSerializer(ctrl)
			mw = NewMockWriter(ctrl)
			mb = Backend{
				Filter:     mf,
				Serializer: ms,
				Writer:     mw,
			}
			amf = NewMockFilter(ctrl)
			ams = NewMockSerializer(ctrl)
			amw = NewMockWriter(ctrl)
			amb = Backend{
				Filter:     amf,
				Serializer: ams,
				Writer:     amw,
			}

			L = NewLogger()
			SetDefault(L)
		})
		AfterEach(func() {
			ctrl.Finish()
		})

		Describe("SetDefault", func() {
			It("should substitute default logger", func() {
				Expect(defaultLogger).To(Equal(L))
			})
		})
		Describe("Threshold", func() {
			Describe("SetThreshold", func() {
				T.DescribeTable("should set valid log level",
					func(level Level) {
						err := SetThreshold(level)
						Expect(err).NotTo(HaveOccurred())
						Expect(defaultLogger.threshold).To(Equal(level))
					},
					T.Entry("EmergLevel", EmergLevel),
					T.Entry("AlertLevel", AlertLevel),
					T.Entry("CritLevel", CritLevel),
					T.Entry("ErrLevel", ErrLevel),
					T.Entry("WarningLevel", WarningLevel),
					T.Entry("NoticeLevel", NoticeLevel),
					T.Entry("InfoLevel", InfoLevel),
					T.Entry("DebugLevel", DebugLevel),
				)
				It("should fail to set invalid log level", func() {
					badLevel := Level(0xBADC0DE)
					err := SetThreshold(badLevel)
					Expect(err).To(Equal(ErrInvalidLogLevel))
					Expect(defaultLogger.threshold).To(Equal(DefaultThreshold))
				})
			})
			Describe("Threshold", func() {
				T.DescribeTable("should return log level",
					func(level Level) {
						err := SetThreshold(level)
						Expect(err).NotTo(HaveOccurred())

						retLevel := Threshold()
						Expect(retLevel).To(Equal(level))
					},
					T.Entry("EmergLevel", EmergLevel),
					T.Entry("AlertLevel", AlertLevel),
					T.Entry("CritLevel", CritLevel),
					T.Entry("ErrLevel", ErrLevel),
					T.Entry("WarningLevel", WarningLevel),
					T.Entry("NoticeLevel", NoticeLevel),
					T.Entry("InfoLevel", InfoLevel),
					T.Entry("DebugLevel", DebugLevel),
				)
			})
		})
		Describe("Backend", func() {
			type Backends map[string]Backend
			expectBackends := func(expected Backends) {
				defaultLogger.mutex.Lock()
				defer defaultLogger.mutex.Unlock()
				Expect(defaultLogger.backends).To(HaveLen(len(expected)))
				for k, v := range expected {
					v.Logger = defaultLogger
					Expect(defaultLogger.backends).To(HaveKeyWithValue(k, v))
				}
			}
			BeforeEach(func() {
				amb.Writer = nil // To differentiate mb and amb.
				expectBackends(Backends{})
			})
			Describe("AddBackend", func() {
				It("should set backend in default Logger", func() {
					AddBackend(backendName, mb)
					expectBackends(Backends{backendName: mb})
				})
				It("should update backend in default Logger", func() {
					AddBackend(backendName, mb)
					AddBackend(backendName, amb)
					expectBackends(Backends{backendName: amb})
				})
				It("should add multiple backends in default Logger", func() {
					AddBackend(backendName, mb)
					AddBackend(anotherBackendName, amb)
					expectBackends(Backends{backendName: mb, anotherBackendName: amb})
				})
			})
			Describe("RemoveBackend", func() {
				BeforeEach(func() {
					AddBackend(backendName, mb)
				})
				It("should remove backend from default Logger", func() {
					err := RemoveBackend(backendName)
					Expect(err).NotTo(HaveOccurred())
					expectBackends(Backends{})
				})
				It("should fail to remove nonexisting backend from default Logger", func() {
					err := RemoveBackend(anotherBackendName)
					Expect(err).To(Equal(ErrInvalidBackendName))
					expectBackends(Backends{backendName: mb})
				})
				It("should fail to remove backend from default Logger 2nd time in a row", func() {
					err := RemoveBackend(backendName)
					Expect(err).NotTo(HaveOccurred())

					err = RemoveBackend(backendName)
					Expect(err).To(Equal(ErrInvalidBackendName))
					expectBackends(Backends{})
				})
			})
			Describe("RemoveAllBackends", func() {
				It("should handle empty backends case", func() {
					RemoveAllBackends()
					expectBackends(Backends{})
				})
				It("should handle single backend case", func() {
					AddBackend(backendName, mb)

					RemoveAllBackends()
					expectBackends(Backends{})
				})
				It("should handle multiple backends case", func() {
					AddBackend(backendName, mb)
					AddBackend(anotherBackendName, amb)

					RemoveAllBackends()
					expectBackends(Backends{})
				})
			})
		})
		Describe("Log functions", func() {
			BeforeEach(func() {
				AddBackend(backendName, mb)
				SetThreshold(DebugLevel)
			})
			Describe("Log", func() {
				T.DescribeTable("should properly build log message and pass to defaults logger's"+
					" backend",
					func(level Level) {
						mf.EXPECT().Verify(gomock.Any()).
							DoAndReturn(func(entry *Entry) (bool, error) {
								Expect(entry.Level).To(Equal(level))
								Expect(entry.Message).To(Equal(testMessage + anotherTestMessage))
								return false, nil
							})
						Log(level, testMessage, anotherTestMessage)
					},
					T.Entry("EmergLevel", EmergLevel),
					T.Entry("AlertLevel", AlertLevel),
					T.Entry("CritLevel", CritLevel),
					T.Entry("ErrLevel", ErrLevel),
					T.Entry("WarningLevel", WarningLevel),
					T.Entry("NoticeLevel", NoticeLevel),
					T.Entry("InfoLevel", InfoLevel),
					T.Entry("DebugLevel", DebugLevel),
				)
				T.DescribeTable("should properly set level and log message and pass to default"+
					" logger's backend",
					func(level Level, testedFunction func(...interface{})) {
						mf.EXPECT().Verify(gomock.Any()).
							DoAndReturn(func(entry *Entry) (bool, error) {
								Expect(entry.Level).To(Equal(level))
								Expect(entry.Message).To(Equal(testMessage + anotherTestMessage))
								return false, nil
							})
						testedFunction(testMessage, anotherTestMessage)
					},
					T.Entry("EmergLevel", EmergLevel, Emergency),
					T.Entry("AlertLevel", AlertLevel, Alert),
					T.Entry("CritLevel", CritLevel, Critical),
					T.Entry("ErrLevel", ErrLevel, Error),
					T.Entry("WarningLevel", WarningLevel, Warning),
					T.Entry("NoticeLevel", NoticeLevel, Notice),
					T.Entry("InfoLevel", InfoLevel, Info),
					T.Entry("DebugLevel", DebugLevel, Debug),
				)
			})
			Describe("Logf", func() {
				T.DescribeTable("should properly build log message and pass to default logger's"+
					" backend",
					func(level Level) {
						mf.EXPECT().Verify(gomock.Any()).
							DoAndReturn(func(entry *Entry) (bool, error) {
								Expect(entry.Level).To(Equal(level))
								Expect(entry.Message).To(Equal(expectedMessage))
								return false, nil
							})
						Logf(level, format, testMessage, anotherTestMessage)
					},
					T.Entry("EmergLevel", EmergLevel),
					T.Entry("AlertLevel", AlertLevel),
					T.Entry("CritLevel", CritLevel),
					T.Entry("ErrLevel", ErrLevel),
					T.Entry("WarningLevel", WarningLevel),
					T.Entry("NoticeLevel", NoticeLevel),
					T.Entry("InfoLevel", InfoLevel),
					T.Entry("DebugLevel", DebugLevel),
				)
				T.DescribeTable("should properly set level and log message and pass to default"+
					" logger's backend",
					func(level Level, testedFunction func(string, ...interface{})) {
						mf.EXPECT().Verify(gomock.Any()).
							DoAndReturn(func(entry *Entry) (bool, error) {
								Expect(entry.Level).To(Equal(level))
								Expect(entry.Message).To(Equal(expectedMessage))
								return false, nil
							})
						testedFunction(format, testMessage, anotherTestMessage)
					},
					T.Entry("EmergLevel", EmergLevel, Emergencyf),
					T.Entry("AlertLevel", AlertLevel, Alertf),
					T.Entry("CritLevel", CritLevel, Criticalf),
					T.Entry("ErrLevel", ErrLevel, Errorf),
					T.Entry("WarningLevel", WarningLevel, Warningf),
					T.Entry("NoticeLevel", NoticeLevel, Noticef),
					T.Entry("InfoLevel", InfoLevel, Infof),
					T.Entry("DebugLevel", DebugLevel, Debugf),
				)
			})
		})
		Describe("Properties", func() {
			const (
				property        = "property"
				value           = "value"
				anotherProperty = "another property"
				anotherValue    = "another value"
			)
			var (
				errorValue = errors.New("error value")
			)
			It("should create a new log message with a property", func() {
				entry := WithProperty(property, value)
				Expect(entry.Properties).To(HaveLen(1))
				Expect(entry.Properties).To(HaveKeyWithValue(property, value))
			})
			It("should create a new log message with multiple properties", func() {
				entry := WithProperties(Properties{
					property:        value,
					anotherProperty: anotherValue,
				})
				Expect(entry.Properties).To(HaveLen(2))
				Expect(entry.Properties).To(HaveKeyWithValue(property, value))
				Expect(entry.Properties).To(HaveKeyWithValue(anotherProperty, anotherValue))
			})
			It("should create a new log message with an error property", func() {
				entry := WithError(errorValue)
				Expect(entry.Properties).To(HaveLen(1))
				Expect(entry.Properties).To(HaveKeyWithValue(ErrorProperty, errorValue.Error()))
			})
		})
	})
})
