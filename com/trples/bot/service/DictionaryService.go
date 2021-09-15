package service

import (
	"gopkg.in/tucnak/telebot.v2"
	"strings"
	"telegram_bot/com/trples/bot/dao"
)



func DictionaryTranslateStart(sender *telebot.User) (error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	return dao.UserStartInput(ctx,client,int64(sender.ID),dao.Translate)
}
func DictionaryLongman(userId int,word string) string  {
	word = strings.ToLower(strings.Trim(word," "))
	return word
}

func DictionaryTranslate(userId int,text string) (string,error)  {
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return text,err
	}
	defer dao.CloseClient(ctx,client)

	_,err =dao.UserGet(ctx,client,int64(userId))
	if err!=nil{
		return text,err
	}
	return text,err
}
