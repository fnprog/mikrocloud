-- +goose Up
-- +goose StatementBegin
ALTER TABLE git_sources ADD COLUMN github_app_id TEXT;
ALTER TABLE git_sources ADD COLUMN github_installation_id TEXT;
ALTER TABLE git_sources ADD COLUMN github_client_id TEXT;
ALTER TABLE git_sources ADD COLUMN github_client_secret TEXT;
ALTER TABLE git_sources ADD COLUMN github_webhook_secret TEXT;
ALTER TABLE git_sources ADD COLUMN github_private_key TEXT;
ALTER TABLE git_sources ADD COLUMN github_app_slug TEXT;
ALTER TABLE git_sources ADD COLUMN is_github_app INTEGER NOT NULL DEFAULT 0;

CREATE INDEX idx_git_sources_github_app ON git_sources(github_app_id);
CREATE INDEX idx_git_sources_github_installation ON git_sources(github_installation_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_git_sources_github_installation;
DROP INDEX IF EXISTS idx_git_sources_github_app;
ALTER TABLE git_sources DROP COLUMN is_github_app;
ALTER TABLE git_sources DROP COLUMN github_app_slug;
ALTER TABLE git_sources DROP COLUMN github_private_key;
ALTER TABLE git_sources DROP COLUMN github_webhook_secret;
ALTER TABLE git_sources DROP COLUMN github_client_secret;
ALTER TABLE git_sources DROP COLUMN github_client_id;
ALTER TABLE git_sources DROP COLUMN github_installation_id;
ALTER TABLE git_sources DROP COLUMN github_app_id;
-- +goose StatementEnd
