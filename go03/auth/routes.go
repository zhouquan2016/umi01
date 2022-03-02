package auth

import "github.com/gin-gonic/gin"

func LoadRoutes(r *gin.Engine) {
	r.POST(LoginPath, Login)
	r.GET(LogoutPath, OutLogin)
	r.GET("/currentUser", CurrentUser)
}
