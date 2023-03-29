package dto

import "mime/multipart"

type ElemesUpdateDTO struct {
	ID         uint64  `json:"id" form:"id" binding:"required"`
	Title      string  `json:"title"`
	NameCourse string  `json:"name_course" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	Category   string  `json:"category"`
	CreatedBy  string  `gorm:"created_by,omitempty" json:"created_by"`
	UpdatedBy  string  `gorm:"updated_by,omitempty" json:"updated_by"`
	DeletedBy  string  `gorm:"deleted_by,omitempty" json:"deleted_by"`
}

type File struct {
	File multipart.File `json:"file,omitempty" validate:"required"`
}

type Url struct {
	Url string `json:"url,omitempty" validate:"required"`
}

type MediaDto struct {
	StatusCode int                    `json:"statusCode"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data"`
}

type ElemesCreateDTO struct {
	Title      string         `json:"title" form:"title" validate:"required"`
	File       multipart.File `json:"file,omitempty" validate:"required"`
	NameCourse string         `json:"name_course" form:"name_course" validate:"required"`
	Price      float64        `json:"price" form:"price" validate:"required"`
	Category   string         `json:"category" form:"category" validate:"required"`
	CreatedBy  string         `gorm:"created_by,omitempty" json:"created_by"`
	UpdatedBy  string         `gorm:"updated_by,omitempty" json:"updated_by"`
	DeletedBy  string         `gorm:"deleted_by,omitempty" json:"deleted_by"`
}
