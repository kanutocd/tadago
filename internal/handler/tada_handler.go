package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/kanutocd/tada/internal/dto"
	"github.com/kanutocd/tada/internal/service"
)

type TadaHandler struct {
	tadaService service.TadaService
}

func NewTadaHandler(tadaService service.TadaService) *TadaHandler {
	return &TadaHandler{tadaService: tadaService}
}

// GetTadas godoc
// @Summary Get tadas with pagination
// @Description Retrieve a paginated list of tadas
// @Tags tadas
// @Accept json
// @Produce json
// @Param cursor query string false "Pagination cursor"
// @Param limit query int false "Items per page (1-100)" minimum(1) maximum(100) default(10)
// @Success 200 {object} dto.PaginationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tadas [get]
func (h *TadaHandler) GetTadas(c *gin.Context) {
	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid pagination parameters",
		})
		return
	}

	response, err := h.tadaService.GetTadas(pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateTada godoc
// @Summary Create a new tada
// @Description Create a new tada task
// @Tags tadas
// @Accept json
// @Produce json
// @Param tada body dto.CreateTadaRequest true "Tada creation data"
// @Success 201 {object} dto.TadaResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /tadas [post]
func (h *TadaHandler) CreateTada(c *gin.Context) {
	var req dto.CreateTadaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	tada, err := h.tadaService.CreateTada(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, tada)
}

// GetTada godoc
// @Summary Get tada by ID
// @Description Get tada details by ID
// @Tags tadas
// @Accept json
// @Produce json
// @Param id path string true "Tada ID"
// @Success 200 {object} dto.TadaResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /tadas/{id} [get]
func (h *TadaHandler) GetTada(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid tada ID",
		})
		return
	}

	tada, err := h.tadaService.GetTadaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Tada not found",
		})
		return
	}

	c.JSON(http.StatusOK, tada)
}

// UpdateTada godoc
// @Summary Update tada
// @Description Update tada details
// @Tags tadas
// @Accept json
// @Produce json
// @Param id path string true "Tada ID"
// @Param tada body dto.UpdateTadaRequest true "Tada update data"
// @Success 200 {object} dto.TadaResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /tadas/{id} [put]
func (h *TadaHandler) UpdateTada(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid tada ID",
		})
		return
	}

	var req dto.UpdateTadaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	tada, err := h.tadaService.UpdateTada(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tada)
}

// DeleteTada godoc
// @Summary Delete tada
// @Description Delete tada by ID
// @Tags tadas
// @Accept json
// @Produce json
// @Param id path string true "Tada ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /tadas/{id} [delete]
func (h *TadaHandler) DeleteTada(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid tada ID",
		})
		return
	}

	err = h.tadaService.DeleteTada(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
