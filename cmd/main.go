package main

import (
	"github.com/ericliao/coupon-system/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化資料庫連線
	config.InitDB()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}
