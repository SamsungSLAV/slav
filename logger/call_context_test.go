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
	"go/build"
	P "path"

	. "github.com/onsi/ginkgo"
	T "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func foo(depth int) *CallContext {
	return getCallContext(depth)
}

func bar(depth int) *CallContext {
	return foo(depth)
}

type CallContextTest struct{}

func (CallContextTest) byValue(depth int) *CallContext {
	return bar(depth)
}

func (c *CallContextTest) byPointer(depth int) *CallContext {
	return c.byValue(depth)
}

var _ = Describe("CallContext", func() {
	const (
		thisFile    = "call_context_test.go"
		thisPackage = "git.tizen.org/tools/slav/logger"
	)
	var (
		thisPath = P.Join(build.Default.GOPATH, "src", thisPackage) + "/"
	)

	Describe("getCallContext", func() {
		T.DescribeTable("should return proper context",
			func(depth, line int, typ, function string) {
				c := CallContextTest{}
				ctx := c.byPointer(depth)

				Expect(ctx).NotTo(BeNil())
				Expect(ctx.Path).To(Equal(thisPath))
				Expect(ctx.File).To(Equal(thisFile))
				Expect(ctx.Line).To(Equal(line))
				Expect(ctx.Package).To(Equal(thisPackage))
				Expect(ctx.Type).To(Equal(typ))
				Expect(ctx.Function).To(Equal(function))
			},
			T.Entry("foo", 0, 29, "", "foo"),
			T.Entry("bar", 1, 33, "", "bar"),
			T.Entry("byValue", 2, 39, "CallContextTest", "byValue"),
			T.Entry("byPointer", 3, 43, "(*CallContextTest)", "byPointer"),
		)
		It("should return nil context if reaching to deep", func() {
			ctx := getCallContext(100000)
			Expect(ctx).To(BeNil())
		})
	})
})
