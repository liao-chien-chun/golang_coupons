package config

import (
	// 格式化字串用
	"fmt"
	// 日誌輸出
	"log"
	// 讀取環境變數
	"os"
	// 引用 GORM model 定義
	"github.com/ericliao/coupon-system/models"

	// GORM MySQL 驅動
	"gorm.io/driver/mysql"
	// GORM 主套件
	"gorm.io/gorm"
)

var DB *gorm.DB // 全域變數DB，做為GORM 的資料庫連線，供其他模組使用

// 初始化資料物連線與建立表 (migration)
func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// 嘗試連線資料庫，如有錯誤 會直接 log.Fatal 結束程式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("無法連接資料庫：", err)
	}

	// 自動建表
	err = db.AutoMigrate(
		&models.User{},
		&models.Coupon{},
		&models.CouponUsage{},
	)
	if err != nil {
		log.Fatal("執行 AutoMigrate 失敗：", err)
	}

	DB = db
	log.Println("成功連線資料庫並執行 migration")
}
