package router

import (
	"restapirian/internal/delivery/http"

	"github.com/gin-gonic/gin"
)

func InitializeRouter() error {
	user := http.NewUserController()
	product := http.NewProductController()

	router := gin.Default()
	// var data = interfaces{"user/list": user.Index, "user/get", "user/post"}
	// routes := map[string]gin.Context{
	// 	"GET /user/list":  user.Index,
	// 	"GET /user/get":   user.Get,
	// 	"POST /user/post": user.Post,
	// }
	userRoutes := router.Group("/api/v1/user")
	productRoutes := router.Group("/api/v1/product")

	userRoutes.GET("/", user.Index)
	userRoutes.GET("/show/:id", user.Index)
	userRoutes.PUT("/", user.Index)
	userRoutes.DELETE("/", user.Index)

	productRoutes.GET("/", product.Index)
	productRoutes.POST("/create", product.Create)

	router.Run(":9090")
	return nil
}
