package router

import (
	"github.com/gin-gonic/gin"
	"humble_blog/logic"
)

func GetCommunityRouters(router *gin.RouterGroup) {
	group := router.Group("/community")
	{
		group.GET("/category", logic.GetCommunityCategory)
	}
}
