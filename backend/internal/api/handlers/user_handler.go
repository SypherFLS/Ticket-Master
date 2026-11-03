package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tmaster/internal/api/maperrors"
	"tmaster/internal/api/utils/helpers"

	"tmaster/internal/dto"
	"tmaster/internal/validation"
)

// LoginHandler handles user authentication.
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginDTO true "User credentials"
// @Success 200 {string} string "JWT token"
// @Failure 400 {object} map[string]string "Validation or parse error"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /auth/login [post]
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
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, token)
}

// RegisterHandler registers a new user.
// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterDTO true "User registration data"
// @Success 201 {string} string "User created"
// @Failure 400 {object} map[string]string "Validation or parse error"
// @Failure 409 {object} map[string]string "User already exists"
// @Router /auth/register [post]
func (h *Handler) RegisterHandler(c *gin.Context) {
	var user dto.RegisterDTO

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(maperrors.StatusFromErr(err), gin.H{
			"error": err.Error(),
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
			"error": err.Error(),
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
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusAccepted)
}
