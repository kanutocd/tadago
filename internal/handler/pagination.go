package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kanutocd/tada/internal/dto"
)

func GetPaginationFromContext(c *gin.Context) dto.PaginationQuery {
	pagination := dto.PaginationQuery{
		Cursor: c.Query("cursor"),
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			pagination.Limit = l
		}
	}

	// Set default limit if not provided
	if pagination.Limit == 0 {
		pagination.Limit = 10
	}

	// Ensure limit is within bounds
	if pagination.Limit > 100 {
		pagination.Limit = 100
	} else if pagination.Limit < 1 {
		pagination.Limit = 1
	}

	return pagination
}
