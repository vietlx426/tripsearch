package hotel

import (
	"context"
	"errors"
	"time"
	"github.com/google/uuid"
	"github.com/vietlx426/tripsearch/db/sqlc"
)

type HotelDTO struct {
    ID             uuid.UUID   `json:"id"`
    HostID         uuid.UUID   `json:"host_id"`
    Name           string      `json:"name"`
    Description    *string     `json:"description"`
    PropertyType   string      `json:"property_type"`
    Address        *string     `json:"address"`
    City           string      `json:"city"`
    Country        string      `json:"country"`
    Latitude       *float64    `json:"latitude"`
    Longitude      *float64    `json:"longitude"`
    PricePerNight  string      `json:"price_per_night"`
    Currency       string      `json:"currency"`
    StarRating     *int32      `json:"star_rating"`
    Status         string      `json:"status"`
    Images         []ImageDTO  `json:"images"`
    CreatedAt      time.Time   `json:"created_at"`
}

type ImageDTO struct {
    ID          uuid.UUID `json:"id"`
    URL         string    `json:"url"`
    Alt         *string   `json:"alt"`
    Description *string   `json:"description"`
    Position    int32     `json:"position"`
}

type Repository interface {
    Insert(ctx context.Context, arg db.InsertHotelParams) (db.Hotel, error)
    FindByID(ctx context.Context, id uuid.UUID) (db.Hotel, error)
    List(ctx context.Context) ([]db.Hotel, error)
    Update(ctx context.Context, arg db.UpdateHotelParams) (db.Hotel, error)
    Delete(ctx context.Context, id uuid.UUID) error

    InsertImage(ctx context.Context, arg db.InsertHotelImageParams) (db.HotelImage, error)
    ListImages(ctx context.Context, hotelID uuid.UUID) ([]db.HotelImage, error)
    DeleteImage(ctx context.Context, id uuid.UUID) error
}

type Service interface {
    Create(ctx context.Context, hostID uuid.UUID, req CreateRequest) (*HotelDTO, error)
    GetByID(ctx context.Context, id uuid.UUID) (*HotelDTO, error)
    List(ctx context.Context) ([]HotelDTO, error)
    Update(ctx context.Context, hostID uuid.UUID, id uuid.UUID, req UpdateRequest) (*HotelDTO, error)
    Delete(ctx context.Context, hostID uuid.UUID, id uuid.UUID) error
}


var (
	ErrHotelNotFound = errors.New("user not found")
	ErrDuplicateHotel = errors.New("duplicate hotel")
)

type CreateRequest struct {
    Name           string   `json:"name"          binding:"required"`
    Description    *string  `json:"description"`
    PropertyType   string   `json:"property_type" binding:"required"`
    Address        *string  `json:"address"`
    City           string   `json:"city"          binding:"required"`
    Country        string   `json:"country"       binding:"required"`
    Latitude       *float64 `json:"latitude"`
    Longitude      *float64 `json:"longitude"`
    PricePerNight  string   `json:"price_per_night" binding:"required"`
    Currency       string   `json:"currency"       binding:"required,len=3"`
    StarRating     *int32   `json:"star_rating"`
    Status         string   `json:"status"`
}

type UpdateRequest struct {
    Name           string   `json:"name"          binding:"required"`
    Description    *string  `json:"description"`
    PropertyType   string   `json:"property_type" binding:"required"`
    Address        *string  `json:"address"`
    City           string   `json:"city"          binding:"required"`
    Country        string   `json:"country"       binding:"required"`
    Latitude       *float64 `json:"latitude"`
    Longitude      *float64 `json:"longitude"`
    PricePerNight  string   `json:"price_per_night" binding:"required"`
    Currency       string   `json:"currency"       binding:"required,len=3"`
    StarRating     *int32   `json:"star_rating"`
    Status         string   `json:"status"`
}
