package filter

import (
	"github.com/gin-gonic/gin"
	"go03/auth"
	"go03/conf"
	"go03/db"
	"net/http"
)

var WhiteList = []string{auth.LoginPath, auth.LogoutPath}

func RolFilter(context *gin.Context) {
	if inWhiteList(context.Request.RequestURI) {
		return
	}
	val, _ := context.Get(gin.AuthUserKey)
	u := val.(db.User)
	var resourceEntity = db.ResourceDao.FindByPath(context.Request.RequestURI)
	if resourceEntity != nil && !db.ResourceDao.ExistsRolePath(u.RoleId, resourceEntity.Path) {
		context.AbortWithStatusJSON(http.StatusOK, conf.ServiceErrorResult("无权限访问"))
	} else {
		context.Next()
	}

}

func inWhiteList(uri string) bool {
	for _, s := range WhiteList {
		if s == uri {
			return true
		}
	}
	return false
}
