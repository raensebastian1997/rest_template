package http

import (
	stdhttp "net/http"

	"restapirian/internal/domain"
	"restapirian/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct{ service service.AuthService }

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{service: authService}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, domain.ErrInvalidInput)
		return
	}
	result, err := ctrl.service.Register(c.Request.Context(), input)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(stdhttp.StatusCreated, gin.H{"data": result})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, domain.ErrInvalidInput)
		return
	}
	result, err := ctrl.service.Login(c.Request.Context(), input)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, gin.H{"data": result})
}
