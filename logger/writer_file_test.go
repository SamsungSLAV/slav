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
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WriterFile", func() {
	const (
		testFile       = "/tmp/test_slav_logger_writerfile.txt"
		impossibleFile = "/it/should/be/impossible/to/create/this/file"
		testMsg        = "testMessage"
		anotherTestMsg = "anotherTestMessage"
		anyLevel       = InfoLevel
		filePerm       = os.FileMode(0600)
	)

	expectFileToContain := func(file, text string) {
		text = text + "\n"
		f, err := os.Open(file)
		Expect(err).NotTo(HaveOccurred())
		defer f.Close()

		buf := make([]byte, 64)
		n, err := f.Read(buf)
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(len(text)))
		Expect(string(buf[0:n])).To(Equal(text))
	}

	BeforeEach(func() {
		// Ensure there is no test file.
		os.Remove(testFile) // Error ignored.
		// Verify, there is no file.
		_, err := os.Stat(testFile)
		Expect(err).To(HaveOccurred())
	})
	AfterEach(func() {
		// Remove file if created during test.
		os.Remove(testFile) // Error ignored.
		// Verify, there is no file.
		_, err := os.Stat(testFile)
		Expect(err).To(HaveOccurred())
	})
	Describe("NewWriterFile", func() {
		It("should create a file and a new empty object", func() {
			w := NewWriterFile(testFile, filePerm)
			Expect(w).NotTo(BeNil())

			_, err := os.Stat(testFile)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should panic if the file cannot be created", func() {
			Expect(func() {
				NewWriterFile(impossibleFile, filePerm)
			}).To(Panic())
		})
	})
	Describe("Write", func() {
		It("should write to newly created file", func() {
			w := NewWriterFile(testFile, filePerm)
			Expect(w).NotTo(BeNil())

			n, err := w.Write(anyLevel, []byte(testMsg))
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(len(testMsg) + 1))
			expectFileToContain(testFile, testMsg)
		})
		It("should append to existing file", func() {
			{ //Write message to newly created file
				p := NewWriterFile(testFile, filePerm)
				Expect(p).NotTo(BeNil())
				n, err := p.Write(anyLevel, []byte(testMsg))
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(len(testMsg) + 1))
				expectFileToContain(testFile, testMsg)
			}
			// Open the file again for appending.
			w := NewWriterFile(testFile, filePerm)
			Expect(w).NotTo(BeNil())

			n, err := w.Write(anyLevel, []byte(anotherTestMsg))
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(len(anotherTestMsg) + 1))
			expectFileToContain(testFile, testMsg+"\n"+anotherTestMsg)
		})
	})
})
