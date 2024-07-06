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
	MemberIDs    string         `json:"member_ids"`
	Year         int            `json:"year"`
	DetainedList string         `json:"detained_list"`
}

type Room struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	RoomType      string         `json:"room_type"`
	Seats         int            `json:"seats"`
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
	StudentID  int            `json:"student_id"`
	SeatNumber int            `json:"seat_number"`
}
