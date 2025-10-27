-- +goose Up

-- Default organization and system user will be created in application code for security

-- +goose Down

DELETE FROM organization_members WHERE id = 'om-system-default-0000';
DELETE FROM user_roles WHERE id = 'ur-system-admin-00000';
DELETE FROM organizations WHERE id = 'org-default-000000000';
DELETE FROM users WHERE id = 'user-system-00000000';
