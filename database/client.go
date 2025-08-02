package database

import (
	"context"
	"log"

	"github.com/geowa4/rara-membership/ent"
	_ "github.com/mattn/go-sqlite3"
)

var Client *ent.Client

func Init() {
	var err error
	Client, err = ent.Open("sqlite3", "file:membership.db?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	ctx := context.Background()
	if err := Client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

func Close() {
	if Client != nil {
		Client.Close()
	}
}
