package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type Class struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	ClassName    string         `json:"class_name"`
	StudentIDs   []string       `json:"student_ids"`
	Year         int            `json:"year"`
	DetainedList string         `json:"detained_list"`
	Branch       string         `json:"branch"`
}

type Item struct {
	ID           int      `json:"id"`
	ClassName    string   `json:"class_name,omitempty"`
	StudentIDs   []string `json:"student_ids,omitempty"`
	Year         int      `json:"year,omitempty"`
	DetainedList []string `json:"detained_list,omitempty"`
	Branch       string   `json:"branch,omitempty"`
	RoomNumber   string   `json:"room_number,omitempty"`
	Capacity     int      `json:"capacity,omitempty"`
	Block        string   `json:"block,omitempty"`
	RoomType     string   `json:"room_type,omitempty"`
}

type Room struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	RoomType      string         `json:"room_type"`
	Capacity      int            `json:"capacity"`
	RoomNumber    string         `json:"room_number"`
	RoomTimetable string         `json:"room_timetable"`
	ClassAssigned string         `json:"class_assigned"`
}

type Assigned struct {
	ID         uint `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	ClassID    int            `json:"class_id"`
	RoomID     int            `json:"room_id"`
	RoomNumber string         `json:"room_number"`
	Year       int            `json:"year"`
}

type ExamAssignment struct {
	ID         uint `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	ExamID     int            `json:"exam_id"`
	RoomID     int            `json:"room_id"`
	RoomNumber string         `json:"room_number"`
	StudentIDs []int          `json:"student_ids"`
}

type AddValuesRequest struct {
	TableName string                 `json:"table_name"`
	Item      map[string]interface{} `json:"item"`
}
type GenerateClassesRequest struct {
	Type   string `json:"type"`
	Params Params `json:"params"`
}
type Params struct {
	Blocks                 []string `json:"blocks"`
	Branches               []string `json:"branches"`
	Years                  []int    `json:"years"`
	SingleChild            bool     `json:"single_child"`
	NumberOfBranchesInRoom int      `json:"number_of_branches"`
	RoomTypes              []string `json:"room_types"`
	ShuffleYears           bool     `json:"shuffle_years"`
}

type Details struct {
	Block         string `json:"block"`
	RoonType      string `json:"classroom"`
	Day           int    `json:"day"`
	HourSegment   int    `json:"hours"`
	NumberofHours int    `json:"no_hours"`
}

type ColumnsData struct {
	Columns map[string]string
}

type Received struct {
	ID      primitive.ObjectID `bson:"_id"`
	RoomNo  string             `bson:"Room_no"`
	DayKey  int                `bson:"Day_key"`
	DayTime string             `bson:"Day/Time"`
	Columns ColumnsData        `bson:",inline"`
}
