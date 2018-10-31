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

var _ = Describe("WriterStderr", func() {
	const (
		testMsg  = "testMessage"
		anyLevel = CritLevel
	)

	var w *WriterStderr

	BeforeEach(func() {
		w = NewWriterStderr()
	})
	Describe("NewWriterStderr", func() {
		It("should create a new empty object", func() {
			Expect(w).NotTo(BeNil())
		})
	})
	Describe("Write", func() {
		It("should write to stderr", func() {
			stderr := withStderrMocked(func() {
				n, err := w.Write(anyLevel, []byte(testMsg))
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(len(testMsg) + 1))
			})
			Expect(stderr).To(Equal(testMsg + "\n"))
		})
	})
})
