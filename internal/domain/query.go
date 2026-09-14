package domain

type ListQuery struct {
	Page       int
	PerPage    int
	Search     string
	SortBy     string
	Order      string
	CategoryID *int32
}

type PageMeta struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	Total      int64  `json:"total"`
	TotalPages int    `json:"total_pages"`
	Search     string `json:"search,omitempty"`
	SortBy     string `json:"sort_by"`
	Order      string `json:"order"`
}

type Page[T any] struct {
	Data []T      `json:"data"`
	Meta PageMeta `json:"meta"`
}
