package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tmaster/internal/api/maperrors"
	"tmaster/internal/api/utils/helpers"
	"tmaster/internal/dto"
	"tmaster/internal/validation"
)

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
			"error" : err.Error(),
		})
		return
	}
	resp := dto.IDResponse{
		ID: id,
	}

	c.JSON(201, resp)
}

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
			"error" : err.Error(),
		})
		return
	}

	c.JSON(200, res)
}

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
			"error" : err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

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
			"error" : err.Error(),
		})
		return
	}
	c.JSON(200, res)
}

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
			"error" : err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
