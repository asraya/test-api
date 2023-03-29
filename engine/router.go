package engine

import (
	"net/http"
	"time"

	"api-test/app/handler"
	"api-test/app/repository"
	"api-test/config"
	"api-test/helpers"
	"api-test/middleware"
	"api-test/service"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

var (
	db         *gorm.DB           = config.SetupDatabaseConnection()
	jwtService service.JWTService = service.NewJWTService()

	elemesRepository repository.ElemesRepository = repository.NewElemesRepository(db)
	elemesService    service.ElemesService       = service.NewElemesService(elemesRepository)
	elemesHandler    handler.ElemesHandler       = handler.NewElemesHandler(elemesService, jwtService)

	userRepository repository.UserRepository = repository.NewUserRepository(db)
	authService    service.AuthService       = service.NewAuthService(userRepository)
	authHandler    handler.AuthHandler       = handler.NewAuthHandler(authService, jwtService)
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(logger(db))
	// CORS
	r.Use(CORSMiddleware())

	// r.POST("/file", elemesHandler.FileUpload())
	// r.POST("/remote", elemesHandler.RemoteUpload())
	// Routes
	v1 := r.Group("api/v1", middleware.AuthorizeJWT(jwtService))
	{
		routes := v1.Group("/")
		{
			n := routes.Group("/elemes")
			{
				n.GET("", func(context *gin.Context) {
					code := http.StatusOK
					pagination := helpers.GeneratePaginationRequest(context)
					response := elemesService.PaginationElemes(elemesRepository, context, pagination)
					if !response.Status {
						code = http.StatusBadRequest
					}
					context.JSON(code, response)
				})
				n.POST("/", elemesHandler.Insert)
				n.GET("/:id", elemesHandler.FindByID)
				n.PUT("/:id", elemesHandler.Update)
				n.DELETE("/:id", elemesHandler.Delete)
			}

		}
		authRoutes := r.Group("api/auth")
		{
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/register", authHandler.Register)
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
