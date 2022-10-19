package models

import "time"

type Log struct {
	Method     string    `gorm:"type:varchar(255)" json:"method"`
	URL        string    `gorm:"type:varchar(255)" json:"url"`
	Code       string    `gorm:"type:varchar(255)" json:"code"`
	Accesstime time.Time `gorm:"type:varchar(255)" json:"accesstime"`
	Handletime time.Time `gorm:"type:varchar(255)" json:"handletime"`
}
