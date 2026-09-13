package domain

import "time"

type User struct {
	Id        int64     `gorm:column="id"`
	Username  string    `gorm:column="username"`
	Password  string    `gorm:column="password"`
	Email     string    `gorm:column="email"`
	Name      string    `gorm:column="name"`
	Updatedat int64     `gorm:column="updated_at"`
	Createdat time.Time `gorm:column="created_at"`
	Createdby int64     `gorm:column="created_by"`
}

func (User) TableName() string {
	return "user"
}
