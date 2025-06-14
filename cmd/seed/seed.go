package main

import (
	"fmt"
	"time"

	"github.com/ericliao/coupon-system/config"
	"github.com/ericliao/coupon-system/models"
)

func main() {
	config.InitDB()
	db := config.DB

	// 建立使用者（僅建立一次，不重複）
	users := []models.User{
		{Name: "Alice"},
		{Name: "Bob"},
	}
	for _, u := range users {
		db.FirstOrCreate(&u, models.User{Name: u.Name})
	}

	// 建立 10 張多樣性優惠券
	coupons := []models.Coupon{
		{Name: "滿100折20", Type: "threshold", Discount: 20, Threshold: 100, Total: 50, StartAt: time.Now().AddDate(0, 0, -2), EndAt: time.Now().AddDate(0, 0, 5)},
		{Name: "滿200折50", Type: "threshold", Discount: 50, Threshold: 200, Total: 30, StartAt: time.Now().AddDate(0, 0, -1), EndAt: time.Now().AddDate(0, 0, 10)},
		{Name: "9折優惠券", Type: "discount", Discount: 0.9, Threshold: 0, Total: 100, StartAt: time.Now().AddDate(0, 0, -3), EndAt: time.Now().AddDate(0, 0, 7)},
		{Name: "8折限時券", Type: "discount", Discount: 0.8, Threshold: 0, Total: 80, StartAt: time.Now().AddDate(0, 0, -5), EndAt: time.Now().AddDate(0, 0, 2)},
		{Name: "滿300折100", Type: "threshold", Discount: 100, Threshold: 300, Total: 40, StartAt: time.Now().AddDate(0, 0, -1), EndAt: time.Now().AddDate(0, 0, 3)},
		{Name: "7折驚喜券", Type: "discount", Discount: 0.7, Threshold: 0, Total: 20, StartAt: time.Now().AddDate(0, 0, -2), EndAt: time.Now().AddDate(0, 0, -1)},        // 過期
		{Name: "滿500折200", Type: "threshold", Discount: 200, Threshold: 500, Total: 60, StartAt: time.Now().AddDate(0, 0, -10), EndAt: time.Now().AddDate(0, 0, -2)}, // 過期
		{Name: "85折春季券", Type: "discount", Discount: 0.85, Threshold: 0, Total: 120, StartAt: time.Now().AddDate(0, 0, -5), EndAt: time.Now().AddDate(0, 0, 8)},
		{Name: "滿150折30", Type: "threshold", Discount: 30, Threshold: 150, Total: 25, StartAt: time.Now().AddDate(0, 0, -2), EndAt: time.Now().AddDate(0, 0, 1)},
		{Name: "95折會員券", Type: "discount", Discount: 0.95, Threshold: 0, Total: 200, StartAt: time.Now().AddDate(0, 0, -3), EndAt: time.Now().AddDate(0, 0, 9)},
	}

	for _, c := range coupons {
		db.Create(&c)
	}

	fmt.Println("成功建立使用者與優惠券種子資料")
}
