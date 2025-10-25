-- +goose Up
ALTER TABLE servers ADD COLUMN ipv6_address TEXT DEFAULT '';

-- +goose Down
ALTER TABLE servers DROP COLUMN ipv6_address;
