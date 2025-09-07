-- +goose Up
-- Migration for new architecture: Projects > Environments > Services
-- This replaces the applications table with a proper hierarchy

-- Drop old tables if they exist (for fresh start)
DROP TABLE IF EXISTS applications;
DROP TABLE IF EXISTS deployments;

-- Environments table
CREATE TABLE IF NOT EXISTS environments (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    project_id TEXT NOT NULL,
    description TEXT,
    variables TEXT, -- JSON encoded map[string]string
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    UNIQUE(project_id, name) -- Environment names must be unique within a project
);

-- Services table (replaces applications)
CREATE TABLE IF NOT EXISTS services (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    project_id TEXT NOT NULL,
    environment_id TEXT NOT NULL,
    
    -- Git configuration
    git_url TEXT NOT NULL,
    git_branch TEXT NOT NULL DEFAULT 'main',
    context_root TEXT,
    
    -- Build configuration
    buildpack_type TEXT NOT NULL CHECK(buildpack_type IN ('nixpacks', 'static', 'dockerfile', 'docker-compose')),
    build_config TEXT, -- JSON encoded build configuration
    
    -- Runtime configuration
    environment_vars TEXT, -- JSON encoded map[string]string
    
    -- Status
    status TEXT NOT NULL DEFAULT 'created' CHECK(status IN ('created', 'building', 'deploying', 'running', 'stopped', 'failed')),
    
    -- Timestamps
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE,
    UNIQUE(project_id, environment_id, name) -- Service names must be unique within an environment
);

-- Deployments table (tracks deployment history)
CREATE TABLE IF NOT EXISTS deployments (
    id TEXT PRIMARY KEY,
    service_id TEXT NOT NULL,
    
    -- Build information
    image_tag TEXT NOT NULL,
    build_logs TEXT,
    
    -- Deployment information
    status TEXT NOT NULL CHECK(status IN ('pending', 'building', 'deploying', 'running', 'failed', 'stopped')),
    container_id TEXT,
    
    -- Timestamps
    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_environments_project_id ON environments(project_id);
CREATE INDEX IF NOT EXISTS idx_services_project_id ON services(project_id);
CREATE INDEX IF NOT EXISTS idx_services_environment_id ON services(environment_id);
CREATE INDEX IF NOT EXISTS idx_services_status ON services(status);
CREATE INDEX IF NOT EXISTS idx_deployments_service_id ON deployments(service_id);
CREATE INDEX IF NOT EXISTS idx_deployments_status ON deployments(status);

-- Insert default production environment for existing projects
INSERT OR IGNORE INTO environments (id, name, project_id, description, variables)
SELECT 
    lower(hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-' || hex(randomblob(2)) || '-' || hex(randomblob(2)) || '-' || hex(randomblob(6))) as id,
    'prod' as name,
    id as project_id,
    'Production environment' as description,
    '{}' as variables
FROM projects;

-- +goose Down
-- Rollback: recreate applications table and remove new tables
DROP TABLE IF EXISTS deployments;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS environments;

-- Recreate the old applications table
CREATE TABLE IF NOT EXISTS applications (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    git_url TEXT NOT NULL,
    git_branch TEXT NOT NULL DEFAULT 'main',
    build_type TEXT NOT NULL DEFAULT 'nixpacks',
    port INTEGER NOT NULL DEFAULT 3000,
    environment_vars TEXT,
    status TEXT NOT NULL DEFAULT 'created',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
