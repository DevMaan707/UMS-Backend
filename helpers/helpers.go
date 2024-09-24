package helpers

import (
	"DevMaan707/UMS/models"
	"fmt"
	"math/rand"
	"net/http"
	"time"

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

	for roomIndex < len(rooms) {
		room := rooms[roomIndex]
		// Dynamically calculate the number of benches based on the room's original capacity
		examCapacity := room.Capacity * 2 / 3 // Convert normal capacity (3 per bench) to exam capacity (2 per bench)
		totalBenches := examCapacity / 2      // Each bench holds 2 students during exams

		branchOrder := shuffleBranches(params.Branches)
		if params.ShuffleYears && len(params.Years) > 1 {
			branchOrder = shuffleBranchesWithYears(params.Branches, params.Years)
		}

		// Track current row and bench position
		row, benchPosition := 1, 1
		roomsCapacity := []map[string]interface{}{}

		// Allocate students in a round-robin fashion across branches
		for totalBenches > 0 {
			branchAllocated := false

			for _, branch := range branchOrder {
				if len(remainingStudents[branch]) == 0 {
					continue
				}

				studentID := remainingStudents[branch][0]
				remainingStudents[branch] = remainingStudents[branch][1:]

				roomsCapacity = append(roomsCapacity, map[string]interface{}{
					"student_id": studentID,
					"row":        row,
					"bench":      benchPosition,
					"column":     (benchPosition-1)%3 + 1, // 3 benches in a row
				})

				branchAllocated = true

				// After seating 2 students per bench, move to the next bench
				if len(roomsCapacity)%2 == 0 {
					benchPosition++
					if benchPosition > 3 {
						benchPosition = 1
						row++
					}
					totalBenches--
				}

				if totalBenches == 0 {
					break
				}
			}

			// If no students were allocated, we are out of students to assign, so break out
			if !branchAllocated {
				break
			}
		}

		if len(roomsCapacity) > 0 {
			assignments = append(assignments, map[string]interface{}{
				"room_id":       room.ID,
				"room_number":   room.RoomNumber,
				"assigned_to":   roomsCapacity,
				"assigned_room": room.RoomNumber,
			})
		}

		roomIndex++
	}

	return assignments
}

// Ensure each bench has students from different branches
func fillBenchesWithDifferentBranches(roomCapacity []map[string]interface{}, numberOfBranchesInRoom int) {
	// Iterate through the room capacity, trying to ensure students from different branches are seated together
	for i := 0; i < len(roomCapacity)-1; i += 2 {
		// Ensure students from different branches sit together on the same bench
		branchA := roomCapacity[i]["student_id"].(string)[:6] // Extract branch identifier from the student ID
		branchB := roomCapacity[i+1]["student_id"].(string)[:6]

		// If two students on the same bench belong to the same branch, swap one with another from a different branch
		if branchA == branchB {
			for j := i + 2; j < len(roomCapacity); j++ {
				branchJ := roomCapacity[j]["student_id"].(string)[:6]
				// Find a student from a different branch and swap
				if branchA != branchJ {
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
	rand.Shuffle(len(shuffledBranches), func(i, j int) {
		shuffledBranches[i], shuffledBranches[j] = shuffledBranches[j], shuffledBranches[i]
	})
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
