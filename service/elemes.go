package service

import (
	"fmt"
	"log"

	"api-test/app/repository"
	"api-test/dto"
	"api-test/helpers"
	"api-test/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mashingan/smapping"
)

// ElemesService is a
type ElemesService interface {
	Insert(b dto.ElemesCreateDTO) models.Elemes
	Update(b dto.ElemesUpdateDTO) models.Elemes
	Delete(b models.Elemes)
	All() []models.Elemes
	FindByID(elemesID uint64) models.Elemes
	IsAllowedToEdit(elemesID uint64) bool
	PaginationElemes(repo repository.ElemesRepository, context *gin.Context, pagination *dto.Pagination) dto.Response
	FileUpload(file dto.ElemesCreateDTO) (string, error)
	RemoteUpload(url dto.ElemesCreateDTO) (string, error)
}

type elemesService struct {
	elemesRepository repository.ElemesRepository
}

var (
	validate = validator.New()
)

// NewElemesService
func NewElemesService(elemesRepo repository.ElemesRepository) ElemesService {
	return &elemesService{
		elemesRepository: elemesRepo,
	}
}

func (service *elemesService) FileUpload(b dto.ElemesCreateDTO) (string, error) {
	//validate
	err := validate.Struct(b)
	if err != nil {
		return "", err
	}

	//upload
	uploadUrl, err := helpers.ImageUploadHelper(b.File)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}

func (service *elemesService) RemoteUpload(b dto.ElemesCreateDTO) (string, error) {
	//validate
	err := validate.Struct(b)
	if err != nil {
		return "", err
	}

	//upload
	uploadUrl, errUrl := helpers.ImageUploadHelper(b)
	if errUrl != nil {
		return "", err
	}
	return uploadUrl, nil
}

func (service *elemesService) Insert(b dto.ElemesCreateDTO) models.Elemes {
	elemes := models.Elemes{}

	err := smapping.FillStruct(&elemes, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}

	res := service.elemesRepository.InsertElemes(elemes)
	return res
}

func (service *elemesService) Update(b dto.ElemesUpdateDTO) models.Elemes {
	elemes := models.Elemes{}
	err := smapping.FillStruct(&elemes, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.elemesRepository.UpdateElemes(elemes)
	return res
}

func (service *elemesService) Delete(b models.Elemes) {
	service.elemesRepository.DeleteElemes(b)
}

func (service *elemesService) All() []models.Elemes {
	return service.elemesRepository.AllElemes()
}

func (service *elemesService) FindByID(elemesID uint64) models.Elemes {
	return service.elemesRepository.FindElemesByID(elemesID)
}

func (service *elemesService) IsAllowedToEdit(elemesID uint64) bool {
	b := service.elemesRepository.FindElemesByID(elemesID)
	id := (b.ID)
	return elemesID == id
}

func (service *elemesService) PaginationElemes(repo repository.ElemesRepository, context *gin.Context, pagination *dto.Pagination) dto.Response {

	operationResult, totalPages := repo.PaginationElemes(pagination)

	if operationResult.Error != nil {
		return dto.Response{Status: false, Message: operationResult.Error.Error()}
	}

	var data = operationResult.Result.(*dto.Pagination)

	// get current url path
	urlPath := context.Request.URL.Path

	// search query params
	searchQueryParams := ""

	for _, search := range pagination.Searchs {
		searchQueryParams += fmt.Sprintf("&%s.%s=%s", search.Column, search.Action, search.Query)
	}

	// set first & last page pagination response
	data.FirstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, 0, pagination.Sort) + searchQueryParams
	data.LastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, totalPages, pagination.Sort) + searchQueryParams

	if data.Page > 0 {
		// set previous page pagination response
		data.PreviousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, data.Page-1, pagination.Sort) + searchQueryParams
	}

	if data.Page < totalPages {
		// set next page pagination response
		data.NextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, data.Page+1, pagination.Sort) + searchQueryParams
	}

	if data.Page > totalPages {
		// reset previous page
		data.PreviousPage = ""
	}

	return dto.Response{Status: true, Data: data}
}
