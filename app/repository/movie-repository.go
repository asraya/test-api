package repository

import (
	"api-test/dto"
	"api-test/models"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

//MovieRepository is a ....
type MovieRepository interface {
	InsertMovie(m models.Movie) models.Movie
	UpdateMovie(m models.Movie) models.Movie
	DeleteMovie(m models.Movie)
	AllMovie() []models.Movie
	FindMovieByID(movieID uint64) models.Movie
	PaginationMovie(pagination *dto.Pagination) (RepositoryResult, int)
}

type movieConnection struct {
	connection *gorm.DB
}

//NewMovieRepository
func NewMovieRepository(dbConn *gorm.DB) MovieRepository {
	return &movieConnection{
		connection: dbConn,
	}
}

func (db *movieConnection) InsertMovie(m models.Movie) models.Movie {
	db.connection.Save(&m)
	db.connection.Preload("Tag").Find(&m)
	return m
}

func (db *movieConnection) UpdateMovie(m models.Movie) models.Movie {
	db.connection.Save(&m)
	db.connection.Preload("Topic").Find(&m)
	return m
}

func (db *movieConnection) DeleteMovie(m models.Movie) {
	db.connection.Delete(&m)
}

func (db *movieConnection) FindMovieByID(movieID uint64) models.Movie {
	var movie models.Movie
	db.connection.Preload("Topic").Find(&movie, movieID)
	return movie
}

func (db *movieConnection) AllMovie() []models.Movie {
	var movies []models.Movie
	db.connection.Preload("Topic").Find(&movies)
	return movies
}

func (db *movieConnection) PaginationMovie(pagination *dto.Pagination) (RepositoryResult, int) {

	var moviey []models.Movie
	var count int64

	totalMovie, totalRows, totalPages, fromRow, toRow, toPro := 0, 0, 0, 0, 0, 0

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

	find = find.Find(&moviey)
	// has error find data
	errFind := find.Error

	if errFind != nil {
		return RepositoryResult{Error: errFind}, totalPages
	}

	pagination.Rows = moviey
	// count all data
	errCount := db.connection.Model(&models.Movie{}).Count(&count).Error

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

	// count all Movie
	errCountMovie := db.connection.Model(&models.Movie{}).Count(&count).Error

	if errCountMovie != nil {
		return RepositoryResult{Error: errCountMovie}, totalMovie
	}

	// calculate total pages
	totalMovie = int(math.Ceil(float64(count)/float64(pagination.Limit))) - 1
	if pagination.Page == 0 {
		// set from & to row on first page
		fromRow = 1
		toRow = pagination.Limit
	} else {
		if pagination.Page <= totalMovie {
			// calculate from & to row
			fromRow = pagination.Page*pagination.Limit + 1
			toRow = (pagination.Page + 1) * pagination.Limit
		}
	}

	if toPro > int(count) {
		// set to row with total rows
		toPro = totalMovie
	}

	pagination.FromRow = fromRow
	pagination.ToRow = toRow

	return RepositoryResult{Result: pagination}, totalPages
}
