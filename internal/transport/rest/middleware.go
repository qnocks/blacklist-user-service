package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const authorizationHeader = "Authorization"

func (h *Handler) authTokenMiddleware(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		buildErrorResponse(c, http.StatusUnauthorized, "empty authorization header")
		return
	}

	splitHeader := strings.Split(header, " ")
	if len(splitHeader) != 2 {
		buildErrorResponse(c, http.StatusUnauthorized, "invalid authorization header")
		return
	}

	_, err := h.services.Auth.ParseToken(splitHeader[1])
	if err != nil {
		buildErrorResponse(c, http.StatusUnauthorized, "error parsing token")
		return
	}
}
