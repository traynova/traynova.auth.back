package middleware

import "github.com/gin-gonic/gin"

func SetupGinLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Skip: func(c *gin.Context) bool {
			return c.Request.Method == "OPTIONS"
		},
	})
}
