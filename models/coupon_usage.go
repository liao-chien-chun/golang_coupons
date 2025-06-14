package models

import "time"

type CouponUsage struct {
	ID        uint       `json:"id" gorm:"primaryKey;comment:使用紀錄ID"`
	UserID    uint       `json:"user_id" gorm:"not null;comment:使用者ID"`
	CouponID  uint       `json:"coupon_id" gorm:"not null;comment:優惠券ID"`
	Status    string     `json:"status" gorm:"type:varchar(20);not null;comment:狀態（unused/used/expired）"`
	UsedAt    *time.Time `json:"used_at" gorm:"comment:實際使用時間"`
	CreatedAt time.Time  `json:"created_at" gorm:"comment:建立時間"`
}
