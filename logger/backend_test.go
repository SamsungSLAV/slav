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
	. "github.com/onsi/gomega"
)

var _ = Describe("Backend", func() {
	const (
		backendName = string("backendName")
	)
	var (
		ctrl      *gomock.Controller
		mf        *MockFilter
		ms        *MockSerializer
		mw        *MockWriter
		mb        Backend
		e         *Entry
		buf       []byte
		testError error
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
		e = &Entry{
			Level:   AlertLevel,
			Message: "message",
		}
		buf = []byte("Lorem ipsum")
		testError = errors.New("Test Error")
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("process", func() {
		It("should pass filter, serialize and write log message", func() {
			gomock.InOrder(
				mf.EXPECT().Verify(e).Return(true, nil),
				ms.EXPECT().Serialize(e).Return(buf, nil),
				mw.EXPECT().Write(e.Level, buf),
			)

			err := mb.process(e)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should pass filter, serialize and write log message but return error when writing fails", func() {
			gomock.InOrder(
				mf.EXPECT().Verify(e).Return(true, nil),
				ms.EXPECT().Serialize(e).Return(buf, nil),
				mw.EXPECT().Write(e.Level, buf).Return(0, testError),
			)

			err := mb.process(e)
			Expect(err).To(Equal(testError))
		})
		It("should pass filter and serialize but return error when serializing fails", func() {
			gomock.InOrder(
				mf.EXPECT().Verify(e).Return(true, nil),
				ms.EXPECT().Serialize(e).Return(nil, testError),
			)

			err := mb.process(e)
			Expect(err).To(Equal(testError))
		})
		It("should not pass filter", func() {
			mf.EXPECT().Verify(e).Return(false, nil)

			err := mb.process(e)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should return error when filtering fails", func() {
			mf.EXPECT().Verify(e).Return(false, testError)

			err := mb.process(e)
			Expect(err).To(Equal(testError))
		})
	})
})
