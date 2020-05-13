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
	Unsupported               MessageType = 0
	VacationSetSelf           MessageType = 1
	VacationSetMention        MessageType = 2
	VacationGetListForMention MessageType = 5
	VacationGetListAll        MessageType = 6
	UserMention               MessageType = 10
	Goodnight                 MessageType = 100
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
			err = bot.saveVacation(m.User, "vacation", dateFrom, dateTo)

			if err != nil {
				m.Text = fmt.Sprint("Error: ", err)
			} else {
				m.Text = "Got it!"
			}

			doAnswer = true
		}
	case VacationSetMention:
		dateFrom, dateTo, err := getDateRangeFromString(m.Text)
		userIdMentioned, err := getMentionedUserId(m.Text)
		if err == nil {
			err = bot.saveVacation(userIdMentioned, "vacation", dateFrom, dateTo)

			if err != nil {
				m.Text = fmt.Sprint("Error: ", err)
			} else {
				m.Text = "Got it!"
			}

			doAnswer = true
		}
	case VacationGetListForMention:
		userIdMentioned, err := getMentionedUserId(m.Text)
		if err == nil {
			lst, err := bot.listVacations(userIdMentioned)

			if err != nil {
				m.Text = fmt.Sprint("Error: ", err)
			} else {
				m.Text = "Here's the list of vacations: " + lst
			}
			doAnswer = true
		}
	case VacationGetListAll:
	case Goodnight:
		m.Text = fmt.Sprintf("Goodnight, <@%s>", m.User)
		doAnswer = true
	default:
	}

	// parse the message into separate sections
	// m.Text = bot.pandaAnswer(strings.TrimPrefix(m.Text, "<@"+id+">"))

	return m, doAnswer
}

func (bot *PandaBot) saveVacation(userId string, vacationType string, dateStart time.Time, dateEnd time.Time) error {
	err := bot.db.Update(func(tx *db.Tx) error {
		v := db.Vacation{
			Tx:        tx,
			UserId:    []byte(userId),
			Type:      vacationType,
			DateStart: dateStart,
			DateEnd:   dateEnd,
		}

		return v.Save()
	})

	return err
}

func (bot *PandaBot) listVacations(userId string) (string, error) {
	var v *db.Vacation
	var err error

	bot.db.View(func(tx *db.Tx) error {
		v, err = tx.Vacation([]byte(userId))
		return err
	})

	// if v == nil {
	// 	return "No vacations found"
	// }

	return fmt.Sprintf("<@%s> is on vacation from %s to %s :panda-roll:", v.UserId, v.DateStart.Format("2006-01-02"), v.DateEnd.Format("2006-01-02")), err
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
			} else if strings.Contains(m.Text, "list vacations for") {
				return VacationGetListForMention
			} else if strings.Contains(m.Text, "list all vacations") {
				return VacationGetListAll
			} else if strings.Contains(strings.ToLower(m.Text), "goodnight") {
				return Goodnight
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

func getMentionedUserId(str string) (string, error) {
	re := regexp.MustCompile(`<@([WU].+?)>`)

	mentions := re.FindAllString(str, -1)

	if len(mentions) < 2 {
		return "", errors.New("no user mention found")
	}

	return mentions[1][2 : len(mentions[1])-1], nil
}
