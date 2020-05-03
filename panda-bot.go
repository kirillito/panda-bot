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

func pandaAnswerMessage(id string, m Message) Message {
	m.Text = pandaAnswer(strings.TrimPrefix(m.Text, "<@"+id+">"))

	return m
}

func save() {
	/*	s.db.Update(func(tx *db.Tx) error {
		r.ParseForm()

		v := db.Vacation{
			Tx:   tx,
			Id: s.getPageName(r),
		}

		v.Data = []byte(strings.TrimSpace(r.FormValue("text")))

		return v.Save()
	})*/

}

var history []string

func pandaAnswer(msg string) string {
	history = append(history, msg)
	return strings.Join(history[:], " ")
}
