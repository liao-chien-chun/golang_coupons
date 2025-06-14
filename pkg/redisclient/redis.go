// pkg/redisclient/redis.go
package redisclient

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client          // 全域 Redis 實例
	Ctx = context.Background() // 全域 Context
)

// InitRedis 初始化 Redis 連線
func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "", // 如果有密碼請設定
		DB:       0,  // 預設用第 0 個 DB
	})

	// 檢查連線
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic("Redis 連線失敗：" + err.Error())
	}

	fmt.Println("Redis 連線成功")
}
