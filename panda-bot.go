/*
panda-bot
Slack bot for team Panda in Xello. Custom features include:
- setting vacation/sick days reminder
- ???
*/

package main

import (
	"strings"
)

func pandaRead(id string, m Message) Message {
	m.Text = pandaAnswer(strings.TrimPrefix(m.Text, "<@"+id+">"))

	return m
}

var history []string

func pandaAnswer(msg string) string {
	history = append(history, msg)
	return strings.Join(history[:], " ")
}
