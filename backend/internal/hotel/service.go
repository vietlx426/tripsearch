package hotel

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	db "github.com/vietlx426/tripsearch/db/sqlc"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, hostID uuid.UUID, req CreateRequest) (*HotelDTO, error) {
	hotel, err := s.repo.Insert(ctx, db.InsertHotelParams{
		HostID:        hostID,
		Name:          req.Name,
		Description:   sql.NullString{String: ptrStr(req.Description), Valid: req.Description != nil},
		PropertyType:  req.PropertyType,
		Address:       sql.NullString{String: ptrStr(req.Address), Valid: req.Address != nil},
		City:          req.City,
		Country:       req.Country,
		Latitude:      sql.NullString{String: ptrFloat(req.Latitude), Valid: req.Latitude != nil},
		Longitude:     sql.NullString{String: ptrFloat(req.Longitude), Valid: req.Longitude != nil},
		PricePerNight: req.PricePerNight,
		Currency:      req.Currency,
		StarRating:    sql.NullInt32{Int32: ptrInt32(req.StarRating), Valid: req.StarRating != nil},
		Status:        statusOrDefault(req.Status),
	})
	if err != nil {
		if errors.Is(err, ErrDuplicateHotel) {
			return nil, ErrDuplicateHotel
		}
		return nil, fmt.Errorf("Create: %w", err)
	}

	return toHotelDTO(hotel, nil), nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*HotelDTO, error) {
	hotel, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrHotelNotFound) {
			return nil, ErrHotelNotFound
		}
		return nil, fmt.Errorf("GetByID: $w", err)
	}
	images, err := s.repo.ListImages(ctx, id)
	if err != nil {
    	return nil, fmt.Errorf("ListImages: %w", err)
	}
	return toHotelDTO(hotel, images), nil
}

func (s *service) List(ctx context.Context) ([]HotelDTO, error) {
	hotels, err := s.repo.List(ctx)
	if err != nil {
    	return nil, fmt.Errorf("List: %w", err)
	}

	dtos := make([]HotelDTO, 0, len(hotels))
	for _, h := range hotels {
		dtos = append(dtos, *toHotelDTO(h, nil))
	}
	return dtos, nil
}

func (s *service) Update(ctx context.Context, hostID uuid.UUID, id uuid.UUID, req UpdateRequest) (*HotelDTO, error) {
	hotel, err := s.repo.FindByID(ctx, id);
	if err != nil {
		if errors.Is(err, ErrHotelNotFound) {
			return nil, ErrHotelNotFound
		}
		return nil, fmt.Errorf("Update: %w", err)
	}
	if hotel.HostID != hostID {
		return nil, ErrUnauthorized
	}
	updated, err := s.repo.Update(ctx, db.UpdateHotelParams{
		ID: id,
		Name:          req.Name,
		Description:   sql.NullString{String: ptrStr(req.Description), Valid: req.Description != nil},
		PropertyType:  req.PropertyType,
		Address:       sql.NullString{String: ptrStr(req.Address), Valid: req.Address != nil},
		City:          req.City,
		Country:       req.Country,
		Latitude:      sql.NullString{String: ptrFloat(req.Latitude), Valid: req.Latitude != nil},
		Longitude:     sql.NullString{String: ptrFloat(req.Longitude), Valid: req.Longitude != nil},
		PricePerNight: req.PricePerNight,
		Currency:      req.Currency,
		StarRating:    sql.NullInt32{Int32: ptrInt32(req.StarRating), Valid: req.StarRating != nil},
		Status:        statusOrDefault(req.Status),
	})
	if err != nil {
		return nil, fmt.Errorf("Update: %w", err)
	}
	return toHotelDTO(updated, nil), nil
}
func (s *service) Delete(ctx context.Context, hostID uuid.UUID, id uuid.UUID) error {
	hotel, err := s.repo.FindByID(ctx, id);
	if err != nil {
		if errors.Is(err, ErrHotelNotFound) {
			return ErrHotelNotFound
		}
		return fmt.Errorf("Delete: %w", err)
	}
	if hotel.HostID != hostID {
		return ErrUnauthorized
	}
	return s.repo.Delete(ctx, id)
}

// helpers

func ptrStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ptrFloat(f *float64) string {
	if f == nil {
		return ""
	}
	return strconv.FormatFloat(*f, 'f', 6, 64)
}

func ptrInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

func statusOrDefault(s string) string {
	if s == "" {
		return "active"
	}
	return s
}

func toHotelDTO(h db.Hotel, images []db.HotelImage) *HotelDTO {
	dto := &HotelDTO{
		ID:            h.ID,
		HostID:        h.HostID,
		Name:          h.Name,
		PropertyType:  h.PropertyType,
		City:          h.City,
		Country:       h.Country,
		PricePerNight: h.PricePerNight,
		Currency:      h.Currency,
		Status:        h.Status,
		CreatedAt:     h.CreatedAt,
	}
	if h.Description.Valid {
		dto.Description = &h.Description.String
	}
	if h.Address.Valid {
		dto.Address = &h.Address.String
	}
	if h.StarRating.Valid {
		dto.StarRating = &h.StarRating.Int32
	}
	for _, img := range images {
		imgDTO := ImageDTO{
			ID:       img.ID,
			URL:      img.Url,
			Position: img.Position,
		}
		if img.Alt.Valid {
			imgDTO.Alt = &img.Alt.String
		}
		if img.Description.Valid {
			imgDTO.Description = &img.Description.String
		}
		dto.Images = append(dto.Images, imgDTO)
	}
	return dto
}

