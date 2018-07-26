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
	"encoding/json"
	"time"
)

const (
	// DefaultSerializerJSONTimestampFormat is the default date and time format.
	DefaultSerializerJSONTimestampFormat = time.RFC3339
)

type serializerJSONRecord struct {
	Level      string     `json:"level"`
	Message    string     `json:"message"`
	Timestamp  string     `json:"timestamp"`
	Properties Properties `json:"properties,omitempty"`
}

// SerializerJSON serializes entry to JSON format.
type SerializerJSON struct {
	// TimestampFormat defines format for displaying date and time.
	// See https://godoc.org/time#Time.Format description for details.
	TimestampFormat string
}

// NewSerializerJSON creates and returns a new default SerializerJSON object.
func NewSerializerJSON() *SerializerJSON {
	return &SerializerJSON{
		TimestampFormat: DefaultSerializerJSONTimestampFormat,
	}
}

// Serialize marshals entry to JSON. It implements loggers' Serializer interface.
func (s *SerializerJSON) Serialize(entry *Entry) ([]byte, error) {
	format := s.TimestampFormat
	if format == "" {
		format = DefaultSerializerJSONTimestampFormat
	}
	record := serializerJSONRecord{
		Level:      entry.Level.String(),
		Message:    entry.Message,
		Timestamp:  entry.Timestamp.UTC().Format(format),
		Properties: entry.Properties,
	}
	return json.Marshal(record)
}
