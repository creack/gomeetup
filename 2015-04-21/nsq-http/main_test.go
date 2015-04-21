package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func MockNSQ(t *testing.T, topic string, data map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if expected, got := "go-client", req.UserAgent(); expected != got {
			t.Fatalf("Unexpected user agent. Expected %q, Got %q\n", expected, got)
		}
		if expected, got := "POST", req.Method; expected != got {
			t.Fatalf("Unexpected Method. Expected %q, Got %q\n", expected, got)
		}
		if expected, got := "application/json", req.Header.Get("Content-Type"); expected != got {
			t.Fatalf("Unexpected content type. Expected %q, Got %q\n", expected, got)
		}
		if expected, got := "/put?topic=testtopic", req.URL.String(); expected != got {
			t.Fatalf("Unexpected Method. Expected %q, Got %q\n", expected, got)
		}

		v := map[string]interface{}{}
		if err := json.NewDecoder(req.Body).Decode(&v); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(v, data) {
			fmt.Fprint(w, "KO")
		}
		fmt.Fprint(w, "OK")
	}
}

func TestFeedNSQ(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
	}
	NSQTopic := "testtopic"

	ts := httptest.NewServer(MockNSQ(t, NSQTopic, m))
	feedNSQ(ts.URL, NSQTopic, m)
	ts.Close()
}
