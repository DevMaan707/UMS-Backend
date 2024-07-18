package helpers

import (
	"DevMaan707/UMS/models"
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"

	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

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

	// Map to hold remaining students per branch
	remainingStudents := make(map[string][]string)
	for _, class := range classes {
		remainingStudents[class.Branch] = append(remainingStudents[class.Branch], class.StudentIDs...)
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
		studentsPerBranch := roomCapacity / params.NumberOfBranchesInRoom

		for branch, students := range remainingStudents {
			if roomCapacity == 0 {
				break
			}
			if len(allocatedBranches) >= params.NumberOfBranchesInRoom {
				break
			}
			if len(students) == 0 {
				continue
			}

			numToAllocate := min(len(students), studentsPerBranch)
			roomsCapacity[roomIndex] = append(roomsCapacity[roomIndex], students[:numToAllocate]...)
			remainingStudents[branch] = students[numToAllocate:]
			roomCapacity -= numToAllocate
			allocatedBranches[branch] = true
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

func EmptyClassGen(c *gin.Context, mongoClient *mongo.Client) {
	const dbName = "TimeTable"

	// const colName = "A_Block"
	const reserve = "Reserve"

	const LoginCol = "Login_Credentials"

	var collectionReserve *mongo.Collection
	//creating instance of the structure
	var payload models.Details

	//binding the structure to the received data
	c.ShouldBindJSON(&payload)

	//Printing received
	fmt.Println("Data Successfully Received from the client")

	//Defining "rooms" var

	var rooms []string

	//Conditional Programming to redirect to different blocks

	if payload.Block == "A" {

		//Getting Collection
		collection := mongoClient.Database(dbName).Collection("A_Block")

		//Getting the search data from mongoDB in "rooms"

		rooms = Find(collectionReserve, collection, payload.HourSegment, payload.Block, payload.Day, payload.NumberofHours)

	} else if payload.Block == "B" {

		//Getting Collection

		collection := mongoClient.Database(dbName).Collection("B_Block")

		//Getting the search data from mongoDB in "rooms"

		rooms = Find(collectionReserve, collection, payload.HourSegment, payload.Block, payload.Day, payload.NumberofHours)

	} else if payload.Block == "C" {

		//Getting COllection

		collection := mongoClient.Database(dbName).Collection("C_Block")

		//Getting the search data from mongoDB in "rooms"

		rooms = Find(collectionReserve, collection, payload.HourSegment, payload.Block, payload.Day, payload.NumberofHours)

	} else if payload.Block == "D" {

		//Getting Collection

		collection := mongoClient.Database(dbName).Collection("D_Block")

		//Getting the search data from mongoDB in "rooms"

		rooms = Find(collectionReserve, collection, payload.HourSegment, payload.Block, payload.Day, payload.NumberofHours)

	} else if payload.Block == "H" {

		//Getting COllection

		collection := mongoClient.Database(dbName).Collection("H_Block")

		//Getting the search data from mongoDB in "rooms"

		rooms = Find(collectionReserve, collection, payload.HourSegment, payload.Block, payload.Day, payload.NumberofHours)
	} else if payload.Block == "E" {

		//Getting Collection

		collection := mongoClient.Database(dbName).Collection("E_Block")

		//Getting the search data from mongoDB in "rooms"

		rooms = Find(collectionReserve, collection, payload.HourSegment, payload.Block, payload.Day, payload.NumberofHours)

	} else if payload.Block == "All" {

		//Getting Collection
		collection := mongoClient.Database(dbName).Collection("All_Block")

		//Getting the search data from mongoDB in "rooms"

		rooms = Find(collectionReserve, collection, payload.HourSegment, payload.Block, payload.Day, payload.NumberofHours)

	}

	//Limiting the search results to only 5 Rooms

	var length = len(rooms)

	//Creating interface which consists of list of rooms and the length of the list

	response := map[string]interface{}{
		"number":    length,
		"classroom": rooms,
	}

	//Finally , sending the data back to the application
	c.JSON(http.StatusOK, response)

	fmt.Println("Response Sent!")
}

func Find(collection_forReserve, collection *mongo.Collection, hour int, Block string, Day int, Num_hours int) []string {

	fmt.Printf("HourSegment = %d\nBlock = %s\nDay= %d Num_Hours = %d\n", hour, Block, Day, Num_hours)

	if err := collection.Database().Client().Ping(context.Background(), nil); err != nil {
		slog.Error("Failed to ping MongoDB:", "Error", err)
	}

	//Describing the filter
	var filter = bson.M{strconv.Itoa(hour): bson.M{"$regex": "(TRAINING|LAB|SPORTS)$"},

		"Day_Key": Day}
	if Num_hours == 1 {
		filter = bson.M{
			strconv.Itoa(hour): bson.M{"$regex": "(TRAINING|LAB|SPORTS)$"},

			"Day_Key": Day,
		}
	} else if Num_hours == 2 {

		if hour <= 5 {
			filter = bson.M{
				strconv.Itoa(hour):     bson.M{"$regex": "(TRAINING|LAB|SPORTS)$"},
				strconv.Itoa(hour + 1): bson.M{"$regex": "(TRAINING|LAB|SPORTS)$"},
				"Day_Key":              Day,
			}
		} else {
			filter = bson.M{
				strconv.Itoa(hour): bson.M{"$regex": "(TRAINING|LAB|SPORTS)$"},

				"Day_Key": Day,
			}
		}
	} else if Num_hours == 3 {
		if hour <= 4 {
			filter = bson.M{
				strconv.Itoa(hour):     bson.M{"$regex": "(TRAINING|LAB|SPORTS)$"},
				strconv.Itoa(hour + 1): bson.M{"$regex": "(TRAINING|LAB|SPORTS)$"},
				strconv.Itoa(hour + 2): bson.M{"$regex": "(TRAINING|LAB|SPORTS)$"},
				"Day_Key":              Day,
			}
		} else if hour <= 5 {
			filter = bson.M{
				strconv.Itoa(hour):     bson.M{"$regex": "(TRAINING|LAB|SPORTS)$"},
				strconv.Itoa(hour + 1): bson.M{"$regex": "(TRAINING|LAB|SPORTS)$"},
				"Day_Key":              Day,
			}

		} else {
			filter = bson.M{
				strconv.Itoa(hour): bson.M{"$regex": "(TRAINING|LAB|SPORTS)$"},

				"Day_Key": Day,
			}
		}
	}

	//Initiating the Find Operation
	fmt.Println("Initiating Filter")
	cursor, err := collection.Find(context.Background(), filter)

	if err != nil {
		slog.Error("error", "Error", err)

	}
	fmt.Println("Got cursor , Searching values: ")

	defer cursor.Close(context.Background())

	//Checking the length of the cursor
	fmt.Println("Cursor Count:", cursor.RemainingBatchLength())

	//Defining slices

	//Iterating through the results
	var rooms []string
	var data models.Received

	//Trying to iterate through the received data
	for cursor.Next(context.Background()) {

		fmt.Println("Decoding the json")
		err := cursor.Decode(&data)

		if err != nil {
			log.Fatal(err)
		}

		//printing the object ID just incase
		fmt.Println(data.ID)

		rooms = append(rooms, data.RoomNo)

	}

	//Iterating through the rooms just incase
	for _, room := range rooms {
		fmt.Println(room)
	}
	//Returning the slice back to routes/PostDetails.go
	return rooms

}
