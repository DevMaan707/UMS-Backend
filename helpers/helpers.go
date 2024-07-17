package helpers

import (
	"DevMaan707/UMS/models"
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
)

func DefineTables(svc *dynamodb.Client) {

	tables := []struct {
		name       string
		keySchema  []types.KeySchemaElement
		attributes []types.AttributeDefinition
	}{
		{
			name: "Class",
			keySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("ID"),
					KeyType:       types.KeyTypeHash,
				},
			},
			attributes: []types.AttributeDefinition{
				{
					AttributeName: aws.String("ID"),
					AttributeType: types.ScalarAttributeTypeN,
				},
			},
		},
		{
			name: "Room",
			keySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("ID"),
					KeyType:       types.KeyTypeHash,
				},
			},
			attributes: []types.AttributeDefinition{
				{
					AttributeName: aws.String("ID"),
					AttributeType: types.ScalarAttributeTypeN,
				},
			},
		},
		{
			name: "Assigned",
			keySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("ID"),
					KeyType:       types.KeyTypeHash,
				},
			},
			attributes: []types.AttributeDefinition{
				{
					AttributeName: aws.String("ID"),
					AttributeType: types.ScalarAttributeTypeN,
				},
			},
		},
		{
			name: "ExamAssignment",
			keySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("ID"),
					KeyType:       types.KeyTypeHash,
				},
			},
			attributes: []types.AttributeDefinition{
				{
					AttributeName: aws.String("ID"),
					AttributeType: types.ScalarAttributeTypeN,
				},
			},
		},
	}

	for _, table := range tables {
		createTable(svc, table.name, table.keySchema, table.attributes)
	}
}

func createTable(svc *dynamodb.Client, tableName string, keySchema []types.KeySchemaElement, attributes []types.AttributeDefinition) {
	input := &dynamodb.CreateTableInput{
		TableName:            aws.String(tableName),
		KeySchema:            keySchema,
		AttributeDefinitions: attributes,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err := svc.CreateTable(context.TODO(), input)
	if err != nil {
		fmt.Printf("Got error calling CreateTable: %s\n", err)
		return
	}

	fmt.Printf("Created the table %s successfully\n", tableName)
}

func AddValues(c *gin.Context, svc *dynamodb.Client) {
	var req models.AddValuesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "TableName is required"})
		return
	}

	item, err := attributevalue.MarshalMap(req.Item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(req.TableName),
		Item:      item,
	}

	_, err = svc.PutItem(context.TODO(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added successfully"})
}

func marshalMap(in map[string]interface{}) (map[string]types.AttributeValue, error) {
	av, err := attributevalue.MarshalMap(in)
	if err != nil {
		return nil, err
	}
	return av, nil
}

func ShuffleStudents(students []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(students), func(i, j int) { students[i], students[j] = students[j], students[i] })
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// func generateExamAssignments(AssignType string, classes []models.Class, rooms []models.Room, params models.Params) []map[string]interface{} {
// 	assignments := make([]map[string]interface{}, 0)
// 	roomIndex := 0

// 	remainingStudents := make(map[string][]string)
// 	for _, class := range classes {
// 		remainingStudents[class.ClassName] = class.StudentIDs
// 	}

// 	roomsCapacity := make([][]string, len(rooms))
// 	for i := range roomsCapacity {
// 		roomsCapacity[i] = make([]string, 0)
// 	}

// 	for roomIndex < len(rooms) {
// 		allocatedBranches := []string{}
// 		roomCapacity := rooms[roomIndex].Capacity

// 		if rooms[roomIndex].RoomType == "classroom" && params.SingleChild {
// 			roomCapacity /= 2
// 		}

// 		for class := range remainingStudents {
// 			ShuffleStudents(remainingStudents[class])
// 		}

// 		for _, class := range classes {
// 			if len(remainingStudents[class.ClassName]) > 0 {
// 				numToAllocate := min(len(remainingStudents[class.ClassName]), roomCapacity/params.NumberOfBranchesInRoom)
// 				roomsCapacity[roomIndex] = append(roomsCapacity[roomIndex], remainingStudents[class.ClassName][:numToAllocate]...)
// 				remainingStudents[class.ClassName] = remainingStudents[class.ClassName][numToAllocate:]
// 				roomCapacity -= numToAllocate
// 				allocatedBranches = append(allocatedBranches, class.Branch)
// 				if len(allocatedBranches) >= params.NumberOfBranchesInRoom {
// 					break
// 				}
// 			}
// 		}

// 		for _, class := range classes {
// 			if roomCapacity == 0 {
// 				break
// 			}
// 			if len(remainingStudents[class.ClassName]) > 0 {
// 				numToAllocate := min(len(remainingStudents[class.ClassName]), roomCapacity)
// 				roomsCapacity[roomIndex] = append(roomsCapacity[roomIndex], remainingStudents[class.ClassName][:numToAllocate]...)
// 				remainingStudents[class.ClassName] = remainingStudents[class.ClassName][numToAllocate:]
// 				roomCapacity -= numToAllocate
// 			}
// 		}

// 		roomIndex++
// 	}

// 	for i, room := range roomsCapacity {
// 		assignment := map[string]interface{}{
// 			"room_id":       rooms[i].ID,
// 			"room_number":   rooms[i].RoomNumber,
// 			"assigned_to":   room,
// 			"assigned_room": rooms[i].RoomNumber,
// 		}
// 		assignments = append(assignments, assignment)
// 	}

// 	return assignments
// }

func generateExamAssignments(assignType string, classes []models.Class, rooms []models.Room, params models.Params) []map[string]interface{} {
	assignments := make([]map[string]interface{}, 0)
	roomIndex := 0

	// Map to hold remaining students per class
	remainingStudents := make(map[string][]string)
	for _, class := range classes {
		remainingStudents[class.ClassName] = class.StudentIDs
	}

	// Array to hold assigned students for each room
	roomsCapacity := make([][]string, len(rooms))
	for i := range roomsCapacity {
		roomsCapacity[i] = make([]string, 0)
	}

	for roomIndex < len(rooms) {
		roomCapacity := rooms[roomIndex].Capacity
		if rooms[roomIndex].RoomType == "classroom" && params.SingleChild {
			roomCapacity /= 2
		}

		allocatedBranches := make(map[string]bool)

		for roomCapacity > 0 && len(allocatedBranches) < params.NumberOfBranchesInRoom {
			for _, class := range classes {
				if len(remainingStudents[class.ClassName]) == 0 {
					continue
				}

				if _, found := allocatedBranches[class.Branch]; !found {
					numToAllocate := min(len(remainingStudents[class.ClassName]), roomCapacity)
					roomsCapacity[roomIndex] = append(roomsCapacity[roomIndex], remainingStudents[class.ClassName][:numToAllocate]...)
					remainingStudents[class.ClassName] = remainingStudents[class.ClassName][numToAllocate:]
					roomCapacity -= numToAllocate

					allocatedBranches[class.Branch] = true
					if len(allocatedBranches) >= params.NumberOfBranchesInRoom {
						break
					}
				}
			}
		}

		roomIndex++
	}

	for i, room := range roomsCapacity {
		assignment := map[string]interface{}{
			"room_id":       rooms[i].ID,
			"room_number":   rooms[i].RoomNumber,
			"assigned_to":   room,
			"assigned_room": rooms[i].RoomNumber,
		}
		assignments = append(assignments, assignment)
	}

	return assignments
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func FetchClassesFromDB(dynamoClient *dynamodb.Client, branches []string, years []int) ([]models.Class, error) {

	slog.Info("In FetchClasses From DB")
	var classes []models.Class
	for _, branch := range branches {
		for _, year := range years {
			input := &dynamodb.ScanInput{
				TableName:        aws.String("Class"),
				FilterExpression: aws.String("#yr = :year AND Branch = :branch"),
				ExpressionAttributeNames: map[string]string{
					"#yr": "Year",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":branch": &types.AttributeValueMemberS{Value: branch},
					":year":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", year)},
				},
			}

			result, err := dynamoClient.Scan(context.TODO(), input)
			if err != nil {
				slog.Error("Error while fetching Class data")
				return nil, err

			}

			var fetchedClasses []models.Class
			err = attributevalue.UnmarshalListOfMaps(result.Items, &fetchedClasses)
			//	fmt.Printf("Result Items => %v", result.Items)
			fmt.Printf("\n\nMarshalled class values => %v\n\n", fetchedClasses)
			if err != nil {
				return nil, err
			}
			classes = append(classes, fetchedClasses...)
			fmt.Printf("\n\nClasses %v", classes)
		}
	}
	return classes, nil
}

func FetchRoomsFromDB(dynamoClient *dynamodb.Client, params models.Params) ([]models.Room, error) {
	slog.Info("In FetchRooms From DB")
	var rooms []models.Room
	for _, block := range params.Blocks {
		for _, roomType := range params.RoomTypes {
			input := &dynamodb.ScanInput{
				TableName:        aws.String("Room"),
				FilterExpression: aws.String("#block = :block AND RoomType = :roomType"),
				ExpressionAttributeNames: map[string]string{
					"#block": "Block",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":block":    &types.AttributeValueMemberS{Value: block},
					":roomType": &types.AttributeValueMemberS{Value: roomType},
				},
			}
			result, err := dynamoClient.Scan(context.TODO(), input)
			if err != nil {
				slog.Error("Error while fetching Room data")
				return nil, err
			}
			var fetchedRooms []models.Room
			err = attributevalue.UnmarshalListOfMaps(result.Items, &fetchedRooms)
			fmt.Printf("\n\nUnMarshalled Room values => %v\n\n", fetchedRooms)
			if err != nil {
				return nil, err
			}
			rooms = append(rooms, fetchedRooms...)
		}
	}
	return rooms, nil
}

func GenerateClassesHandler(c *gin.Context, dynamoClient *dynamodb.Client) {
	var req models.GenerateClassesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	classes, err := FetchClassesFromDB(dynamoClient, req.Params.Branches, req.Params.Years)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	rooms, err := FetchRoomsFromDB(dynamoClient, req.Params)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	assignments := generateExamAssignments(req.Type, classes, rooms, req.Params)
	c.JSON(200, assignments)
}
