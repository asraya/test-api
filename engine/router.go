package engine

import (
	"net/http"
	"time"

	"api-test/app/handler"
	"api-test/app/repository"
	"api-test/config"
	"api-test/helpers"
	"api-test/service"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

var (
	db *gorm.DB = config.SetupDatabaseConnection()

	movieRepository repository.MovieRepository = repository.NewMovieRepository(db)
	movieService    service.MovieService       = service.NewMovieService(movieRepository)
	movieHandler    handler.MovieHandler       = handler.NewMovieHandler(movieService)
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(logger(db))
	// CORS
	r.Use(CORSMiddleware())

	// Routes
	v1 := r.Group("api/v1")
	{
		routes := v1.Group("/")
		{
			n := routes.Group("/movie")
			{
				n.GET("", func(context *gin.Context) {
					code := http.StatusOK
					pagination := helpers.GeneratePaginationRequest(context)
					response := movieService.PaginationMovie(movieRepository, context, pagination)
					if !response.Status {
						code = http.StatusBadRequest
					}
					context.JSON(code, response)
				})
				n.POST("/", movieHandler.Insert)
				n.GET("/:id", movieHandler.FindByID)
				n.PUT("/:id", movieHandler.Update)
				n.DELETE("/:id", movieHandler.Delete)
			}

		}
	}
	return r
}

// CORSMiddleware ..
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func logger(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now()
		ctx.Next()
		db.Exec("INSERT INTO logs(method,url,code,accesstime,handletime) VALUE(?,?,?,?,?)", ctx.Request.Method, ctx.Request.URL.Path, ctx.Writer.Status(), time.Now().Format(time.RFC822), time.Now().Sub(now).String())
	}
}
