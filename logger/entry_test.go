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
	"fmt"
	"runtime"
	"time"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	T "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Entry", func() {
	const (
		backendName        = string("backendName")
		testMessage        = string("Test Message")
		anotherTestMessage = string("Another Test Message")
		format             = string("%s >>> %s")
		expectedMessage    = string(testMessage + " >>> " + anotherTestMessage)
		thisFile           = string("entry_test.go")
		thisPackage        = string("git.tizen.org/tools/slav/logger")
	)
	var (
		ctrl  *gomock.Controller
		mf    *MockFilter
		ms    *MockSerializer
		mw    *MockWriter
		mb    Backend
		L     *Logger
		entry *Entry
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

		L = NewLogger()
		L.AddBackend(backendName, mb)
		entry = &Entry{
			Logger:     L,
			Properties: make(Properties),
		}
		L.SetThreshold(DebugLevel)
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("process", func() {
		It("should set level, message and timestamp, and pass entry to Logger's backends", func() {
			before := time.Now()
			mf.EXPECT().Verify(entry).DoAndReturn(func(entry *Entry) (bool, error) {
				_, _, line, _ := runtime.Caller(0)
				Expect(entry.Level).To(Equal(WarningLevel))
				Expect(entry.Message).To(Equal(testMessage))
				Expect(entry.Timestamp).To(BeTemporally(">", before))
				Expect(entry.Timestamp).To(BeTemporally("<", time.Now()))
				Expect(entry.CallContext).NotTo(BeNil())
				Expect(entry.CallContext.Line).To(Equal(line + 11))
				Expect(entry.CallContext.File).To(Equal(thisFile))
				Expect(entry.CallContext.Package).To(Equal(thisPackage))
				return false, nil
			})
			entry.process(WarningLevel, testMessage)
		})
		It("should not set CallContext it it cannot be get", func() {
			before := time.Now()
			mf.EXPECT().Verify(entry).DoAndReturn(func(entry *Entry) (bool, error) {
				Expect(entry.Level).To(Equal(WarningLevel))
				Expect(entry.Message).To(Equal(testMessage))
				Expect(entry.Timestamp).To(BeTemporally(">", before))
				Expect(entry.Timestamp).To(BeTemporally("<", time.Now()))
				Expect(entry.CallContext).To(BeNil())
				return false, nil
			})
			// Ask for context deeper than call stack depth.
			entry.IncDepth(100000)
			entry.process(WarningLevel, testMessage)
		})
		It("should not set anything and return if level doesn't pass threshold", func() {
			L.SetThreshold(ErrLevel)
			entry.process(WarningLevel, testMessage)
			Expect(entry.Level).To(BeZero())
			Expect(entry.Message).To(BeZero())
			Expect(entry.Timestamp).To(BeZero())
			Expect(entry.CallContext).To(BeNil())
		})
	})
	Describe("Log", func() {
		T.DescribeTable("should properly build log message and pass to backends",
			func(level Level) {
				mf.EXPECT().Verify(entry).Return(false, nil)
				_, _, line, _ := runtime.Caller(0)
				entry.Log(level, testMessage, anotherTestMessage)
				Expect(entry.Level).To(Equal(level))
				Expect(entry.Message).To(Equal(testMessage + anotherTestMessage))
				Expect(entry.CallContext).NotTo(BeNil())
				Expect(entry.CallContext.Line).To(Equal(line + 1))
				Expect(entry.CallContext.File).To(Equal(thisFile))
				Expect(entry.CallContext.Package).To(Equal(thisPackage))
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
		T.DescribeTable("should properly set level and log message and pass to logger's backend",
			func(level Level, testedFunction func(*Entry, ...interface{})) {
				mf.EXPECT().Verify(entry).Return(false, nil)
				_, _, line, _ := runtime.Caller(0)
				testedFunction(entry, testMessage, anotherTestMessage)
				Expect(entry.Level).To(Equal(level))
				Expect(entry.Message).To(Equal(testMessage + anotherTestMessage))
				Expect(entry.CallContext).NotTo(BeNil())
				Expect(entry.CallContext.Line).To(Equal(line + 1))
				Expect(entry.CallContext.File).To(Equal(thisFile))
				Expect(entry.CallContext.Package).To(Equal(thisPackage))
			},
			T.Entry("EmergLevel", EmergLevel, (*Entry).Emergency),
			T.Entry("AlertLevel", AlertLevel, (*Entry).Alert),
			T.Entry("CritLevel", CritLevel, (*Entry).Critical),
			T.Entry("ErrLevel", ErrLevel, (*Entry).Error),
			T.Entry("WarningLevel", WarningLevel, (*Entry).Warning),
			T.Entry("NoticeLevel", NoticeLevel, (*Entry).Notice),
			T.Entry("InfoLevel", InfoLevel, (*Entry).Info),
			T.Entry("DebugLevel", DebugLevel, (*Entry).Debug),
		)
	})
	Describe("Logf", func() {
		T.DescribeTable("should properly build log message and pass to logger's backend",
			func(level Level) {
				mf.EXPECT().Verify(entry).Return(false, nil)
				_, _, line, _ := runtime.Caller(0)
				entry.Logf(level, format, testMessage, anotherTestMessage)
				Expect(entry.Level).To(Equal(level))
				Expect(entry.Message).To(Equal(expectedMessage))
				Expect(entry.CallContext).NotTo(BeNil())
				Expect(entry.CallContext.Line).To(Equal(line + 1))
				Expect(entry.CallContext.File).To(Equal(thisFile))
				Expect(entry.CallContext.Package).To(Equal(thisPackage))
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
		T.DescribeTable("should properly set level and log message and pass to logger's backend",
			func(level Level, testedFunction func(*Entry, string, ...interface{})) {
				mf.EXPECT().Verify(entry).Return(false, nil)
				_, _, line, _ := runtime.Caller(0)
				testedFunction(entry, format, testMessage, anotherTestMessage)
				Expect(entry.Level).To(Equal(level))
				Expect(entry.Message).To(Equal(expectedMessage))
				Expect(entry.CallContext).NotTo(BeNil())
				Expect(entry.CallContext.Line).To(Equal(line + 1))
				Expect(entry.CallContext.File).To(Equal(thisFile))
				Expect(entry.CallContext.Package).To(Equal(thisPackage))
			},
			T.Entry("EmergLevel", EmergLevel, (*Entry).Emergencyf),
			T.Entry("AlertLevel", AlertLevel, (*Entry).Alertf),
			T.Entry("CritLevel", CritLevel, (*Entry).Criticalf),
			T.Entry("ErrLevel", ErrLevel, (*Entry).Errorf),
			T.Entry("WarningLevel", WarningLevel, (*Entry).Warningf),
			T.Entry("NoticeLevel", NoticeLevel, (*Entry).Noticef),
			T.Entry("InfoLevel", InfoLevel, (*Entry).Infof),
			T.Entry("DebugLevel", DebugLevel, (*Entry).Debugf),
		)
	})
	Describe("Properties", func() {
		const (
			property                  = "property"
			value                     = "value"
			anotherProperty           = "another property"
			anotherValue              = "another value"
			completelyAnotherProperty = "completely another property"
			completelyAnotherValue    = "completely another value"
		)
		var (
			errorValue        = errors.New("error value")
			anotherErrorValue = errors.New("another error value")
		)
		It("should add a property", func() {
			Expect(entry.Properties).To(BeEmpty())

			e := entry.WithProperty(property, value)
			Expect(e).To(Equal(entry))
			Expect(entry.Properties).To(HaveLen(1))
			Expect(entry.Properties).To(HaveKeyWithValue(property, value))
		})
		It("should overwrite a property", func() {
			Expect(entry.Properties).To(BeEmpty())
			entry.WithProperty(property, value)

			e := entry.WithProperty(property, anotherValue)
			Expect(e).To(Equal(entry))
			Expect(entry.Properties).To(HaveLen(1))
			Expect(entry.Properties).To(HaveKeyWithValue(property, anotherValue))
		})
		It("should add properties", func() {
			Expect(entry.Properties).To(BeEmpty())

			e := entry.WithProperties(Properties{
				property:        value,
				anotherProperty: anotherValue,
			})
			Expect(e).To(Equal(entry))
			Expect(entry.Properties).To(HaveLen(2))
			Expect(entry.Properties).To(HaveKeyWithValue(property, value))
			Expect(entry.Properties).To(HaveKeyWithValue(anotherProperty, anotherValue))
		})
		It("should update properties", func() {
			Expect(entry.Properties).To(BeEmpty())
			entry.WithProperties(Properties{
				property:        value,
				anotherProperty: anotherValue,
			})

			e := entry.WithProperties(Properties{
				property:                  anotherValue,
				completelyAnotherProperty: completelyAnotherValue,
			})
			Expect(e).To(Equal(entry))
			Expect(entry.Properties).To(HaveLen(3))
			// Updated by 2nd call.
			Expect(entry.Properties).To(HaveKeyWithValue(property, anotherValue))
			// Not changed since 1st call.
			Expect(entry.Properties).To(HaveKeyWithValue(anotherProperty, anotherValue))
			// Added by a 2nd call.
			Expect(entry.Properties).To(HaveKeyWithValue(completelyAnotherProperty,
				completelyAnotherValue))
		})
		It("should add an error property", func() {
			Expect(entry.Properties).To(BeEmpty())

			e := entry.WithError(errorValue)
			Expect(e).To(Equal(entry))
			Expect(entry.Properties).To(HaveLen(1))
			Expect(entry.Properties).To(HaveKeyWithValue(ErrorProperty, errorValue.Error()))
		})
		It("should overwrite an error property", func() {
			Expect(entry.Properties).To(BeEmpty())
			entry.WithError(errorValue)

			e := entry.WithError(anotherErrorValue)
			Expect(e).To(Equal(entry))
			Expect(entry.Properties).To(HaveLen(1))
			Expect(entry.Properties).To(HaveKeyWithValue(ErrorProperty, anotherErrorValue.Error()))
		})
		It("should create properties map if nil", func() {
			Expect(entry.Properties).NotTo(BeNil())
			entry.Properties = nil
			Expect(entry.Properties).To(BeNil())

			e := entry.WithProperties(Properties{
				property:        value,
				anotherProperty: anotherValue,
			})
			Expect(e).To(Equal(entry))
			Expect(entry.Properties).NotTo(BeNil())
			Expect(entry.Properties).To(HaveLen(2))
			Expect(entry.Properties).To(HaveKeyWithValue(property, value))
			Expect(entry.Properties).To(HaveKeyWithValue(anotherProperty, anotherValue))
		})
	})
	Describe("IncDepth", func() {
		const dep = 67
		It("should increase depth of entry call stack", func() {
			Expect(entry.depth).To(BeZero())

			for i := 0; i < 3; i++ {
				By(fmt.Sprintf("i = %d", i))
				e := entry.IncDepth(dep)
				Expect(e).To(Equal(entry))
				Expect(entry.depth).To(Equal((i + 1) * dep))
			}
		})
	})
})
