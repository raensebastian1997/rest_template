package domain

import "time"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type Role struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;size:40;not null;uniqueIndex" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Role) TableName() string { return "role" }
