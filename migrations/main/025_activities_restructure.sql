-- +goose Up
-- Rename activity_type to event_type
ALTER TABLE activities RENAME COLUMN activity_type TO event_type;

-- Add level column with default 'info'
ALTER TABLE activities ADD COLUMN level TEXT NOT NULL DEFAULT 'info';

-- Remove initiator_name column (data will come from user join)
ALTER TABLE activities DROP COLUMN initiator_name;

-- Update indexes
DROP INDEX IF EXISTS idx_activities_activity_type;
CREATE INDEX idx_activities_event_type ON activities(event_type);
CREATE INDEX idx_activities_level ON activities(level);

-- +goose Down
-- Reverse the changes
ALTER TABLE activities ADD COLUMN initiator_name TEXT NOT NULL DEFAULT '';
ALTER TABLE activities DROP COLUMN level;
ALTER TABLE activities RENAME COLUMN event_type TO activity_type;

DROP INDEX IF EXISTS idx_activities_level;
DROP INDEX IF EXISTS idx_activities_event_type;
CREATE INDEX idx_activities_activity_type ON activities(activity_type);