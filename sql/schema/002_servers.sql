-- +goose Up
CREATE TABLE servers (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    server_name TEXT NOT NULL,
    description TEXT,
    icon_url TEXT,
    banner_url TEXT,
    is_public BOOLEAN DEFAULT TRUE,
    member_count INTEGER DEFAULT 1,
    server_level INTEGER DEFAULT 0,
    max_members INTEGER DEFAULT 50,
    invite_code TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE user_servers (
    user_id UUID NOT NULL,
    server_id UUID NOT NULL,
    PRIMARY KEY (user_id, server_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE user_servers;
DROP TABLE servers;
