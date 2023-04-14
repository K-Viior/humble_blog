package router

import (
	"github.com/gin-gonic/gin"
	"humble_blog/controller"
)

func GetUserRouters(router *gin.RouterGroup) {
	group := router.Group("/user")
	{
		group.POST("/register", controller.Register)
		group.POST("/login", controller.Login)
	}
}
