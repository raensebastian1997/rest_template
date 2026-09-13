package http

import (
	"restapirian/config"
	"restapirian/internal/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController() *UserController {
	db, _ := config.ConnectDb()
	return NewUserControllerWithDB(db)
}

func NewUserControllerWithDB(db *gorm.DB) *UserController {
	return &UserController{
		db: db,
	}
}

func (ctrl *UserController) Index(ctx *gin.Context) {
	if ctrl.db == nil {
		ctx.JSON(500, gin.H{"error": "database unavailable"})
		return
	}

	var data []domain.User

	err := ctrl.db.Find(&data).Error
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to fetch users"})
		return
	}

	ctx.JSON(200, gin.H{"data": data})

}

func (ctrl *UserController) Update(ctx *gin.Context) {

}

func (ctrl *UserController) Create(ctx *gin.Context) {

}

func (ctrl *UserController) Delete(ctx *gin.Context) {

}
