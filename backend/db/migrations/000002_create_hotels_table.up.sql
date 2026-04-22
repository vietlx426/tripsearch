CREATE TABLE hotels (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    host_id         UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    property_type   VARCHAR(50) NOT NULL DEFAULT 'hotel',
    address         TEXT,
    city            VARCHAR(100) NOT NULL,
    country         VARCHAR(100) NOT NULL,
    latitude        NUMERIC(9,6),
    longitude       NUMERIC(9,6),
    price_per_night NUMERIC(10,2) NOT NULL,
    currency        CHAR(3) NOT NULL DEFAULT 'USD' CHECK (currency = UPPER(currency)),
    star_rating     INTEGER CHECK (star_rating BETWEEN 1 AND 5),
    status          VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'pending')),
    search_vector   TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('english', name || ' ' || COALESCE(description, '') || ' ' || city || ' ' || country)
    ) STORED,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_hotels_host_id ON hotels(host_id);
CREATE INDEX idx_hotels_status ON hotels(status);
CREATE INDEX idx_hotels_search_vector ON hotels USING GIN(search_vector);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER hotels_set_updated_at
    BEFORE UPDATE ON hotels
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
