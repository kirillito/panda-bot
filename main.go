package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	//"./db"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: panda-bot slack-bot-token\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	ws, id := slackConnect(os.Args[1])
	fmt.Println("panda-bot ready, ^C exits")

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
				m = pandaRead(id, m)

				postMessage(ws, m)
			}(m)
			// NOTE: the Message object is copied, this is intentional
		}

	}
}
