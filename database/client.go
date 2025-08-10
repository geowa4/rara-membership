package database

import (
	"context"
	"log"
	"os"

	"github.com/geowa4/rara-membership/ent"
	_ "github.com/mattn/go-sqlite3"
)

var Client *ent.Client

func Init() {
	var err error
	
	// Use environment variable for database path, default to membership.db
	dbPath := os.Getenv("RARA_DB_PATH")
	if dbPath == "" {
		dbPath = "membership.db"
	}
	
	dsn := "file:" + dbPath + "?cache=shared&_fk=1"
	Client, err = ent.Open("sqlite3", dsn)
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
