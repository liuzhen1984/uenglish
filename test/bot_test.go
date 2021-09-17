package test

import (
	"fmt"
	"strings"
	"telegram_bot/com/trples/bot/service"
	"testing"
)

func TestSend(t *testing.T){
	service.GetBot()
	service.Send(1470647544,"test")
}
func TestStrings(t *testing.T){
	s:=strings.SplitN("asdfa:bsfasdf:cdfasdf",":",0)
	print(s)
	s1:=strings.SplitN("asdfa:bsfasdf:cdfasdf",":",1)
	print(s1)
	s2:=strings.SplitN("asdfa:bsfasdf:cdfasdf",":",2)
	print(s2)
	s3:=strings.SplitN("asdfa:bsfasdf:cdfasdf",":",3)
	print(s3)
	s3=strings.SplitN("asdfa:bsfasdf:cdfasdf",":",4)
	print(s3)
}
func print(sList []string){
	for i,v:=range sList{
		fmt.Printf("index=%d,value=%s\n",i,v)
	}
	fmt.Printf("=============\n\n")

}

