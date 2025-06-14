// controllers/coupon_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ericliao/coupon-system/config"
	"github.com/ericliao/coupon-system/models"

	"github.com/gin-gonic/gin"
)

// RedeemCoupon handles POST /coupons/:id/redeem
func RedeemCoupon(c *gin.Context) {
	couponIDStr := c.Param("id")
	couponID, err := strconv.Atoi(couponIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "優惠券 ID 無效"})
		return
	}

	// 模擬使用者 ID，實務上應由 JWT 或登入 session 取得
	userID := uint(1)

	db := config.DB

	// 檢查是否已領取過
	var existing models.CouponUsage
	db.Where("user_id = ? AND coupon_id = ?", userID, couponID).First(&existing)
	if existing.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "已領取過該優惠券"})
		return
	}

	// 查詢優惠券
	var coupon models.Coupon
	db.First(&coupon, couponID)
	if coupon.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到優惠券"})
		return
	}

	now := time.Now()
	if now.Before(coupon.StartAt) || now.After(coupon.EndAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "優惠券已過期或未開始"})
		return
	}

	if coupon.Redeemed >= coupon.Total {
		c.JSON(http.StatusBadRequest, gin.H{"error": "優惠券已被領完"})
		return
	}

	// 更新已領取數量
	db.Model(&coupon).Update("redeemed", coupon.Redeemed+1)

	// 建立領取記錄
	usage := models.CouponUsage{
		UserID:    userID,
		CouponID:  coupon.ID,
		Status:    "unused",
		CreatedAt: now,
	}
	db.Create(&usage)

	c.JSON(http.StatusOK, gin.H{"message": "優惠券領取成功"})
}
