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
	"math"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SerializerJSON", func() {
	var (
		e *Entry
		s *SerializerJSON
	)

	BeforeEach(func() {
		e = &Entry{
			Level:     ErrLevel,
			Message:   "message",
			Timestamp: time.Unix(1234567890, 0),
			CallContext: &CallContext{
				Path:     "somePath",
				File:     "someFile",
				Line:     1234567,
				Package:  "somePackage",
				Type:     "someType",
				Function: "someFunction",
			},
		}
		s = NewSerializerJSON()
	})
	Describe("NewSerializerJSON", func() {
		It("should create a new object with default configuration", func() {
			Expect(s).NotTo(BeNil())
			Expect(s.TimestampFormat).To(Equal(DefaultSerializerJSONTimestampFormat))
		})
	})
	Describe("Serialize", func() {
		It("should serialize simple message without properties", func() {
			buf, err := s.Serialize(e)
			Expect(err).NotTo(HaveOccurred())
			expected := []byte(`{"level":"error","message":"message","timestamp":` +
				`"2009-02-13T23:31:30Z","callcontext":{"path":"somePath","file":"someFile",` +
				`"line":1234567,"package":"somePackage","type":"someType","function":` +
				`"someFunction"}}`)
			Expect(buf).To(Equal(expected))
		})
		It("should serialize message with properties", func() {
			e.WithProperties(Properties{
				"name": "Alice",
				"male": true,
				"age":  37,
				"skills": Properties{
					"cooking": "expert",
					"coding":  7,
				},
				"issues": "",
			})
			buf, err := s.Serialize(e)
			Expect(err).NotTo(HaveOccurred())
			expected := []byte(`{"level":"error","message":"message","timestamp":` +
				`"2009-02-13T23:31:30Z","callcontext":{"path":"somePath","file":"someFile",` +
				`"line":1234567,"package":"somePackage","type":"someType","function":` +
				`"someFunction"},"properties":{"age":37,"issues":"","male":true,"name":"Alice",` +
				`"skills":{"coding":7,"cooking":"expert"}}}`)
			Expect(buf).To(Equal(expected))
		})
		It("should serialize message with error", func() {
			err := errors.New("test error")
			e.WithError(err)
			buf, err := s.Serialize(e)
			Expect(err).NotTo(HaveOccurred())
			expected := []byte(`{"level":"error","message":"message","timestamp":` +
				`"2009-02-13T23:31:30Z","callcontext":{"path":"somePath","file":` +
				`"someFile","line":1234567,"package":"somePackage","type":"someType",` +
				`"function":"someFunction"},"properties":{"error":"test error"}}`)
			Expect(buf).To(Equal(expected))
		})
		It("should return error if serialization is not possible", func() {
			e.WithProperties(Properties{
				"power": math.Inf(1),
			})
			buf, err := s.Serialize(e)
			Expect(err.Error()).To(Equal("json: unsupported value: +Inf"))
			Expect(buf).To(BeNil())
		})
		It("should serialize with default timestamp format if format is not initialized", func() {
			s = &SerializerJSON{}
			buf, err := s.Serialize(e)
			Expect(err).NotTo(HaveOccurred())
			expected := []byte(`{"level":"error","message":"message","timestamp":` +
				`"2009-02-13T23:31:30Z","callcontext":{"path":"somePath","file":"someFile",` +
				`"line":1234567,"package":"somePackage","type":"someType","function":` +
				`"someFunction"}}`)
			Expect(buf).To(Equal(expected))
		})
		It("should serialize with different timestamp format", func() {
			s.TimestampFormat = time.ANSIC
			buf, err := s.Serialize(e)
			Expect(err).NotTo(HaveOccurred())
			expected := []byte(`{"level":"error","message":"message","timestamp":` +
				`"Fri Feb 13 23:31:30 2009","callcontext":{"path":"somePath","file":"someFile",` +
				`"line":1234567,"package":"somePackage","type":"someType","function":` +
				`"someFunction"}}`)
			Expect(buf).To(Equal(expected))
		})
		It("should serialize if context's type is missing", func() {
			e.CallContext.Type = ""
			buf, err := s.Serialize(e)
			Expect(err).NotTo(HaveOccurred())
			expected := []byte(`{"level":"error","message":"message","timestamp":` +
				`"2009-02-13T23:31:30Z","callcontext":{"path":"somePath","file":"someFile",` +
				`"line":1234567,"package":"somePackage","function":"someFunction"}}`)
			Expect(buf).To(Equal(expected))
		})
		It("should serialize if there is no context at all", func() {
			e.CallContext = nil
			buf, err := s.Serialize(e)
			Expect(err).NotTo(HaveOccurred())
			expected := []byte(`{"level":"error","message":"message","timestamp":` +
				`"2009-02-13T23:31:30Z"}`)
			Expect(buf).To(Equal(expected))
		})
	})
})
