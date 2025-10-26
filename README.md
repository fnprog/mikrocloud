 <img src="docs/assets/og.png" alt="mikrocloud" width="100">

<br/>
<br/>

> **Mikrocloud** is a modern, self-hostable **Platform as a Service (PaaS)** that simplifies deploying, managing, and scaling applications and databases — on your own servers.

## Quick Start

### One-Line Installation

```bash
curl -sSL https://cdn.mikrocloud.io/install.sh | sh
```

This installs:
- Docker (if not already installed)
- Mikrocloud CLI binary
- Default configuration at `~/.config/mikrocloud/mikrocloud.toml`
- Data directory at `~/.local/share/mikrocloud`

### Start Mikrocloud

```bash
mikrocloud-cli start
```

This launches:
- **mikrocloud-queue** - Dragonfly message queue (port 6379)
- **mikrocloud-proxy** - Traefik reverse proxy (ports 80, 443, 8080)
- **mikrocloud** - Main server (port 3000)
- **mikrocloud-metrics** - Prometheus (port 9090, if enabled)
- **mikrocloud-cloudflared** - Cloudflare Tunnel (if enabled)

Access the dashboard at: **http://localhost:3000**

### Manage Mikrocloud

```bash
mikrocloud-cli status    # Check service status
mikrocloud-cli stop      # Stop all services
mikrocloud-cli start     # Start all services
```

## Configuration

Edit `~/.config/mikrocloud/mikrocloud.toml` to customize:

```toml
[server]
port = 3000
data_dir = "${HOME}/.local/share/mikrocloud"

[proxy]
enabled = true
auto_start = true
http_port = 80
https_port = 443

[queue]
enabled = true
auto_start = true

[metrics]
enabled = false    # Set to true to enable Prometheus
auto_start = false

[tunnel]
enabled = false    # Set to true for Cloudflare Tunnel
token = ""         # Add your Cloudflare Tunnel token
```

## Docker Compose (Alternative)

For production deployments, use docker-compose:

```bash
docker-compose up -d
```

Customize via environment variables:
```bash
ACME_EMAIL=admin@example.com \
MIKROCLOUD_DOMAIN=mikrocloud.example.com \
docker-compose up -d
```

For metrics and tunnel profiles:
```bash
docker-compose --profile metrics --profile tunnel up -d
```

## Requirements

- **Docker** or Podman
- **Linux** (amd64, arm64, armv7)
- Ports 80, 443, 3000 available

## Development

```bash
make build        # Build backend
make build-web    # Build frontend
make build-full   # Build both with embedded frontend
make test         # Run tests
make lint         # Run linters
```

## License

MIT License - see LICENSE file for details
