# MC-Manager

A web-based management tool for Minecraft Bedrock dedicated servers running in Docker.

## Overview

**MC-Manager** provides a modern, responsive web interface to manage multiple Minecraft Bedrock server worlds. It enables creating, configuring, launching, and switching between different worlds without manually editing configuration files or running command-line scripts. The application is built with Go and Docker to containerize game server instances, ensuring clean isolation and simple management.

It includes a **Play Timer** feature, ideal for parents or server administrators who want to limit gameplay session durations. The server automatically shuts down after a set duration, broadcasting periodic in-game warning messages to connected players.

---

## Features

- **WebAssembly PWA UI**: Single Page Application (PWA) built with `go-app` (v10), compiling Go UI components directly into WebAssembly (`app.wasm`).
- **Multi-World Orchestration**: Configure and store multiple world settings. MC-Manager enforces a single active container rule, stopping running instances before launching a new world.
- **Docker SDK Integration**: Direct orchestration of Minecraft Bedrock server containers via the Docker Engine API.
- **Dynamic Server Properties**: Modify game mode (`survival`, `creative`, etc.), level names, seeds, and server settings directly from the web interface.
- **Play Timer & In-Game Broadcasts**: Set session timers with automated countdown warnings sent into the Minecraft console via `docker exec`.
- **Favorite Worlds**: Mark favorite worlds for quick access and deletion protection.
- **Forward Authentication Integration**: Supports session cookie checks and redirect authentication proxies (Authelia, Authentik).
- **Responsive Bootstrap 5 Design**: Desktop and mobile-friendly UI layout.

---

## Technology Stack

- **Primary Language**: Go (1.25+)
- **Frontend Framework**: [go-app (v10)](https://github.com/maxence-charriere/go-app/v10) (Go compiled to WebAssembly)
- **Backend Framework**: [Echo (v4)](https://github.com/labstack/echo/v4)
- **Database**: MongoDB (`go.mongodb.org/mongo-driver/v2`)
- **Containerization**: Docker Engine API (`github.com/docker/docker`)
- **Configuration**: Viper (`config.yaml` with `APP_*` environment variable overrides)
- **UI & Styling**: Bootstrap 5

---

## Getting Started

### Prerequisites

* **Go**: 1.25 or newer
* **Docker Engine / Docker Desktop**: Running with socket or named pipe access
* **MongoDB**: Instance running locally or accessible via network

### Installation & Execution

1. **Clone the repository:**
   ```bash
   git clone https://github.com/ajaxe/mc-manager.git
   cd mc-manager
   ```

2. **Configure the application:**
   Copy or edit `config.yaml` in the root directory. See the **Configuration** section below for field descriptions.

3. **Build the application:**
   Build the WebAssembly frontend asset and native backend server binary into `./tmp/`:
   ```bash
   make build
   ```

4. **Run the server:**
   ```bash
   make run
   ```
   Access the web interface at `http://localhost:8000` (or your configured port).

5. **Clean build artifacts:**
   ```bash
   make clean
   ```

---

## Configuration

The application reads configuration from `config.yaml` using Viper, supporting environment variable overrides prefixed with `APP_` (e.g., `APP_SERVER_PORT=8080`).

### Example `config.yaml`

```yaml
is_dev: true

server:
  # Port for the REST API and WASM server to listen on
  port: 8000
  # Optional TLS certificate and key paths
  cert_file: ""
  key_file: ""
  # Docker host URL (e.g. unix:///var/run/docker.sock or npipe:////./pipe/docker_engine)
  docker_host_url: ""
  # Forward auth settings
  auth_token: ""
  auth_server_enabled: false
  auth_server_url: ""
  auth_redirect_path: "/login/check"
  auth_introspect_path: "/login/introspect"
  auth_cookie_name: ".fwd-auth-custom"
  admins: "admin"

database:
  connection_uri: "mongodb://localhost:27017"
  db_name: "mc_manager"

game_server:
  # Host directory where server world files are hosted
  hosting_dir: "/opt/mc-manager/servers"
  # Target directory inside container for world data
  world_dir: "/minecraft/worlds"
  # Minecraft Bedrock server container image
  image_name: "itzg/minecraft-bedrock-server"
  # Default container environment variables
  env_vars:
    - "EULA=TRUE"
    - "GAMEMODE=survival"
    - "DIFFICULTY=normal"
  # Additional volume mappings
  volumes: []
  # Docker container labels
  labels:
    - "traefik.enable=true"
  # Docker network names
  networks:
    - "public-proxy"
  # Logging configuration
  logging:
    driver: "json-file"
    options:
      max-size: "10m"
      max-file: "3"
```

---

## Development & Testing

### Live Development (Watch Mode)

Use `wgo` for automatic rebuilding and hot-reloading on Go/JS/CSS changes:

```pwsh
wgo -xdir tmp -file .go -file .js -file .css make run
```

### Go Build Tags

MC-Manager uses Go build tags to isolate WebAssembly client code from native Docker SDK dependencies:

* `//go:build wasm`: Compiles `app.wasm` for the browser frontend target. Excludes Docker SDK and OS system calls.
* `//go:build !wasm` / `//go:build unix || windows`: Compiles the native OS backend server binary with Docker Engine API support.

### Testing Commands

* **Run all package unit tests:**
  ```bash
  go test ./...
  ```

* **Run game server package unit tests:**
  ```bash
  go test -v ./internal/gameserver
  ```

* **Build integration test Docker image:**
  ```bash
  docker build . --tag apogee-dev/mc-manager-tests:local -f Dockerfile.test
  ```

* **Execute containerized integration tests:**
  ```bash
  docker run --volume=/var/run/docker.sock:/var/run/docker.sock apogee-dev/mc-manager-tests:local
  ```

### Production Docker Image

To build the production-ready Docker image:

```bash
make prod
```
or manually:
```bash
docker build . --tag apogee-dev/mc-manager:local -f Dockerfile
```

