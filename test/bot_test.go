package test

import (
	"telegram_bot/com/trples/bot/service"
	"testing"
)

func TestSend(t *testing.T){
	service.GetBot()
	service.Send(1470647544,"test")
}
