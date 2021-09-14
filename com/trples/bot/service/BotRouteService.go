package service

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strings"
	domain "telegram_bot/com/trples/bot/config"
	"telegram_bot/com/trples/bot/dao"
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
		b.Send(m.Sender,fmt.Sprintf("Welcome %s %s to add learning english bot.",m.Sender.FirstName,m.Sender.LastName))
	})


	b.Handle("/stop", func(m *tb.Message) {
		err := UserStop(m.Sender)
		if err!=nil {
			b.Send(m.Sender,fmt.Sprintf("You haven't stopped all reminder"))
			return
		}
		b.Send(m.Sender,fmt.Sprintf("You have stopped all reminder"))
	})

	b.Handle("/add", func(m *tb.Message) {
		err:=VocabularyAdd(m.Sender)
		if err!=nil{
			b.Send(m.Sender, fmt.Sprintf("Server error, please retry later: %s ",err))
			return
		}
		b.Send(m.Sender, "Begin add your input [word:sentence]:")
	})
	b.Handle("/end", func(m *tb.Message) {
		//	//Receive new words, and sentences | update words and sentences | Review words
		message:=strings.ReplaceAll(m.Text,"/end","")
		user:=UserGet(int64(m.Sender.ID))
		switch user.WaitType {
			case dao.Review:
				if message == ""{
					b.Send(m.Sender, "Please enter the word you want to end")
					return
				}
				result,err:=VocabularyEndReview(m.Sender,message)
				if err!=nil{
					b.Send(m.Sender, fmt.Sprintf("error: %s, total:%d, pass:%d",err.Error(),result.Total,result.Pass))
					return
				}
				b.Send(m.Sender, fmt.Sprintf("Review word: %s completed, total:%d, pass:%d",message,result.Total,result.Pass))
			default:
				err:=VocabularyEnd(m.Sender)
				if err!=nil{
					b.Send(m.Sender, fmt.Sprintf("Server error, please retry later: %s ",err))
				}
				b.Send(m.Sender, "===End Successful===")
		}

	})
	b.Handle("/review", func(m *tb.Message) {
		message:=strings.ReplaceAll(m.Text,"/review","")
		err:=VocabularyReview(m.Sender,message)
		if err!=nil{
			b.Send(m.Sender, fmt.Sprintf("Server error, please retry later %s",err.Error()))
			return
		}
		b.Send(m.Sender, "Begin review your input [word:sentence]:")
	})
	b.Handle("/update", func(m *tb.Message) {
		err:=VocabularyUpdate(m.Sender)
		if err!=nil{
			b.Send(m.Sender, "Server error, please retry later")
		}
		b.Send(m.Sender, "Begin update your [word:sentence]")
	})
	b.Handle("/get", func(m *tb.Message) {
		message:=strings.ReplaceAll(m.Text,"/get","")
		vocabulary,err:=VocabularyGet(m.Sender.ID,message)
		if err!=nil {
			b.Send(m.Sender, fmt.Sprintf("%s doesn't exist",m.Text))
			return
		}
		b.Send(m.Sender,fmt.Sprintf("%s, the status is %s, review is %s",vocabulary.Word,vocabulary.LearnStatus,vocabulary.ReviewStatus))
		sentences,err:=SentenceFindByWord(m.Sender.ID,message)
		if err!=nil {
			b.Send(m.Sender, fmt.Sprintf("%s doesn't have sentences",m.Text))
			return
		}
		for _,v:=range sentences{
			b.Send(m.Sender,fmt.Sprintf("%s : %s",v.Word,v.Sentence))
		}
	})
	b.Handle("/schedule", func(m *tb.Message) {
		b.Send(m.Sender, "Hello schedule!")
	})

	b.Handle("/delete", func(m *tb.Message) {
		message:=strings.ReplaceAll(m.Text,"/delete","")

		fmt.Printf(" delete %s\n",message)
		VocabularyDeleteByWord(m.Sender.ID,message)
		b.Send(m.Sender, "You will delete vocabulary "+m.Text)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		user:=UserGet(int64(m.Sender.ID))
		if !user.IsInput {
			b.Send(m.Sender, fmt.Sprintf("Your status is error waitType=%s",user.WaitType))
		}
		switch user.WaitType {
			case dao.Add:
				err:=VocabularyAddReceive(m.Sender,m.Text)
				if err!=nil{
					b.Send(m.Sender, "Server error, please retry later")
					return
				}
				b.Send(m.Sender, fmt.Sprintf("Add vocabulary %s successful",m.Text))
			case dao.Update:
				err:=VocabularyUpdateReceive(m.Sender,m.Text)
				if err!=nil{
					b.Send(m.Sender, fmt.Sprintf("Server error, please retry later %s",err))
					return
				}
				b.Send(m.Sender, fmt.Sprintf("Add vocabulary %s successful",m.Text))
			case dao.Review:
				result,err:=VocabularyReviewReceive(m.Sender,m.Text)
				if err!=nil{
					b.Send(m.Sender, fmt.Sprintf("Review word: %s %t, total:%d, pass:%d",result.Word,result.Result,result.Total,result.Pass))
				}else{
					b.Send(m.Sender, fmt.Sprintf("Review word: %s %t, total:%d, pass:%d",result.Word,result.Result,result.Total,result.Pass))
				}
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
