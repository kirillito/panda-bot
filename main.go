package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"./db"
	"./panda"
)

var dbFile string

func newDB(fn string) (*db.DB, error) {
	db := &db.DB{}

	if err := db.Open(dbFile, 0600); err != nil {
		return nil, err
	}

	return db, nil
}

func setup() (*db.DB, *log.Logger) {
	// connect DB
	flag.StringVar(&dbFile, "db", ".\\tmp\\panda.db", "Path to the BoltDB file")
	flag.Parse()

	// Setup the logger
	logger := log.New(os.Stdout, "", 0)

	// Setup the database
	db, err := newDB(dbFile)
	if err != nil {
		logger.Fatal(err)
	}

	return db, logger
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: panda-bot slack-bot-token\n")
		os.Exit(1)
	}

	db, logger := setup()
	defer db.Close()

	bot := panda.Create(logger, db)
	bot.Run(os.Args[1])
}
