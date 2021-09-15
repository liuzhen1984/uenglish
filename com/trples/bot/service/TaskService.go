package service

import (
	"context"
	"fmt"
	"github.com/robfig/cron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	domain "telegram_bot/com/trples/bot/config"
	"telegram_bot/com/trples/bot/dao"
	"time"
)
func TaskRun(){
	run:=cron.New()
	run.AddFunc(domain.LoadProperties().BotSchedule, UserChecking)
	fmt.Println("start task check ... "+domain.LoadProperties().BotSchedule)
	run.Start()
}

func UserChecking() {
	fmt.Println("start check...")
	ctx,client,err:=dao.GetClient()
	if err!=nil{
		fmt.Println(err)
		return
	}
	defer dao.CloseClient(ctx,client)
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	currentTime := time.Now().UnixMilli()

	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(dao.Collection_user)

	opts := options.Find()
	opts.SetSort(bson.D{{"delay_to_time", -1}})

	where:= bson.D{
		{
			"$or",bson.A{
				bson.D{{"is_bot",false}},
				bson.D{{"is_bot",nil}},
			},
		},
		{"is_enable",true},
		{
			"$or",bson.A{
				bson.D{{"reminder_at",nil}},
				bson.D{{"reminder_at",bson.D{{"$lt",time.Now().UnixMilli()}}}},
			},
		},
		{
			"$or",bson.A{
			bson.D{{"is_delay", nil}},
			bson.D{{"is_delay", false}},
			bson.D{{"is_delay", true}, {"delay_to_time", bson.D{{"$gt", 0}}}, {"delay_to_time", bson.D{{"$lt", currentTime}}},
			},
		}},
	}
	cursor,err:=collection.Find(ctx,where,opts)

	var count int64
	if(err!=nil){
		fmt.Println(err)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var userConfig dao.UserConfig
		if err := cursor.Decode(&userConfig); err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Printf("start check the user %s %s\n",userConfig.FirstName,userConfig.LastName)
		vCount,err:=checkVocabulary(ctx,client,userConfig.UserID)
		if err!=nil{
			continue
		}
		if vCount>0{
			dao.UserUpdateByUserId(ctx,client,userConfig.UserID,bson.M{"reminder_at":time.Now().UnixMilli()+60*10*1000})
			Send(int(userConfig.UserID),fmt.Sprintf("You have %d vocabularies to review",vCount))
		}
		count = count +1
	}
	fmt.Printf("user count %d \n",count)
	return
}

func checkVocabulary(ctx context.Context, client *mongo.Client,userId int64)(int64,error){
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(dao.Collection_Vocabulary)
	ctime:=time.Now().UnixMilli()
	where:= bson.D{
			{"user_id",userId},
			{
				"$or",bson.A{
					bson.D{{"learn_status",dao.Waiting}},
					bson.D{{"learn_status",dao.Finished}},
					bson.D{{"is_remember",nil}},
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
	result,err:=collection.UpdateMany(ctx,where,bson.D{{"$set",bson.M{"learn_status":dao.Learning}}})
	if(err!=nil){
		fmt.Println(err)
		return 0,err
	}
	return result.ModifiedCount,err
}

