package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type Vocabulary struct {
	Id   			 primitive.ObjectID		`bson:"_id,omitempty"`
	Word 		 	 string  				`bson:"word,omitempty"`  //unique
	LearnStatus		 string 				`bson:"learn_status,omitempty"`  // waiting, learning, finished
	ReviewStatus     string 				`bson:"review_status,omitempty"` //pass, fail,nil
	IsRemember		 bool 					`bson:"is_remember,omitempty"`
	UserId 			 int64 					`bson:"user_id,omitempty"`
	Period			 int64 					`bson:"period,omitempty"`  //unit is hour
	ReminderCount	 int64  				`bson:"reminder_count,omitempty"`
	CreateAt	 	 int64					`bson:"create_at,omitempty"`
	UpdatedAt		 int64					`bson:"updated_at,omitempty"`
}

type Sentences struct {
	Id   			primitive.ObjectID		`bson:"_id,omitempty"`
	Word			string					`bson:"word,omitempty"`
	UserId 			int64 					`bson:"user_id,omitempty"`
	Sentence		string					`bson:"sentence,omitempty"`
	Status			string					`bson:"status"` //pass, fail
	ReviewCount		int64					`bson:"review_count"`
	CreateAt	 	int64					`bson:"create_at,omitempty"`
	UpdatedAt		int64					`bson:"updated_at,omitempty"`
}

func SentencesDeleteByWord(ctx context.Context, client *mongo.Client,userId int64,word string) (int64,error){
	database := client.Database("telegram_bot")
	collection := database.Collection("sentences")
	result,error := collection.DeleteMany(ctx,bson.M{"user_id":userId,"word":word})
	if error!=nil{
		return 0,nil
	}
	return result.DeletedCount,nil
}

func SentenceFindByWord(ctx context.Context, client *mongo.Client,userId int64,word string) ([]Sentences,error){
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	database := client.Database("telegram_bot")
	collection := database.Collection("sentences")

	cursor,err:=collection.Find(ctx,bson.M{"user_id":userId,"word":word})
	var sentences []Sentences
	if err!=nil{
		return sentences,err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var sentence Sentences
		if err := cursor.Decode(&sentence); err != nil {
			log.Fatal(err)
			continue
		}
		sentences = append(sentences, sentence)
		if len(sentences)>10{
			break
		}
	}
	return sentences,nil
}
func SentenceSave(ctx context.Context, client *mongo.Client,sentence Sentences) (interface{},error){
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	sentence.CreateAt = time.Now().UnixMilli()
	sentence.UpdatedAt = sentence.CreateAt
	//userConfig.IsEnable = true
	//userConfig.IsDelay = false
	//userConfig.DelayToTime = -1
	database := client.Database("telegram_bot")
	collection := database.Collection("sentences")
	result,err:=collection.InsertOne(ctx,sentence)
	if(err!=nil){
		fmt.Println(err)
		return nil,err
	}
	return result.InsertedID,nil
}


func VocabularyGet(ctx context.Context, client *mongo.Client,userId int64,word string) (Vocabulary,error){
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	database := client.Database("telegram_bot")
	collection := database.Collection("vocabulary")

	result:=collection.FindOne(ctx,bson.M{"user_id":userId,"word":word})
	var vocabulary Vocabulary
	err:=result.Decode(&vocabulary)
	return vocabulary,err
}

func VocabularySave(ctx context.Context, client *mongo.Client,vocabulary Vocabulary) (interface{},error){
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	vocabulary.CreateAt = time.Now().UnixMilli()
	vocabulary.UpdatedAt = vocabulary.CreateAt
	//userConfig.IsEnable = true
	//userConfig.IsDelay = false
	//userConfig.DelayToTime = -1
	database := client.Database("telegram_bot")
	collection := database.Collection("vocabulary")
	result,err:=collection.InsertOne(ctx,vocabulary)
	if(err!=nil){
		fmt.Println(err)
		return nil,err
	}
	return result.InsertedID,nil
}

func VocabularyDeleteByWord(ctx context.Context, client *mongo.Client,userId int64,word string) (int64,error){
	database := client.Database("telegram_bot")
	collection := database.Collection("vocabulary")
	result,error := collection.DeleteMany(ctx,bson.M{"user_id":userId,"word":word})
	if error!=nil{
		return 0,nil
	}
	return result.DeletedCount,nil
}