CREATE TABLE IF NOT EXISTS auth.role (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(20) UNIQUE NOT NULL,
    description VARCHAR(100)
);