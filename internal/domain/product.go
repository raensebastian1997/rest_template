package domain

import "time"

type Product struct {
	ID         int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	Code       string    `gorm:"column:code" json:"code"`
	Price      float64   `gorm:"column:price" json:"price"`
	CategoryID *int32    `gorm:"column:category_id;index" json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"category,omitempty"`
	CreatedBy  int64     `gorm:"column:created_by" json:"created_by"`
	UpdatedBy  int64     `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Product) TableName() string { return "product" }
