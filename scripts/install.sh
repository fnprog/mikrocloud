#!/bin/bash
set -e
set -o pipefail

CDN="https://cdn.mikrocloud.io"
GITHUB_REPO="mikrocloud/mikrocloud"
VERSION="latest"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="${HOME}/.config/mikrocloud"
DATA_DIR="${HOME}/.local/share/mikrocloud"

MIKROCLOUD_PORT=${MIKROCLOUD_PORT:-3000}
ROOT_USERNAME=${ROOT_USERNAME:-}
ROOT_USER_EMAIL=${ROOT_USER_EMAIL:-}
ROOT_USER_PASSWORD=${ROOT_USER_PASSWORD:-}

if [ $EUID != 0 ]; then
    echo "⚠️  This script should be run as root or with sudo for system-wide installation"
    echo "   Continuing with user-level installation..."
    INSTALL_DIR="${HOME}/.local/bin"
fi

echo "╔════════════════════════════════════════════════════════════╗"
echo "║         Welcome to Mikrocloud Installer v0.1.0            ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "This script will install Mikrocloud on your system."
echo "Source: https://github.com/${GITHUB_REPO}"
echo ""

OS_TYPE=$(grep -w "ID" /etc/os-release | cut -d "=" -f 2 | tr -d '"' 2>/dev/null || echo "unknown")
echo "📋 Detected OS: ${OS_TYPE}"

command_exists() {
    command -v "$@" > /dev/null 2>&1
}

is_proxmox_lxc() {
    if [ -n "$container" ] && [ "$container" = "lxc" ]; then
        return 0
    fi
    if grep -q "container=lxc" /proc/1/environ 2>/dev/null; then
        return 0
    fi
    return 1
}

check_port() {
    local port=$1
    if command_exists ss; then
        if ss -tulnp | grep ":${port} " >/dev/null 2>&1; then
            return 1
        fi
    elif command_exists netstat; then
        if netstat -tuln | grep ":${port} " >/dev/null 2>&1; then
            return 1
        fi
    fi
    return 0
}

echo ""
echo "🔍 Checking prerequisites..."

if [ -f /.dockerenv ]; then
    echo "❌ Error: This script cannot be run inside a Docker container"
    exit 1
fi

echo "   ✓ Not running in container"

for port in 80 443 ${MIKROCLOUD_PORT}; do
    if ! check_port ${port}; then
        echo "❌ Error: Port ${port} is already in use"
        echo "   Please free up this port before continuing"
        exit 1
    fi
done

echo "   ✓ Required ports (80, 443, ${MIKROCLOUD_PORT}) are available"

if is_proxmox_lxc; then
    echo "⚠️  WARNING: Proxmox LXC container detected"
    echo "   Docker Swarm mode may require special configuration"
    echo "   Continuing in 3 seconds..."
    sleep 3
fi

echo ""
echo "🐳 Checking Docker installation..."

if command_exists docker; then
    DOCKER_VERSION=$(docker version --format '{{.Server.Version}}' 2>/dev/null || echo "unknown")
    echo "   ✓ Docker already installed (version: ${DOCKER_VERSION})"
else
    echo "   Docker not found. Installing Docker..."
    
    if [ "$EUID" -ne 0 ]; then
        echo "❌ Error: Docker installation requires root privileges"
        echo "   Please run: curl -sSL https://get.docker.com | sudo sh"
        exit 1
    fi
    
    curl -fsSL https://get.docker.com | sh
    
    if [ $? -eq 0 ]; then
        echo "   ✓ Docker installed successfully"
    else
        echo "❌ Error: Failed to install Docker"
        exit 1
    fi
    
    if command_exists systemctl; then
        systemctl enable docker
        systemctl start docker
    fi
    
    if [ -n "$SUDO_USER" ]; then
        usermod -aG docker "$SUDO_USER" 2>/dev/null || true
        echo "   ✓ Added $SUDO_USER to docker group"
        echo "   ⚠️  Please log out and back in for group changes to take effect"
    fi
fi

if ! docker info >/dev/null 2>&1; then
    echo "❌ Error: Docker daemon is not running"
    echo "   Please start Docker: sudo systemctl start docker"
    exit 1
fi

echo "   ✓ Docker daemon is running"

echo ""
echo "📦 Downloading Mikrocloud CLI..."

ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    armv7l)
        ARCH="armv7"
        ;;
    *)
        echo "❌ Error: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

OS=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ ! -d "$INSTALL_DIR" ]; then
    mkdir -p "$INSTALL_DIR"
fi

BINARY_URL="https://github.com/${GITHUB_REPO}/releases/latest/download/mikrocloud-cli-${OS}-${ARCH}"
TEMP_BINARY="/tmp/mikrocloud-cli-$$"

echo "   Downloading from: ${BINARY_URL}"

if command_exists curl; then
    curl -fSL "${BINARY_URL}" -o "${TEMP_BINARY}" 2>/dev/null || {
        echo "⚠️  GitHub release not available, building from source..."
        echo "   This requires Go 1.24+ to be installed"
        
        if ! command_exists go; then
            echo "❌ Error: Go is not installed. Please install Go 1.24+ or wait for binary releases"
            exit 1
        fi
        
        TEMP_DIR="/tmp/mikrocloud-build-$$"
        git clone "https://github.com/${GITHUB_REPO}.git" "${TEMP_DIR}"
        cd "${TEMP_DIR}"
        make build-cli
        cp bin/mikrocloud-cli "${TEMP_BINARY}"
        cd -
        rm -rf "${TEMP_DIR}"
    }
elif command_exists wget; then
    wget -q "${BINARY_URL}" -O "${TEMP_BINARY}" 2>/dev/null || {
        echo "❌ Error: Failed to download CLI binary and curl/wget not available"
        exit 1
    }
else
    echo "❌ Error: Neither curl nor wget found. Please install one of them."
    exit 1
fi

chmod +x "${TEMP_BINARY}"

if [ -w "$INSTALL_DIR" ] || [ "$EUID" -eq 0 ]; then
    mv "${TEMP_BINARY}" "${INSTALL_DIR}/mikrocloud-cli"
    echo "   ✓ Installed mikrocloud-cli to ${INSTALL_DIR}/mikrocloud-cli"
else
    echo "❌ Error: No write permission to ${INSTALL_DIR}"
    echo "   Please run with sudo or choose a different installation directory"
    rm -f "${TEMP_BINARY}"
    exit 1
fi

if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo "⚠️  Warning: ${INSTALL_DIR} is not in your PATH"
    echo "   Add this to your ~/.bashrc or ~/.zshrc:"
    echo "   export PATH=\"${INSTALL_DIR}:\$PATH\""
fi

echo ""
echo "📝 Creating configuration..."

mkdir -p "${CONFIG_DIR}"
mkdir -p "${DATA_DIR}"

if [ ! -f "${CONFIG_DIR}/mikrocloud.toml" ]; then
    cat > "${CONFIG_DIR}/mikrocloud.toml" << 'CONFIGEOF'
[server]
host = "0.0.0.0"
port = 3000
data_dir = "${HOME}/.local/share/mikrocloud"
log_level = "info"

[database]
type = "sqlite"
url = "${HOME}/.local/share/mikrocloud/mikrocloud.db"

[analytics]
type = "duckdb"
url = "${HOME}/.local/share/mikrocloud/analytics.duckdb"

[queue]
type = "dragonfly"
enabled = true
auto_start = true
url = "redis://mikrocloud-queue:6379/0"

[docker]
runtime = "docker"
socket_path = "/var/run/docker.sock"
rootless = false
network_mode = "bridge"

[proxy]
enabled = true
image = "traefik:v3.0"
dashboard_enabled = true
dashboard_port = 8080

[metrics]
enabled = false
image = "prom/prometheus:latest"

[tunnel]
enabled = false
image = "cloudflare/cloudflared:latest"
token = ""

[ssl]
enabled = false
acme_email = ""
staging = true
certs_dir = "${HOME}/.local/share/mikrocloud/certs"

[auth]
enabled = true
jwt_secret = ""
CONFIGEOF

    JWT_SECRET=$(openssl rand -hex 32 2>/dev/null || cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 64 | head -n 1)
    
    if [ "$(uname)" = "Darwin" ]; then
        sed -i '' "s/jwt_secret = \"\"/jwt_secret = \"${JWT_SECRET}\"/" "${CONFIG_DIR}/mikrocloud.toml"
    else
        sed -i "s/jwt_secret = \"\"/jwt_secret = \"${JWT_SECRET}\"/" "${CONFIG_DIR}/mikrocloud.toml"
    fi
    
    echo "   ✓ Created configuration at ${CONFIG_DIR}/mikrocloud.toml"
    echo "   ✓ Generated JWT secret"
else
    echo "   ✓ Configuration already exists at ${CONFIG_DIR}/mikrocloud.toml"
fi

echo "   ✓ Created data directory at ${DATA_DIR}"

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║            Installation completed successfully!            ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "🚀 Next steps:"
echo ""
echo "   1. Start Mikrocloud:"
echo "      $ mikrocloud-cli start"
echo ""
echo "   2. Access the dashboard:"
echo "      http://localhost:${MIKROCLOUD_PORT}"
echo ""
echo "   3. Check status anytime:"
echo "      $ mikrocloud-cli status"
echo ""
echo "   4. View logs:"
echo "      $ docker logs -f mikrocloud"
echo ""
echo "📚 Documentation: https://mikrocloud.io/docs"
echo "💬 Community: https://discord.gg/mikrocloud"
echo "🐛 Issues: https://github.com/${GITHUB_REPO}/issues"
echo ""

if [ -n "$ROOT_USERNAME" ] && [ -n "$ROOT_USER_EMAIL" ] && [ -n "$ROOT_USER_PASSWORD" ]; then
    echo "ℹ️  Predefined credentials detected (will be used on first start)"
fi
