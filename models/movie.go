package models

import (
	"time"

	"gorm.io/gorm"
)

//User represents users table in database
type Movie struct {
	ID        uint64         `gorm:"primary_key:auto_increment" json:"id"`
	Year      string         `json:"year" form:"year" binding:"required"`
	Title     string         `gorm:"type:varchar(255)" json:"title"  binding:"required"`
	Type      string         `gorm:"type:text" json:"type" binding:"required"`
	Poster    string         `json:"poster" form:"poster"`
	CreatedBy string         `gorm:"created_by,omitempty" json:"created_by"`
	UpdatedBy string         `gorm:"updated_by,omitempty" json:"updated_by"`
	DeletedBy string         `gorm:"deleted_by,omitempty" json:"deleted_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
