package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)


type UserConfig struct {
	Id   		 primitive.ObjectID		`bson:"_id,omitempty"`
	UserID 		 int64  `bson:"user_id,omitempty"`
	FirstName    string `bson:"first_name,omitempty"`
	LastName     string `bson:"last_name,omitempty"`
	Username     string `bson:"username,omitempty"`
	Email	     string `bson:"email,omitempty"`
	LanguageCode string `bson:"language_code,omitempty"`
	IsBot        bool   `bson:"is_bot,omitempty"`
	IsEnable	 bool	`bson:"is_enable,omitempty"`
	IsDelay		 bool	`bson:"is_delay,omitempty"`
	DelayToTime	 int64	`bson:"delay_to_time,omitempty"`  //Scheduled time for next activation
	Schedule	 string	`bson:"schedule_time,omitempty"`  //day|hour
	IsInput		 bool	`bson:"is_input,omitempty"`
	//Receive new words, and sentences | update words and sentences | Review words
	WaitType	string	  `bson:"wait_type,omitempty"`
	InputUpdatedAt	int64 `bson:"input_updated_at,omitempty"`
	CreateAt	 int64	`bson:"create_at,omitempty"`
	UpdatedAt	 int64	`bson:"updated_at,omitempty"`
}

func UserFindByDelay(ctx context.Context, client *mongo.Client) ([]UserConfig,error){
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	currentTime := time.Now().UnixMilli()
	database := client.Database("telegram_bot")
	collection := database.Collection("user_config")

	opts := options.Find()
	opts.SetSort(bson.D{{"delay_to_time", -1}})

	cursor,err:=collection.Find(ctx,bson.D{{"is_delay",true},{"delay_to_time",bson.D{{"$gt",0}}},{"delay_to_time",bson.D{{"$lt",currentTime}}}},opts)

	var userConfigList []UserConfig
	if(err!=nil){
		fmt.Println(err)
		return userConfigList,err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var userConfig UserConfig
		if err := cursor.Decode(&userConfig); err != nil {
			log.Fatal(err)
			continue
		}
		userConfigList = append(userConfigList, userConfig)
	}
	return userConfigList,err
}


func UserGet(ctx context.Context, client *mongo.Client,userId int64) (UserConfig,error){
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	database := client.Database("telegram_bot")
	collection := database.Collection("user_config")

	result:=collection.FindOne(ctx,bson.M{"user_id":userId})

	var userConfig UserConfig

	err:=result.Decode(&userConfig)
	return userConfig,err
}

func UserSave(ctx context.Context, client *mongo.Client,userConfig UserConfig) (interface{},error){
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	userConfig.CreateAt = time.Now().UnixMilli()
	userConfig.UpdatedAt = userConfig.CreateAt
	//userConfig.IsEnable = true
	//userConfig.IsDelay = false
	//userConfig.DelayToTime = -1
	database := client.Database("telegram_bot")
	collection := database.Collection("user_config")
	result,err:=collection.InsertOne(ctx,userConfig)
	if(err!=nil){
		fmt.Println(err)
		return nil,err
	}
	return result.InsertedID,nil
}

func UserUpdateByUserId(ctx context.Context, client *mongo.Client,userId int64,userConfig bson.M) (int64,error){
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database("telegram_bot")
	collection := database.Collection("user_config")
	userConfig["updated_at"] = time.Now().UnixMilli()
	result,err:=collection.UpdateOne(ctx,bson.M{"user_id":userId},bson.D{{"$set",userConfig}})
	if(err!=nil){
		fmt.Println(err)
		return 0,err
	}
	return result.ModifiedCount,nil
}

func UserSetEmail(ctx context.Context, client *mongo.Client,userId int64,email string) error{
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database("telegram_bot")
	collection := database.Collection("user_config")
	_,err:=collection.UpdateOne(ctx,bson.M{"user_id":userId},bson.D{{"$set",bson.M{"email":email,"updated_at":time.Now().UnixMilli()}}})
	if(err!=nil){
		fmt.Println(err)
	}
	return err
}


func UserStopped(ctx context.Context, client *mongo.Client,userId int64)  error{
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database("telegram_bot")
	collection := database.Collection("user_config")
	_,err:=collection.UpdateOne(ctx,bson.M{"user_id":userId},bson.D{{"$set",bson.M{"is_enable":false,"updated_at":time.Now().UnixMilli()}}})
	if(err!=nil){
		fmt.Println(err)
	}
	return err
}

func UserStart(ctx context.Context, client *mongo.Client,userId int64)  error{
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database("telegram_bot")
	collection := database.Collection("user_config")
	_,err:=collection.UpdateOne(ctx,bson.M{"user_id":userId},bson.D{{"$set",bson.M{"is_enable":true,"updated_at":time.Now().UnixMilli()}}})
	if(err!=nil){
		fmt.Println(err)
	}
	return err
}

// update delay
func UserDelay(ctx context.Context, client *mongo.Client,userId int64,delay bool,delayTime int64)  error{
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database("telegram_bot")
	collection := database.Collection("user_config")
	_,err:=collection.UpdateOne(ctx,bson.M{"user_id":userId},bson.D{{"$set",bson.M{"is_delay":delay,"delay_to_time":delayTime,"updated_at":time.Now().UnixMilli()}}})
	if(err!=nil){
		fmt.Println(err)
	}
	return err
}

//Receive new words, and sentences | update words and sentences | Review words
func UserStartInput(ctx context.Context, client *mongo.Client,userId int64,waitType string)  error{
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database("telegram_bot")
	collection := database.Collection("user_config")
	_,err:=collection.UpdateOne(ctx,bson.M{"user_id":userId},bson.D{{"$set",bson.M{"is_input":true,"wait_type":waitType,"input_updated_at":time.Now().UnixMilli()}}})
	if(err!=nil){
		fmt.Println(err)
	}
	return err
}

func UserEndInput(ctx context.Context, client *mongo.Client,userId int64)  error{
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database("telegram_bot")
	collection := database.Collection("user_config")
	_,err:=collection.UpdateOne(ctx,bson.M{"user_id":userId},bson.D{{"$set",bson.M{"is_input":false,"wait_type":""}}})
	if(err!=nil){
		fmt.Println(err)
	}
	return err
}