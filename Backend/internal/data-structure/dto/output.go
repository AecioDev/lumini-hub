package dto

import "simple-erp-service/internal/data-structure/models"

// ApiPagination representa informações de paginação
type ApiPagination struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Sort       string `json:"sort"`
	Order      string `json:"order"`
	TotalRows  int64  `json:"totalRows"`
	TotalPages int    `json:"totalPages"`
}

// ApiPaginationFromModel converte um objeto de paginação para PaginationDTO
func ApiPaginationFromModel(pagination *models.Pagination) *ApiPagination {
	if pagination == nil {
		return nil
	}

	return &ApiPagination{
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		Sort:       pagination.Sort,
		Order:      pagination.Order,
		TotalRows:  pagination.TotalRows,
		TotalPages: pagination.TotalPages,
	}
}
