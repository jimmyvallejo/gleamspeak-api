-- +goose Up
ALTER TABLE servers DROP COLUMN IF EXISTS invite_code;

ALTER TABLE servers ADD COLUMN invite_code TEXT NOT NULL;


-- +goose Down
DROP COLUMN IF EXISTS invite_code;