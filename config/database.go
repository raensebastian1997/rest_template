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

	if err := cnn.AutoMigrate(&domain.User{}, &domain.Product{}); err != nil {
		return nil,
			fmt.Errorf("failed to migrate database: %w", err)
	}

	// for i := 1; i < 100000; i++ {
	// 	// domain.User{

	// 	// }
	// 	newUser := domain.User{Name: "Alice", Email: "riaan@gmaola" + string(i), Createdat: time.Now(), Updatedat: time.Now().Unix(), Createdby: int64(i)}
	// 	result := cnn.Create(&newUser) // Pass a pointer to backfill the generated ID
	// 	// 4. Handle errors and metadata
	// 	if result.Error != nil {
	// 		panic(result.Error)
	// 	}

	// }

	return cnn, nil
}
