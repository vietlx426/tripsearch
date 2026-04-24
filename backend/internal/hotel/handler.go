package hotel

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vietlx426/tripsearch/pkg/middleware"
)


type handler struct {
	svc Service
}


func NewHandler(svc Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) Create(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hotel, err := h.svc.Create(c.Request.Context(), claims.UserID, req)
	if err != nil {
		if errors.Is(err, ErrDuplicateHotel) {
			c.JSON(http.StatusConflict, gin.H{"error": "hotel already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // TODO: remove before prod
		return
	}

	c.JSON(http.StatusCreated, hotel)
}

func (h *handler) List(c *gin.Context) {
	hotels, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, hotels)
}

func (h *handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hotel id"})
		return
	}

	hotel, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrHotelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "hotel not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, hotel)
}

func (h *handler) Update(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hotel id"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hotel, err := h.svc.Update(c.Request.Context(), claims.UserID, id, req)
	if err != nil {
		if errors.Is(err, ErrHotelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "hotel not found"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, hotel)
}

func (h *handler) Delete(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hotel id"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), claims.UserID, id); err != nil {
		if errors.Is(err, ErrHotelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "hotel not found"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

