package helpers

import (
	"DevMaan707/UMS/models"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Room struct {
	ID            uint   `json:"id"`
	RoomType      string `json:"room_type"`
	Capacity      int    `json:"capacity"`
	RoomNumber    string `json:"room_number"`
	AssignedClass string `json:"assigned_class"`
	Rows          int    `json:"rows"`
	Columns       int    `json:"columns"`
}

type Class struct {
	ID         uint     `json:"id"`
	ClassName  string   `json:"class_name"`
	Year       int      `json:"year"`
	Branch     string   `json:"branch"`
	StudentIDs []string `json:"student_ids"`
}

func GenerateTestData() (map[string][]Room, map[string][]Class) {
	const numRows = 8
	const numColumns = 3

	branches := map[string]struct {
		Name   string
		Prefix map[int]string
	}{
		"A": {"CSE", map[int]string{1: "24EG105", 2: "23EG105", 3: "22EG105", 4: "21EG105"}},
		"B": {"AIML", map[int]string{1: "24EG106", 2: "23EG106", 3: "22EG106", 4: "21EG106"}},
		"C": {"CS", map[int]string{1: "24EG107", 2: "23EG107", 3: "22EG107", 4: "21EG107"}},
		"D": {"ECE", map[int]string{1: "24EG108", 2: "23EG108", 3: "22EG108", 4: "21EG108"}},
	}

	blocks := map[string][]Room{}
	classes := map[string][]Class{}

	for block, branchData := range branches {
		blockRooms := []Room{}
		branchClasses := []Class{}

		for i := 1; i <= 5; i++ {
			roomType := "classroom"
			roomNumber := fmt.Sprintf("%s-%02d", block, i)
			capacity := numRows * numColumns * 2

			room := Room{
				ID:         uint(len(blockRooms) + 1),
				RoomType:   roomType,
				Capacity:   capacity,
				RoomNumber: roomNumber,
				Rows:       numRows,
				Columns:    numColumns,
			}
			blockRooms = append(blockRooms, room)
		}

		blocks[block] = blockRooms

		sections := []string{"A", "B", "C", "D", "E"}
		for year := 1; year <= 4; year++ {
			for _, section := range sections {
				className := fmt.Sprintf("%s-%s", branchData.Name, section)
				studentIDs := generateStudentIDs(branchData.Prefix[year], section, rand.Intn(12)+55)

				class := Class{
					ID:         uint(len(branchClasses) + 1),
					ClassName:  className,
					Year:       year,
					Branch:     branchData.Name,
					StudentIDs: studentIDs,
				}
				branchClasses = append(branchClasses, class)
			}
		}

		classes[branchData.Name] = branchClasses
	}

	return blocks, classes
}

func generateStudentIDs(prefix, section string, capacity int) []string {
	studentIDs := []string{}
	for i := 1; i <= capacity; i++ {
		studentID := fmt.Sprintf("%s%s%02d", prefix, section, i)
		studentIDs = append(studentIDs, studentID)
	}
	return studentIDs
}

func AssignRoomsForExams(c *gin.Context) {
	var params models.Params
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	blocks, classes := GenerateTestData()

	selectedRooms := []Room{}
	for _, block := range params.Blocks {
		if rooms, found := blocks[block]; found {
			selectedRooms = append(selectedRooms, rooms...)
		}
	}

	selectedStudents := make(map[string][]string)
	for _, branch := range params.Branches {
		if branchClasses, found := classes[branch]; found {
			for _, class := range branchClasses {
				if containsInt(params.Years, class.Year) {
					selectedStudents[branch] = append(selectedStudents[branch], class.StudentIDs...)
				}
			}
		}
	}

	if params.InternalShuffle {
		for branch := range selectedStudents {
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(selectedStudents[branch]), func(i, j int) {
				selectedStudents[branch][i], selectedStudents[branch][j] = selectedStudents[branch][j], selectedStudents[branch][i]
			})
		}
	}

	assignments := generateExamAssignments("exam", selectedRooms, selectedStudents, params)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Exam Room Assignments",
		"assignments": assignments,
	})
}

func generateExamAssignments(assignType string, rooms []Room, students map[string][]string, params models.Params) []map[string]interface{} {
	assignments := make([]map[string]interface{}, 0)
	roomIndex := 0

	remainingStudents := map[string]int{}
	for branch, studentList := range students {
		remainingStudents[branch] = len(studentList)
	}

	for roomIndex < len(rooms) && len(remainingStudents) > 0 {
		room := rooms[roomIndex]

		totalBenches := room.Rows * room.Columns
		assignedStudents := []map[string]interface{}{}
		benchIndex := 0
		sections := map[string]int{}
		years := map[int]int{}
		totalStudents := 0

		if params.NumberOfBranchesInRoom == 1 {
			currentBranch := ""
			for _, branch := range params.Branches {
				if remainingStudents[branch] > 0 {
					currentBranch = branch
					break
				}
			}

			if currentBranch == "" {
				break
			}

			for benchIndex < totalBenches && remainingStudents[currentBranch] > 0 {
				var row, column int
				if params.RowWise {
					row = (benchIndex / room.Columns) + 1
					column = (benchIndex % room.Columns) + 1
				} else {
					column = (benchIndex / room.Rows) + 1
					row = (benchIndex % room.Rows) + 1
				}

				if remainingStudents[currentBranch] > 0 {
					studentID := students[currentBranch][0]
					students[currentBranch] = students[currentBranch][1:]
					remainingStudents[currentBranch]--

					year := extractYear(studentID)
					section := extractSection(studentID)

					// Update counts
					years[year]++
					sections[section]++
					totalStudents++

					assignedStudents = append(assignedStudents, map[string]interface{}{
						"student_id": studentID,
						"row":        row,
						"column":     column,
						"side":       "left",
					})
				}

				if remainingStudents[currentBranch] > 0 {
					studentID := students[currentBranch][0]
					students[currentBranch] = students[currentBranch][1:]
					remainingStudents[currentBranch]--

					year := extractYear(studentID)
					section := extractSection(studentID)

					years[year]++
					sections[section]++
					totalStudents++

					assignedStudents = append(assignedStudents, map[string]interface{}{
						"student_id": studentID,
						"row":        row,
						"column":     column,
						"side":       "right",
					})
				}

				benchIndex++
			}
		} else if params.NumberOfBranchesInRoom == 2 && !params.SingleChild {
			branchIndex := 0
			branchOrder := params.Branches

			for benchIndex < totalBenches && len(remainingStudents) > 0 {
				for branchIndex < len(branchOrder) && remainingStudents[branchOrder[branchIndex]] == 0 {
					branchIndex++
				}

				if branchIndex >= len(branchOrder) {
					break
				}

				branchA := branchOrder[branchIndex]
				var row, column int
				if params.RowWise {
					row = (benchIndex / room.Columns) + 1
					column = (benchIndex % room.Columns) + 1
				} else {
					column = (benchIndex / room.Rows) + 1
					row = (benchIndex % room.Rows) + 1
				}

				if remainingStudents[branchA] > 0 {
					studentID := students[branchA][0]
					students[branchA] = students[branchA][1:]
					remainingStudents[branchA]--

					year := extractYear(studentID)
					section := extractSection(studentID)

					years[year]++
					sections[section]++
					totalStudents++

					assignedStudents = append(assignedStudents, map[string]interface{}{
						"student_id": studentID,
						"row":        row,
						"column":     column,
						"side":       "left",
					})
				}
				branchIndex++
				if branchIndex >= len(branchOrder) {
					branchIndex = 0
				}
				for branchIndex < len(branchOrder) && remainingStudents[branchOrder[branchIndex]] == 0 {
					branchIndex++
				}

				if branchIndex < len(branchOrder) {
					branchB := branchOrder[branchIndex]

					if remainingStudents[branchB] > 0 {
						studentID := students[branchB][0]
						students[branchB] = students[branchB][1:]
						remainingStudents[branchB]--

						year := extractYear(studentID)
						section := extractSection(studentID)

						years[year]++
						sections[section]++
						totalStudents++

						assignedStudents = append(assignedStudents, map[string]interface{}{
							"student_id": studentID,
							"row":        row,
							"column":     column,
							"side":       "right",
						})
					}
				}

				benchIndex++
			}
		}

		if len(assignedStudents) > 0 {
			assignments = append(assignments, map[string]interface{}{
				"room_id":        room.ID,
				"room_number":    room.RoomNumber,
				"total_students": totalStudents,
				"sections":       sections,
				"years":          years,
				"assigned_to":    assignedStudents,
			})
		}

		roomIndex++
	}

	return assignments
}

func extractYear(studentID string) int {
	yearStr := studentID[:2]
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0
	}
	return year
}

func extractSection(studentID string) string {

	return string(studentID[7])
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsInt(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
