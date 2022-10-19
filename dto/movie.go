package dto

type MovieUpdateDTO struct {
	ID        uint64 `json:"id" form:"id" binding:"required"`
	Title     string `gorm:"type:varchar(255)" json:"title"`
	Year      string `gorm:"type:varchar(255)" json:"year"`
	ImdbID    string `gorm:"type:text" json:"imdbID"`
	Type      string `json:"type" form:"type" binding:"required"`
	Poster    string `json:"poster" form:"poster" binding:"required"`
	CreatedBy string `gorm:"created_by,omitempty" json:"created_by"`
	UpdatedBy string `gorm:"updated_by,omitempty" json:"updated_by"`
	DeletedBy string `gorm:"deleted_by,omitempty" json:"deleted_by"`
}

type MovieCreateDTO struct {
	ID        uint64 `gorm:"id" json:"id"`
	Title     string `gorm:"type:varchar(255)" json:"title"`
	Year      string `gorm:"type:varchar(255)" json:"year"`
	ImdbID    string `gorm:"type:text" json:"imdbID"`
	Type      string `json:"type" form:"type" binding:"required"`
	Poster    string `json:"poster" form:"poster" binding:"required"`
	CreatedBy string `gorm:"created_by,omitempty" json:"created_by"`
	UpdatedBy string `gorm:"updated_by,omitempty" json:"updated_by"`
	DeletedBy string `gorm:"deleted_by,omitempty" json:"deleted_by"`
}
