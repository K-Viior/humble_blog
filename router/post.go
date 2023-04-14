package router

import (
	"github.com/gin-gonic/gin"
	"humble_blog/controller"
	"humble_blog/middleware"
)

func GetPostRouters(router *gin.RouterGroup) {
	group := router.Group("/post")
	{
		group.GET("", controller.GetPostList)
		group.GET("/:postId", controller.GetPostById)
	}
	group.Use(middleware.AuthRequired())
	{
		group.POST("", controller.CreatePost)
		group.POST("/vote", controller.PostVote)
	}
}
