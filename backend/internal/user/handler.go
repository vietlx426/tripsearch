package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vietlx426/tripsearch/pkg/middleware"
)

type handler struct {
	svc Service
}

func NewHandler(svc Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) GetMe(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.svc.GetMe(c.Request.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
