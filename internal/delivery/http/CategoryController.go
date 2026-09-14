package http

import (
	stdhttp "net/http"

	"restapirian/internal/domain"
	"restapirian/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryController struct{ service service.CategoryService }

func NewCategoryController(categoryService service.CategoryService) *CategoryController {
	return &CategoryController{service: categoryService}
}

func (ctrl *CategoryController) Index(c *gin.Context) {
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

func (ctrl *CategoryController) Show(c *gin.Context) {
	id, err := parseInt32ID(c)
	if err != nil {
		writeError(c, err)
		return
	}
	category, err := ctrl.service.Get(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, gin.H{"data": category})
}

func (ctrl *CategoryController) Create(c *gin.Context) {
	var input service.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, domain.ErrInvalidInput)
		return
	}
	category, err := ctrl.service.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(stdhttp.StatusCreated, gin.H{"data": category})
}

func (ctrl *CategoryController) Update(c *gin.Context) {
	id, err := parseInt32ID(c)
	if err != nil {
		writeError(c, err)
		return
	}
	var input service.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, domain.ErrInvalidInput)
		return
	}
	category, err := ctrl.service.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, gin.H{"data": category})
}

func (ctrl *CategoryController) Delete(c *gin.Context) {
	id, err := parseInt32ID(c)
	if err != nil {
		writeError(c, err)
		return
	}
	if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	c.Status(stdhttp.StatusNoContent)
}
