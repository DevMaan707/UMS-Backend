package main

import (
	"DevMaan707/UMS/helpers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/test/generate-classes", func(c *gin.Context) {
		helpers.AssignRoomsForExams(c)
	})
	router.GET("/assignments", helpers.GetAllAssignments)
	router.GET("/assignments/:student_id", helpers.GetStudentSpecificAssignment)
	router.Run()
}
