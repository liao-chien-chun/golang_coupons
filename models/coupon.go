package models

import "time"

type Coupon struct {
	ID        uint      `json:"id" gorm:"primaryKey;comment:優惠券ID"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null;comment:優惠券名稱"`
	Type      string    `json:"type" gorm:"type:varchar(20);not null;comment:優惠券類型"`
	Discount  float64   `json:"discount" gorm:"comment:折扣金額或百分比"`
	Threshold float64   `json:"threshold" gorm:"comment:滿額門檻"`
	Total     int       `json:"total" gorm:"comment:發行總數"`
	Redeemed  int       `json:"redeemed" gorm:"comment:已被領取數"`
	StartAt   time.Time `json:"start_at" gorm:"comment:開始時間"`
	EndAt     time.Time `json:"end_at" gorm:"comment:結束時間"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:建立時間"`
}
