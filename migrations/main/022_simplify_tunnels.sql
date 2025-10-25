-- +goose Up
DROP INDEX IF EXISTS idx_cloudflare_tunnels_target;
ALTER TABLE cloudflare_tunnels DROP COLUMN target_type;
ALTER TABLE cloudflare_tunnels DROP COLUMN target_id;

-- +goose Down
ALTER TABLE cloudflare_tunnels ADD COLUMN target_type TEXT NOT NULL DEFAULT 'proxy' CHECK(target_type IN ('proxy', 'application'));
ALTER TABLE cloudflare_tunnels ADD COLUMN target_id TEXT;
CREATE INDEX IF NOT EXISTS idx_cloudflare_tunnels_target ON cloudflare_tunnels(target_type, target_id);
