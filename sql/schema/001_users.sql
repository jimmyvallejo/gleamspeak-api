-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    handle TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    bio TEXT,
    avatar_url TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_handle ON users(handle);
CREATE TABLE roles (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);
-- +goose Down
DROP TABLE users;
DROP TABLE roles;
DROP TABLE user_roles;