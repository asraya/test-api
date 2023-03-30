package repository

import (
	"api-test/dto"
	"api-test/models"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

// ElemesRepository is a ....
type ElemesRepository interface {
	InsertElemes(m models.Elemes) models.Elemes
	UpdateElemes(m models.Elemes) models.Elemes
	DeleteElemes(m models.Elemes)
	AllElemes() []models.Elemes
	FindElemesByID(elemesID uint64) models.Elemes
	PaginationElemes(pagination *dto.Pagination) (RepositoryResult, int)
}

type elemesConnection struct {
	connection *gorm.DB
}

// NewElemesRepository
func NewElemesRepository(dbConn *gorm.DB) ElemesRepository {
	return &elemesConnection{
		connection: dbConn,
	}
}

func (db *elemesConnection) InsertElemes(m models.Elemes) models.Elemes {
	db.connection.Save(&m)
	db.connection.Preload("Tag").Find(&m)
	return m
}

func (db *elemesConnection) UpdateElemes(m models.Elemes) models.Elemes {
	db.connection.Save(&m)
	db.connection.Preload("Elemes").Find(&m)
	return m
}

func (db *elemesConnection) DeleteElemes(m models.Elemes) {
	db.connection.Delete(&m)
}

func (db *elemesConnection) FindElemesByID(elemesID uint64) models.Elemes {
	var elemes models.Elemes
	db.connection.Preload("Elemes").Find(&elemes, elemesID)
	return elemes
}

func (db *elemesConnection) GetCategoryCourse(elemesID uint64) models.Elemes {
	var elemes models.Elemes
	db.connection.Preload("Elemes").Find(&elemes, elemesID)
	return elemes
}

func (db *elemesConnection) AllElemes() []models.Elemes {
	var elemess []models.Elemes
	db.connection.Preload("Elemes").Find(&elemess)
	return elemess
}

func (db *elemesConnection) PaginationElemes(pagination *dto.Pagination) (RepositoryResult, int) {
	var contacts models.Elemes
	var count int64

	totalRows, totalPages, fromRow, toRow := 0, 0, 0, 0

	offset := pagination.Page * pagination.Limit

	// get data with limit, offset & order
	find := db.connection.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)

	// generate where query
	searchs := pagination.Searchs

	if searchs != nil {
		for _, value := range searchs {
			column := value.Column
			action := value.Action
			query := value.Query

			switch action {
			case "equals":
				whereQuery := fmt.Sprintf("%s = ?", column)
				find = find.Where(whereQuery, query)
				break
			case "contains":
				whereQuery := fmt.Sprintf("%s LIKE ?", column)
				find = find.Where(whereQuery, "%"+query+"%")
				break
			case "in":
				whereQuery := fmt.Sprintf("%s IN (?)", column)
				queryArray := strings.Split(query, ",")
				find = find.Where(whereQuery, queryArray)
				break

			}
		}
	}

	find = find.Find(&contacts)

	// has error find data
	errFind := find.Error

	if errFind != nil {
		return RepositoryResult{Error: errFind}, totalPages
	}

	pagination.Rows = contacts

	// count all data
	errCount := db.connection.Model(&models.Elemes{}).Count(&count).Error

	if errCount != nil {
		return RepositoryResult{Error: errCount}, totalPages
	}

	pagination.TotalRows = int(count)
	// calculate total pages
	totalPages = int(math.Ceil(float64(count)/float64(pagination.Limit))) - 1

	if pagination.Page == 0 {
		// set from & to row on first page
		fromRow = 1
		toRow = pagination.Limit
	} else {
		if pagination.Page <= totalPages {
			// calculate from & to row
			fromRow = pagination.Page*pagination.Limit + 1
			toRow = (pagination.Page + 1) * pagination.Limit
		}
	}

	if toRow > int(count) {
		// set to row with total rows
		toRow = totalRows
	}

	pagination.FromRow = fromRow
	pagination.ToRow = toRow

	return RepositoryResult{Result: pagination}, totalPages
}
