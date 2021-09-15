package main

import (
	"fmt"
	"os"
	domain "telegram_bot/com/trples/bot/config"
	"telegram_bot/com/trples/bot/service"
)

func main(){
	fileName := "/Users/zliu/work/golang/uenglish/resource/config.properties"
	if len(os.Args)>1 {
		fileName = os.Args[1]
		_,err:=os.ReadFile(fileName)
		if err!=nil{
			fmt.Println("this file don't exist")
		}
	}
	domain.LoadProperties(fileName)
	service.TaskRun()
	service.BotRoute()
	//web.WebBoot()
}
