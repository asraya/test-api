package handler

import (
	"net/http"
	"strconv"

	"api-test/dto"
	"api-test/helpers"
	"api-test/models"
	"api-test/service"

	"github.com/gin-gonic/gin"
)

//MovieHandler is a ...
type MovieHandler interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type movieHandler struct {
	movieService service.MovieService
}

func NewMovieHandler(movieServ service.MovieService) MovieHandler {
	return &movieHandler{
		movieService: movieServ,
	}
}

func (c *movieHandler) All(context *gin.Context) {
	var movies []models.Movie = c.movieService.All()
	res := helpers.BuildResponse(true, "OK", movies)
	context.JSON(http.StatusOK, res)
}

func (c *movieHandler) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helpers.BuildErrorResponse("No param id was found", err.Error(), helpers.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var movie models.Movie = c.movieService.FindByID(id)
	if &movie != &movie {
		res := helpers.BuildErrorResponse("Data not found", "No data with given id", helpers.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helpers.BuildResponse(true, "OK", movie)
		context.JSON(http.StatusOK, res)
	}
}

func (c *movieHandler) Insert(context *gin.Context) {
	var movieCreateDTO dto.MovieCreateDTO
	errDTO := context.ShouldBind(&movieCreateDTO)
	if errDTO != nil {
		res := helpers.BuildErrorResponse("Failed to process request", errDTO.Error(), helpers.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		result := c.movieService.Insert(movieCreateDTO)
		response := helpers.BuildResponse(true, "OK", result)
		context.JSON(http.StatusCreated, response)
	}
}

func (c *movieHandler) Update(context *gin.Context) {
	var movieUpdateDTO dto.MovieUpdateDTO

	errDTO := context.ShouldBind(&movieUpdateDTO)
	if errDTO != nil {
		res := helpers.BuildErrorResponse("Failed to process request", errDTO.Error(), helpers.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	if c.movieService.IsAllowedToEdit(movieUpdateDTO.ID) {
		result := c.movieService.Update(movieUpdateDTO)
		response := helpers.BuildResponse(true, "OK", result)
		context.JSON(http.StatusOK, response)
	}
}

func (c *movieHandler) Delete(context *gin.Context) {
	var movie models.Movie
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helpers.BuildErrorResponse("Failed tou get id", "No param id were found", helpers.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	movie.ID = id
	if c.movieService.IsAllowedToEdit(movie.ID) {
		c.movieService.Delete(movie)
		res := helpers.BuildResponse(true, "Deleted", helpers.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		response := helpers.BuildErrorResponse("You dont have permission", "You are not the owner", helpers.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}
