package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/app/middleware"
	"github.com/ngoctd314/c72-api-server/app/usecase"
	"github.com/ngoctd314/c72-api-server/app/usecase/company"
	"github.com/ngoctd314/c72-api-server/app/usecase/department"
	"github.com/ngoctd314/c72-api-server/app/usecase/laundry"
	"github.com/ngoctd314/c72-api-server/app/usecase/lending"
	"github.com/ngoctd314/c72-api-server/app/usecase/setting"
	"github.com/ngoctd314/c72-api-server/app/usecase/tag"
	"github.com/ngoctd314/c72-api-server/app/usecase/tagname"
	"github.com/ngoctd314/c72-api-server/app/usecase/txlog"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
)

func Handler(repo *repository.Laundry) *gin.Engine {
	gin.DisableBindValidation()
	mux := gin.New()

	mux.NoRoute(noRouteHandleFunc)
	mux.Use(middleware.Cors())
	mux.Use(gin.Recovery())

	api := mux.Group("/api/v1")

	api.POST("/tags", ghttp.GinHandleFunc(tag.AssignTag(repo)))
	api.GET("/tags", ghttp.GinHandleFunc(tag.List(repo)))
	api.GET("/tags/active", ghttp.GinHandleFunc(tag.GetActiveTagsByIDs(repo)))
	api.PATCH("/tags/:id", ghttp.GinHandleFunc(tag.UpdateTagName(repo)))
	api.DELETE("/tags/:id", ghttp.GinHandleFunc(tag.DeleteByID(repo)))
	api.PATCH("/tags/by-name", ghttp.GinHandleFunc(tag.UpdateTagNameByName(repo)))

	api.POST("/departments", ghttp.GinHandleFunc(department.Create(repo)))
	api.POST("/departments/upload", ghttp.GinHandleFunc(department.CreateByUpload(repo)))
	api.GET("/departments", ghttp.GinHandleFunc(department.List(repo)))
	api.DELETE("/departments", ghttp.GinHandleFunc(department.DeleteBatch(repo)))
	api.PATCH("/departments/by-name", ghttp.GinHandleFunc(department.Change(repo)))

	api.POST("/companies/upload", ghttp.GinHandleFunc(company.CreateByUpload(repo)))
	api.GET("/companies", ghttp.GinHandleFunc(company.List(repo)))
	api.DELETE("/companies", ghttp.GinHandleFunc(company.DeleteBatch(repo)))
	api.PATCH("/companies/by-name", ghttp.GinHandleFunc(company.Change(repo)))

	api.POST("/tag-names/upload", ghttp.GinHandleFunc(tagname.CreateByUpload(repo)))
	api.DELETE("/tag-names", ghttp.GinHandleFunc(tagname.DeleteBatch(repo)))
	api.GET("/tag-names", ghttp.GinHandleFunc(tagname.List(repo)))
	api.PATCH("/tag-names", ghttp.GinHandleFunc(tagname.Change(repo)))

	api.POST("/settings", ghttp.GinHandleFunc(setting.Create(repo)))
	api.GET("/settings", ghttp.GinHandleFunc(setting.List(repo)))
	api.PATCH("/settings", ghttp.GinHandleFunc(setting.Update(repo)))
	api.DELETE("/settings", ghttp.GinHandleFunc(setting.Delete(repo)))

	api.POST("/tx-log", ghttp.GinHandleFunc(txlog.Create(repo)))
	api.GET("/tx-log/departments", ghttp.GinHandleFunc(txlog.ListDept(repo)))
	api.GET("/tx-log/departments/:id", ghttp.GinHandleFunc(txlog.GetDept(repo)))
	api.GET("/tx-log/companies", ghttp.GinHandleFunc(txlog.ListCompany(repo)))
	api.GET("/tx-log/companies/:id", ghttp.GinHandleFunc(txlog.GetCompany(repo)))

	api.GET("/ping", ghttp.GinHandleFunc(usecase.Ping()))

	api.POST("/lending", ghttp.GinHandleFunc(lending.DoLending(repo)))
	api.GET("/lending/:id/tags", ghttp.GinHandleFunc(lending.GetTags(repo)))
	api.GET("/lending", ghttp.GinHandleFunc(lending.List(repo)))
	api.PATCH("/lending/return-dirty", ghttp.GinHandleFunc(lending.ReturnDirty(repo)))

	api.POST("/laundry", ghttp.GinHandleFunc(laundry.DoLaundry(repo)))
	api.GET("/laundry/:id", ghttp.GinHandleFunc(laundry.Get(repo)))
	api.GET("/washing/:id/tags", ghttp.GinHandleFunc(laundry.GetTags(repo)))
	api.GET("/laundry", ghttp.GinHandleFunc(laundry.List(repo)))
	api.PATCH("/laundry/return-clean", ghttp.GinHandleFunc(laundry.ReturnClean(repo)))

	return mux
}

func noRouteHandleFunc(c *gin.Context) {
	c.JSON(http.StatusNotFound, dto.Response{
		Success: false,
		Message: "PAGE_NOT_FOUND",
	})
}
