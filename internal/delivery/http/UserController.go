package http

import (
	stdhttp "net/http"

	"restapirian/internal/domain"
	"restapirian/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct{ service service.UserService }

func NewUserController(userService service.UserService) *UserController {
	return &UserController{service: userService}
}

func (ctrl *UserController) Index(c *gin.Context) {
	if ctrl.service == nil {
		c.JSON(stdhttp.StatusInternalServerError, gin.H{"error": "database unavailable"})
		return
	}
	query, err := parseListQuery(c)
	if err != nil {
		writeError(c, err)
		return
	}
	page, err := ctrl.service.List(c.Request.Context(), query)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, page)
}

func (ctrl *UserController) Show(c *gin.Context) {
	id, err := parseInt64ID(c)
	if err != nil {
		writeError(c, err)
		return
	}
	user, err := ctrl.service.Get(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, gin.H{"data": user})
}

func (ctrl *UserController) Create(c *gin.Context) {
	var input service.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, domain.ErrInvalidInput)
		return
	}
	user, err := ctrl.service.Create(c.Request.Context(), input, currentUserID(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(stdhttp.StatusCreated, gin.H{"data": user})
}

func (ctrl *UserController) Update(c *gin.Context) {
	id, err := parseInt64ID(c)
	if err != nil {
		writeError(c, err)
		return
	}
	var input service.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, domain.ErrInvalidInput)
		return
	}
	user, err := ctrl.service.Update(c.Request.Context(), id, input, currentUserID(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, gin.H{"data": user})
}

func (ctrl *UserController) Delete(c *gin.Context) {
	id, err := parseInt64ID(c)
	if err != nil {
		writeError(c, err)
		return
	}
	if err := ctrl.service.Delete(c.Request.Context(), id, currentUserID(c)); err != nil {
		writeError(c, err)
		return
	}
	c.Status(stdhttp.StatusNoContent)
}
