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
	"io/ioutil"
	"os"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	T "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func withStderrMocked(testFunction func()) string {
	r, w, _ := os.Pipe()

	tmp := os.Stderr
	defer func() {
		os.Stderr = tmp
	}()
	os.Stderr = w

	go func() {
		testFunction()
		w.Close()
	}()

	buffer, _ := ioutil.ReadAll(r)
	return string(buffer)
}

var _ = Describe("Logger", func() {
	const (
		backendName        = string("backendName")
		anotherBackendName = string("anotherBackendName")
	)
	var (
		ctrl    *gomock.Controller
		mf, amf *MockFilter
		ms, ams *MockSerializer
		mw, amw *MockWriter
		mb, amb Backend

		L *Logger
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
		amf = NewMockFilter(ctrl)
		ams = NewMockSerializer(ctrl)
		amw = NewMockWriter(ctrl)
		amb = Backend{
			Filter:     amf,
			Serializer: ams,
			Writer:     amw,
		}

		L = NewLogger()
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("NewLogger", func() {
		It("should create a new default logger object", func() {
			Expect(L).NotTo(BeNil())
			Expect(L.mutex).NotTo(BeNil())
			L.mutex.Lock()
			defer L.mutex.Unlock()
			Expect(L.threshold).To(Equal(DefaultThreshold))
			Expect(L.backends).To(BeEmpty())
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
	Describe("Backend", func() {
		type Backends map[string]Backend
		expectBackends := func(expected Backends) {
			L.mutex.Lock()
			defer L.mutex.Unlock()
			Expect(L.backends).To(HaveLen(len(expected)))
			for k, v := range expected {
				v.Logger = L
				Expect(L.backends).To(HaveKeyWithValue(k, v))
			}
		}
		BeforeEach(func() {
			amb.Writer = nil // To differentiate mb and amb.
			expectBackends(Backends{})
		})
		Describe("AddBackend", func() {
			It("should set backend in Logger", func() {
				L.AddBackend(backendName, mb)
				expectBackends(Backends{backendName: mb})
			})
			It("should update backend in Logger", func() {
				L.AddBackend(backendName, mb)
				L.AddBackend(backendName, amb)
				expectBackends(Backends{backendName: amb})
			})
			It("should add multiple backends in Logger", func() {
				L.AddBackend(backendName, mb)
				L.AddBackend(anotherBackendName, amb)
				expectBackends(Backends{backendName: mb, anotherBackendName: amb})
			})
		})
		Describe("RemoveBackend", func() {
			BeforeEach(func() {
				L.AddBackend(backendName, mb)
			})
			It("should remove backend from Logger", func() {
				err := L.RemoveBackend(backendName)
				Expect(err).NotTo(HaveOccurred())
				expectBackends(Backends{})
			})
			It("should fail to remove nonexisting backend from Logger", func() {
				err := L.RemoveBackend(anotherBackendName)
				Expect(err).To(Equal(ErrInvalidBackendName))
				expectBackends(Backends{backendName: mb})
			})
			It("should fail to remove backend from Logger 2nd time in a row", func() {
				err := L.RemoveBackend(backendName)
				Expect(err).NotTo(HaveOccurred())
				expectBackends(Backends{})

				err = L.RemoveBackend(backendName)
				Expect(err).To(Equal(ErrInvalidBackendName))
				expectBackends(Backends{})
			})
		})
		Describe("RemoveAllBackends", func() {
			It("should handle empty backends case", func() {
				L.RemoveAllBackends()
				expectBackends(Backends{})
			})
			It("should handle single backend case", func() {
				L.AddBackend(backendName, mb)

				L.RemoveAllBackends()
				expectBackends(Backends{})
			})
			It("should handle multiple backends case", func() {
				L.AddBackend(backendName, mb)
				L.AddBackend(anotherBackendName, amb)

				L.RemoveAllBackends()
				expectBackends(Backends{})
			})
		})
	})
	Describe("process", func() {
		entry := &Entry{
			Level:   ErrLevel,
			Message: "Message",
		}
		buf := []byte("Lorem ipsum")
		testError := errors.New("Test Error")
		buildError := func(err error, name string) string {
			return "Error <" + err.Error() + "> printing log message to <" + name + "> backend.\n"
		}
		BeforeEach(func() {
			L.AddBackend(backendName, mb)
			L.AddBackend(anotherBackendName, amb)
		})
		It("should log to all backends", func() {
			gomock.InOrder(
				mf.EXPECT().Verify(entry).Return(true, nil),
				ms.EXPECT().Serialize(entry).Return(buf, nil),
				mw.EXPECT().Write(entry.Level, buf),
			)
			gomock.InOrder(
				amf.EXPECT().Verify(entry).Return(true, nil),
				ams.EXPECT().Serialize(entry).Return(buf, nil),
				amw.EXPECT().Write(entry.Level, buf),
			)
			L.process(entry)
		})
		It("should print errors to stderr", func() {
			mf.EXPECT().Verify(entry).Return(false, testError)
			gomock.InOrder(
				amf.EXPECT().Verify(entry).Return(true, nil),
				ams.EXPECT().Serialize(entry).Return([]byte{}, testError),
			)
			stderr := withStderrMocked(func() {
				L.process(entry)
			})
			Expect(stderr).To(ContainSubstring(buildError(testError, backendName)))
			Expect(stderr).To(ContainSubstring(buildError(testError, anotherBackendName)))
		})
	})
})
