package service

import (
	"fmt"
	"log"

	"api-test/app/repository"
	"api-test/dto"
	"api-test/models"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

//MovieService is a
type MovieService interface {
	Insert(b dto.MovieCreateDTO) models.Movie
	Update(b dto.MovieUpdateDTO) models.Movie
	Delete(b models.Movie)
	All() []models.Movie
	FindByID(movieID uint64) models.Movie
	IsAllowedToEdit(movieID uint64) bool
	PaginationMovie(repo repository.MovieRepository, context *gin.Context, pagination *dto.Pagination) dto.Response
}

type movieService struct {
	movieRepository repository.MovieRepository
}

//NewMovieService
func NewMovieService(movieRepo repository.MovieRepository) MovieService {
	return &movieService{
		movieRepository: movieRepo,
	}
}

func (service *movieService) Insert(b dto.MovieCreateDTO) models.Movie {
	movie := models.Movie{}
	err := smapping.FillStruct(&movie, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.movieRepository.InsertMovie(movie)
	return res
}

func (service *movieService) Update(b dto.MovieUpdateDTO) models.Movie {
	movie := models.Movie{}
	err := smapping.FillStruct(&movie, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.movieRepository.UpdateMovie(movie)
	return res
}

func (service *movieService) Delete(b models.Movie) {
	service.movieRepository.DeleteMovie(b)
}

func (service *movieService) All() []models.Movie {
	return service.movieRepository.AllMovie()
}

func (service *movieService) FindByID(movieID uint64) models.Movie {
	return service.movieRepository.FindMovieByID(movieID)
}

func (service *movieService) IsAllowedToEdit(movieID uint64) bool {
	b := service.movieRepository.FindMovieByID(movieID)
	id := (b.ID)
	return movieID == id
}

func (service *movieService) PaginationMovie(repo repository.MovieRepository, context *gin.Context, pagination *dto.Pagination) dto.Response {

	operationResult, totalPages := repo.PaginationMovie(pagination)

	if operationResult.Error != nil {
		return dto.Response{Status: false, Message: operationResult.Error.Error()}
	}

	var data = operationResult.Result.(*dto.Pagination)

	urlPath := context.Request.URL.Path

	searchQueryParams := ""

	for _, search := range pagination.Searchs {
		searchQueryParams += fmt.Sprintf("&%s.%s=%s", search.Column, search.Action, search.Query)
	}

	data.FirstPage = fmt.Sprintf("%s?limit=%d&page=%d&apikey=%s&sort=%s", urlPath, pagination.Limit, 0, pagination.Apikey, pagination.Sort) + searchQueryParams
	data.LastPage = fmt.Sprintf("%s?limit=%d&page=%d&apikey=%s&sort=%s", urlPath, pagination.Limit, totalPages, pagination.Apikey, pagination.Sort) + searchQueryParams

	if data.Page > 0 {
		data.PreviousPage = fmt.Sprintf("%s?limit=%d&page=%d&apikey=%s&sort=%s", urlPath, pagination.Limit, data.Page-1, pagination.Apikey, pagination.Sort) + searchQueryParams
	}

	if data.Page < totalPages {
		data.NextPage = fmt.Sprintf("%s?limit=%d&page=%d&apikey=%s&sort=%s", urlPath, pagination.Limit, data.Page+1, pagination.Apikey, pagination.Sort) + searchQueryParams
	}

	if data.Page > totalPages {
		data.PreviousPage = ""
	}

	return dto.Response{Status: true, Data: data}
}
