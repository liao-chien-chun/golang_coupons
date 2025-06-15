package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ericliao/coupon-system/config"
	"github.com/gin-gonic/gin"
)

type CouponWithStatus struct {
	ID      uint      `json:"id"`
	Name    string    `json:"name"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at`
	Status  string    `json:"status"` // unused, used, expired
}

// 查詢使用者所有優惠券狀態
func GetUserCoupons(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"status":  http.StatusBadRequest,
			"message": "查詢失敗",
			"error":   "使用者 ID 無效",
		})
		return
	}

	db := config.DB
	var results []CouponWithStatus

	now := time.Now()

	// 使用 JOIN 撈取優惠券與使用紀錄
	rows, err := db.
		Table("coupon_usages").
		Select("coupons.id, coupons.name, coupons.start_at, coupons.end_at, coupon_usages.status").
		Joins("JOIN coupons ON coupon_usages.coupon_id = coupons.id").
		Where("coupon_usages.user_id = ?", userId).
		Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"status":  http.StatusInternalServerError,
			"message": "查詢失敗",
			"error":   "資料庫錯誤",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r CouponWithStatus
		if err := rows.Scan(&r.ID, &r.Name, &r.StartAt, &r.EndAt, &r.Status); err != nil {
			continue
		}

		if r.Status == "unused" && now.After(r.EndAt) {
			r.Status = "expired"
		}

		results = append(results, r)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  http.StatusOK,
		"message": "查詢成功",
		"data":    results,
	})
}
