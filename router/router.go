package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	server := gin.Default()
	server.LoadHTMLFiles("template/index.html")
	server.Static("/static", "./static")
	server.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})
	router := server.Group("/api")
	{
		GetUserRouters(router)
		GetCommunityRouters(router)
		GetPostRouters(router)
	}
	return server
}
