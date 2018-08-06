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
	"fmt"
	"log/syslog"
	"net"
	"os"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	T "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

const (
	listenAddr = "127.0.0.1:0"
	protocol   = "udp"
)

type serverUDP struct {
	conn  net.PacketConn
	data  []string
	mutex sync.Locker
}

func newServerUDP() *serverUDP {
	return &serverUDP{
		data:  make([]string, 0),
		mutex: new(sync.Mutex),
	}
}

func (server *serverUDP) addr() string {
	return server.conn.LocalAddr().String()
}

func (server *serverUDP) start() {
	var err error

	server.conn, err = net.ListenPacket(protocol, listenAddr)
	Expect(err).NotTo(HaveOccurred())

	go func() {
		var buf [4096]byte

		for {
			var n int
			var err error
			n, _, err = server.conn.ReadFrom(buf[:])
			if err != nil {
				return
			}
			server.mutex.Lock()
			server.data = append(server.data, string(buf[:n]))
			server.mutex.Unlock()
		}
	}()
}

func (server *serverUDP) get() []string {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	ret := make([]string, len(server.data))
	copy(ret, server.data)
	return ret
}

func (server *serverUDP) count() int {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	return len(server.data)
}

func (server *serverUDP) stop() {
	server.conn.Close()
}

var _ = Describe("WriterSyslog", func() {
	const (
		testMsg      = "testMessage"
		testTag      = "testTag"
		badAddr      = "987.654.321.098:765"
		testSeverity = syslog.LOG_USER
		badLevel     = Level(0xBADC0DE)
	)
	var (
		srv *serverUDP
	)

	expectLog := func(index int, priority Level, severity syslog.Priority, before, after time.Time,
		tag string, msg string) {
		Eventually(srv.count).Should(BeNumerically(">", index))
		logs := srv.get()
		Expect(len(logs)).To(BeNumerically(">", index))

		var prio int
		var date, host, tagpid, m string
		fmt.Sscanf(logs[index], "<%d>%s %s %s %s", &prio, &date, &host, &tagpid, &m)

		Expect(prio).To(Equal(int(priority) + int(severity)))

		t, err := time.Parse(time.RFC3339, date)
		Expect(err).NotTo(HaveOccurred())
		Expect(t).To(BeTemporally(">=", before.Truncate(time.Second)))
		Expect(t).To(BeTemporally("<", after))

		hostname, _ := os.Hostname()
		Expect(host).To(Equal(hostname))

		Expect(tagpid).To(Equal(fmt.Sprintf("%s[%d]:", tag, os.Getpid())))

		Expect(m).To(Equal(msg))
	}

	BeforeEach(func() {
		srv = newServerUDP()
		srv.start()
	})
	AfterEach(func() {
		srv.stop()
	})

	Describe("NewWriterSyslog", func() {
		It("should create a new object connected to log daemon", func() {
			w := NewWriterSyslog(protocol, srv.addr(), testSeverity, testTag)
			Expect(w).NotTo(BeNil())
			Expect(w.syslogClient).NotTo(BeNil())
		})
		It("should panic if cannot connect to daemon", func() {
			Expect(func() {
				NewWriterSyslog(protocol, badAddr, testSeverity, testTag)
			}).To(Panic())
		})
		T.DescribeTable("should connect using different severities",
			func(severity syslog.Priority) {
				before := time.Now()
				NewWriterSyslog(protocol, srv.addr(), severity, testTag).
					Write(EmergLevel, []byte(testMsg))
				after := time.Now()
				expectLog(0, EmergLevel, severity, before, after, testTag, testMsg)
			},
			T.Entry("LOG_USER", syslog.LOG_USER),
			T.Entry("LOG_DAEMON", syslog.LOG_DAEMON),
			T.Entry("LOG_NEWS", syslog.LOG_NEWS),
		)
		T.DescribeTable("should connect using different tags",
			func(tag string) {
				before := time.Now()
				NewWriterSyslog(protocol, srv.addr(), testSeverity, tag).
					Write(EmergLevel, []byte(testMsg))
				after := time.Now()
				expectLog(0, EmergLevel, testSeverity, before, after, tag, testMsg)
			},
			T.Entry("tag", "tag"),
			T.Entry("anotherTag", "anotherTag"),
			T.Entry("completelyAnotherTag", "completelyAnotherTag"),
		)
		It("should connect with empty tag using argv[0] as tag", func() {
			before := time.Now()
			NewWriterSyslog(protocol, srv.addr(), testSeverity, "").
				Write(EmergLevel, []byte(testMsg))
			after := time.Now()
			expectLog(0, EmergLevel, testSeverity, before, after, os.Args[0], testMsg)
		})
	})
	Describe("Write", func() {
		It("should write multiple messages", func() {
			w := NewWriterSyslog(protocol, srv.addr(), testSeverity, testTag)
			messages := []string{"To", "log", "or", "not", "to", "log?", "That's", "a", "question"}
			before := time.Now()
			for _, s := range messages {
				w.Write(EmergLevel, []byte(s))
			}
			after := time.Now()
			for i, s := range messages {
				expectLog(i, EmergLevel, testSeverity, before, after, testTag, s)
			}
		})
		T.DescribeTable("should log using different levels",
			func(level Level) {
				before := time.Now()
				n, err := NewWriterSyslog(protocol, srv.addr(), testSeverity, testTag).
					Write(level, []byte(testMsg))
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(0))
				after := time.Now()
				expectLog(0, level, testSeverity, before, after, testTag, testMsg)
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
		It("should return error if log level is unknown", func() {
			n, err := NewWriterSyslog(protocol, srv.addr(), testSeverity, testTag).
				Write(badLevel, []byte(testMsg))
			Expect(err).To(Equal(ErrInvalidLogLevel))
			Expect(n).To(Equal(0))
		})
	})
})
