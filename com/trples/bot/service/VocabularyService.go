package service

import (
	"gopkg.in/tucnak/telebot.v2"
	"strings"
	"telegram_bot/com/trples/bot/dao"
)

func VocabularyDeleteByWord(userId int,word string)  {
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		panic(err)
	}
	defer dao.CloseClient(ctx,client)
	dao.SentencesDeleteByWord(ctx,client,int64(userId),word)
	dao.VocabularyDeleteByWord(ctx,client,int64(userId),word)
}

func VocabularyGet(userId int,word string) dao.Vocabulary{
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		panic(err)
	}
	defer dao.CloseClient(ctx,client)
	vocabulary,_:=dao.VocabularyGet(ctx,client,int64(userId),word)
	return vocabulary
}
// vocabulary: sentences
func VocabularyAdd(sender *telebot.User,message string) (error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	words := strings.Split(message,":")
	vocabulary,err:=dao.VocabularyGet(ctx,client,int64(sender.ID),words[0])
	if err != nil{
		vocabulary=dao.Vocabulary{}
		vocabulary.Word = words[0]
		vocabulary.IsRemember = false
		vocabulary.LearnStatus = "waiting"
		vocabulary.Period = 24
		vocabulary.ReminderCount = 0
		vocabulary.UserId = int64(sender.ID)
		_,err:=dao.VocabularySave(ctx,client,vocabulary)
		if err != nil{
			return err
		}
		dao.SentencesDeleteByWord(ctx,client,int64(sender.ID),words[0])
		sentence:=dao.Sentences{}
		sentence.Word = words[0]
		sentence.Sentence = words[1]
		sentence.ReviewCount = 0
		sentence.UserId = int64(sender.ID)
		sentence.Status = "pass"
		_,err = dao.SentenceSave(ctx,client,sentence)
		return err
	}

	sentences,err:=dao.SentenceFindByWord(ctx,client,int64(sender.ID),words[0])

	if err != nil{
		return err
	}
	exist:= false
	for _,v:= range sentences{
		if strings.ToLower(strings.Trim(v.Sentence," "))==strings.ToLower(strings.Trim(words[1], " ")) {
			exist = true
			break
		}
	}
	if !exist {
		sentence:=dao.Sentences{}
		sentence.Word = words[0]
		sentence.Sentence = words[1]
		sentence.ReviewCount = 0
		sentence.Status = "pass"
		sentence.UserId = int64(sender.ID)
		dao.SentenceSave(ctx,client,sentence)
	}
	return nil
}

func VocabularyStart(sender *telebot.User) (error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	return dao.UserStartInput(ctx,client,int64(sender.ID),"receive")
}

func VocabularyEnd(sender *telebot.User) (error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	return dao.UserEndInput(ctx,client,int64(sender.ID))
}