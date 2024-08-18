-- +goose Up
CREATE TABLE text_messages (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    channel_id UUID NOT NULL,
    message TEXT NOT NULL,
    image TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(id),
    CONSTRAINT fk_channel FOREIGN KEY (channel_id) REFERENCES text_channels(id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE IF EXISTS text_messages;