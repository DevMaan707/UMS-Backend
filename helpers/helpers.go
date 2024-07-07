package helpers

import (
	"DevMaan707/UMS/models"
	"context"
	"fmt"
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

func ShuffleStudents(students []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(students), func(i, j int) { students[i], students[j] = students[j], students[i] })
}

func generateExamAssignments(assignType string, classes []models.Class, rooms []models.Room) []map[string]interface{} {
	assignments := make([]map[string]interface{}, 0)
	roomIndex := 0

	remainingStudents := make(map[string][]string)

	for _, class := range classes {
		remainingStudents[class.ClassName] = class.StudentIDs
	}

	roomsCapacity := make([][]string, len(rooms))
	for i := range roomsCapacity {
		roomsCapacity[i] = make([]string, 0)
	}

	for roomIndex < len(rooms) {
		allocatedBranches := []string{}
		roomCapacity := rooms[roomIndex].Capacity

		for class := range remainingStudents {
			ShuffleStudents(remainingStudents[class])
		}

		for _, class := range classes {
			if len(remainingStudents[class.ClassName]) > 0 {
				numToAllocate := min(len(remainingStudents[class.ClassName]), roomCapacity/2)
				roomsCapacity[roomIndex] = append(roomsCapacity[roomIndex], remainingStudents[class.ClassName][:numToAllocate]...)
				remainingStudents[class.ClassName] = remainingStudents[class.ClassName][numToAllocate:]
				roomCapacity -= numToAllocate
				allocatedBranches = append(allocatedBranches, class.Branch)
				if len(allocatedBranches) >= 2 {
					break
				}
			}
		}

		for _, class := range classes {
			if roomCapacity == 0 {
				break
			}
			if len(remainingStudents[class.ClassName]) > 0 {
				numToAllocate := min(len(remainingStudents[class.ClassName]), roomCapacity)
				roomsCapacity[roomIndex] = append(roomsCapacity[roomIndex], remainingStudents[class.ClassName][:numToAllocate]...)
				remainingStudents[class.ClassName] = remainingStudents[class.ClassName][numToAllocate:]
				roomCapacity -= numToAllocate
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

func fetchClassesFromDB(dynamoClient *dynamodb.Client, branches []string, years []int) ([]models.Class, error) {
	var classes []models.Class
	// Add code to fetch from DynamoDB here
	return classes, nil
}

func fetchRoomsFromDB(dynamoClient *dynamodb.Client, years []int) ([]models.Room, error) {
	var rooms []models.Room
	// Add code to fetch from DynamoDB here
	return rooms, nil
}

func generateClassesHandler(c *gin.Context, dynamoClient *dynamodb.Client) {
	var req models.GenerateClassesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	classes, err := fetchClassesFromDB(dynamoClient, req.Params.Branches, req.Params.Years)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	rooms, err := fetchRoomsFromDB(dynamoClient, req.Params.Years)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	assignments := generateExamAssignments(req.Type, classes, rooms)
	c.JSON(200, assignments)
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

func GenerateClasses(c *gin.Context, dynamoClient *dynamodb.Client) {

}

func AddValues(c *gin.Context, svc *dynamodb.Client) {
	var req models.AddValuesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := marshalMap(req.Item)
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
