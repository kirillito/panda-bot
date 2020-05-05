/*
panda-bot
Slack bot for team Panda in Xello. Custom features include:
- setting vacation/sick days reminder
- ???
*/

package panda

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"../db"
)

// Message types
type MessageType int

const (
	Unsupported        MessageType = 0
	VacationSetSelf    MessageType = 1
	VacationSetMention MessageType = 2
	VacationGetList    MessageType = 5
	UserMention        MessageType = 10
)

type PandaBot struct {
	logger *log.Logger
	db     *db.DB
}

func Create(logger *log.Logger, db *db.DB) *PandaBot {
	return &PandaBot{logger: logger, db: db}
}

func (bot *PandaBot) Run(token string) {
	// start a websocket-based Real Time API session
	ws, id := slackConnect(token)
	fmt.Println("panda-bot ready, ^C exits")

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}
		// see if we're supposed to answer
		go func(m Message) {
			var doAnswer = false

			m, doAnswer = bot.pandaAnswerMessage(id, m)

			if doAnswer {
				postMessage(ws, m)
			}
		}(m)
		// NOTE: the Message object is copied, this is intentional

	}
}

func (bot *PandaBot) pandaAnswerMessage(id string, m Message) (Message, bool) {
	var doAnswer bool = false

	// determine the type of the message
	var msgType = bot.getMessageType(id, m)

	switch msgType {
	case VacationSetSelf:
		dateFrom, dateTo, err := getDateRangeFromString(m.Text)
		if err == nil {
			bot.saveVacation(m.User, "vacation", dateFrom, dateTo)
			m.Text = "Got it!"
			doAnswer = true
		}
	case VacationSetMention:
		doAnswer = true
	case VacationGetList:
		lst := bot.listVacations(m.User)
		m.Text = "Here's the list of vacations: " + lst
		doAnswer = true
	default:
	}

	// parse the message into separate sections
	// m.Text = bot.pandaAnswer(strings.TrimPrefix(m.Text, "<@"+id+">"))

	return m, doAnswer
}

func (bot *PandaBot) saveVacation(userId string, vacationType string, dateStart time.Time, dateEnd time.Time) {
	bot.db.Update(func(tx *db.Tx) error {
		v := db.Vacation{
			Tx:        tx,
			UserId:    []byte(userId),
			Type:      vacationType,
			DateStart: dateStart,
			DateEnd:   dateEnd,
		}

		return v.Save()
	})

	return
}

func (bot *PandaBot) listVacations(userId string) string {
	var v = db.Vacation{
		UserId: []byte(userId),
	}
	//var err error

	bot.db.View(func(tx *db.Tx) error {
		v.Load()

		return nil
	})

	// if v == nil {
	// 	return "No vacations found"
	// }

	return string(v.UserId) + " is on vacation from " + time.Time.String(v.DateStart) + " to " + time.Time.String(v.DateEnd)
}

func (bot *PandaBot) getMessageType(id string, m Message) MessageType {
	if m.Type == "message" {

		// messages targeting the bot
		if strings.HasPrefix(m.Text, "<@"+id+">") {
			if strings.Contains(m.Text, "on vacation") {
				if strings.Contains(m.Text, "I'm") || strings.Contains(m.Text, "I am") || strings.Contains(m.Text, m.User) {
					return VacationSetSelf
				} else {
					return VacationSetMention
				}
			} else if strings.Contains(m.Text, "vacation list") {
				return VacationGetList
			}
		}
	}

	return Unsupported
}

func (bot *PandaBot) save() {
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

// var history []string

// func (bot *PandaBot) pandaAnswer(msg string) string {
// 	history = append(history, msg)
// 	return strings.Join(history[:], " ")
// }

func getDateRangeFromString(str string) (time.Time, time.Time, error) {
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

	dates := re.FindAllString(str, -1)

	if len(dates) < 2 {
		return time.Now(), time.Now(), errors.New("no date range found")
	}

	layout := "2006-01-02"
	date1, err := time.Parse(layout, dates[0])
	if err != nil {
		return time.Now(), time.Now(), errors.New("no date range found")
	}

	date2, err := time.Parse(layout, dates[1])
	if err != nil {
		return time.Now(), time.Now(), errors.New("no date range found")
	}

	return date1, date2, nil
}
