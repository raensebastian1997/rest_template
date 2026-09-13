package config

import (
	"fmt"
	"restapirian/internal/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDb() (*gorm.DB, error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/cmd_api_db?charset=utf8mb4&parseTime=True&loc=Local"
	cnn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil,
			fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := cnn.AutoMigrate(domain.User{}); err != nil {
		return nil,
			fmt.Errorf("failed to migrate database: %w", err)
	}

	return cnn, nil
}
