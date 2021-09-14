package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	domain "telegram_bot/com/trples/bot/config"
	"time"
)

var Collection_Vocabulary = "vocabulary"
var Collection_Sentences = "sentences"

type WaitTypeEnum string
const (
	Add		WaitTypeEnum = "add"
	Review	WaitTypeEnum = "review"
	Update	WaitTypeEnum = "update"
	Translate	WaitTypeEnum = "translate"
)

type ReviewStatus string
const (
	PASS	ReviewStatus = "pass"
	FAIL 	ReviewStatus = "fail"
)

type LearnStatus string
const (
	Waiting		LearnStatus = "waiting"
	Learning 	LearnStatus = "learning"
	Finished 	LearnStatus = "finished"
)

type Vocabulary struct {
	Id   			 primitive.ObjectID		`bson:"_id,omitempty"`
	Word 		 	 string  				`bson:"word,omitempty"`  //unique
	LearnStatus		 LearnStatus 			`bson:"learn_status,omitempty"`  // waiting, learning, finished
	ReviewStatus     ReviewStatus 			`bson:"review_status,omitempty"` //pass, fail,nil
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
	Status			ReviewStatus			`bson:"status"` //pass, fail
	ReviewCount		int64					`bson:"review_count"`
	CreateAt	 	int64					`bson:"create_at,omitempty"`
	UpdatedAt		int64					`bson:"updated_at,omitempty"`
}
func SentencesDeleteBySentence(ctx context.Context, client *mongo.Client,userId int64,word string,sentence string) (int64,error){
	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Sentences)
	result,error := collection.DeleteMany(ctx,bson.M{"user_id":userId,"word":word,"sentence":sentence})
	if error!=nil{
		return 0,nil
	}
	return result.DeletedCount,nil
}
func SentencesDeleteByWord(ctx context.Context, client *mongo.Client,userId int64,word string) (int64,error){
	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Sentences)
	result,error := collection.DeleteMany(ctx,bson.M{"user_id":userId,"word":word})
	if error!=nil{
		return 0,nil
	}
	return result.DeletedCount,nil
}

func SentenceFindBySentence(ctx context.Context, client *mongo.Client,userId int64,word string,sentence string) ([]Sentences,error){
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Sentences)
	cursor,err:=collection.Find(ctx,bson.M{"user_id":userId,"word":word,"sentence":sentence})
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

func SentenceFindByWord(ctx context.Context, client *mongo.Client,userId int64,word string) ([]Sentences,error){
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Sentences)

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
	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Sentences)
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

	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Vocabulary)

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
	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Vocabulary)
	result,err:=collection.InsertOne(ctx,vocabulary)
	if(err!=nil){
		fmt.Println(err)
		return nil,err
	}
	return result.InsertedID,nil
}

func VocabularyDeleteByWord(ctx context.Context, client *mongo.Client,userId int64,word string) (int64,error){
	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Vocabulary)
	result,error := collection.DeleteMany(ctx,bson.M{"user_id":userId,"word":word})
	if error!=nil{
		return 0,nil
	}
	return result.DeletedCount,nil
}


//Receive new words, and sentences | update words and sentences | Review words
func SentenceUpdateStatus(ctx context.Context, client *mongo.Client,id primitive.ObjectID, status ReviewStatus)  error{
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Sentences)
	_,err:=collection.UpdateByID(ctx,id,bson.D{{"$set",bson.M{"status":status,"input_updated_at":time.Now().UnixMilli()}}})
	if(err!=nil){
		fmt.Println(err)
	}
	return err
}
//Receive new words, and sentences | update words and sentences | Review words
func SentenceUpdateStatusByWord(ctx context.Context, client *mongo.Client,userId int64,word string, status ReviewStatus)  error{
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Sentences)
	_,err:=collection.UpdateMany(ctx,bson.M{"user_id":userId,"word":word},bson.D{{"$set",bson.M{"status":status,"input_updated_at":time.Now().UnixMilli()}}})
	if(err!=nil){
		fmt.Println(err)
	}
	return err
}

func VocabularyUpdateStatus(ctx context.Context, client *mongo.Client,userId int64,word string, status ReviewStatus)  error{
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Vocabulary)
	_,err:=collection.UpdateOne(ctx,bson.M{"user_id":userId,"word":word},bson.D{{"$set",bson.M{"review_status":status,"input_updated_at":time.Now().UnixMilli()}}})
	if(err!=nil){
		fmt.Println(err)
	}
	return err
}

func VocabularyUpdateLearnStatus(ctx context.Context, client *mongo.Client,userId int64,word string, status LearnStatus)  error{
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Vocabulary)
	_,err:=collection.UpdateOne(ctx,bson.M{"user_id":userId,"word":word},bson.D{{"$set",bson.M{"learn_status":status,"input_updated_at":time.Now().UnixMilli()}}})
	if(err!=nil){
		fmt.Println(err)
	}
	return err
}
func VocabularyRemember(ctx context.Context, client *mongo.Client,userId int64,word string,isRemember bool)  error{
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Vocabulary)
	_,err:=collection.UpdateOne(ctx,bson.M{"user_id":userId,"word":word},bson.D{{"$set",bson.M{"is_remember":isRemember,"input_updated_at":time.Now().UnixMilli()}}})
	if(err!=nil){
		fmt.Println(err)
	}
	return err
}
/**
Period			 int64 					`bson:"period,omitempty"`  //unit is hour
ReminderCount	 int64  				`bson:"reminder_count,omitempty"`
 */
func VocabularyUpdatePeriod(ctx context.Context, client *mongo.Client,userId int64,word string, period int64,count int64)  error{
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database(domain.LoadProperties().MongodbDatase)
	collection := database.Collection(Collection_Vocabulary)
	_,err:=collection.UpdateOne(ctx,bson.M{"user_id":userId,"word":word},bson.D{{"$set",bson.M{"reminder_count":count,"period":period,"input_updated_at":time.Now().UnixMilli()}}})
	if(err!=nil){
		fmt.Println(err)
	}
	return err
}