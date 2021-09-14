package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	domain "telegram_bot/com/trples/bot/config"
	"time"
)

func GetClient() (context.Context,*mongo.Client,error){
	config:=domain.LoadProperties()
	// Replace the uri string with your MongoDB deployment's connection string.
	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/test?w=majority",config.MongodbUsername,config.MongodbPassword,config.MongodbHost)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return ctx,client,err
}

func CloseClient(ctx context.Context,client *mongo.Client) error {
	return client.Disconnect(ctx)
}