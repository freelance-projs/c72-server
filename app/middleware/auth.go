package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/pkg/model"
)

func Authn() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func Authz(role model.ERole) gin.HandlerFunc {
	return func(c *gin.Context) {}
}
