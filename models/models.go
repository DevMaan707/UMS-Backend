package models

import (
	"time"

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

type Room struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	RoomType      string         `json:"room_type"`
	Capacity      int            `json:"capacity"`
	RoomNumber    string         `json:"room_number"`
	RoomTimetable string         `json:"room_timetable"`
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
	Blocks   []string `json:"blocks"`
	Branches []string `json:"branches"`
	Years    []int    `json:"years"`
}
