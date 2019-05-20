package main

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
)

func fetch() error {
	resp, err := http.Get("https://www.google.com")
	if err != nil {
		return errors.Wrap(err, "error fetching google page")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected status code returned by google: %d", resp.StatusCode)
	}
	println("Success!")
	return nil
}

func main() {
	if err := fetch(); err != nil {
		log.Fatal(err)
	}
}
