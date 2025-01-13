package route

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/app/usecase/tag"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/helper"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/env"
)

func Handler(tagRepo *repository.Tag) *gin.Engine {
	mux := gin.New()

	// register no route
	mux.NoRoute(noRouteHandleFunc)
	// register global middlewares
	mux.Use(
		gin.Recovery(),
	)
	mux.Use(cors.New(cors.Config{
		AllowOrigins:     env.GetStringSlice("http.cors.allowOrigins"),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	apiGroup := mux.Group("/api")

	apiGroup.POST("/tags", tag.ScanInBatches(tagRepo))
	apiGroup.GET("/tags/:id", tag.GetByID(tagRepo))
	apiGroup.GET("/tags/scan-histories", tag.NewTagScanHistory(tagRepo))
	apiGroup.PATCH("/tags", tag.UpdateTagName(tagRepo))

	apiGroup.GET("/ips", func(c *gin.Context) {
		ip, err := helper.LocalIP()
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, dto.Response{
			Success: true,
			Data:    ip,
		})

	})
	apiGroup.POST("/logs", func(c *gin.Context) {
		var req dto.LogRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			slog.Info("received log", "err", err.Error())

			c.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: "INVALID_REQUEST",
			})
			return
		}
		slog.Info("received log", "message", req.Message, "data", req.Data)

		c.JSON(http.StatusOK, dto.Response{
			Success: true,
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
