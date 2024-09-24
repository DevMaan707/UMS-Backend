package main

import (
	"DevMaan707/UMS/helpers"

	"github.com/gin-gonic/gin"
)

func main() {
	//dbConnection := db.ConnectMySQL()

	router := gin.Default()
	// router.GET("/live/generate-classes", func(c *gin.Context) {
	// 	helpers.GenerateClassesLive(c, dbConnection)
	// })

	router.POST("/test/generate-classes", func(c *gin.Context) {
		helpers.GenerateClassesTest(c)
	})
	router.Run()
}
