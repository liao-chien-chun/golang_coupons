package main

import (
	"fmt"
	"time"

	"github.com/ericliao/coupon-system/config"
	"github.com/ericliao/coupon-system/models"
	"gorm.io/gorm"
)

func main() {
	config.InitDB()
	db := config.DB

	// 清空資料表並重置 auto_increment
	resetTables(db)

	// 建立使用者
	users := []models.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	for _, u := range users {
		db.Create(&u)
	}

	// // 建立使用者（僅建立一次，不重複）
	// users := []models.User{
	// 	{Name: "Alice"},
	// 	{Name: "Bob"},
	// }
	// for _, u := range users {
	// 	db.FirstOrCreate(&u, models.User{Name: u.Name})
	// }

	// 建立 10 種多樣性優惠券
	coupons := []models.Coupon{
		{ID: 1, Name: "滿100折20", Type: "threshold", Discount: 20, Threshold: 100, Total: 50, Redeemed: 5, StartAt: time.Now().AddDate(0, 0, -2), EndAt: time.Now().AddDate(0, 0, 5)},
		{ID: 2, Name: "滿200折50", Type: "threshold", Discount: 50, Threshold: 200, Total: 30, Redeemed: 30, StartAt: time.Now().AddDate(0, 0, -1), EndAt: time.Now().AddDate(0, 0, 10)}, // 已領完
		{ID: 3, Name: "9折優惠券", Type: "discount", Discount: 0.9, Threshold: 0, Total: 100, Redeemed: 20, StartAt: time.Now().AddDate(0, 0, -3), EndAt: time.Now().AddDate(0, 0, 7)},
		{ID: 4, Name: "8折限時券", Type: "discount", Discount: 0.8, Threshold: 0, Total: 80, Redeemed: 0, StartAt: time.Now().AddDate(0, 0, -5), EndAt: time.Now().AddDate(0, 0, 2)},
		{ID: 5, Name: "滿300折100", Type: "threshold", Discount: 100, Threshold: 300, Total: 40, Redeemed: 10, StartAt: time.Now().AddDate(0, 0, -1), EndAt: time.Now().AddDate(0, 0, 3)},
		{ID: 6, Name: "7折驚喜券", Type: "discount", Discount: 0.7, Threshold: 0, Total: 20, Redeemed: 5, StartAt: time.Now().AddDate(0, 0, -4), EndAt: time.Now().AddDate(0, 0, -1)},         // 過期
		{ID: 7, Name: "滿500折200", Type: "threshold", Discount: 200, Threshold: 500, Total: 60, Redeemed: 60, StartAt: time.Now().AddDate(0, 0, -10), EndAt: time.Now().AddDate(0, 0, -2)}, // 過期領完
		{ID: 8, Name: "85折春季券", Type: "discount", Discount: 0.85, Threshold: 0, Total: 120, Redeemed: 50, StartAt: time.Now().AddDate(0, 0, -5), EndAt: time.Now().AddDate(0, 0, 8)},
		{ID: 9, Name: "滿150折30", Type: "threshold", Discount: 30, Threshold: 150, Total: 25, Redeemed: 0, StartAt: time.Now().AddDate(0, 0, -2), EndAt: time.Now().AddDate(0, 0, 1)},
		{ID: 10, Name: "95折會員券", Type: "discount", Discount: 0.95, Threshold: 0, Total: 200, Redeemed: 200, StartAt: time.Now().AddDate(0, 0, -3), EndAt: time.Now().AddDate(0, 0, 9)},
	}
	for _, c := range coupons {
		db.Create(&c)
	}

	// 建立使用者1的優惠券使用紀錄
	usages := []models.CouponUsage{
		{UserID: 1, CouponID: 1, Status: "unused", CreatedAt: time.Now()},
		{UserID: 1, CouponID: 2, Status: "used", UsedAt: ptrTime(time.Now().AddDate(0, 0, -1)), CreatedAt: time.Now().AddDate(0, 0, -1)},
		{UserID: 1, CouponID: 6, Status: "expired", CreatedAt: time.Now().AddDate(0, 0, -3)},
	}
	for _, u := range usages {
		db.Create(&u)
	}

	fmt.Println("成功建立使用者與優惠券種子資料")
}

// 清空資料表
func resetTables(db *gorm.DB) {
	db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	db.Exec("TRUNCATE TABLE coupon_usages;")
	db.Exec("TRUNCATE TABLE coupons;")
	db.Exec("TRUNCATE TABLE users;")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1;")
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
