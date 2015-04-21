package main

import (
	"log"

	"github.com/simplereach/platform/connectors/vertica"
	"github.com/simplereach/platform/util/proxydumper"
)

func main() {
	l, errChan, errC := proxydumper.Start("tcp", "127.0.0.1:8085", "54.85.149.150:5433", "client.log", "server.log")
	if err := <-errC; err != nil {
		log.Fatal(err)
	}

	c, err := vertica.NewODBC("vertica.json", "dbadmin", "simplereach_analytics")
	if err != nil {
		log.Fatal(err)
	}
	c.Query("SELECT NOW();")
	c.Close()

	l.Close()
	for err := range errChan {
		if err != nil {
			log.Printf("-> %s\n", err)
		}
	}
}
