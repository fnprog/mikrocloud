# MikroCloud API v1 - New Architecture

This document describes the new API architecture based on the hierarchy: **Projects → Environments → Services**

## Architecture Overview

```
Project
├── Environment (prod, dev, staging, etc.)
│   ├── Service 1 (e.g., frontend)
│   ├── Service 2 (e.g., backend)
│   └── Service 3 (e.g., database)
└── Environment (another env)
    ├── Service 1
    └── Service 2
```

## Quick Start: One-off Service Deployment

For simple use cases, you can create a service with auto-generated project and environment:

```bash
POST /api/v1/services/quick
{
  "name": "my-app",
  "git_url": "https://github.com/user/my-app.git",
  "git_branch": "main",
  "buildpack_type": "nixpacks",
  "environment": {
    "NODE_ENV": "production"
  }
}
```

This will:

1. Create a project named "my-app"
2. Create a "prod" environment in that project
3. Create the service in the prod environment
4. Return the service with project and environment details

## Full Workflow: Project → Environment → Service

### 1. Create a Project

```bash
POST /api/v1/projects
{
  "name": "my-company-website",
  "description": "Company website and blog"
}
```

Response:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "my-company-website",
  "description": "Company website and blog",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "default_environment": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "prod"
  }
}
```

### 2. Create Additional Environments

```bash
POST /api/v1/projects/550e8400-e29b-41d4-a716-446655440000/environments
{
  "name": "dev",
  "description": "Development environment",
  "variables": {
    "DEBUG": "true",
    "LOG_LEVEL": "debug"
  }
}
```

### 3. Create Services in Environments

#### Frontend Service (Static Site)

```bash
POST /api/v1/projects/550e8400-e29b-41d4-a716-446655440000/environments/660e8400-e29b-41d4-a716-446655440001/services
{
  "name": "frontend",
  "git_url": "https://github.com/company/website-frontend.git",
  "git_branch": "main",
  "context_root": "",
  "buildpack_type": "static",
  "environment": {
    "NODE_ENV": "production"
  },
  "build_config": {
    "static": {
      "build_command": "npm run build",
      "output_dir": "dist"
    }
  }
}
```

#### Backend Service (Node.js with Dockerfile)

```bash
POST /api/v1/projects/550e8400-e29b-41d4-a716-446655440000/environments/660e8400-e29b-41d4-a716-446655440001/services
{
  "name": "backend",
  "git_url": "https://github.com/company/website-backend.git",
  "git_branch": "main",
  "buildpack_type": "dockerfile",
  "environment": {
    "PORT": "3000",
    "NODE_ENV": "production"
  },
  "build_config": {
    "dockerfile": {
      "dockerfile_path": "Dockerfile",
      "build_args": {
        "NODE_VERSION": "18"
      }
    }
  }
}
```

#### Database Service (Docker Compose)

```bash
POST /api/v1/projects/550e8400-e29b-41d4-a716-446655440000/environments/660e8400-e29b-41d4-a716-446655440001/services
{
  "name": "database",
  "git_url": "https://github.com/company/database-setup.git",
  "git_branch": "main",
  "buildpack_type": "docker-compose",
  "environment": {
    "POSTGRES_DB": "website",
    "POSTGRES_USER": "app",
    "POSTGRES_PASSWORD": "secret123"
  },
  "build_config": {
    "compose": {
      "compose_file": "docker-compose.yml",
      "service": "postgres"
    }
  }
}
```

## Service Operations

### Deploy a Service

```bash
POST /api/v1/projects/{project_id}/environments/{environment_id}/services/{service_id}/deploy
```

This will:

1. Clone the Git repository
2. Build the image using the specified buildpack
3. Deploy the container
4. Update service status

### Stop a Service

```bash
POST /api/v1/projects/{project_id}/environments/{environment_id}/services/{service_id}/stop
```

### Restart a Service

```bash
POST /api/v1/projects/{project_id}/environments/{environment_id}/services/{service_id}/restart
```

### Delete a Service

```bash
DELETE /api/v1/projects/{project_id}/environments/{environment_id}/services/{service_id}
```

## List Resources

### List All Projects

```bash
GET /api/v1/projects
```

### List Environments in a Project

```bash
GET /api/v1/projects/{project_id}/environments
```

### List Services in an Environment

```bash
GET /api/v1/projects/{project_id}/environments/{environment_id}/services
```

## Buildpack Types

### 1. Nixpacks (Auto-detection)

```json
{
  "buildpack_type": "nixpacks",
  "build_config": {
    "nixpacks": {
      "start_command": "npm start",
      "build_command": "npm run build",
      "variables": {
        "NODE_VERSION": "18"
      }
    }
  }
}
```

### 2. Static Sites

```json
{
  "buildpack_type": "static",
  "build_config": {
    "static": {
      "build_command": "npm run build",
      "output_dir": "dist",
      "nginx_config": "nginx.conf"
    }
  }
}
```

### 3. Dockerfile

```json
{
  "buildpack_type": "dockerfile",
  "build_config": {
    "dockerfile": {
      "dockerfile_path": "Dockerfile.prod",
      "build_args": {
        "NODE_VERSION": "18",
        "ENV": "production"
      },
      "target": "production"
    }
  }
}
```

### 4. Docker Compose

```json
{
  "buildpack_type": "docker-compose",
  "build_config": {
    "compose": {
      "compose_file": "docker-compose.prod.yml",
      "service": "web"
    }
  }
}
```

## Example: Full Stack Application

Here's how you would deploy a full-stack application with frontend, backend, and database:

```bash
# 1. Create project
POST /api/v1/projects
{
  "name": "my-saas-app",
  "description": "SaaS application with React frontend and Node.js backend"
}

# 2. Create staging environment (prod is created automatically)
POST /api/v1/projects/{project_id}/environments
{
  "name": "staging",
  "description": "Staging environment for testing"
}

# 3. Deploy frontend (React app)
POST /api/v1/projects/{project_id}/environments/{prod_env_id}/services
{
  "name": "frontend",
  "git_url": "https://github.com/company/saas-frontend.git",
  "buildpack_type": "static",
  "build_config": {
    "static": {
      "build_command": "npm run build",
      "output_dir": "build"
    }
  }
}

# 4. Deploy backend (Node.js API)
POST /api/v1/projects/{project_id}/environments/{prod_env_id}/services
{
  "name": "api",
  "git_url": "https://github.com/company/saas-backend.git",
  "buildpack_type": "nixpacks",
  "environment": {
    "PORT": "3000",
    "NODE_ENV": "production",
    "DATABASE_URL": "postgresql://user:pass@db:5432/saas"
  }
}

# 5. Deploy database (PostgreSQL)
POST /api/v1/projects/{project_id}/environments/{prod_env_id}/services
{
  "name": "database",
  "git_url": "https://github.com/company/db-setup.git",
  "buildpack_type": "docker-compose",
  "environment": {
    "POSTGRES_DB": "saas",
    "POSTGRES_USER": "user",
    "POSTGRES_PASSWORD": "secure_password"
  }
}

# 6. Deploy all services
POST /api/v1/projects/{project_id}/environments/{prod_env_id}/services/{frontend_id}/deploy
POST /api/v1/projects/{project_id}/environments/{prod_env_id}/services/{api_id}/deploy
POST /api/v1/projects/{project_id}/environments/{prod_env_id}/services/{db_id}/deploy
```

This new architecture provides:

- **Clear Organization**: Projects contain environments, environments contain services
- **Environment Separation**: Easy to have dev, staging, prod environments
- **Service Isolation**: Each service can be deployed, scaled, and managed independently
- **Quick Deployment**: One-off service creation for simple use cases
- **Full Control**: Complete project/environment/service hierarchy for complex applications
