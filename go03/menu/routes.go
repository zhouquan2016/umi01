package menu

import "github.com/gin-gonic/gin"

func LoadRoute(r *gin.Engine) {
	r.GET("/menu/:id", GetById)
	r.GET("/menu/existsByPath", ExistsByPath)
	r.POST("/menu/children", Children)
	r.POST("/admin/menu/add", Add)
	r.POST("/admin/menu/edit", Edit)
	r.POST("/admin/menu/delete", Delete)
	r.POST("/admin/menu/list", List)
	r.POST("/admin/menu/addResource", AddResource)
	r.GET("/admin/menu/tree", Tree)
}
