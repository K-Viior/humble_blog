package main

import (
	"humble_blog/config"
	"humble_blog/pkg/snowflake"
	"humble_blog/router"
)

func main() {
	//初始化数据库
	config.InitDB()
	//初始化redis
	config.InitRedis()
	//初始化雪花算法
	snowflake.Init("2023-01-01", 1)
	//初始化路由
	server := router.InitRouter()

	server.Run(":8081")
}
