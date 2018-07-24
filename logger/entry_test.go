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
		entry = &Entry{Logger: L}
		L.SetThreshold(DebugLevel)
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("process", func() {
		It("should set level and message and pass entry to Logger's backends", func() {
			mf.EXPECT().Verify(entry).DoAndReturn(func(entry *Entry) (bool, error) {
				Expect(entry.Level).To(Equal(WarningLevel))
				Expect(entry.Message).To(Equal(testMessage))
				return false, nil
			})
			entry.process(WarningLevel, testMessage)
		})
		It("should not set anything and return if level doesn't pass threshold", func() {
			L.SetThreshold(ErrLevel)
			entry.process(WarningLevel, testMessage)
			Expect(entry.Level).To(BeZero())
			Expect(entry.Message).To(BeZero())
		})
	})
	Describe("Log", func() {
		T.DescribeTable("should properly build log message and pass to backends",
			func(level Level) {
				mf.EXPECT().Verify(entry).Return(false, nil)
				entry.Log(level, testMessage, anotherTestMessage)
				Expect(entry.Level).To(Equal(level))
				Expect(entry.Message).To(Equal(testMessage + anotherTestMessage))
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
				testedFunction(entry, testMessage, anotherTestMessage)
				Expect(entry.Level).To(Equal(level))
				Expect(entry.Message).To(Equal(testMessage + anotherTestMessage))
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
				entry.Logf(level, format, testMessage, anotherTestMessage)
				Expect(entry.Level).To(Equal(level))
				Expect(entry.Message).To(Equal(expectedMessage))
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
				testedFunction(entry, format, testMessage, anotherTestMessage)
				Expect(entry.Level).To(Equal(level))
				Expect(entry.Message).To(Equal(expectedMessage))
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
})
