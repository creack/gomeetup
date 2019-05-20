package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Publish writes the given body to NSQ.
func Publish(addr string, body *bytes.Buffer) error {
	req, err := http.NewRequest("POST", addr, body)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "go-client")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
	if !bytes.Equal(data, []byte("OK")) {
		return fmt.Errorf("resp %s != OK", data)
	}
	return nil
}

func feedNSQ(addr, topic string, data map[string]interface{}) {
	destURL := fmt.Sprintf("%s/put?topic=%s", addr, topic)

	buf, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: failed to marshal %+v", data)
		return
	}

	if err := Publish(destURL, bytes.NewBuffer(buf)); err != nil {
		log.Printf("ERROR: failed to publish %s", err)
	}

}
