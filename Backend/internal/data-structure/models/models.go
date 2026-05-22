package models

// Pagination representa os parâmetros de paginação
type Pagination struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Sort       string `json:"sort"`
	Order      string `json:"order"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
}
