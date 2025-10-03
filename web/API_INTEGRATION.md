# Frontend API Integration

## Overview

The frontend uses a relative path (`/api`) for all API requests. This works in both development and production:

### Production
- Frontend is embedded in the Go binary
- All requests to `/api/*` are handled by the backend API routes
- No CORS issues since everything is same-origin

### Development
- Vite dev server runs on port 5173 (default)
- Backend runs on port 8080
- Vite proxy forwards `/api/*` requests to `http://localhost:8080`
- See `vite.config.ts` for proxy configuration

## API Client

Located in `/web/src/lib/api/`:

- `client.ts` - Generic HTTP client with auth token handling
- `auth.ts` - Authentication API methods (login, register, logout)
- `index.ts` - Barrel exports for clean imports

### Usage

```typescript
import { authApi } from '$lib/api';

// Login
const { token, user } = await authApi.login({ 
  email: 'user@example.com', 
  password: 'password' 
});

// Register
const { token, user } = await authApi.register({
  name: 'John Doe',
  email: 'john@example.com',
  password: 'securepassword'
});

// Check authentication
if (authApi.isAuthenticated()) {
  // User is logged in
}

// Logout
authApi.logout();
```

## Authentication Flow

1. User submits login/register form
2. API client sends request to `/api/auth/login` or `/api/auth/register`
3. Backend returns JWT token and user data
4. Token stored in localStorage
5. Token automatically added to all subsequent requests via `Authorization: Bearer <token>` header
6. Protected routes check auth status and redirect to login if needed

## Protected Routes

Dashboard routes are protected in `/web/src/routes/dashboard/+layout.svelte`:

```typescript
onMount(() => {
  authStore.initialize();
  if (!authStore.isAuthenticated()) {
    goto('/login');
  }
});
```

## Development Setup

1. Start backend: `make run` or `go run main.go`
2. Start frontend: `cd web && pnpm dev`
3. Open browser to `http://localhost:5173`
4. API requests will be proxied to backend automatically

## Environment Variables

Create `web/.env` for custom configuration (optional):

```bash
# Override API URL for development (not usually needed)
# VITE_API_URL=http://localhost:8080
```

The proxy handles everything by default, so you typically don't need any environment variables.
