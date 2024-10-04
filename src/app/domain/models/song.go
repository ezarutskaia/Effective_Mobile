package models

import (
	"time"
	"gorm.io/gorm"
)

type Song struct {
	gorm.Model				 `json:"-"`
	ID int                   `gorm:"primaryKey" json:"-"`
	Name string              `json:"name"`
	ReleaseDate *time.Time   `json:"releasedate"`
	Link string              `json:"link"`
	GroupID int				 `json:"-"`
	Group Group				 `json:"-"`
	Verses []Verse			 `json:"-"`
}