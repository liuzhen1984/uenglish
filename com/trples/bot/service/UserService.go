package service

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/tucnak/telebot.v2"
	"telegram_bot/com/trples/bot/dao"
)

func UserGet(userId int64) dao.UserConfig{
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		panic(err)
	}
	defer dao.CloseClient(ctx,client)
	user,_:=dao.UserGet(ctx,client,userId)
	return user
}

/**
schedule = day|hour
 */
func UserUpdateSchedule(sender *telebot.User,schedule string)  error{
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	_,err=dao.UserUpdateByUserId(ctx,client,int64(sender.ID),bson.M{"schedule":schedule})
	return err
}

func UserUpdateLang(sender *telebot.User,lang string)  error{
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	err=dao.UserSetLang(ctx,client,int64(sender.ID),lang)
	return err
}

func UserStop(sender *telebot.User)  error{
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	return dao.UserStopped(ctx,client,int64(sender.ID))
}
func UserReStart(sender *telebot.User) error {
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	return dao.UserStart(ctx,client,int64(sender.ID))
}

func UserNewOrExisting(sender *telebot.User) (error){
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		return err
	}
	defer dao.CloseClient(ctx,client)
	user,err:=dao.UserGet(ctx,client,int64(sender.ID))
	if err != nil{
		fmt.Printf("new account id=%d username=%s fname,lname=%s,%s,isBot=%b",sender.ID,sender.Username,sender.FirstName,sender.LastName,sender.IsBot)
		user=dao.UserConfig{}
		user.UserID = int64(sender.ID)
		user.Username = sender.Username
		user.LanguageCode = sender.LanguageCode
		user.LastName = sender.LastName
		user.FirstName = sender.FirstName
		user.IsBot = sender.IsBot
		user.IsEnable = true
		user.IsDelay = false
		user.DelayToTime = -1
		user.ReminderAt = -1
		dao.UserSave(ctx,client,user)
		return err
	}
	return dao.UserStart(ctx,client,int64(sender.ID))
}