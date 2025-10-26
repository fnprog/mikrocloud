-- +goose Up

INSERT INTO users (
    id, 
    email, 
    password_hash, 
    username, 
    name,
    status, 
    email_verified_at,
    timezone,
    created_at, 
    updated_at
) VALUES (
    'user-system-00000000',
    'system@mikrocloud.local',
    '$2a$10$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx',
    'mikrocloud',
    'System User',
    'active',
    CURRENT_TIMESTAMP,
    'UTC',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

INSERT INTO organizations (
    id,
    name,
    slug,
    description,
    owner_id,
    billing_email,
    plan,
    status,
    created_at,
    updated_at
) VALUES (
    'org-default-000000000',
    'Default Organization',
    'default',
    'Default organization for mikrocloud',
    'user-system-00000000',
    'system@mikrocloud.local',
    'free',
    'active',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

INSERT INTO user_roles (
    id,
    user_id,
    role_id,
    granted_by,
    granted_at
) VALUES (
    'ur-system-admin-00000',
    'user-system-00000000',
    'role-admin-00000000',
    'user-system-00000000',
    CURRENT_TIMESTAMP
);

INSERT INTO organization_members (
    id,
    organization_id,
    user_id,
    role,
    invited_by,
    invited_at,
    joined_at,
    status
) VALUES (
    'om-system-default-0000',
    'org-default-000000000',
    'user-system-00000000',
    'owner',
    'user-system-00000000',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    'active'
);

-- +goose Down

DELETE FROM organization_members WHERE id = 'om-system-default-0000';
DELETE FROM user_roles WHERE id = 'ur-system-admin-00000';
DELETE FROM organizations WHERE id = 'org-default-000000000';
DELETE FROM users WHERE id = 'user-system-00000000';
