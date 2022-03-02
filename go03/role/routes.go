package role

import (
	"github.com/gin-gonic/gin"
)

func LoadRoutes(r *gin.Engine) {
	r.POST("/admin/role/list", List)
	r.GET("/admin/role/exists", Exists)
	r.POST("/admin/role/add", Add)
	r.GET("/admin/role/delete", Delete)
	r.POST("/admin/role/edit", Update)
	r.GET("/admin/role/getById", GetById)
	r.GET("/admin/role/getAll", GetAll)
}
