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

// FilterPassAll is a dummy implementation of Filter interface which accepts all entries.
type FilterPassAll struct{}

// NewFilterPassAll creates and returns a new FilterPassAll object.
func NewFilterPassAll() *FilterPassAll {
	return &FilterPassAll{}
}

// Verify accepts all entries and returns true.
// It implements Filter interface in FilterPassAll type.
func (*FilterPassAll) Verify(*Entry) (bool, error) {
	return true, nil
}
