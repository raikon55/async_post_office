package main

import (
	"file_processor/src/webservice"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.POST("/message", webservice.TriggerMessages)
	router.GET("/document", webservice.ReturnSomething)

	router.Run("127.0.0.1:8080")
}
