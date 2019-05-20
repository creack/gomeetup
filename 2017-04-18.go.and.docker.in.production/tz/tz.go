package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

func tz() error {
	austinLoc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		return errors.Wrap(err, "error loading Austin timezone location")
	}

	fmt.Printf("%s\n", time.Now().In(austinLoc))
	return nil
}

func main() {
	if err := tz(); err != nil {
		log.Fatal(err)
	}
}
