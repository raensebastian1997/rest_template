package service

import (
	"strings"

	"restapirian/internal/domain"
)

func normalizeListQuery(query domain.ListQuery, defaultSort string) domain.ListQuery {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 {
		query.PerPage = 10
	}
	if query.PerPage > 100 {
		query.PerPage = 100
	}
	query.Search = strings.TrimSpace(query.Search)
	query.SortBy = strings.ToLower(strings.TrimSpace(query.SortBy))
	if query.SortBy == "" {
		query.SortBy = defaultSort
	}
	query.Order = strings.ToLower(strings.TrimSpace(query.Order))
	if query.Order != "desc" {
		query.Order = "asc"
	}
	return query
}

func pageMeta(query domain.ListQuery, total int64) domain.PageMeta {
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(query.PerPage) - 1) / int64(query.PerPage))
	}
	return domain.PageMeta{
		Page: query.Page, PerPage: query.PerPage, Total: total, TotalPages: totalPages,
		Search: query.Search, SortBy: query.SortBy, Order: query.Order,
	}
}
