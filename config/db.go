package config

import (
	"github.com/dnopihananto/gin-full-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/learn_gin?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect")
	}
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Article{})

	DB.Model(&models.User{}).Association("Article")
}
