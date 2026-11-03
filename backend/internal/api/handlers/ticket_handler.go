package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tmaster/internal/api/maperrors"
	"tmaster/internal/api/utils/helpers"
	"tmaster/internal/dto"
	"tmaster/internal/validation"
)

// CreateTicketHandler creates a new ticket for the authenticated user.
// @Summary Create ticket
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.TicketDTO true "Ticket payload"
// @Success 201 {object} dto.IDResponse
// @Failure 400 {object} map[string]string "Validation or parse error"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /api/create [post]
func (h *Handler) CreateTicketHandler(c *gin.Context) {
	var ticket dto.TicketDTO

	user_id, user_role, err := helpers.GetAllParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := validation.Validate(ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := h.service.CreateTicketService(c.Request.Context(), ticket, user_id, user_role)
	if err != nil {
		c.JSON(maperrors.StatusFromErr(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	resp := dto.IDResponse{
		ID: id,
	}

	c.JSON(http.StatusCreated, resp)
}

// GetOwnTicketsHandler returns all tickets owned by the authenticated user.
// @Summary List own tickets
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.UserTicketResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /api/tickets [get]
func (h *Handler) GetOwnTicketsHandler(c *gin.Context) {
	user_id, user_role, err := helpers.GetAllParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	pag_data := helpers.GetPagData(c)

	res, err := h.service.GetOwnTicketsService(c.Request.Context(), user_id, user_role, pag_data)

	if err != nil {
		c.JSON(maperrors.StatusFromErr(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetNewTicketsHandler returns the queue of new tickets available for operators.
// @Summary List new tickets queue
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.OperatorTicketResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /api/tickets_queue [get]
func (h *Handler) GetNewTicketsHandler(c *gin.Context) {
	_, user_role, err := helpers.GetAllParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	pag_data := helpers.GetPagData(c)

	res, err := h.service.GetNewTicketsService(c.Request.Context(), user_role, pag_data)
	if err != nil {
		c.JSON(maperrors.StatusFromErr(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

// ClaimNextTicketHandler assigns the next available ticket to the authenticated operator.
// @Summary Claim next ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.OperatorTicketResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /api/claim [post]
func (h *Handler) ClaimNextTicketHandler(c *gin.Context) {
	user_id, user_role, err := helpers.GetAllParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.service.ClaimNextTicketService(c.Request.Context(), user_id, user_role)
	if err != nil {
		c.JSON(maperrors.StatusFromErr(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

// CloseTicketHandler closes a ticket assigned to the authenticated operator.
// @Summary Close ticket
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Success 204 {string} string "Ticket closed"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /api/close [patch]
func (h *Handler) CloseTicketHandler(c *gin.Context) {
	user_id, user_role, err := helpers.GetAllParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.CloseTicketService(c.Request.Context(), user_id, 0, user_role); err != nil {
		c.JSON(maperrors.StatusFromErr(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
