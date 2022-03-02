package history

import "github.com/gin-gonic/gin"

func LoadRoutes(r *gin.Engine) {
	r.POST("/admin/history/list", List)
	r.GET("/admin/history/indexAll", IndexAll)
}
