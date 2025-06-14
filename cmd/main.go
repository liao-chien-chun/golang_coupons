package main

import (
	"github.com/ericliao/coupon-system/config"
	"github.com/ericliao/coupon-system/pkg/redisclient"
	"github.com/ericliao/coupon-system/routes"
)

func main() {
	// 初始化資料庫連線
	config.InitDB()
	redisclient.InitRedis() // 初始化 Redis

	// r := gin.Default()

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "pong"})
	// })

	// 使用自訂的 router
	r := routes.SetupRouter()

	r.Run(":8080")
}
