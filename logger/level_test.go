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
		T.Entry("EmergLevel", EmergLevel, "emergency"),
		T.Entry("AlertLevel", AlertLevel, "alert"),
		T.Entry("CritLevel", CritLevel, "critical"),
		T.Entry("ErrLevel", ErrLevel, "error"),
		T.Entry("WarningLevel", WarningLevel, "warning"),
		T.Entry("NoticeLevel", NoticeLevel, "notice"),
		T.Entry("InfoLevel", InfoLevel, "info"),
		T.Entry("DebugLevel", DebugLevel, "debug"),
		T.Entry("Unknown level", Level(0xBADC0DE), "unknown"),
	)
})
