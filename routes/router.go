package routes

import (
	"github.com/ericliao/coupon-system/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 領取優惠券
	r.POST("/coupons/:id/redeem", controllers.RedeemCoupon)

	// 使用優惠券
	r.POST("/coupons/:id/use", controllers.UseCoupon)

	// 查詢使用者所有優惠願
	r.GET("/users/:id/coupons", controllers.GetUserCoupons)

	return r
}
