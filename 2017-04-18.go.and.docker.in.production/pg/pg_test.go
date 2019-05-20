package main

import (
	"flag"
	"testing"
)

var pgSource = flag.String("pg", "", "Postgres datasource. Ex: `postgres://postgres@pg/postgres`")

func TestPGTimeSuccess(t *testing.T) {
	if err := pgTime(*pgSource); err != nil {
		t.Fatal(err)
	}
}

func TestPGTimeFailure(t *testing.T) {
	// Missing data source.
	if err := pgTime(""); err == nil {
		t.Fatal("Expected error upon empty data source")
	}
	// Wrong data source.
	if err := pgTime("hello world"); err == nil {
		t.Fatal("Expected error upon wrong data source")
	}
}
