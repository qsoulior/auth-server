CREATE TABLE IF NOT EXISTS auth.token (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expires_at TIMESTAMP NOT NULL,
    fingerprint BYTEA NOT NULL,
    is_session BOOLEAN DEFAULT FALSE,
    user_id UUID REFERENCES auth.user(id) ON DELETE CASCADE NOT NULL
);