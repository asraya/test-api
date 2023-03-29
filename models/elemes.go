package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents users table in database
type Elemes struct {
	ID         uint64         `gorm:"primary_key:auto_increment" json:"id"`
	NameCourse string         `gorm:"name_course" json:"name_course" binding:"required"`
	Title      string         `gorm:"type:varchar(255)" json:"title"  binding:"required"`
	Price      float64        `gorm:"price" json:"price"`
	Category   string         `gorm:"category" json:"category"`
	CreatedBy  string         `gorm:"created_by,omitempty" json:"created_by"`
	UpdatedBy  string         `gorm:"updated_by,omitempty" json:"updated_by"`
	DeletedBy  string         `gorm:"deleted_by,omitempty" json:"deleted_by"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
