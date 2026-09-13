package router

import (
	"restapirian/internal/delivery/http"

	"github.com/gin-gonic/gin"
)

func InitializeRouter() error {
	user := http.NewUserController()
	router := gin.Default()
	// var data = interfaces{"user/list": user.Index, "user/get", "user/post"}
	// routes := map[string]gin.Context{
	// 	"GET /user/list":  user.Index,
	// 	"GET /user/get":   user.Get,
	// 	"POST /user/post": user.Post,
	// }
	v1 := router.Group("/api/v1")
	v1.GET("/", user.Index)
	v1.GET("/show/:id", user.Index)
	v1.PUT("/", user.Index)
	v1.DELETE("/", user.Index)

	router.Run(":9090")
	return nil
}
