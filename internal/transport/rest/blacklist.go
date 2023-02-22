package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/qnocks/blacklist-user-service/internal/entity"
	"github.com/qnocks/blacklist-user-service/pkg/util"
	"net/http"
	"strconv"
)

type createBlacklistUserInput struct {
	Phone    string `json:"phone" binding:"required"`
	Username string `json:"username" binding:"required"`
	Cause    string `json:"cause" binding:"required"`
	CausedBy string `json:"caused_by" binding:"required"`
}

type userResponse struct {
	ID        int    `json:"id"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
	Cause     string `json:"cause"`
	Timestamp int64  `json:"timestamp"`
	CausedBy  string `json:"caused_by"`
}

// @Summary Save blacklisted user
// @Security ApiKeyAuth
// @Tags blacklist
// @Description Store blacklisted user
// @Accept  json
// @Param input body createBlacklistUserInput true "blacklisted user"
// @Success 201
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/blacklist/ [post]
func (h *Handler) saveBlacklistedUser(c *gin.Context) {
	var requestBody createBlacklistUserInput
	if err := c.BindJSON(&requestBody); err != nil {
		buildErrorResponse(c, http.StatusBadRequest, "error parsing request body")
		return
	}

	if err := h.services.Blacklist.Save(entity.BlacklistedUser{
		Phone:    requestBody.Phone,
		Username: requestBody.Username,
		Cause:    requestBody.Cause,
		CausedBy: requestBody.CausedBy,
	}); err != nil {
		buildErrorResponse(c, http.StatusInternalServerError, "error persisting blacklisted user")
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary Delete blacklisted user
// @Security ApiKeyAuth
// @Tags blacklist
// @Description Delete blacklisted user by provided id
// @Accept  json
// @Param id path string true "id"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/blacklist/{id} [delete]
func (h *Handler) deleteBlacklistedUser(c *gin.Context) {
	idParam := c.Params.ByName("id")
	if len(idParam) == 0 {
		buildErrorResponse(c, http.StatusBadRequest, "missing [id] parameter")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		buildErrorResponse(c, http.StatusBadRequest, "failed parsing [id] parameter")
		return
	}

	if err := h.services.Blacklist.Delete(id); err != nil {
		buildErrorResponse(c, http.StatusInternalServerError, "error deleting blacklisted user")
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get blacklisted users
// @Security ApiKeyAuth
// @Tags blacklist
// @Description Get blacklisted users by provided phone or username
// @Accept json
// @Produce json
// @Param phone query string false "phone to search"
// @Param username query string false "username to search"
// @Success 200 {object} userResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/blacklist/ [get]
func (h *Handler) getBlacklistedUsers(c *gin.Context) {
	phone := c.Query("phone")
	username := c.Query("username")
	if len(phone) == 0 && len(username) == 0 {
		buildErrorResponse(c, http.StatusBadRequest, "missing [phone] and [username] parameters")
		return
	}

	users, err := h.services.Blacklist.Find(phone, username)
	if err != nil {
		buildErrorResponse(c, http.StatusInternalServerError, "error selecting blacklisted users")
		return
	}

	c.JSON(http.StatusOK, util.Map(users, func(entity entity.BlacklistedUser) userResponse {
		return userResponse{
			ID:        entity.ID,
			Phone:     entity.Phone,
			Username:  entity.Username,
			Cause:     entity.Cause,
			Timestamp: entity.Timestamp.Unix(),
			CausedBy:  entity.CausedBy,
		}
	}))
}
