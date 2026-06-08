CREATE TABLE IF NOT EXISTS authors (
    id          VARCHAR(50) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    country     VARCHAR(255) NOT NULL,
    genres      TEXT[] NOT NULL DEFAULT '{}',
    birth_day   VARCHAR(20) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
