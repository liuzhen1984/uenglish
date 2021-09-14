package domain

import (
	"github.com/magiconair/properties"
	"sync"
)

type config struct {
	MongodbHost			string
	MongodbUsername		string
	MongodbPassword		string
	MongodbDatase		string
	BotToken			string
	BotSchedule			string
}

/**

mongodb.address=triple-cluster-db.eszal.mongodb.net
mongodb.username=triples
mongodb.password=triples123456
mongodb.database=stock
mongodb.max.pool.size=200

telegram.bot.token=1956917595:AAF3sOP-l6xgPwlKuUZbtdrfgGTu0bdRxXg

telegram.bot.schedule= 0 * * * * *

 */

var globalConfig *config
var once sync.Once


func LoadProperties(propertyFile ...string) *config {
	once.Do(func() {
		globalConfig = &config{}
		properties:= properties.MustLoadFile(propertyFile[0],properties.UTF8)
		globalConfig.MongodbPassword = properties.GetString("mongodb.password","")
		globalConfig.MongodbHost = properties.GetString("mongodb.address","")
		globalConfig.MongodbUsername = properties.GetString("mongodb.username","")
		globalConfig.MongodbDatase = properties.GetString("mongodb.database","")

		globalConfig.BotToken = properties.GetString("telegram.bot.token","")
		globalConfig.BotSchedule = properties.GetString("telegram.bot.schedule","")
	})
	return globalConfig
}
