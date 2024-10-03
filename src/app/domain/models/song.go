package models

import (
	"time"
	"gorm.io/gorm"
)

type Song struct {
	gorm.Model
	ID int                   `gorm:"primaryKey" json:"-"`
	Name string              `json:"name"`
	ReleaseDate *time.Time   `json:"releasedate"`
	Link string              `json:"link"`
	GroupID int
	Group Group
	Verses []Verse
}