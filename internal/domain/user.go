package domain

import "time"

type User struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"column:username" json:"username"`
	Password  string    `gorm:"column:password" json:"-"`
	Email     string    `gorm:"column:email" json:"email"`
	Name      string    `gorm:"column:name" json:"name"`
	RoleID    *int64    `gorm:"column:role_id;index" json:"role_id"`
	Role      *Role     `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"role,omitempty"`
	CreatedBy int64     `gorm:"column:created_by" json:"created_by"`
	UpdatedBy int64     `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt int64     `gorm:"column:updated_at;autoUpdateTime:milli" json:"updated_at"`
}

func (User) TableName() string { return "user" }
