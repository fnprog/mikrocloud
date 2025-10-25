-- +goose Up
CREATE TABLE IF NOT EXISTS cloudflare_tunnels (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	project_id TEXT REFERENCES projects(id) ON DELETE CASCADE,
	organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
	tunnel_token TEXT NOT NULL,
	container_id TEXT,
	target_type TEXT NOT NULL CHECK(target_type IN ('proxy', 'application')),
	target_id TEXT,
	status TEXT NOT NULL DEFAULT 'stopped' CHECK(status IN ('stopped', 'starting', 'running', 'error', 'stopping')),
	last_health_check DATETIME,
	health_status TEXT CHECK(health_status IN ('healthy', 'unhealthy', 'unknown')),
	error_message TEXT,
	config TEXT DEFAULT '{}',
	created_by TEXT NOT NULL REFERENCES users(id),
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(organization_id, name)
);

CREATE INDEX IF NOT EXISTS idx_cloudflare_tunnels_project_id ON cloudflare_tunnels(project_id);
CREATE INDEX IF NOT EXISTS idx_cloudflare_tunnels_organization_id ON cloudflare_tunnels(organization_id);
CREATE INDEX IF NOT EXISTS idx_cloudflare_tunnels_status ON cloudflare_tunnels(status);
CREATE INDEX IF NOT EXISTS idx_cloudflare_tunnels_target ON cloudflare_tunnels(target_type, target_id);

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS update_cloudflare_tunnels_timestamp
	AFTER UPDATE ON cloudflare_tunnels
	FOR EACH ROW
	WHEN NEW.updated_at = OLD.updated_at
BEGIN
	UPDATE cloudflare_tunnels SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS update_cloudflare_tunnels_timestamp;
DROP INDEX IF EXISTS idx_cloudflare_tunnels_target;
DROP INDEX IF EXISTS idx_cloudflare_tunnels_status;
DROP INDEX IF EXISTS idx_cloudflare_tunnels_organization_id;
DROP INDEX IF EXISTS idx_cloudflare_tunnels_project_id;
DROP TABLE IF EXISTS cloudflare_tunnels;
