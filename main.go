package main

import (
	"DevMaan707/UMS/db"
	helpers "DevMaan707/UMS/helpers"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	dynamoClient := db.ConnectDynamoDB()
	mongoClient, err := db.ConnectMongoDB()

	if err != nil {
		log.Fatal("Error occurred while connecting to MongoDB")
	}
	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
		fmt.Println("Success")
	})

	router.GET(
		"/create-tables",
		func(c *gin.Context) {
			helpers.DefineTables(dynamoClient)
		},
	)

	router.Run()
}
