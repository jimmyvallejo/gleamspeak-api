-- +goose Up

CREATE TABLE languages (
    id UUID PRIMARY KEY,
    language TEXT NOT NULL
);


INSERT INTO languages (id, language)
VALUES 
    (gen_random_uuid(), 'english'),
    (gen_random_uuid(), 'spanish'),
    (gen_random_uuid(), 'french');


CREATE TABLE text_channels (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    server_id UUID NOT NULL,
    language_id UUID NOT NULL,
    channel_name TEXT NOT NULL,
    last_active TIMESTAMP,
    is_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(id),
    CONSTRAINT fk_server FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE,
    CONSTRAINT fk_language FOREIGN KEY (language_id) REFERENCES languages(id)
);

-- +goose Down
DROP TABLE IF EXISTS text_channels;
DELETE FROM languages
WHERE language IN ('english', 'french', 'spanish');
DROP TABLE IF EXISTS languages;