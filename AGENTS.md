# Agent Instructions for mikrocloud-2

## Build/Test Commands
- **Backend build**: `make build` (or `go build -o bin/mikrocloud ./main.go`)
- **Frontend build**: `make build-web` (or `cd web && pnpm run build`)
- **Full build**: `make build-full` (builds both with embedded frontend)
- **Backend tests**: `make test` (or `go test -v ./...`) 
- **Frontend tests**: `cd web && pnpm test` (single run) or `cd web && pnpm test:unit` (watch mode)
- **Lint backend**: `make lint` (requires golangci-lint)
- **Lint frontend**: `cd web && pnpm lint`
- **Format code**: `make fmt` (Go) and `cd web && pnpm format` (frontend)

## Code Style

### Go Backend
- Use standard Go conventions: camelCase for private, PascalCase for public
- Follow DDD patterns: domain models in `internal/domain/`, separate value objects
- Error handling: return errors, don't panic; use `fmt.Errorf()` for wrapping
- Use Go Chi for API handlers with proper context and structured responses
- Database: SQLite with goose migrations, use Bob ORM patterns
- Imports: group std library, third-party, then local packages

### Frontend (Svelte/TypeScript)
- Use tabs for indentation, single quotes, no trailing commas (per .prettierrc)
- TypeScript strict mode, prefer explicit types over `any`
- Use Svelte 5 syntax and features
- UI components follow shadcn-svelte patterns in `lib/components/ui/`
- Import order: Svelte imports first, then libraries, then local
- Use TanStack Query for data fetching
- CSS: Tailwind classes, avoid custom CSS unless necessary

### Misc

You are able to use the Svelte MCP server, where you have access to comprehensive Svelte 5 and SvelteKit documentation. Here's how to use the available tools effectively:

## Available MCP Tools:

### 1. list-sections

Use this FIRST to discover all available documentation sections. Returns a structured list with titles, use_cases, and paths.
When asked about Svelte or SvelteKit topics, ALWAYS use this tool at the start of the chat to find relevant sections.

### 2. get-documentation

Retrieves full documentation content for specific sections. Accepts single or multiple sections.
After calling the list-sections tool, you MUST analyze the returned documentation sections (especially the use_cases field) and then use the get-documentation tool to fetch ALL documentation sections that are relevant for the user's task.

### 3. svelte-autofixer

Analyzes Svelte code and returns issues and suggestions.
You MUST use this tool whenever writing Svelte code before sending it to the user. Keep calling it until no issues or suggestions are returned.

### 4. playground-link

Generates a Svelte Playground link with the provided code.
After completing the code, ask the user if they want a playground link. Only call this tool after user confirmation and NEVER if code was written to files in their project.
