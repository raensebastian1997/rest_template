package http

import (
	"restapirian/config"
	"restapirian/internal/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// func
type ProductController struct {
	db *gorm.DB
}

func NewProductController() *ProductController {
	db, _ := config.ConnectDb()
	return &ProductController{db: db}
}

func (ctrl *ProductController) Index(c *gin.Context) {
	var data []domain.Product

	if err := ctrl.db.Find(&data); err != nil {
		c.JSON(200, gin.H{"error": "database unavailables"})
		return
	}
	c.JSON(200, gin.H{"error": "database unavailable", "data": data})
	// return
}

func (ctrl *ProductController) Create(c *gin.Context) {
	var data domain.Product

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Save db unavailable", "data": nil})
		return
	}
	if result := ctrl.db.Table("product").Create(&data); result.Error != nil {
		c.JSON(400, gin.H{"error": "cnapt save because ", "data": nil})
		return
	}
	c.JSON(200, gin.H{"error": "save successfully", "data": data})

}
