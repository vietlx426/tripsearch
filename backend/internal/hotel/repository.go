package hotel
import (
	"context"
	"fmt"
	"errors"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/vietlx426/tripsearch/db/sqlc"
)
type repository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) Repository {
	return &repository{q: q}
}

func (r *repository) Insert(ctx context.Context, args db.InsertHotelParams) (db.Hotel, error) {
	hotel, err := r.q.InsertHotel(ctx, args)
	if err != nil {
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) && pgErr.Code == "23505" {
        return db.Hotel{}, ErrDuplicateHotel
    }
    return db.Hotel{}, fmt.Errorf("Insert: %w", err)
}

	return hotel, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (db.Hotel, error) {
	hotel, err := r.q.FindHotelByID(ctx, id)
	if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return db.Hotel{}, ErrHotelNotFound
    }
    	return db.Hotel{}, fmt.Errorf("FindByID: %w", err)
	}

	return hotel, nil
}

func (r *repository) List(ctx context.Context) ([]db.Hotel, error) {
	list, err := r.q.ListHotels(ctx)
	if err != nil {
    	return nil, fmt.Errorf("List: %w", err)
	}

	return list, nil
}

func (r *repository) Update(ctx context.Context, arg db.UpdateHotelParams) (db.Hotel, error) {
	hotel, err := r.q.UpdateHotel(ctx, arg)
	if err != nil {
    	return db.Hotel{}, fmt.Errorf("Update: %w", err)
	}
	return hotel, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteHotel(ctx, id)
}

func (r *repository) InsertImage(ctx context.Context, arg db.InsertHotelImageParams) (db.HotelImage, error) {
	hotelimage, err := r.q.InsertHotelImage(ctx, arg)
	if err != nil {
    	return db.HotelImage{}, fmt.Errorf("InsertImage: %w", err)
	}
	return hotelimage, nil
}

func (r *repository) ListImages(ctx context.Context, hotelID uuid.UUID) ([]db.HotelImage, error) {
	list, err := r.q.ListImagesByHotelID(ctx, hotelID)
	if err != nil {
    	return nil, fmt.Errorf("ListImages: %w", err)
	}
	return list, nil
}

func (r *repository) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteHotelImage(ctx, id)
}
