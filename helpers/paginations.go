package helpers

import (
	"strconv"
	"strings"

	"api-test/dto"

	"github.com/gin-gonic/gin"
)

func GeneratePaginationRequest(context *gin.Context) *dto.Pagination {
	limit := 10
	page := 1
	apikey := "faf7e5bb"
	sort := "created_at desc"

	var searchs []dto.Search

	query := context.Request.URL.Query()

	for key, value := range query {
		queryValue := value[len(value)-1]

		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		case "apikey":
			apikey = queryValue
			break
		case "sort":
			sort = queryValue
			break
		}
		if strings.Contains(key, ".") {
			searchKeys := strings.Split(key, ".")
			search := dto.Search{Column: searchKeys[0], Action: searchKeys[1], Query: queryValue}
			searchs = append(searchs, search)
		}
	}

	return &dto.Pagination{Limit: limit, Page: page, Apikey: apikey, Sort: sort, Searchs: searchs}
}
