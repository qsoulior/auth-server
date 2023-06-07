CREATE TABLE IF NOT EXISTS auth.token (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expires_at TIMESTAMP NOT NULL,
    fingerprint CHAR(64) NOT NULL,
    user_id UUID REFERENCES auth.user(id) ON DELETE CASCADE NOT NULL
);