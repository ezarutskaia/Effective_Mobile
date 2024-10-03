package models

import (
	"gorm.io/gorm"
)

type Verse struct {
	gorm.Model
	ID int       `gorm:"primaryKey"`
	Text string
	SongID int
}