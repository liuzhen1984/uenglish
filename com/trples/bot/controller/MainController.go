package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func WebBoot()  {
	r := gin.Default()
	r.Delims("{{", "}}")
	r.Static("/assets", "./")
	r.LoadHTMLFiles("./html/main.html")
	r.GET("/v1/bot", hello)
	r.POST("/v1/bot", hello)
	r.Run("0.0.0.0:8081") // listen and serve on 0.0.0.0:8080
}


func hello(c *gin.Context){
	fmt.Println(c)
	c.JSON(200, gin.H{
		"data": "hello you",
	})
}

