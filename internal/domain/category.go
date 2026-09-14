package domain

import "time"

type Category struct {
	ID        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;size:120;not null;uniqueIndex" json:"name"`
	Code      string    `gorm:"column:code;size:50;not null;uniqueIndex" json:"code"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Category) TableName() string { return "category" }
