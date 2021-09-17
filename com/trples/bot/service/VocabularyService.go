package service

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"strings"
	domain "telegram_bot/com/trples/bot/config"
	"telegram_bot/com/trples/bot/dao"
	"time"
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
		vocabulary.Period = 2
		vocabulary.ReminderCount = 0
		vocabulary.LatestReviewAt = time.Now().UnixMilli() + 2*60*60*1000
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



func VocabularyReview(userId int,word string) ([]string,error){
	word = strings.ToLower(strings.Trim(word," "))
	results:=[]string{word}
	if word==""{
		return VocabularyReviewAll(userId)
	}
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return results,err
	}
	defer dao.CloseClient(ctx,client)
	vocabulary,err:=dao.VocabularyGet(ctx,client,int64(userId),word)
	if err!=nil{
		return results,err
	}
	if vocabulary.LearnStatus == dao.Learning {
		return results,errors.New(fmt.Sprintf("%s is reviewing\n",word))
	}
	dao.VocabularyUpdateLearnStatus(ctx,client,int64(userId),word,dao.Learning)
	dao.SentenceUpdateStatusByWord(ctx,client,int64(userId),word,dao.FAIL)
	return results,dao.UserStartInput(ctx,client,int64(userId),dao.Review)
}

func VocabularyReviewAll(userId int) ([]string,error){
	ctx,client,err:=dao.GetClient()
	var vLits []string

	if err!=nil{
		return vLits,err
	}
	defer dao.CloseClient(ctx,client)
	err= dao.UserStartInput(ctx,client,int64(userId),dao.Review)
	result,err:=dao.VocabularyFindByReview(ctx,client,int64(userId))
	if err==nil{
		for _,v:=range result {
			dao.SentenceUpdateStatusByWord(ctx,client,int64(userId),v.Word,dao.FAIL)
			vLits = append(vLits,v.Word)
		}
	}

	return vLits,err
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
func VocabularyUpdateReceive(sender *telebot.User,message string) (int64,error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return 0,err
	}
	defer dao.CloseClient(ctx,client)
	words := strings.Split(message,":")
	word:=strings.ToLower(strings.Trim(words[0]," "))
	sent:=strings.ToLower(strings.Trim(words[1]," "))
	_,err=dao.VocabularyGet(ctx,client,int64(sender.ID),word)
	if err != nil{
		return 0,err
	}
	count,err:=dao.SentencesDeleteBySentence(ctx,client,int64(sender.ID),word,sent)

	return count,err
}

func VocabularyEnd(userId int) (error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	return dao.UserEndInput(ctx,client,int64(userId))
}
func VocabularyEndReview(userId int,word string) (ReviewResult,error){
	word = strings.ToLower(strings.Trim(word," "))

	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return ReviewResult{},err
	}
	defer dao.CloseClient(ctx,client)

	vocab,err:=dao.VocabularyGet(ctx,client,int64(userId),word)
	if err != nil{
		return ReviewResult{},err
	}

	sentence,err:=dao.SentenceFindByWord(ctx,client,int64(userId),word)
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
		dao.VocabularyReviewCompleted(ctx,client,vocab)
		err= dao.UserEndInput(ctx,client,int64(userId))
		return reviewResult,err
	}
	err= errors.New("You need to pass this review ["+word+"]")
	return reviewResult,err
}

func VocabularyEndAllReview(userId int) (ReviewResult,error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return ReviewResult{},err
	}
	defer dao.CloseClient(ctx,client)

	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(dao.Collection_Vocabulary)
	ctime:=time.Now().UnixMilli()
	where:= bson.D{
		{"user_id",int64(userId)},
		{
			"$or",bson.A{
			bson.D{{"learn_status",dao.Learning}},
			bson.D{{"learn_status",nil}},
		},
		},
		{
			"$or",bson.A{
			bson.D{{"is_remember",false}},
			bson.D{{"is_remember",nil}},
		},
		},
		{"latest_review_at",bson.D{{"$lt", ctime}},
		},
	}
	cursor,err:=collection.Find(ctx,where)
	var count int64
	reviewResult := ReviewResult{}
	if(err!=nil){
		fmt.Println(err)
		return reviewResult,err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		reviewResult.Total = reviewResult.Total+1
		var vocabulary dao.Vocabulary
		if err := cursor.Decode(&vocabulary); err != nil {
			log.Fatal(err)
			continue
		}
		_,err=VocabularyEndReview(userId,vocabulary.Word)
		if err!=nil{
			reviewResult.Pass = reviewResult.Pass+1
			continue
		}
		if reviewResult.Word!=""{
			reviewResult.Word = reviewResult.Word +", "+vocabulary.Word
		} else{
			reviewResult.Word = vocabulary.Word
		}
	}
	fmt.Printf("user count %d \n",count)
	reviewResult.Pass = reviewResult.Total-reviewResult.Pass
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


func VocabularyFindByUserId(userId int) ([]dao.Vocabulary,error){
	ctx,client,err:=dao.GetClient()
	vocabularyList:=[]dao.Vocabulary{}
	if err!=nil{
		return vocabularyList,err
	}
	defer dao.CloseClient(ctx,client)
	where:= bson.D{
		{"user_id",userId},
		//{
		//	"$or",bson.A{
		//	bson.D{{"is_remember",false}},
		//	bson.D{{"is_remember",nil}},
		//},
		//},
		//{"learn_status",Learning},
	}
	return dao.VocabularyFind(ctx,client,where)
}

func VocabularyRemember(userId int,word string) error{
	word = strings.ToLower(strings.Trim(word," "))

	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	_,err=dao.VocabularyGet(ctx,client,int64(userId),word)
	if err!=nil{
		return err
	}
	return dao.VocabularyRemember(ctx,client,int64(userId),word,true)
}

func VocabularyReset(userId int,word string) error{
	word = strings.ToLower(strings.Trim(word," "))

	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	vocabulary,err:=dao.VocabularyGet(ctx,client,int64(userId),word)
	if err!=nil{
		return err
	}
	return dao.VocabularyUpdatePeriod(ctx,client,int64(userId),word,vocabulary.Period,0)
}

func VocabularyCheck(userId int) ([]string,error){

	ctx,client,err:=dao.GetClient()
	var vLits []string

	if err!=nil{
		return vLits,err
	}
	defer dao.CloseClient(ctx,client)
	result,err:=dao.VocabularyFindByReview(ctx,client,int64(userId))
	if err!=nil{
		return vLits,err
	}
	for _,v:=range result {
		vLits = append(vLits,v.Word)
	}
	return vLits,err
}