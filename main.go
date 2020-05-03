package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"./db"
	"golang.org/x/net/websocket"
)

var dbFile string

func newDB(fn string) (*db.DB, error) {
	db := &db.DB{}

	if err := db.Open(dbFile, 0600); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: panda-bot slack-bot-token\n")
		os.Exit(1)
	}

	setup()

	// start a websocket-based Real Time API session
	ws, id := slackConnect(os.Args[1])
	fmt.Println("panda-bot ready, ^C exits")

	runBot(ws, id)
}

func setup() {
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
	defer db.Close()
}

func runBot(ws *websocket.Conn, id string) {
	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}
		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			// if so try to parse if
			//parts := strings.Fields(m.Text)
			// looks good, get the quote and reply with the result
			go func(m Message) {
				m = pandaAnswerMessage(id, m)
				postMessage(ws, m)
			}(m)
			// NOTE: the Message object is copied, this is intentional
		}

	}
}
