package utils

import (
	"simple-erp-service/internal/data-structure/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetPaginationParams extrai os parâmetros de paginação da requisição
func GetPaginationParams(c *gin.Context) models.Pagination {
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

	return models.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Order: order,
	}
}

// Paginate aplica a paginação a uma consulta GORM
func Paginate(value interface{}, pagination *models.Pagination, db *gorm.DB) (*gorm.DB, error) {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	pagination.TotalPages = int(totalRows) / pagination.Limit
	if int(totalRows)%pagination.Limit > 0 {
		pagination.TotalPages++
	}

	offset := (pagination.Page - 1) * pagination.Limit
	return db.Offset(offset).Limit(pagination.Limit).Order(pagination.Sort + " " + pagination.Order), nil
}
