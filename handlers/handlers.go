package handlers

import (
	"DevMaan707/UMS/helpers"
	"DevMaan707/UMS/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AssignRoomsForExams(c *gin.Context) {
	var params models.Params

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	toe, err := time.Parse(time.RFC3339, params.TOE)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Time of Exam format"})
		return
	}
	doe, err := time.ParseDuration(params.DOE)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Duration of Exam format"})
		return
	}

	blocks, classes := helpers.GenerateTestData()

	selectedRooms := []helpers.Room{}
	for _, block := range params.Blocks {
		if rooms, found := blocks[block]; found {
			selectedRooms = append(selectedRooms, rooms...)
		}
	}

	selectedStudents := make(map[string][]string)
	for _, branch := range params.Branches {
		if branchClasses, found := classes[branch]; found {
			for _, class := range branchClasses {
				if helpers.ContainsInt(params.Years, class.Year) {
					selectedStudents[branch] = append(selectedStudents[branch], class.StudentIDs...)
				}
			}
		}
	}

	if params.InternalShuffle {
		helpers.ShuffleStudents(selectedStudents)
	}

	assignments := helpers.GenerateExamAssignments("exam", selectedRooms, selectedStudents, params, toe, doe)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Exam Room Assignments",
		"assignments": assignments,
	})
}

func GetAllAssignments(c *gin.Context) {

	toeStr := c.Query("toe")
	toe, err := time.Parse(time.RFC3339, toeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Time of Exam format"})
		return
	}

	filteredAssignments, _ := helpers.FetchAssignmentsByTime(toe)
	c.JSON(http.StatusOK, gin.H{
		"assignments": filteredAssignments,
	})
}

func GetStudentSpecificAssignment(c *gin.Context) {
	toeStr := c.Query("toe")
	studentID := c.Param("student_id")
	toe, err := time.Parse(time.RFC3339, toeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Time of Exam format"})
		return
	}

	roomAssignment, err := helpers.GetStudentAssignment(studentID, toe)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Assignment not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"room_number": roomAssignment.RoomNumber,
			"details":     roomAssignment.Details,
			"toe":         roomAssignment.Toe,
			"block":       roomAssignment.Block,
		})
	}
}
func GeneratePDFByTOE(c *gin.Context) {
	toeStr := c.Query("toe")
	toe, err := time.Parse(time.RFC3339, toeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Time of Exam format"})
		return
	}

	assignments, err := helpers.FetchAssignmentsByTime(toe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
		return
	}

	pdfPath, err := helpers.GeneratePDF(assignments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	c.File(pdfPath)
}
