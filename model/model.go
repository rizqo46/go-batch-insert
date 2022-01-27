package model

import "time"

type Student struct {
	ID          int `gorm:"primaryKey;autoIncrement"`
	Name        string
	DateOfBirth time.Time
	Grade       int
}
