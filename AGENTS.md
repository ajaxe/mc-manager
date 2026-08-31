# AGENTS.md: Minecraft Server Manager (mc-manager)

## 1. Project Mission
**MC-Manager** is a lightweight, web-based platform designed to simplify the lifecycle management of **Minecraft Bedrock server instances**. It leverages **Docker** and **Go WebAssembly** to provide a robust, isolated environment for creating, configuring, launching, monitoring, and playtime-limiting game worlds with ease.

## 2. Tech Stack & Core Dependencies
- **Primary Language:** **Go (1.25+)** - Powers both the backend REST API and frontend WebAssembly binary.
- **Frontend Framework:** **go-app (v10)** (`github.com/maxence-charriere/go-app/v10`) - Compiles Go UI components into a Progressive Web App (PWA) running in WebAssembly (WASM).
- **Backend Framework:** **Echo (v4)** (`github.com/labstack/echo/v4`) - High-performance HTTP server handling REST API endpoints and static WASM asset distribution.
- **Database:** **MongoDB** (`go.mongodb.org/mongo-driver/v2`) - Stores persistent data for world configurations, launch histories, and play timer session records.
- **Infrastructure & Orchestration:** **Docker Engine API** (`github.com/docker/docker`, `github.com/docker/cli`) - Orchestrates game server containers, manages host volume bindings, controls server properties, and routes logs.
- **Configuration Management:** **Viper** (`github.com/spf13/viper`) & **fsnotify** - Manages settings via `config.yaml` with environment variable overrides (`APP_*` prefix) and hot-reloading support.
- **UI & Styling:** **Bootstrap 5** - Provides responsive layouts, CSS styling, icons, and components rendered declaratively via `go-app`.

## 3. Architecture & Design Patterns
- **Hybrid Target Compilation:** The codebase shares domain models (`internal/models`) across targets while using Go build tags (`//go:build wasm` vs `//go:build unix || windows`) to decouple client WASM execution from native OS and Docker SDK dependencies.
- **Single Page Application (WASM PWA):** The browser frontend runs `app.wasm` compiled from `cmd/webapp/`, communicating with the backend exclusively through JSON REST endpoints via `internal/client/api_client.go`.
- **Repository Pattern:** Database interactions in `internal/db` encapsulate MongoDB collection operations (`worlds`, `launches`, `playtimers`) into strongly typed interfaces used by Echo handlers.
- **Container Orchestration Layer:** The `internal/gameserver` package translates user actions and world settings into Docker container creations (`ContainerCreate`), starts (`ContainerStart`), stops (`ContainerStop`), and console commands (`ContainerExecAttach`).
- **Background Play Timer Monitor:** A background goroutine managed in `internal/job` (`StartMonitor`) continuously polls active play timer sessions, sends in-game countdown alerts via Docker exec, and automatically stops server containers upon session expiration.
- **Forward Auth Integration:** Backend auth handlers (`internal/handlers/login_handler.go`) check session cookies (`auth_cookie_name`) and support external forward-authentication proxy redirect loops when enabled.

## 4. Directory Mental Model
- **`cmd/webapp/`**: Main application entry points (`main.go`, `backend-entry.go`, `frontend-entry.go`) coordinating WASM client routing (`Frontend()`) and Echo backend initialization (`Backend()`).
- **`internal/client/`**: Frontend-only logic including WASM HTTP client wrappers (`api_client.go`, `fetch.go`), global reactive application context (`AppContext`), and state keys (`state_keys.go`).
- **`internal/components/`**: Reusable declarative `go-app` UI components (e.g., `WorldItemCard`, `SidebarMenu`, `PlayTimerCountdown`, `WorldDeleteConfirm`, `StatusToast`).
- **`internal/config/`**: Configuration structs (`AppConfig`) and Viper initialization logic supporting `config.yaml` loading, env binding, and file watching.
- **`internal/db/`**: MongoDB connection client (`NewClient`) and repository CRUD logic for worlds, launches, and play timers.
- **`internal/gameserver/`**: Core logic for Docker daemon client management, server container configuration generators (`serverconfig.go`), and in-game console messaging (`console.go`).
- **`internal/handlers/`**: Echo HTTP handlers implementing the backend REST API (`worlds_handlers.go`, `launch_handlers.go`, `play_timer_handler.go`, `login_handler.go`, `healthcheck_handler.go`).
- **`internal/http/`**: Common HTTP status constants and response definitions.
- **`internal/job/`**: Asynchronous background jobs, including play timer monitoring loops (`monitor.go`) and server shutdown orchestration (`tracking.go`).
- **`internal/models/`**: Shared domain models (`WorldItem`, `LaunchItem`, `PlayTimerItem`, `GameServerDetail`), API payload wrappers (`ApiResult`), custom Echo context (`AppContext`), and error structs (`AppError`).
- **`internal/pages/`**: Top-level WASM PWA page components (e.g., `HomePage`, `WorldAddPage`, `WorldLaunchPage`, `PlayTimerPage`, `SetupPage`).
- **`internal/server/`**: Server initialization (`NewBackendApi`), middleware pipeline configuration, error handling rules, and graceful OS signal shutdown (`Start`).
- **`web/`**: Static assets (CSS, JS, icons) and output destination for compiled `app.wasm`.

## 5. Development Standards
- **Naming Conventions:**
  - **Go Files:** Use `snake_case.go` for backend and domain files (e.g., `worlds_handlers.go`, `api_client.go`); use `PascalCase.go` or `snake_case.go` for UI components.
  - **UI Components:** Use **PascalCase** for `go-app` component structs (e.g., `WorldItemCard`, `LaunchItemCard`) implementing `app.Composer`.
  - **Database Methods:** Use explicit action-based method names (e.g., `GetWorldByID`, `InsertWorld`, `DeleteWorldByID`, `UpdateWorldFav`).
  - **REST Endpoints:** Group handlers under the `/api` route group with lowercase pluralized resources (e.g., `/api/worlds`, `/api/launches`, `/api/playtimer`).
- **Error Handling & Logging Strategies:**
  - Handlers must use `models.AppError` to encapsulate HTTP status codes, user-facing error messages, and internal error details.
  - Centralized Echo error middleware (`AppErrorHandler`) converts errors into standard `AppError` JSON responses.
  - Structured logging uses standard `slog` with a JSON handler (`slog.NewJSONHandler`) for backend operations, combined with Echo logger instances.
- **Implementation Patterns:**
  - **Build Tag Isolation:** Keep Docker Engine API and native OS dependencies constrained under `//go:build unix || windows` or `//go:build !wasm`. Stub or separate WASM implementation files using `//go:build wasm` (e.g., `game_service_js.go`).
  - **WASM State Management:** Manage global frontend state via `internal/client/AppContext` and trigger reactive UI updates using `SetState` or `app.Context.Handle()`.
  - **WASM API Calls:** Always use helper wrappers in `internal/client/api_client.go` (`httpGet`, `httpPost`, `httpDelete`) to ensure proper `fetch` headers, credential handling, and response deserialization.
  - **MongoDB Handling:** Use MongoDB Driver v2 with `ObjectIDAsHexString: true` BSON options. Enforce explicit context timeouts for all database operations (`2*time.Second` ping, `30*time.Second` query operations).
  - **Authentication Checks:** Read session cookies via `cc.AuthCookieValue()`. When `auth_server_enabled` is true, redirect unauthenticated users to the configured `AuthRedirectURL()`.
  - **Go Test Coverage:** Every change or addition to Go source code (`.go`) must include creating new unit tests or updating existing unit tests in the corresponding `*_test.go` file.

## 6. Hard Constraints & Anti-Patterns
- **DO NOT Import Docker Packages into WASM Code:** Never import `github.com/docker/docker` or `internal/gameserver` into code compiled under `//go:build wasm`.
- **DO NOT Perform Direct DOM Manipulation:** In frontend components, strictly use `go-app`'s declarative Go syntax (`app.Div()`, `app.Button()`); avoid raw JavaScript DOM calls unless wrapped in encapsulated helpers.
- **DO NOT Bypass `models.AppError`:** API handlers must not return plain error strings or raw unhandled errors; always format responses through `AppError` or standard `ApiResult`.
- **DO NOT Allow Deletion of Favorite Worlds:** Worlds with `IsFavorite == true` must be strictly protected from deletion in both API backend handlers and UI client modals.
- **DO NOT Leave Concurrent Container Instances:** When launching a server world, ensure all existing game server container instances are stopped (`stopAllInstances()`) so only one container runs at any time.
- **DO NOT Hardcode Host Paths or Secret Tokens:** Pass all host volume paths (`HOSTING_DIR`), database connection strings, and authentication secret tokens via `config.yaml` or `APP_*` environment variables.
- **DO NOT Skip Build Synchronization:** Run `make build` whenever frontend components or pages are modified; otherwise, the server will continue serving an outdated `app.wasm` from `./tmp/web/`.
- **DO NOT Modify Go Code Without Updating Tests:** Never submit or complete a Go (`.go`) code modification without creating a new test or updating existing unit tests in the package.

## 7. Operational Commands
- **Local Building & Running:**
  - `make build`: Compiles `app.wasm` (`GOARCH=wasm GOOS=js`) and native server binary (`GOARCH=amd64`) into `./tmp/`, copying static web assets and `config.yaml`.
  - `make run`: Executes `make build` and starts the compiled server binary (`./tmp/server.exe`).
  - `make clean`: Removes the `./tmp/` build directory and compiled output artifacts.
- **Testing:**
  - `go test ./...`: Runs all package unit tests across the repository.
  - `go test -v ./internal/gameserver`: Runs unit tests specifically for game server container operations.
  - `docker build . --tag apogee-dev/mc-manager-tests:local -f Dockerfile.test`: Builds the containerized integration test suite image.
  - `docker run --volume=/var/run/docker.sock:/var/run/docker.sock apogee-dev/mc-manager-tests:local`: Runs containerized tests with mounted Docker socket.
- **Production Packaging:**
  - `make prod`: Packages the web application and server into a production Docker image `apogee-dev/mc-manager:local`.
  - `docker build . --tag apogee-dev/mc-manager:local -f Dockerfile`: Manually builds the production Docker image.
- **Live Development Watch Mode:**
  - `wgo -xdir tmp -file .go -file .js -file .css make run`: Watches Go/JS/CSS source files and automatically rebuilds and restarts the server on changes.
