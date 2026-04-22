-- name: InsertHotelImage :one
INSERT INTO hotel_images (hotel_id, url, alt, description, position)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteHotelImage :exec
DELETE FROM hotel_images
WHERE id = $1;

-- name: ListImagesByHotelID :many
SELECT * FROM hotel_images
WHERE hotel_id = $1
ORDER BY position DESC;