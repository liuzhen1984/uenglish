package service

import (
	"errors"
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
	"telegram_bot/com/trples/bot/dao"
)


type ReviewResult struct {
	Word	 string `json:"word"`
	Total	 int    `json:"total"`
	Pass	 int    `json:"pass"`
	Result   bool  	`json:"result"`
}

func VocabularyDeleteByWord(userId int,word string)  {
	word = strings.ToLower(strings.Trim(word," "))
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		panic(err)
	}
	defer dao.CloseClient(ctx,client)
	dao.SentencesDeleteByWord(ctx,client,int64(userId),word)
	dao.VocabularyDeleteByWord(ctx,client,int64(userId),word)
}

func VocabularyGet(userId int,word string) (dao.Vocabulary,error){
	word = strings.ToLower(strings.Trim(word," "))
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return dao.Vocabulary{},err
	}
	defer dao.CloseClient(ctx,client)
	vocabulary,err:=dao.VocabularyGet(ctx,client,int64(userId),word)
	return vocabulary,err
}
func VocabularyAdd(sender *telebot.User) (error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	return dao.UserStartInput(ctx,client,int64(sender.ID),dao.Add)
}
// vocabulary: sentences
func VocabularyAddReceive(sender *telebot.User,message string) (error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	words := strings.Split(message,":")
	word:= strings.ToLower(strings.Trim(words[0]," "))
	sent:=strings.ToLower(strings.Trim(words[1]," "))
	vocabulary,err:=dao.VocabularyGet(ctx,client,int64(sender.ID),word)
	if err != nil{
		vocabulary=dao.Vocabulary{}
		vocabulary.Word = word
		vocabulary.IsRemember = false
		vocabulary.LearnStatus = dao.Waiting
		vocabulary.ReviewStatus = dao.FAIL
		vocabulary.Period = 24
		vocabulary.ReminderCount = 0
		vocabulary.UserId = int64(sender.ID)
		_,err:=dao.VocabularySave(ctx,client,vocabulary)
		if err != nil{
			return err
		}
		dao.SentencesDeleteByWord(ctx,client,int64(sender.ID),word)
		sentence:=dao.Sentences{}
		sentence.Word = word
		sentence.Sentence = sent
		sentence.ReviewCount = 0
		sentence.Status = dao.FAIL
		sentence.UserId = int64(sender.ID)
		sentence.Status = "pass"
		_,err = dao.SentenceSave(ctx,client,sentence)
		return err
	}

	sentences,err:=dao.SentenceFindByWord(ctx,client,int64(sender.ID),word)

	if err != nil{
		return err
	}
	exist:= false
	for _,v:= range sentences{
		if strings.ToLower(strings.Trim(v.Sentence," "))==sent {
			exist = true
			break
		}
	}
	if !exist {
		sentence:=dao.Sentences{}
		sentence.Word = word
		sentence.Sentence = sent
		sentence.ReviewCount = 0
		sentence.Status = dao.FAIL
		sentence.UserId = int64(sender.ID)
		dao.SentenceSave(ctx,client,sentence)
	}
	return nil
}



func VocabularyReview(sender *telebot.User,word string) (error){
	word = strings.ToLower(strings.Trim(word," "))
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	vocabulary,err:=dao.VocabularyGet(ctx,client,int64(sender.ID),word)
	if err!=nil{
		return err
	}
	if vocabulary.LearnStatus == dao.Learning {
		return errors.New(fmt.Sprintf("%s is reviewing\n",word))
	}
	dao.VocabularyUpdateLearnStatus(ctx,client,int64(sender.ID),word,dao.Learning)
	dao.SentenceUpdateStatusByWord(ctx,client,int64(sender.ID),word,dao.FAIL)
	return dao.UserStartInput(ctx,client,int64(sender.ID),dao.Review)
}


//only delete the sentence
func VocabularyReviewReceive(sender *telebot.User,message string) (ReviewResult,error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return ReviewResult{},err
	}
	defer dao.CloseClient(ctx,client)
	words := strings.Split(message,":")
	word:=strings.ToLower(strings.Trim(words[0]," "))
	sent:=strings.ToLower(strings.Trim(words[1]," "))

	_,err=dao.VocabularyGet(ctx,client,int64(sender.ID),word)
	if err != nil{
		return ReviewResult{},err
	}

	sentence,err:=dao.SentenceFindByWord(ctx,client,int64(sender.ID),word)
	if err != nil{
		return ReviewResult{},err
	}

	reviewResult := ReviewResult{}
	reviewResult.Word = word
	reviewResult.Total = len(sentence)
	reviewResult.Result = false

	for _,v:=range sentence {
		if strings.ToLower(strings.Trim(v.Sentence," ")) == sent {
			dao.SentenceUpdateStatus(ctx,client,v.Id,dao.PASS)
			reviewResult.Pass = reviewResult.Pass+1
			reviewResult.Result = true
			continue
		}
		if v.Status == dao.PASS {
			reviewResult.Pass = reviewResult.Pass+1
		}
	}

	if reviewResult.Result {
		return reviewResult,nil
	}
	return reviewResult,errors.New(fmt.Sprintf("Review word %s failed",word))
}

func VocabularyUpdate(sender *telebot.User) (error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	return dao.UserStartInput(ctx,client,int64(sender.ID),dao.Update)
}

//only delete the sentence
func VocabularyUpdateReceive(sender *telebot.User,message string) (error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	words := strings.Split(message,":")
	word:=strings.ToLower(strings.Trim(words[0]," "))
	sent:=strings.ToLower(strings.Trim(words[1]," "))
	_,err=dao.VocabularyGet(ctx,client,int64(sender.ID),word)
	if err != nil{
		return err
	}
	_,err=dao.SentencesDeleteBySentence(ctx,client,int64(sender.ID),word,sent)
	if err != nil{
		return err
	}
	return nil
}

func VocabularyEnd(sender *telebot.User) (error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	return dao.UserEndInput(ctx,client,int64(sender.ID))
}
func VocabularyEndReview(sender *telebot.User,word string) (ReviewResult,error){
	word = strings.ToLower(strings.Trim(word," "))

	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return ReviewResult{},err
	}
	defer dao.CloseClient(ctx,client)

	_,err=dao.VocabularyGet(ctx,client,int64(sender.ID),word)
	if err != nil{
		return ReviewResult{},err
	}

	sentence,err:=dao.SentenceFindByWord(ctx,client,int64(sender.ID),word)
	if err != nil{
		return ReviewResult{},err
	}

	reviewResult := ReviewResult{}
	reviewResult.Word = word
	reviewResult.Total = len(sentence)

	for _,v:=range sentence{
		if v.Status == dao.PASS {
			reviewResult.Pass = reviewResult.Pass+1
		}
	}

	if reviewResult.Total == reviewResult.Pass {
		dao.VocabularyUpdateLearnStatus(ctx,client,int64(sender.ID),word,dao.Finished)
		dao.VocabularyUpdateStatus(ctx,client,int64(sender.ID),word,dao.PASS)
		err= dao.UserEndInput(ctx,client,int64(sender.ID))
		return reviewResult,err
	}
	err= errors.New("You need to pass this review ["+word+"]")
	return reviewResult,err
}

func SentenceFindByWord(userId int, word string) ([]dao.Sentences,error){
	word = strings.ToLower(strings.Trim(word," "))
	ctx,client,err:=dao.GetClient()
	sentenceList:=[]dao.Sentences{}
	if err!=nil{
		return sentenceList,err
	}
	defer dao.CloseClient(ctx,client)

	_,err=dao.VocabularyGet(ctx,client,int64(userId),word)
	if err != nil{
		return sentenceList,err
	}

	return dao.SentenceFindByWord(ctx,client,int64(userId),word)
}