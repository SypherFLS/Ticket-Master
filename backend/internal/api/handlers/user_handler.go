package handlers

import (
	"net/http"
	"tmaster/internal/api/maperrors"
	"github.com/gin-gonic/gin"
	"tmaster/internal/api/utils/helpers"

	"tmaster/internal/dto"
	"tmaster/internal/validation"
)

func (h *Handler) LoginHandler(c *gin.Context) {
	var logData dto.LoginDTO

	if err := c.ShouldBindJSON(&logData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := validation.Validate(logData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.service.LoginService(c.Request.Context(), logData)

	if err != nil {
		c.JSON(maperrors.StatusFromErr(err), gin.H{
			"error" : err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, token)
}

func (h *Handler) RegisterHandler(c *gin.Context) {
	var user dto.RegisterDTO

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(maperrors.StatusFromErr(err), gin.H{
			"error" : err.Error(),
		})
		return
	}

	if err := validation.Validate(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.RegisterService(c.Request.Context(), user)
	if err != nil {
		c.JSON(maperrors.StatusFromErr(err), gin.H{
			"error" : err.Error(),
		})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) SetRoleHandler(c *gin.Context) {
	var ru dto.RequestUser

	_, user_role, err := helpers.GetAllParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := c.ShouldBindJSON(&ru); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.SetRoleService(c.Request.Context(), user_role, ru.Email, ru.Role); err != nil {
		c.JSON(maperrors.StatusFromErr(err), gin.H{
			"error" : err.Error(),
		})
		return
	}

	c.Status(http.StatusAccepted)
}
