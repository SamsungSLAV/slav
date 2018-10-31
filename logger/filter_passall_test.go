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
	. "github.com/onsi/gomega"
)

var _ = Describe("FilterPassAll", func() {
	var (
		f *FilterPassAll
	)

	BeforeEach(func() {
		f = NewFilterPassAll()
	})
	Describe("NewFilterPassAll", func() {
		It("should create a new empty object", func() {
			Expect(f).NotTo(BeNil())
		})
	})
	Describe("Verify", func() {
		It("should return true for any Entry", func() {
			any := &Entry{
				Level:   EmergLevel,
				Message: "AnyMessage",
			}
			ret, err := f.Verify(any)
			Expect(err).NotTo(HaveOccurred())
			Expect(ret).To(BeTrue())
		})
	})
})
