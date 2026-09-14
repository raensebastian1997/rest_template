package http

import (
	stdhttp "net/http"

	"restapirian/internal/service"

	"github.com/gin-gonic/gin"
)

type RoleController struct{ service service.RoleService }

func NewRoleController(roleService service.RoleService) *RoleController {
	return &RoleController{service: roleService}
}

func (ctrl *RoleController) Index(c *gin.Context) {
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
