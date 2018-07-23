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
	. "github.com/onsi/gomega"
)

var _ = Describe("Entry", func() {
	const (
		backendName = string("backendName")
		testMessage = string("Test Message")
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
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("process", func() {
		It("should set level and message and pass entry to Logger's backends", func() {
			L.SetThreshold(DebugLevel)
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
})
