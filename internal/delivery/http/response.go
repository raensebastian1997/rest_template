package http

import (
	"errors"
	stdhttp "net/http"
	"strconv"
	"strings"

	"restapirian/internal/domain"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, err error) {
	status := stdhttp.StatusInternalServerError
	message := "internal server error"
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		status, message = stdhttp.StatusBadRequest, err.Error()
	case errors.Is(err, domain.ErrUnauthorized):
		status, message = stdhttp.StatusUnauthorized, domain.ErrUnauthorized.Error()
	case errors.Is(err, domain.ErrForbidden):
		status, message = stdhttp.StatusForbidden, domain.ErrForbidden.Error()
	case errors.Is(err, domain.ErrNotFound):
		status, message = stdhttp.StatusNotFound, err.Error()
	case errors.Is(err, domain.ErrConflict):
		status, message = stdhttp.StatusConflict, domain.ErrConflict.Error()
	}
	c.JSON(status, gin.H{"error": message})
}

func parseListQuery(c *gin.Context) (domain.ListQuery, error) {
	query := domain.ListQuery{
		Page: 1, PerPage: 10, Search: c.Query("search"),
		SortBy: c.Query("sort_by"), Order: c.Query("order"),
	}
	var err error
	if value := c.Query("page"); value != "" {
		query.Page, err = strconv.Atoi(value)
		if err != nil || query.Page < 1 {
			return query, domain.ErrInvalidInput
		}
	}
	if value := c.Query("per_page"); value != "" {
		query.PerPage, err = strconv.Atoi(value)
		if err != nil || query.PerPage < 1 {
			return query, domain.ErrInvalidInput
		}
	}
	if value := c.Query("category_id"); value != "" {
		parsed, parseErr := strconv.ParseInt(value, 10, 32)
		if parseErr != nil || parsed < 1 {
			return query, domain.ErrInvalidInput
		}
		categoryID := int32(parsed)
		query.CategoryID = &categoryID
	}
	query.Order = strings.ToLower(query.Order)
	if query.Order != "" && query.Order != "asc" && query.Order != "desc" {
		return query, domain.ErrInvalidInput
	}
	return query, nil
}

func parseInt32ID(c *gin.Context) (int32, error) {
	value, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil || value < 1 {
		return 0, domain.ErrInvalidInput
	}
	return int32(value), nil
}

func parseInt64ID(c *gin.Context) (int64, error) {
	value, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || value < 1 {
		return 0, domain.ErrInvalidInput
	}
	return value, nil
}

func currentUserID(c *gin.Context) int64 {
	value, exists := c.Get("auth.user_id")
	if !exists {
		return 0
	}
	userID, _ := value.(int64)
	return userID
}
