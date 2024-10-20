package helpers

import (
	"DevMaan707/UMS/models"
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
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

func LogAssignments(assignments []map[string]interface{}, logFileName string, toe time.Time) error {
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening log file: %w", err)
	}
	defer file.Close()
	logEntry := map[string]interface{}{
		"time":        toe.Format(time.RFC3339),
		"room":        assignments[0]["room"],
		"assignments": assignments,
	}

	entryData, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("error marshalling log entry to JSON: %w", err)
	}

	if _, err := file.Write(append(entryData, '\n')); err != nil {
		return fmt.Errorf("error writing to log file: %w", err)
	}

	return nil
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

func GenerateExamAssignments(assignType string, rooms []Room, students map[string][]string, params models.Params, toe time.Time, doe time.Duration) []map[string]interface{} {
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

					years[year]++
					sections[section]++
					totalStudents++

					assignedStudents = append(assignedStudents, map[string]interface{}{
						"student_id": studentID,
						"row":        row,
						"column":     column,
						"side":       "left",
						"toe":        toe.Format(time.RFC3339),
						"doe":        doe.String(),
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
						"toe":        toe.Format(time.RFC3339),
						"doe":        doe.String(),
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
						"toe":        toe.Format(time.RFC3339),
						"doe":        doe.String(),
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
							"toe":        toe.Format(time.RFC3339),
							"doe":        doe.String(),
						})
					}
				}

				benchIndex++
			}
		}

		roomIndex++
		assignments = append(assignments, map[string]interface{}{
			"room":        room.RoomNumber,
			"assignments": assignedStudents,
		})
	}
	logFileName := "exam_assignments.log"
	err := LogAssignments(assignments, logFileName, toe)
	if err != nil {
		fmt.Printf("Failed to log assignments: %v\n", err)
	}

	return assignments
}

func extractYear(studentID string) int {
	year, _ := strconv.Atoi(string(studentID[0:2]))
	return year
}

func extractSection(studentID string) string {
	return string(studentID[2])
}

type ExamAssignment struct {
	StudentID  string
	RoomNumber string
	TOE        time.Time
}

var Assignments []ExamAssignment

func FindStudentAssignment(assignments []ExamAssignment, studentID string, toe time.Time) string {
	for _, assignment := range assignments {
		if assignment.StudentID == studentID && assignment.TOE.Equal(toe) {
			return assignment.RoomNumber
		}
	}
	return ""
}

func FilterAssignmentsByTOE(assignments []ExamAssignment, toe time.Time) []ExamAssignment {
	filtered := []ExamAssignment{}
	for _, assignment := range assignments {
		if assignment.TOE.Equal(toe) {
			filtered = append(filtered, assignment)
		}
	}
	return filtered
}

func ShuffleStudents(students map[string][]string) {
	rand.Seed(time.Now().UnixNano())
	for branch := range students {
		rand.Shuffle(len(students[branch]), func(i, j int) {
			students[branch][i], students[branch][j] = students[branch][j], students[branch][i]
		})
	}
}

func ContainsInt(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
func FetchAssignmentsByTime(targetTime time.Time) ([]map[string]interface{}, error) {
	file, err := os.Open("exam_assignments.log")
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %w", err)
	}
	defer file.Close()

	var assignments []map[string]interface{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		var logEntry map[string]interface{}
		err := json.Unmarshal([]byte(line), &logEntry)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
		}
		if logTimeStr, exists := logEntry["time"]; exists {
			logTime, err := time.Parse(time.RFC3339, logTimeStr.(string))
			if err == nil && logTime.Equal(targetTime) {
				if assignmentArray, exists := logEntry["assignments"]; exists {
					if assignmentsList, ok := assignmentArray.([]interface{}); ok {
						for _, assign := range assignmentsList {
							if assignmentMap, ok := assign.(map[string]interface{}); ok {
								assignments = append(assignments, assignmentMap)
							}
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from file: %w", err)
	}

	return assignments, nil
}

type StudentAssignmentResponse struct {
	RoomNumber string `json:"room_number"`
	Details    string `json:"details"`
	Toe        string `json:"toe"`
	Block      string `json:"block"`
}

func GetStudentAssignment(studentID string, toe time.Time) (StudentAssignmentResponse, error) {
	assignments, err := FetchAssignmentsByTime(toe)
	if err != nil {
		return StudentAssignmentResponse{}, err
	}
	for _, assignment := range assignments {
		if innerAssignments, exists := assignment["assignments"]; exists {
			if assignmentList, ok := innerAssignments.([]interface{}); ok {
				for _, assign := range assignmentList {
					assignmentData, ok := assign.(map[string]interface{})
					if !ok {
						continue
					}
					if assignmentData["student_id"] == studentID {
						if assignmentToeStr, exists := assignmentData["toe"]; exists {
							if assignmentToeStr == toe.Format(time.RFC3339) {
								row, rowOk := assignmentData["row"].(float64)
								side, sideOk := assignmentData["side"].(string)

								if rowOk && sideOk {
									response := StudentAssignmentResponse{
										RoomNumber: assignment["room"].(string),
										Details:    fmt.Sprintf("%s - Row: %d", side, int(row)),
										Toe:        assignmentToeStr.(string),
										Block:      "",
									}
									return response, nil
								} else {
									return StudentAssignmentResponse{}, fmt.Errorf("invalid row or side format for student %s", studentID)
								}
							}
						} else {
							return StudentAssignmentResponse{}, fmt.Errorf("toe field not found for student %s", studentID)
						}
					}
				}
			}
		}
	}

	return StudentAssignmentResponse{}, fmt.Errorf("no assignment found for student %s at time %s", studentID, toe.Format(time.RFC3339))
}
func GeneratePDF(assignments []map[string]interface{}) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "B", 16)

	shiftRight := 8.0

	for _, assignment := range assignments {
		roomNumber := assignment["room"].(string)
		pdf.AddPage()

		pdf.CellFormat(0, 10, "Room: "+roomNumber, "", 1, "C", false, 0, "")

		if assignmentList, ok := assignment["assignments"].([]interface{}); ok {
			const (
				benchWidth  = 60.0
				benchHeight = 15.0
				padding     = 5.0
			)

			for _, assign := range assignmentList {
				assignMap := assign.(map[string]interface{})
				studentID := assignMap["student_id"].(string)
				row := assignMap["row"].(float64)
				column := assignMap["column"].(float64)
				side := assignMap["side"].(string)

				xPosition := shiftRight + (benchWidth+padding)*(column-1)
				yPosition := float64(row) * 20

				pdf.Rect(xPosition, yPosition, benchWidth, benchHeight, "D")
				pdf.SetFont("Arial", "B", 10)

				if side == "left" {
					pdf.Text(xPosition+5, yPosition+10, studentID)
				} else if side == "right" {
					pdf.Text(xPosition+benchWidth/2+5, yPosition+10, studentID)
				}
			}
		}
	}

	pdfFilePath := "assignments.pdf"
	err := pdf.OutputFileAndClose(pdfFilePath)
	if err != nil {
		return "", err
	}

	return pdfFilePath, nil
}
