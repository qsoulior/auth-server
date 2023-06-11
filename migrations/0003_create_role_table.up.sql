CREATE TABLE IF NOT EXISTS auth.role (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(20) UNIQUE NOT NULL,
    description VARCHAR(100) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS auth.user_role (
    role_id UUID REFERENCES auth.role(id) ON DELETE CASCADE,
    user_id UUID REFERENCES auth.user(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, user_id)
);