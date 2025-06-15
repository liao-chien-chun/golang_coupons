// controllers/coupon_controller.go
package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ericliao/coupon-system/config"
	"github.com/ericliao/coupon-system/models"
	"github.com/ericliao/coupon-system/pkg/redisclient"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// RedeemCoupon handles POST /coupons/:id/redeem
func RedeemCoupon(c *gin.Context) {
	couponIDStr := c.Param("id")
	couponID, err := strconv.Atoi(couponIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"status":  http.StatusBadRequest,
			"message": "領券失敗",
			"error":   "優惠券 ID 無效",
		})
		return
	}

	// 模擬使用者 ID，實務上應由 JWT 或登入 session 取得
	userID := uint(1)

	db := config.DB

	// Redis 上鎖（悲觀鎖）避免同時領券
	lockKey := fmt.Sprintf("lock:coupon:%d", couponID)
	lockSuccess, err := redisclient.Rdb.SetNX(redisclient.Ctx, lockKey, "locked", 3*time.Second).Result()
	if err != nil || !lockSuccess {
		log.Printf("Redis 鎖定失敗：%v", err)
		c.JSON(http.StatusTooManyRequests, gin.H{
			"success": false,
			"status":  http.StatusTooManyRequests,
			"message": "領取失敗",
			"error":   "系統繁忙，請稍後再試",
		})
		return
	}
	defer func() {
		log.Printf("釋放 Redis 鎖: %s", lockKey)
		redisclient.Rdb.Del(redisclient.Ctx, lockKey) // 確保最後會釋放鎖
	}()

	// 檢查是否已領取過
	var existing models.CouponUsage
	db.Where("user_id = ? AND coupon_id = ?", userID, couponID).First(&existing)
	if existing.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"status":  http.StatusConflict,
			"message": "領取失敗",
			"error":   "已領取過該優惠券",
		})
		return
	}

	// 先從 Redis 讀取快取資料
	var coupon models.Coupon
	redisKey := fmt.Sprintf("coupon_data:%d", couponID)
	cachedData, err := redisclient.Rdb.Get(redisclient.Ctx, redisKey).Result()
	if err == nil {
		log.Printf("從 Redis 快取取得 coupon #%d", couponID)
		// 成功從 Redis 取得快取，反序列化
		if err := json.Unmarshal([]byte(cachedData), &coupon); err != nil {
			log.Printf("反序列化快取失敗，fallback 查 DB: %v", err)
			// 快取壞掉 fallback 查資料庫
			db.First(&coupon, couponID)
		}
	} else if err == redis.Nil {
		log.Printf("Redis 無資料，從 DB 查 coupon #%d 並寫入快取", couponID)
		// Redis 無此快取 → 從資料庫撈
		db.First(&coupon, couponID)
		if coupon.ID != 0 {
			// 寫入 Redis 快取，預設 10 分鐘 TTL
			if jsonBytes, err := json.Marshal(coupon); err == nil {
				redisclient.Rdb.Set(redisclient.Ctx, redisKey, jsonBytes, 10*time.Minute)
			}
		}
	} else {
		log.Printf("Redis 讀取錯誤，fallback 查 DB: %v", err)
		// Redis 有問題 fallback DB
		db.First(&coupon, couponID)
	}

	if coupon.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"status":  http.StatusNotFound,
			"message": "領取失敗",
			"error":   "找不到優惠券",
		})
		return
	}

	now := time.Now()
	if now.Before(coupon.StartAt) || now.After(coupon.EndAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"status":  http.StatusBadRequest,
			"message": "領取失敗",
			"error":   "優惠券已過期或未開始",
		})
		return
	}

	if coupon.Redeemed >= coupon.Total {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"status":  http.StatusBadRequest,
			"message": "領取失敗",
			"error":   "優惠券已被領完",
		})
		return
	}

	// 更新已領取數量（建議包在交易中）
	db.Model(&coupon).Update("redeemed", coupon.Redeemed+1)

	usage := models.CouponUsage{
		UserID:    userID,
		CouponID:  coupon.ID,
		Status:    "unused",
		CreatedAt: now,
	}
	db.Create(&usage)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  http.StatusOK,
		"message": "優惠券領取成功",
	})
}

// 使用優惠券 API POST /coupons/:id/use
func UseCoupon(c *gin.Context) {
	// 模擬登入的使用者 id
	userId := 1

	// 從 URL 參數取得優惠券 ID
	couponId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"status":  http.StatusBadRequest,
			"message": "請求格式錯誤",
			"error":   "優惠券 ID 應為數字",
		})
		return
	}

	db := config.DB

	// 查詢該使用者的該張優惠券使用紀錄，且狀態為未使用
	var usage models.CouponUsage
	if err := db.Where("user_id = ? AND coupon_id =  ?", userId, couponId).First(&usage).Error; err != nil {
		status := http.StatusInternalServerError
		msg := "查詢使用紀錄失敗"
		if err == gorm.ErrRecordNotFound {
			status = http.StatusBadRequest
			msg = "你尚未領取此優惠券"
		}
		c.JSON(status, gin.H{
			"success": false,
			"status":  status,
			"message": "使用失敗",
			"error":   msg,
		})
		return
	}

	// 查詢優惠券資料，檢查是否過期
	var coupon models.Coupon
	if err := db.First(&coupon, couponId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"status":  http.StatusBadRequest,
			"message": "使用失敗",
			"error":   "找不到優惠券資料",
		})
		return
	}

	now := time.Now()
	if now.Before(coupon.StartAt) || now.After(coupon.EndAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"status":  http.StatusBadRequest,
			"message": "使用失敗",
			"error":   "優惠券已過期或尚未開始",
		})
		return
	}

	// 狀態不是 unused 就不能使用
	if usage.Status != "unused" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"status":  http.StatusBadRequest,
			"message": "使用失敗",
			"error":   "此優惠券已使用或無法再使用",
		})
		return
	}

	// 更新使用紀錄為已使用
	usage.Status = "used"
	usage.UsedAt = &now
	if err := db.Save(&usage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"status":  http.StatusInternalServerError,
			"message": "使用失敗",
			"error":   "更新優惠券狀態失敗",
		})
		return
	}

	// 成功回傳
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  http.StatusOK,
		"message": "優惠券使用成功",
	})
}
