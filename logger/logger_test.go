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
	. "github.com/onsi/ginkgo"
	T "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger", func() {
	var L *Logger

	BeforeEach(func() {
		L = NewLogger()
	})

	Describe("NewLogger", func() {
		It("should create a new default logger object", func() {
			Expect(L).NotTo(BeNil())
			Expect(L.mutex).NotTo(BeNil())
			L.mutex.Lock()
			defer L.mutex.Unlock()
			Expect(L.threshold).To(Equal(DefaultThreshold))
		})
	})
	Describe("Threshold", func() {
		Describe("SetThreshold", func() {
			T.DescribeTable("should set valid log level",
				func(level Level) {
					err := L.SetThreshold(level)
					Expect(err).NotTo(HaveOccurred())
					L.mutex.Lock()
					defer L.mutex.Unlock()
					Expect(L.threshold).To(Equal(level))
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
				err := L.SetThreshold(badLevel)
				Expect(err).To(Equal(ErrInvalidLogLevel))
				Expect(L.threshold).To(Equal(DefaultThreshold))
			})
		})
		Describe("Threshold", func() {
			T.DescribeTable("should return log level",
				func(level Level) {
					err := L.SetThreshold(level)
					Expect(err).NotTo(HaveOccurred())

					retLevel := L.Threshold()
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
		Describe("PassThreshold", func() {
			It("should pass threshold for high priority levels", func() {
				levels := []Level{EmergLevel, AlertLevel, CritLevel, ErrLevel, WarningLevel, NoticeLevel, InfoLevel, DebugLevel, Level(0xBADC0DE)}
				for ti, thr := range levels[:len(levels)-1] {
					err := L.SetThreshold(thr)
					Expect(err).NotTo(HaveOccurred(), "setting threshold to %s", thr)

					for _, lvl := range levels[:ti] {
						Expect(L.PassThreshold(lvl)).To(BeTrue(), "threshold %s; level %s", thr, lvl)
					}
					for _, lvl := range levels[ti+1:] {
						Expect(L.PassThreshold(lvl)).To(BeFalse(), "threshold %s; level %s", thr, lvl)
					}
				}
			})
		})
	})
})
