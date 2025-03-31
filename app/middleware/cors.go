package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/common/env"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     env.GetStringSlice("http.cors.allowOrigins"),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, MaxAge: 12 * time.Hour,
	})

}
