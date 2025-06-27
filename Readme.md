# MC-Manager

A web-based management tool for Minecraft Bedrock dedicated servers running in Docker.

## Overview

MC-Manager provides a user-friendly web interface to manage multiple Minecraft Bedrock server worlds. It allows you to easily configure, launch, and switch between different worlds without manually editing files or using the command line. The application is built with Go and uses Docker to containerize and run the Minecraft servers, ensuring a clean and isolated environment for each world.

It also includes a "Play Timer" feature, perfect for parents or server admins who want to limit session playtime. The server will automatically shut down after a set duration, with periodic in-game warnings for the players.

## Features

*   **Web-based UI**: Manage your server from anywhere using a modern, responsive web interface built with `go-app` (Go compiled to WebAssembly).
*   **Multi-world Management**: Create and store configurations for multiple worlds. The application ensures that only one Minecraft server Docker container is active at a time, allowing you to seamlessly switch between your defined worlds with a single click.
*   **Docker Integration**: Runs Minecraft servers in isolated Docker containers for better security and resource management.
*   **Dynamic Configuration**: Modify world settings like Game Mode (`survival`, `creative`, etc.) directly from the UI.
*   **Play Timer**: Set a time limit for a gameplay session. The server will automatically stop, and players will receive countdown notifications.
*   **Favorite Worlds**: Mark your most-used worlds for quick access.
*   **Secure**: Supports integration with an authentication proxy for secure access.
*   **PWA Support**: The frontend is a Progressive Web App that can be "installed" on your desktop or mobile device and notifies you when updates are available.

## Technology Stack

*   **Backend**: Go
*   **Frontend**: Go (compiled to Wasm) using go-app
*   **Database**: MongoDB
*   **Containerization**: Docker
*   **UI Library**: Bootstrap

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

*   Go (version 1.18 or newer)
*   Docker
*   MongoDB

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/ajaxe/mc-manager.git
    cd mc-manager
    ```

2.  **Configure the application:**
    Create a `config.yml` file in the root of the project. You can use the example below as a starting point. See the **Configuration** section for more details on each option.

3.  **Build the application:**
    ```bash
    # This will build the server binary
    go build -o mc-manager ./cmd/server
    ```

4.  **Run the application:**
    ```bash
    ./mc-manager
    ```
    The server will start, and you can access the web UI at `http://localhost:8080` (or the port you specified in your config).

## Configuration

The application is configured via a `config.yml` file. The structure is based on the `AppConfig` struct in `internal/config/config.go`.

Here is an example with explanations for each field:

```yaml
server:
  # Port for the web server to listen on.
  port: "8080"

  # (Optional) For enabling HTTPS on the development server.
  # cert_file: "/path/to/cert.pem"
  # key_file: "/path/to/key.pem"

  # (Optional) Required if running Docker on Windows/macOS via Docker Desktop.
  # Use "unix:///var/run/docker.sock" on Linux.
  docker_host_url: "npipe:////./pipe/docker_engine"

  # (Optional) If using an external authentication provider (e.g., Authelia, Authentik),
  # this is the URL the user will be redirected to for login.
  auth_redirect_url: "https://auth.example.com/"

  # The name of the cookie to check for an authenticated session.
  auth_cookie_name: "my_auth_cookie"

  # A secret token sent to the auth provider to identify this service.
  auth_token: "your-secret-auth-token"

database:
  # Connection string for your MongoDB instance.
  connection_uri: "mongodb://localhost:27017"
  # The name of the database to use.
  db_name: "mc_manager"

game_server:
  # The base directory on the host where all game server files will be stored.
  hosting_dir: "/opt/mc-manager/servers"

  # The directory inside the container where world data is stored.
  # This should generally not be changed.
  world_dir: "/data/worlds"

  # The Docker image for the Minecraft Bedrock server.
  # Example: "itzg/minecraft-bedrock-server"
  image_name: "itzg/minecraft-bedrock-server"

  # Environment variables to pass to the container.
  # These are often used to configure the itzg/minecraft-bedrock-server image.
  env_vars:
    - "EULA=TRUE"
    - "GAMEMODE=survival"
    - "DIFFICULTY=normal"

  # Host-to-container volume mappings.
  # The application automatically manages the world data volume.
  volumes: []

  # Docker labels to apply to the container for discovery or management (e.g., for Traefik).
  labels:
    - "traefik.enable=true"

  # Docker networks to connect the container to.
  networks:
    - "proxy_network"

  # (Optional) Configure Docker container logging.
  logging:
    driver: "json-file"
    options:
      max-size: "10m"
      max-file: "3"
