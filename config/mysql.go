package config

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"humble_blog/model"
	"time"
)

var DB *gorm.DB

func InitDB() {

	dsn := "root:1234@tcp(127.0.0.1:3306)/humble_blog?charset=utf8mb4&parseTime=True&loc=Local"
	//连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		zap.L().Error("failed to connect database ,err: " + err.Error())
	}
	//获得数据库连接池
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Error("failed to get server ,err :" + err.Error())
	}
	sqlDB.SetMaxIdleConns(10)           //设置连接池中空闲的最大数量
	sqlDB.SetMaxOpenConns(100)          //设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) //设置连接可复用的最大时间
	//数据库迁移
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Category{})
	db.AutoMigrate(&model.Post{})
	db.Create(&model.Category{
		BaseModel:   model.BaseModel{ID: 1, CreateTime: time.Now(), UpdateTime: time.Now()},
		Name:        "Golang",
		Description: "Go社区",
	})
	db.Create(&model.Category{
		BaseModel:   model.BaseModel{ID: 2, CreateTime: time.Now(), UpdateTime: time.Now()},
		Name:        "java",
		Description: "java社区",
	})
	db.Create(&model.Category{
		BaseModel:   model.BaseModel{ID: 3, CreateTime: time.Now(), UpdateTime: time.Now()},
		Name:        "LeetCode",
		Description: "LeetCode社区",
	})
	zap.L().Info("database init success")
	DB = db
}

func GetDB() *gorm.DB {
	return DB
}
