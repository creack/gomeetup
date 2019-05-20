package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func pgTime(pgSource string) error {
	if pgSource == "" {
		return errors.New("missing pg data source")
	}
	db, err := sqlx.Connect("postgres", pgSource)
	if err != nil {
		return errors.Wrap(err, "error connecting to the db")
	}
	defer func() { _ = db.Close() }() // Best effort.

	var t time.Time
	if err := db.Get(&t, "SELECT NOW()"); err != nil {
		return errors.Wrap(err, "error querying time in db")
	}
	fmt.Printf("%s\n", t)
	return nil
}

func main() {
	pgSource := flag.String("pg", "", "Postgres datasource. Ex: `postgres://postgres@pg/postgres`")
	flag.Parse()
	if err := pgTime(*pgSource); err != nil {
		log.Fatal(err)
	}
}
