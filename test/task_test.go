package test

import (
	"telegram_bot/com/trples/bot/dao"
	"telegram_bot/com/trples/bot/service"
	"testing"
)

func TestCheckTask(t *testing.T){
	service.UserChecking()
}

func TestReviewCompleted(t *testing.T){
	ctx,client,_:=dao.GetClient()

	defer dao.CloseClient(ctx,client)

	vocab,_:=dao.VocabularyGet(ctx,client,1470647544,"ass")
	dao.VocabularyReviewCompleted(ctx,client,vocab)
}

func TestEndReviewAll(t *testing.T){
	service.VocabularyEndAllReview(1470647544)
}