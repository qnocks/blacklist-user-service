package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/qnocks/blacklist-user-service/internal/entity"
	"net/http"
)

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

// @Summary Login
// @Tags auth
// @Description Login with provided credentials
// @Accept json
// @Produce json
// @Param input body signInInput true "user credentials"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /auth/login/ [post]
func (h *Handler) login(c *gin.Context) {
	var requestBody signInInput

	if err := c.BindJSON(&requestBody); err != nil {
		buildErrorResponse(c, http.StatusBadRequest, "error parsing request body")
		return
	}

	token, err := h.services.Auth.Login(entity.User{
		Username: requestBody.Username,
		Password: requestBody.Password,
	})
	if err != nil {
		buildErrorResponse(c, http.StatusBadRequest, "error verifying credentials")
		return
	}

	c.JSON(http.StatusOK, tokenResponse{Token: token})
}
