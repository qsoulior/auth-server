CREATE TABLE IF NOT EXISTS auth.user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(20) UNIQUE NOT NULL,
    password BYTEA NOT NULL
);