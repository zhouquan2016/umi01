package user

import "github.com/gin-gonic/gin"

func LoadRoutes(r *gin.Engine) {
	r.GET("/menu/tree", FindMenus)
	r.POST("/user/list", List)
	r.POST("/admin/user/save", Save)
	r.GET("/admin/user/getById", GetById)
	r.GET("/admin/user/deleteById", DeleteById)
	r.POST("/admin/user/edit", Edit)
}
