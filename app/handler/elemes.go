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

// ElemesHandler is a ...
type ElemesHandler interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type elemesHandler struct {
	elemesService service.ElemesService
	jwtService    service.JWTService
}

func NewElemesHandler(elemesServ service.ElemesService, jwtServ service.JWTService) ElemesHandler {
	return &elemesHandler{
		elemesService: elemesServ,
		jwtService:    jwtServ,
	}

}

func (c *elemesHandler) All(context *gin.Context) {
	var elemess []models.Elemes = c.elemesService.All()
	res := helpers.BuildResponse(true, "OK", elemess)
	context.JSON(http.StatusOK, res)
}

func (c *elemesHandler) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helpers.BuildErrorResponse("No param id was found", err.Error(), helpers.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var elemes models.Elemes = c.elemesService.FindByID(id)
	if &elemes != &elemes {
		res := helpers.BuildErrorResponse("Data not found", "No data with given id", helpers.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helpers.BuildResponse(true, "OK", elemes)
		context.JSON(http.StatusOK, res)
	}
}

func (c *elemesHandler) Insert(context *gin.Context) {
	var elemesCreateDTO dto.ElemesCreateDTO
	errDTO := context.ShouldBind(&elemesCreateDTO)
	//upload
	formFile, _, err := context.Request.FormFile("file")
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			dto.MediaDto{
				StatusCode: http.StatusInternalServerError,
				Message:    "error",
				Data:       map[string]interface{}{"data": "Select a file to upload"},
			})
		return
	}
	uploadUrl, err := c.elemesService.FileUpload(dto.ElemesCreateDTO{
		Title:      elemesCreateDTO.Title,
		File:       formFile,
		Price:      elemesCreateDTO.Price,
		NameCourse: elemesCreateDTO.NameCourse,
	})
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			dto.MediaDto{
				StatusCode: http.StatusInternalServerError,
				Message:    "error",
				Data:       map[string]interface{}{"data": "Error uploading file"},
			})
		return
	}
	context.JSON(
		http.StatusOK,
		dto.MediaDto{
			StatusCode: http.StatusOK,
			Message:    "success",
			Data: map[string]interface{}{
				"title":       elemesCreateDTO.Title,
				"name_course": elemesCreateDTO.NameCourse,
				"price":       elemesCreateDTO.Price,
				"category":    elemesCreateDTO.Category,
				"file":        uploadUrl,
			},
		})
	if errDTO != nil {
		res := helpers.BuildErrorResponse("Failed to process request", errDTO.Error(), helpers.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
}

func (c *elemesHandler) Update(context *gin.Context) {
	var elemesUpdateDTO dto.ElemesUpdateDTO

	errDTO := context.ShouldBind(&elemesUpdateDTO)
	if errDTO != nil {
		res := helpers.BuildErrorResponse("Failed to process request", errDTO.Error(), helpers.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	if c.elemesService.IsAllowedToEdit(elemesUpdateDTO.ID) {
		result := c.elemesService.Update(elemesUpdateDTO)
		response := helpers.BuildResponse(true, "OK", result)
		context.JSON(http.StatusOK, response)
	}
}

func (c *elemesHandler) Delete(context *gin.Context) {
	var elemes models.Elemes
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helpers.BuildErrorResponse("Failed tou get id", "No param id were found", helpers.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	elemes.ID = id
	if c.elemesService.IsAllowedToEdit(elemes.ID) {
		c.elemesService.Delete(elemes)
		res := helpers.BuildResponse(true, "Deleted", helpers.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		response := helpers.BuildErrorResponse("You dont have permission", "You are not the owner", helpers.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}
