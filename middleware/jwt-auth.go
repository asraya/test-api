package middleware

import (
	"strings"

	"api-test/service"

	"github.com/gin-gonic/gin"
)

// AuthorizeJWT validates the token user given, return 401 if not valid
func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(403, "No Authorization header provided")
			c.Abort()
			return
		}

		extractedToken := strings.Split(authHeader, "Bearer ")

		if len(extractedToken) == 2 {
			authHeader = strings.TrimSpace(extractedToken[1])
		} else {
			c.JSON(400, "Incorrect Format of Authorization Token")
			c.Abort()

			return
		}

	}
}
