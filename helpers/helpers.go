package helpers

import (
	"DevMaan707/UMS/models"
	"math/rand"
	"net/http"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"
)

func GenerateClassesTest(c *gin.Context) {
	var params models.Params
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	branches := map[string]struct {
		Name   string
		Prefix map[int]string
	}{
		"A": {"CSE", map[int]string{1: "24EG105", 2: "23EG105", 3: "22EG105", 4: "21EG105"}},
		"B": {"AIML", map[int]string{1: "24EG106", 2: "23EG106", 3: "22EG106", 4: "21EG106"}},
		"C": {"DS", map[int]string{1: "24EG107", 2: "23EG107", 3: "22EG107", 4: "21EG107"}},
		"D": {"ECE", map[int]string{1: "24EG108", 2: "23EG108", 3: "22EG108", 4: "21EG108"}},
	}

	testRooms := []models.Room{}
	testClasses := []models.Class{}

	for block, branchData := range branches {
		if !contains(params.Blocks, block) {
			continue
		}

		branchName := branchData.Name
		studentIDPrefixes := branchData.Prefix

		if !contains(params.Branches, branchName) {
			continue
		}

		sections := []string{"A", "B", "C", "D"}

		for year := 1; year <= 4; year++ {
			if !containsInt(params.Years, year) {
				continue
			}

			for i, section := range sections {
				roomNumber := fmt.Sprintf("%s-%d0%d", block, year, i+1)
				capacity := rand.Intn(6) + 55

				room := models.Room{
					ID:            uint(len(testRooms) + 1),
					RoomType:      "classroom",
					Capacity:      capacity,
					RoomNumber:    roomNumber,
					ClassAssigned: fmt.Sprintf("%s-%s", branchName, section),
				}
				testRooms = append(testRooms, room)

				studentIDs := generateStudentIDs(studentIDPrefixes[year], section, capacity)

				class := models.Class{
					ID:         uint(len(testClasses) + 1),
					ClassName:  fmt.Sprintf("%s-%s", branchName, section),
					Year:       year,
					Branch:     branchName,
					StudentIDs: studentIDs,
				}
				testClasses = append(testClasses, class)
			}
		}

		for labNum := 1; labNum <= 3; labNum++ {
			lab := models.Room{
				ID:            uint(len(testRooms) + 1),
				RoomType:      "lab",
				Capacity:      30,
				RoomNumber:    fmt.Sprintf("%s-00%d", block, labNum),
				ClassAssigned: "Shared across all years",
			}
			testRooms = append(testRooms, lab)
		}
	}

	assignments := generateExamAssignments("test", testClasses, testRooms, params)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Simulated Test Data",
		"assignments": assignments,
	})
}

func generateStudentIDs(prefix, section string, capacity int) []string {
	studentIDs := []string{}
	for i := 1; i <= capacity; i++ {
		studentID := fmt.Sprintf("%s%s%02d", prefix, section, i)
		studentIDs = append(studentIDs, studentID)
	}
	return studentIDs
}
func generateExamAssignments(assignType string, classes []models.Class, rooms []models.Room, params models.Params) []map[string]interface{} {
	assignments := make([]map[string]interface{}, 0)
	roomIndex := 0

	// Ensure params.NumberOfBranchesInRoom is not zero to avoid division by zero
	if params.NumberOfBranchesInRoom == 0 {
		params.NumberOfBranchesInRoom = 1
	}

	// Map to hold remaining students per branch
	remainingStudents := make(map[string][]string)
	for _, class := range classes {
		if contains(params.Branches, class.Branch) && containsInt(params.Years, class.Year) {
			remainingStudents[class.Branch] = append(remainingStudents[class.Branch], class.StudentIDs...)
		}
	}

	// Array to hold assigned students with row and column positions for each room
	roomsCapacity := make([][]map[string]interface{}, len(rooms))

	for roomIndex < len(rooms) {
		roomCapacity := rooms[roomIndex].Capacity
		if rooms[roomIndex].RoomType == "classroom" && params.SingleChild {
			roomCapacity /= 2
		}

		// Calculate number of benches based on room capacity
		//totalBenches := roomCapacity / 2
		allocatedBranches := make(map[string]bool)
		studentsPerBranch := roomCapacity / params.NumberOfBranchesInRoom

		// Shuffle branches if necessary
		branchOrder := shuffleBranches(params.Branches)
		if params.ShuffleYears && len(params.Years) > 1 {
			branchOrder = shuffleBranchesWithYears(params.Branches, params.Years)
		}

		// Track current row and column
		row, column := 1, 1

		// Loop over each branch
		for _, branch := range branchOrder {
			students := remainingStudents[branch]

			// Only assign students if we haven't exceeded room capacity
			if roomCapacity == 0 || len(allocatedBranches) >= params.NumberOfBranchesInRoom {
				break
			}
			if len(students) == 0 {
				continue
			}

			numToAllocate := Min(len(students), studentsPerBranch)
			allocatedBranches[branch] = true

			for i := 0; i < numToAllocate; i++ {
				// Assign the student to a bench and column
				roomsCapacity[roomIndex] = append(roomsCapacity[roomIndex], map[string]interface{}{
					"student_id": students[i],
					"row":        row,
					"column":     column,
				})

				// Move to the next seat in the room
				column++
				if column > 3 { // After 3 columns, move to the next row
					column = 1
					row++
				}
			}

			remainingStudents[branch] = students[numToAllocate:]
			roomCapacity -= numToAllocate
		}

		// Ensure all benches are filled according to the rules
		fillBenchesWithDifferentBranches(roomsCapacity[roomIndex], params.NumberOfBranchesInRoom)

		roomIndex++
	}

	// Create assignment output, excluding empty rooms
	for i, room := range roomsCapacity {
		if len(room) > 0 { // Only include rooms with assigned students
			assignment := map[string]interface{}{
				"room_id":       rooms[i].ID,
				"room_number":   rooms[i].RoomNumber,
				"assigned_to":   room,
				"assigned_room": rooms[i].RoomNumber,
			}
			assignments = append(assignments, assignment)
		}
	}

	return assignments
}

func fillBenchesWithDifferentBranches(roomCapacity []map[string]interface{}, numberOfBranchesInRoom int) {
	// Logic to ensure students from different branches sit together on benches if NumberOfBranchesInRoom > 1
	for i := 0; i < len(roomCapacity)-1; i += 2 {
		// If there are two students per bench, ensure they are from different branches if necessary
		if numberOfBranchesInRoom > 1 && roomCapacity[i]["student_id"].(string)[:6] == roomCapacity[i+1]["student_id"].(string)[:6] {
			// Shuffle them with another branch student if possible
			for j := i + 2; j < len(roomCapacity); j++ {
				if roomCapacity[i]["student_id"].(string)[:6] != roomCapacity[j]["student_id"].(string)[:6] {
					// Swap the students
					roomCapacity[i+1], roomCapacity[j] = roomCapacity[j], roomCapacity[i+1]
					break
				}
			}
		}
	}
}

func shuffleBranchesWithYears(branches []string, years []int) []string {
	rand.Seed(time.Now().UnixNano())
	shuffledBranches := []string{}
	for _, year := range years {
		for _, branch := range branches {
			shuffledBranches = append(shuffledBranches, fmt.Sprintf("%s-%d", branch, year))
		}
	}
	rand.Shuffle(len(shuffledBranches), func(i, j int) { shuffledBranches[i], shuffledBranches[j] = shuffledBranches[j], shuffledBranches[i] })
	return shuffledBranches
}

func shuffleBranches(branches []string) []string {
	rand.Seed(time.Now().UnixNano())
	shuffledBranches := make([]string, len(branches))
	copy(shuffledBranches, branches)
	rand.Shuffle(len(shuffledBranches), func(i, j int) { shuffledBranches[i], shuffledBranches[j] = shuffledBranches[j], shuffledBranches[i] })
	return shuffledBranches
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
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