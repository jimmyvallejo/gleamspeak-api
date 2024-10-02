-- +goose Up
CREATE TABLE voice_channels (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    server_id UUID NOT NULL,
    language_id UUID NOT NULL,
    channel_name TEXT NOT NULL,
    channel_id TEXT NOT NULL,
    last_active TIMESTAMP,
    is_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(id),
    CONSTRAINT fk_server FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE,
    CONSTRAINT fk_language FOREIGN KEY (language_id) REFERENCES languages(id)
);
CREATE TABLE voice_channel_members (
    user_id UUID NOT NULL,
    channel_id UUID NOT NULL,
    server_id UUID NOT NULL,
    PRIMARY KEY (user_id, channel_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (channel_id) REFERENCES voice_channels(id) ON DELETE CASCADE,
    FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE IF EXISTS voice_channel_members;
DROP TABLE IF EXISTS voice_channels;