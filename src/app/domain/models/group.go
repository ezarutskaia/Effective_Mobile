package models

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	ID int                   `gorm:"primaryKey" json: -`
	Name string              `gorm:"unique" json:"name"`
	Songs []Song             `json:"-"`
}