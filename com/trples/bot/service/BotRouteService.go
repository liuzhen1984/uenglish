package service

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	domain "telegram_bot/com/trples/bot/config"
	"time"
)

func BotRoute()  {
	config:=domain.LoadProperties()

	b, err := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		//URL: "http://195.129.111.17:8012",

		Token:  config.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			fmt.Printf("username= %s,%s,id=%d, is private \n",m.Sender.FirstName,m.Sender.LastName,m.Sender.ID)
			return
		}
		err := UserNewOrExisting(m.Sender)
		if err!=nil {
			b.Send(m.Sender,fmt.Sprintf("Welcome %s %s to add learning english bot, but init user failed",m.Sender.FirstName,m.Sender.LastName))
			return
		}
		b.Send(m.Sender,fmt.Sprintf("Welcome %s %s to add learning english bot",m.Sender.FirstName,m.Sender.LastName))
	})

	b.Handle("/stop", func(m *tb.Message) {
		err := UserStop(m.Sender)
		if err!=nil {
			b.Send(m.Sender,fmt.Sprintf("You haven't stopped all reminder"))
			return
		}
		b.Send(m.Sender,fmt.Sprintf("You have stopped all reminder"))
	})

	b.Handle("/daily", func(m *tb.Message) {
		err:=VocabularyStart(m.Sender)
		if err!=nil{
			b.Send(m.Sender, "Server error, please retry later")
		}
		b.Send(m.Sender, "Begin receive your input [word:sentence]")
		b.Send(m.Sender, "=======================================")
	})
	b.Handle("/end", func(m *tb.Message) {
		//	//Receive new words, and sentences | update words and sentences | Review words
		err:=VocabularyEnd(m.Sender)
		if err!=nil{
			b.Send(m.Sender, "Server error, please retry later")
		}
		b.Send(m.Sender, "=======================================")
	})
	b.Handle("/review", func(m *tb.Message) {
		b.Send(m.Sender, "Hello review!")
	})

	b.Handle("/schedule", func(m *tb.Message) {
		b.Send(m.Sender, "Hello schedule!")
	})

	b.Handle("/delete", func(m *tb.Message) {
		fmt.Printf(" delete %s\n",m.Text)
		VocabularyDeleteByWord(m.Sender.ID,m.Text)
		b.Send(m.Sender, "You will delete vocabulary "+m.Text)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		user:=UserGet(int64(m.Sender.ID))
		if !user.IsInput {
			b.Send(m.Sender, fmt.Sprintf("Your status is error waitType=%s",user.WaitType))
		}
		switch user.WaitType {
			case "receive":
				err:=VocabularyAdd(m.Sender,m.Text)
				if err!=nil{
					b.Send(m.Sender, "Server error, please retry later")
				}
				b.Send(m.Sender, fmt.Sprintf("Add vocabulary %s successful",m.Text))
			case "update":
				//todo
				b.Send(m.Sender, fmt.Sprintf("Update vocabulary %s successful",m.Text))
			case "review":
				//todo
				b.Send(m.Sender, fmt.Sprintf("Review vocabulary %s successful",m.Text))
			default:
				b.Send(m.Sender, "User wait type value is error, please retry")
				VocabularyEnd(m.Sender)
		}
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		// photos only
	})

	b.Handle(tb.OnChannelPost, func (m *tb.Message) {
		// channel posts only
	})

	b.Handle(tb.OnQuery, func (q *tb.Query) {
	})

	b.Start()
}
