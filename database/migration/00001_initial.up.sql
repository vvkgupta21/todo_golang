CREATE TABLE IF NOT EXISTS users
(
    id          SERIAL PRIMARY KEY,
    name        TEXT                     NOT NULL,
    email       TEXT                     NOT NULL,
    password    TEXT                     NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);