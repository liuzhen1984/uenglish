package service

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strings"
	"sync"
	domain "telegram_bot/com/trples/bot/config"
	"telegram_bot/com/trples/bot/dao"
	"time"
)

var bot *tb.Bot
var once sync.Once

func GetBot() *tb.Bot{
	config:=domain.LoadProperties()
	once.Do(func() {
		bot = &tb.Bot{}
		var err error
		bot,err = tb.NewBot(tb.Settings{
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
	})
	return bot
}
func Send(userId int, message string)  {
	user := tb.User{}
	user.ID=userId
	bot.Send(&user,message)
}

func BotRoute()  {
	GetBot()
	bot.Handle("/test", func(m *tb.Message) {
		Send(1470647544,fmt.Sprintf("Test sender %s %s, send message %s",m.Sender.FirstName,m.Sender.LastName,m.Text))
	})
	bot.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			fmt.Printf("username= %s,%s,id=%d, is private \n",m.Sender.FirstName,m.Sender.LastName,m.Sender.ID)
			return
		}
		err := UserNewOrExisting(m.Sender)
		if err!=nil {
			bot.Send(m.Sender,fmt.Sprintf("Welcome %s %s to add learning english bot, but init user failed",m.Sender.FirstName,m.Sender.LastName))
			return
		}
		bot.Send(m.Sender,fmt.Sprintf("Welcome %s %s to add learning english bot.",m.Sender.FirstName,m.Sender.LastName))
	})

	bot.Handle("/list",func(m *tb.Message) {
		vList,err:=VocabularyFindByUserId(m.Sender.ID)
		if err!=nil{
			bot.Send(m.Sender,fmt.Sprintf("Cannot get anything vocabularies."))
			return
		}
		bot.Send(m.Sender,fmt.Sprintf("You have added vocabularies:"))
		for _,v:=range vList{
			bot.Send(m.Sender,fmt.Sprintf("%s , review: %s, review count: %d",v.Word,v.ReviewStatus,v.ReminderCount))
		}
	})

	bot.Handle("/stop", func(m *tb.Message) {
		err := UserStop(m.Sender)
		if err!=nil {
			bot.Send(m.Sender,fmt.Sprintf("You haven't stopped all reminder"))
			return
		}
		bot.Send(m.Sender,fmt.Sprintf("You have stopped all reminder"))
	})

	bot.Handle("/add", func(m *tb.Message) {
		err:=VocabularyAdd(m.Sender)
		if err!=nil{
			bot.Send(m.Sender, fmt.Sprintf("Server error, please retry later: %s ",err))
			return
		}
		bot.Send(m.Sender, "Begin add your input [word:sentence]:")
	})
	bot.Handle("/end", func(m *tb.Message) {
		//	//Receive new words, and sentences | update words and sentences | Review words
		message:=strings.ReplaceAll(m.Text,"/end","")
		user:=UserGet(int64(m.Sender.ID))
		switch user.WaitType {
			case dao.Review:
				message = strings.ToLower(strings.Trim(message," "))
				var result ReviewResult
				var err error
				if message == ""{
					result,err =VocabularyEndAllReview(m.Sender.ID)
				}else{
					result,err =VocabularyEndReview(m.Sender.ID,message)
				}
				if err!=nil{
					bot.Send(m.Sender, fmt.Sprintf("error: %s, total:%d, pass:%d",err.Error(),result.Total,result.Pass))
					return
				}
				if message == ""{
					bot.Send(m.Sender, fmt.Sprintf("Review word: %s completed, total:%d, pass:%d",result.Word,result.Total,result.Pass))
				} else{
					bot.Send(m.Sender, fmt.Sprintf("Review %s 's sentences [total:%d, pass:%d]",result.Word,result.Total,result.Pass))
				}
			default:
				err:=VocabularyEnd(m.Sender.ID)
				if err!=nil{
					bot.Send(m.Sender, fmt.Sprintf("Server error, please retry later: %s ",err))
				}
				if user.WaitType!=""{
					bot.Send(m.Sender, fmt.Sprintf("Operation [%s] finished",user.WaitType))
				}else{
					bot.Send(m.Sender, "Nothing...")
				}
		}

	})
	bot.Handle("/review", func(m *tb.Message) {
		message:=strings.ReplaceAll(m.Text,"/review","")
		vList,err:=VocabularyReview(m.Sender.ID,message)
		if err!=nil{
			bot.Send(m.Sender, fmt.Sprintf("Server error, please retry later %s",err.Error()))
			return
		}
		if len(vList)<=0 {
			bot.Send(m.Sender,"No need to review vocabularies")
			err:=VocabularyEnd(m.Sender.ID)
			if err!=nil{
				bot.Send(m.Sender, fmt.Sprintf("End all review error: %s ",err))
			}
			return
		}
		bot.Send(m.Sender, "Begin review these vocabularies:")
		for _,v:=range vList{
			bot.Send(m.Sender,v)
		}
	})
	bot.Handle("/update", func(m *tb.Message) {
		err:=VocabularyUpdate(m.Sender)
		if err!=nil{
			bot.Send(m.Sender, "Server error, please retry later")
		}
		bot.Send(m.Sender, "Begin update your [word:sentence]")
	})
	bot.Handle("/get", func(m *tb.Message) {
		message:=strings.ReplaceAll(m.Text,"/get","")
		vocabulary,err:=VocabularyGet(m.Sender.ID,message)
		if err!=nil {
			bot.Send(m.Sender, fmt.Sprintf("%s doesn't exist",m.Text))
			return
		}
		bot.Send(m.Sender,fmt.Sprintf("%s, the status is %s, review is %s",vocabulary.Word,vocabulary.LearnStatus,vocabulary.ReviewStatus))
		sentences,err:=SentenceFindByWord(m.Sender.ID,message)
		if err!=nil {
			bot.Send(m.Sender, fmt.Sprintf("%s doesn't have sentences",m.Text))
			return
		}
		for _,v:=range sentences{
			bot.Send(m.Sender,fmt.Sprintf("%s : %s",v.Word,v.Sentence))
		}
	})
	bot.Handle("/schedule", func(m *tb.Message) {
		bot.Send(m.Sender, "Hello schedule!")
	})

	bot.Handle("/delete", func(m *tb.Message) {
		message:=strings.ReplaceAll(m.Text,"/delete","")

		fmt.Printf(" delete %s\n",message)
		VocabularyDeleteByWord(m.Sender.ID,message)
		bot.Send(m.Sender, "You will delete vocabulary "+m.Text)
	})
	bot.Handle("/lang", func(m *tb.Message) {
		message:=strings.ReplaceAll(m.Text,"/lang","")

		fmt.Printf(" lang %s\n",message)
		UserUpdateLang(m.Sender,message)
		bot.Send(m.Sender, "You will delete vocabulary "+m.Text)
	})

	bot.Handle("/t", func(m *tb.Message) {
		err:=DictionaryTranslateStart(m.Sender)
		if err!=nil{
			bot.Send(m.Sender, fmt.Sprintf("Server error, please retry later: %s ",err))
			return
		}
		bot.Send(m.Sender, "Begin add your input sentence:")
	})

	bot.Handle("/longman", func(m *tb.Message) {
		message:=strings.ReplaceAll(m.Text,"/longman","")

		fmt.Printf(" longman %s\n",message)
		result:=DictionaryLongman(m.Sender.ID,message)
		bot.Send(m.Sender, "From longman dictionary : " + result)
		bot.Send(m.Sender, "https://www.ldoceonline.com/dictionary/" + strings.ToLower(strings.Trim(result, " ")))
	})

	bot.Handle("/remember", func(m *tb.Message) {
		message:=strings.ReplaceAll(m.Text,"/remember","")

		fmt.Printf(" remember %s\n",message)
		result:=VocabularyRemember(m.Sender.ID,message)
		if result!=nil {
			bot.Send(m.Sender, fmt.Sprintf("Remember the vocabulary [%s] operate failed %s",message,result))
			return
		}
		bot.Send(m.Sender, fmt.Sprintf("Remember the vocabulary [%s] operate successful",message))
	})

	bot.Handle("/reset", func(m *tb.Message) {
		message:=strings.ReplaceAll(m.Text,"/remember","")

		fmt.Printf(" remember %s\n",message)
		result:=VocabularyReset(m.Sender.ID,message)
		if result!=nil {
			bot.Send(m.Sender, fmt.Sprintf("Reset the vocabulary [%s] operate failed %s",message,result))
			return
		}
		bot.Send(m.Sender, fmt.Sprintf("Reset the vocabulary [%s] operate successful",message))
	})

	bot.Handle("/check", func(m *tb.Message) {
		vList,err:=VocabularyCheck(m.Sender.ID)
		if err!=nil{
			bot.Send(m.Sender, fmt.Sprintf("Server error, please retry later %s",err.Error()))
			return
		}
		if len(vList)<=0 {
			bot.Send(m.Sender,"No need to review vocabularies")
			return
		}
		bot.Send(m.Sender, "These vocabularies need to be reviewed:")
		for _,v:=range vList{
			bot.Send(m.Sender,v)
		}
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		user:=UserGet(int64(m.Sender.ID))
		if !user.IsInput {
			bot.Send(m.Sender, fmt.Sprintf("Your status is error waitType=%s",user.WaitType))
		}
		switch user.WaitType {
			case dao.Add:
				err:=VocabularyAddReceive(m.Sender,m.Text)
				if err!=nil{
					bot.Send(m.Sender, "Server error, please retry later")
					return
				}
				bot.Send(m.Sender, fmt.Sprintf("Add vocabulary %s successful",m.Text))
			case dao.Update:
				count,err:=VocabularyUpdateReceive(m.Sender,m.Text)
				if err!=nil{
					bot.Send(m.Sender, fmt.Sprintf("Server error, please retry later %s",err))
					return
				}
				if count>0{
					bot.Send(m.Sender, fmt.Sprintf("Updated the vocabulary %s successful",m.Text))
				} else{
					bot.Send(m.Sender, fmt.Sprintf("Didn't update for the vocabulary %s",m.Text))
				}
			case dao.Review:
				result,err:=VocabularyReviewReceive(m.Sender,m.Text)
				if err!=nil{
					bot.Send(m.Sender, fmt.Sprintf("Review word: %s %t, total:%d, pass:%d",result.Word,result.Result,result.Total,result.Pass))
				}else{
					bot.Send(m.Sender, fmt.Sprintf("Review word: %s %t, total:%d, pass:%d",result.Word,result.Result,result.Total,result.Pass))
				}
			case dao.Translate:
				result,err:=DictionaryTranslate(m.Sender.ID,m.Text)
				if err!=nil{
					bot.Send(m.Sender, fmt.Sprintf("Translate failed from google : %s, error %s" , result,err))
				}else{
					bot.Send(m.Sender, "From google translate : " )
					bot.Send(m.Sender, result )
				}
				VocabularyEnd(m.Sender.ID)
			default:
				bot.Send(m.Sender, "User wait type value is error, please retry")
				VocabularyEnd(m.Sender.ID)
		}
	})

	bot.Handle(tb.OnPhoto, func(m *tb.Message) {
		// photos only
	})

	bot.Handle(tb.OnChannelPost, func (m *tb.Message) {
		// channel posts only
	})

	bot.Handle(tb.OnQuery, func (q *tb.Query) {
	})

	bot.Start()
}


