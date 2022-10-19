package dto

import "time"

type LogCreateDTO struct {
	Method     string    `json:"method" form:"method"`
	URL        string    `json:"url" form:"url"`
	Code       string    `json:"code" form:"code"`
	Accesstime time.Time `json:"accesstime" form:"accesstime"`
	Handletime time.Time `json:"handletime" form:"handletime"`
}
