package main

import (
	"DevMaan707/UMS/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//router.Use(middleware.JWTAuthMiddleware())

	router.POST("/test/generate-classes", handlers.AssignRoomsForExams)
	router.GET("/assignments", handlers.GetAllAssignments)
	router.GET("/assignments/:student_id", handlers.GetStudentSpecificAssignment)
	router.GET("/generatepdfbytoe", handlers.GeneratePDFByTOE)

	router.Run()
}
