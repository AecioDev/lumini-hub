package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Pagination representa os parâmetros de paginação
type Pagination struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Sort       string `json:"sort"`
	Order      string `json:"order"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
}

// GetPaginationParams extrai os parâmetros de paginação da requisição
func GetPaginationParams(c *gin.Context) Pagination {
	// Valores padrão
	page := 1
	limit := 10
	sort := "created_at"
	order := "desc"

	// Extrair parâmetros da query
	if pageParam := c.Query("page"); pageParam != "" {
		pageInt, err := strconv.Atoi(pageParam)
		if err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	if limitParam := c.Query("limit"); limitParam != "" {
		limitInt, err := strconv.Atoi(limitParam)
		if err == nil && limitInt > 0 {
			limit = limitInt
		}
	}

	if sortParam := c.Query("sort"); sortParam != "" {
		sort = sortParam
	}

	if orderParam := c.Query("order"); orderParam != "" {
		if orderParam == "asc" || orderParam == "desc" {
			order = orderParam
		}
	}

	return Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Order: order,
	}
}

// Paginate aplica a paginação a uma consulta GORM. A ordenação passa por safeOrderBy: só colunas
// reais do modelo, nunca texto da requisição direto no SQL (e pagination.Sort/Order são
// corrigidos para o que foi de fato usado).
func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) (*gorm.DB, error) {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	pagination.TotalPages = int(totalRows) / pagination.Limit
	if int(totalRows)%pagination.Limit > 0 {
		pagination.TotalPages++
	}

	offset := (pagination.Page - 1) * pagination.Limit
	paginated := db.Offset(offset).Limit(pagination.Limit)
	if orderBy, ok := safeOrderBy(value, db, pagination.Sort, pagination.Order); ok {
		paginated = paginated.Order(orderBy)
		pagination.Sort = orderBy.Column.Name
		pagination.Order = "asc"
		if orderBy.Desc {
			pagination.Order = "desc"
		}
	}
	return paginated, nil
}

// PaginationFromFilter monta a paginação de uma busca POST .../filter (page_no, page_size,
// order_by_column, is_asc). Sem página vira 1, sem tamanho 10 e sem coluna `defaultSort` em ordem
// decrescente; a coluna enviada só é aceita por Paginate se for uma coluna real do modelo.
func PaginationFromFilter(pageNo, pageSize *int, orderByColumn string, isAsc *bool, defaultSort string) Pagination {
	pagination := Pagination{Page: 1, Limit: 10, Sort: defaultSort, Order: "desc"}
	if pageNo != nil && *pageNo > 0 {
		pagination.Page = *pageNo
	}
	if pageSize != nil && *pageSize > 0 {
		pagination.Limit = *pageSize
	}
	if orderByColumn != "" {
		pagination.Sort = orderByColumn
		pagination.Order = "asc"
		if isAsc != nil && !*isAsc {
			pagination.Order = "desc"
		}
	}
	return pagination
}

// ApiPagination representa informações de paginação para a API (usado em DTOs)
type ApiPagination struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Sort       string `json:"sort"`
	Order      string `json:"order"`
	TotalRows  int64  `json:"totalRows"`
	TotalPages int    `json:"totalPages"`
}

// ApiPaginationFromModel converte um objeto de paginação interno para ApiPagination
func ApiPaginationFromModel(pagination *Pagination) *ApiPagination {
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

