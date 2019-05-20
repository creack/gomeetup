package main

import (
	"fmt"
	"net"
	"testing"
)

type redisTestQueryData struct {
	request, response string
}

type testCacheDataSet struct {
	query  []redisTestQueryData
	status string
	data   string
}

var (
	testLookupCacheData = map[string]testCacheDataSet{
		"first_time": {
			status: "not-started",
			data:   "",
			query: []redisTestQueryData{
				{
					request: "*1\r\n$5\r\nMULTI\r\n" +
						"*3\r\n$5\r\nSETNX\r\n$47\r\nstatus:reports:99914b932bd37a50b983c5e7c90ae93b\r\n$11\r\nin-progress\r\n" +
						"*2\r\n$3\r\nGET\r\n$47\r\nstatus:reports:99914b932bd37a50b983c5e7c90ae93b\r\n" +
						"*4\r\n$6\r\nLRANGE\r\n$46\r\ncache:reports:99914b932bd37a50b983c5e7c90ae93b\r\n$1\r\n0\r\n$2\r\n-1\r\n" +
						"*1\r\n$4\r\nEXEC\r\n",
					response: "+OK\r\n+QUEUED\r\n+QUEUED\r\n+QUEUED\r\n" +
						"*3\r\n:1\r\n$11\r\nin-progress\r\n*0\r\n",
				}, {
					request:  "*2\r\n$3\r\nDEL\r\n$46\r\ncache:reports:99914b932bd37a50b983c5e7c90ae93b\r\n",
					response: ":0\r\n",
				}, {
					request:  "*3\r\n$6\r\nEXPIRE\r\n$47\r\nstatus:reports:99914b932bd37a50b983c5e7c90ae93b\r\n$3\r\n120\r\n",
					response: ":1\r\n",
				},
			},
		},
		"second_time": {
			status: "in-pgrogress",
			data:   "",
			query: []redisTestQueryData{
				{
					request: "*1\r\n$5\r\nMULTI\r\n" +
						"*3\r\n$5\r\nSETNX\r\n$47\r\nstatus:reports:99914b932bd37a50b983c5e7c90ae93b\r\n$11\r\nin-progress\r\n" +
						"*2\r\n$3\r\nGET\r\n$47\r\nstatus:reports:99914b932bd37a50b983c5e7c90ae93b\r\n" +
						"*4\r\n$6\r\nLRANGE\r\n$46\r\ncache:reports:99914b932bd37a50b983c5e7c90ae93b\r\n$1\r\n0\r\n$2\r\n-1\r\n" +
						"*1\r\n$4\r\nEXEC\r\n",
					response: "+OK\r\n+QUEUED\r\n+QUEUED\r\n+QUEUED\r\n" +
						"*3\r\n:0\r\n$11\r\nin-progress\r\n*0\r\n",
				},
			},
		},
	}
)

func TestLookupSuccess(t *testing.T) {
	for step, testData := range testLookupCacheData {
		ts := newTestServer(t, func(c net.Conn) {
			for _, elem := range testData.query {
				buf := make([]byte, 2048)
				n, err := c.Read(buf)
				if err != nil {
					t.Errorf("[%s] %s", step, err)
					break
				}
				if expect, got := elem.request, string(buf[:n]); expect != got {
					t.Errorf("[%s] Unexpected result.\nExpect:\t%q\nGot:\t%q", step, expect, got)
					break
				}
				fmt.Fprint(c, elem.response)
			}
			c.Close()
		})
		_ = ts
		// Call redis

		// if status, data, err := cacheManager.Lookup(); err != nil {
		// 	t.Fatalf("[%s] %s", step, err)
		// } else if expect, got := testData.status, status; expect != got {
		// 	t.Fatalf("[%s] Unexpected result.\nExpect:\t%s\nGot:\t%s", step, expect, got)
		// } else if expect, got := testData.data, string(bytes.Join(data, []byte(","))); expect != got {
		// 	t.Fatalf("[%s] Unexpected result.\nExpect:\t%s\nGot:\t%s", step, expect, got)
		// }
	}
}
