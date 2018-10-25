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

var _ = Describe("Level", func() {
	T.DescribeTable("should properly stringify log level",
		func(level Level, expected string) {
			Expect(level.String()).To(Equal(expected))
		},
		T.Entry("EmergLevel", EmergLevel, EmergLevelStr),
		T.Entry("AlertLevel", AlertLevel, AlertLevelStr),
		T.Entry("CritLevel", CritLevel, CritLevelStr),
		T.Entry("ErrLevel", ErrLevel, ErrLevelStr),
		T.Entry("WarningLevel", WarningLevel, WarningLevelStr),
		T.Entry("NoticeLevel", NoticeLevel, NoticeLevelStr),
		T.Entry("InfoLevel", InfoLevel, InfoLevelStr),
		T.Entry("DebugLevel", DebugLevel, DebugLevelStr),
		T.Entry("Unknown level", Level(0xBADC0DE), UnknownLevelStr),
	)
	Describe("IsValid", func() {
		T.DescribeTable("should treat known log level as valid",
			func(level Level) {
				Expect(level.IsValid()).To(BeTrue())
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
		It("should treat unknown log level as invalid", func() {
			badLevel := Level(0xBADC0DE)
			Expect(badLevel.IsValid()).To(BeFalse())
		})
	})
	Describe("StringToLevel", func() {
		T.DescribeTable("should return proper log levels",
			func(input string, expected Level) {
				l, err := StringToLevel(input)
				Expect(err).NotTo(HaveOccurred())
				Expect(l).To(Equal(expected))
			},
			T.Entry("emergency->EmergLevel", EmergLevelStr, EmergLevel),
			T.Entry("alert->AlertLevel", AlertLevelStr, AlertLevel),
			T.Entry("critical->CritLevel", CritLevelStr, CritLevel),
			T.Entry("error->ErrLevel", ErrLevelStr, ErrLevel),
			T.Entry("warning->WarningLevel", WarningLevelStr, WarningLevel),
			T.Entry("notice->NoticeLevel", NoticeLevelStr, NoticeLevel),
			T.Entry("info->InfoLevel", InfoLevelStr, InfoLevel),
			T.Entry("debug->DebugLevel", DebugLevelStr, DebugLevel),
		)
		It("should return invalid log level and ErrInvalidLogLevel on unmatched input", func() {
			badInput := "!#@mam12"
			l, err := StringToLevel(badInput)
			Expect(err).To(Equal(ErrInvalidLogLevel))
			Expect(l.IsValid()).To(BeFalse())
		})
	})
})
