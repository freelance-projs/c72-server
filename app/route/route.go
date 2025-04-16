package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/c72-api-server/app/middleware"
	"github.com/ngoctd314/c72-api-server/app/usecase"
	"github.com/ngoctd314/c72-api-server/app/usecase/company"
	"github.com/ngoctd314/c72-api-server/app/usecase/department"
	"github.com/ngoctd314/c72-api-server/app/usecase/setting"
	"github.com/ngoctd314/c72-api-server/app/usecase/stat"
	"github.com/ngoctd314/c72-api-server/app/usecase/tag"
	"github.com/ngoctd314/c72-api-server/app/usecase/tagname"
	"github.com/ngoctd314/c72-api-server/app/usecase/txlog"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
)

func Handler(repo *repository.Repository) *gin.Engine {
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

	api.POST("/tag-names/upload", ghttp.GinHandleFunc(tagname.CreateByUpload(repo)))
	api.DELETE("/tag-names", ghttp.GinHandleFunc(tagname.DeleteBatch(repo)))
	api.GET("/tag-names", ghttp.GinHandleFunc(tagname.List(repo)))
	api.PATCH("/tag-names", ghttp.GinHandleFunc(tagname.Change(repo)))

	api.POST("/tag-departments", ghttp.GinHandleFunc(tag.CreateTagDepartment(repo)))
	api.GET("/tag-departments", ghttp.GinHandleFunc(tag.ListTagDepartment(repo)))
	api.POST("/tag-companies", ghttp.GinHandleFunc(tag.CreateTagCompany(repo)))
	api.GET("/tag-companies", ghttp.GinHandleFunc(tag.ListTagCompanies(repo)))

	api.POST("/departments", ghttp.GinHandleFunc(department.Create(repo)))
	api.POST("/departments/upload", ghttp.GinHandleFunc(department.CreateByUpload(repo)))
	api.GET("/departments", ghttp.GinHandleFunc(department.List(repo)))
	api.DELETE("/departments", ghttp.GinHandleFunc(department.DeleteBatch(repo)))
	api.PATCH("/departments/by-name", ghttp.GinHandleFunc(department.Change(repo)))

	api.POST("/companies/upload", ghttp.GinHandleFunc(company.CreateByUpload(repo)))
	api.GET("/companies", ghttp.GinHandleFunc(company.List(repo)))
	api.DELETE("/companies", ghttp.GinHandleFunc(company.DeleteBatch(repo)))
	api.PATCH("/companies/by-name", ghttp.GinHandleFunc(company.Change(repo)))

	api.GET("/settings", ghttp.GinHandleFunc(setting.Get(repo)))
	api.PATCH("/settings", ghttp.GinHandleFunc(setting.Update(repo)))

	api.POST("/tx-log", ghttp.GinHandleFunc(txlog.Create(repo)))
	api.GET("/tx-log/departments", ghttp.GinHandleFunc(txlog.ListDept(repo)))
	api.GET("/tx-log/departments/:id", ghttp.GinHandleFunc(txlog.GetDept(repo)))
	api.GET("/tx-log/companies", ghttp.GinHandleFunc(txlog.ListCompany(repo)))
	api.GET("/tx-log/companies/:id", ghttp.GinHandleFunc(txlog.GetCompany(repo)))

	api.GET("/stats/departments", ghttp.GinHandleFunc(stat.ListDepartment(repo)))
	api.GET("/stats/departments/:department", ghttp.GinHandleFunc(stat.GetDepartment(repo)))

	api.GET("/stats/tags", ghttp.GinHandleFunc(stat.ListTag(repo)))
	api.GET("/stats/tags/:tag_name", ghttp.GinHandleFunc(stat.GetTag(repo)))

	api.GET("/stats/companies", ghttp.GinHandleFunc(stat.ListCompany(repo)))
	api.GET("/stats/companies/:company", ghttp.GinHandleFunc(stat.GetCompany(repo)))

	api.GET("/ping", ghttp.GinHandleFunc(usecase.Ping()))

	return mux
}

func noRouteHandleFunc(c *gin.Context) {
	c.JSON(http.StatusNotFound, dto.Response{
		Success: false,
		Message: "PAGE_NOT_FOUND",
	})
}
