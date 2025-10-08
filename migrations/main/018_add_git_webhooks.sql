-- +goose Up
-- +goose StatementBegin
ALTER TABLE git_sources ADD COLUMN webhook_url TEXT;
ALTER TABLE git_sources ADD COLUMN allow_preview_deployments INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE git_sources DROP COLUMN allow_preview_deployments;
ALTER TABLE git_sources DROP COLUMN webhook_url;
-- +goose StatementEnd
