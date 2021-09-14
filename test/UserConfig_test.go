package test

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	domain "telegram_bot/com/trples/bot/config"
	"telegram_bot/com/trples/bot/dao"
	"testing"
	"time"
)

func TestSaveUser(t *testing.T){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		panic(err)
	}
	defer dao.CloseClient(ctx,client)
	userConfig := dao.UserConfig{}
	userConfig.UserID=123
	userConfig.Username = "zhen"
	userConfig.FirstName ="z"
	userConfig.LastName ="liu"
	userConfig.LanguageCode ="us"
	userConfig.IsEnable =true
	dao.UserSave(ctx,client,userConfig)
}

func TestUpdateUser(t *testing.T)  {
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		panic(err)
	}
	defer dao.CloseClient(ctx,client)

	dao.UserUpdateByUserId(ctx,client,123,bson.M{"first_name":"zhen"})
}

func TestUserStop(t *testing.T)  {
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		panic(err)
	}
	defer dao.CloseClient(ctx,client)

	dao.UserStopped(ctx,client,123)
}

func TestUserStart(t *testing.T)  {
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		panic(err)
	}
	defer dao.CloseClient(ctx,client)

	dao.UserStart(ctx,client,123)
}

func TestUserDelay(t *testing.T)  {
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		panic(err)
	}
	defer dao.CloseClient(ctx,client)

	dao.UserDelay(ctx,client,123,true,time.Now().UnixMilli())
}
func TestUserDiableDelay(t *testing.T)  {
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		panic(err)
	}
	defer dao.CloseClient(ctx,client)
	dao.UserDelay(ctx,client,123,false,-1)
}

func TestUserUpdateEmail(t *testing.T)  {
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		panic(err)
	}
	defer dao.CloseClient(ctx,client)
	dao.UserSetEmail(ctx,client,123,"liuzhen1984@gmail.com")
	userConfig,err:=dao.UserGet(ctx,client,123)
	if err!=nil{
		panic(err)
	}

	if userConfig.Email=="liuzhen1984@gmail.com" {
		fmt.Println("ok")
	}else{
		panic(errors.New(userConfig.Email))
	}

}

func TestFindUserByDelay(t *testing.T){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		panic(err)
	}
	defer dao.CloseClient(ctx,client)
	fmt.Println(dao.UserFindByDelay(ctx,client))
}

func TestDeleteVocabulary(t *testing.T){
	fileName := "/Users/zliu/work/golang/telegram_bot/resource/config.properties"
	domain.LoadProperties(fileName)
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		panic(err)
	}
	defer dao.CloseClient(ctx,client)
	dao.VocabularyDeleteByWord(ctx,client,1470647544,"e")
	dao.SentencesDeleteByWord(ctx,client,1470647544,"e")
}