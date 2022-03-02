package main

import (
	"github.com/gin-gonic/gin"
	"go03/auth"
	"go03/conf"
	"go03/filter"
	"go03/history"
	"go03/menu"
	"go03/role"
	"go03/user"
	"log"
)

func main() {
	r := gin.New()
	r.Use(filter.HistoryFilter)
	r.Use(conf.ErrorFilter)
	r.Use(auth.Filter)
	r.Use(filter.RolFilter)
	auth.LoadRoutes(r)
	menu.LoadRoute(r)
	user.LoadRoutes(r)
	role.LoadRoutes(r)
	history.LoadRoutes(r)
	log.Panicln(r.Run(":8282"))
}
