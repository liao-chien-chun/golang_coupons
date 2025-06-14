package routes

import (
	"github.com/ericliao/coupon-system/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 領取優惠券
	r.POST("/coupons/:id/redeem", controllers.RedeemCoupon)

	return r
}
