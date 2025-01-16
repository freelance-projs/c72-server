package route

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/app/usecase/tag"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/env"
)

func Handler(tagRepo *repository.Tag) *gin.Engine {
	gin.DisableBindValidation()
	mux := gin.New()

	// register no route
	mux.NoRoute(noRouteHandleFunc)
	// register global middlewares
	mux.Use(
		gin.Recovery(),
	)
	mux.Use(cors.New(cors.Config{
		AllowOrigins:     env.GetStringSlice("http.cors.allowOrigins"),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	apiGroup := mux.Group("/api")

	apiGroup.POST("/tags", tag.ScanInBatches(tagRepo))
	apiGroup.POST("/tags-mapping/upload", tag.TagMappingUpload(tagRepo))
	apiGroup.GET("/tags/:id", tag.GetByID(tagRepo))
	apiGroup.GET("/tags", tag.List(tagRepo))
	apiGroup.GET("/tags/scan-histories", tag.NewTagScanHistory(tagRepo))
	apiGroup.PATCH("/tags", tag.UpdateTagName(tagRepo))
	apiGroup.DELETE("/tags/:id", tag.DeleteByID(tagRepo))
	apiGroup.PATCH("/tags/by-name", tag.UpdateTagNameByName(tagRepo))
	apiGroup.DELETE("/tags/by-name/:name", tag.DeleteByName(tagRepo))

	apiGroup.GET("/host-ips", func(c *gin.Context) {
		c.JSON(http.StatusOK, dto.Response{
			Success: true,
			Data:    os.Getenv("HOST_IP"),
		})
	})
	apiGroup.GET("/ping", func(c *gin.Context) {
		slog.Info("/ping")
		c.JSON(http.StatusOK, dto.Response{
			Success: true,
			Data:    "pong",
		})
	})

	return mux
}

func noRouteHandleFunc(c *gin.Context) {
	c.JSON(http.StatusNotFound, dto.Response{
		Success: false,
		Message: "PAGE_NOT_FOUND",
	})
}
