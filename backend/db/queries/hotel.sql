-- name: InsertHotel :one
INSERT INTO hotels (host_id, name, description, property_type, address, city, country, latitude, longitude, price_per_night, currency, star_rating, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: FindHotelByID :one
SELECT * FROM hotels
WHERE id = $1
LIMIT 1;

-- name: ListHotels :many
SELECT * FROM hotels
WHERE status = 'active'
ORDER BY created_at DESC;

-- name: UpdateHotel :one
UPDATE hotels
SET
    name            = $2,
    description     = $3,
    property_type   = $4,
    address         = $5,
    city            = $6,
    country         = $7,
    latitude        = $8,
    longitude       = $9,
    price_per_night = $10,
    currency        = $11,
    star_rating     = $12,
    status          = $13
WHERE id = $1
RETURNING *;

-- name: DeleteHotel :exec
DELETE FROM hotels
WHERE id = $1;