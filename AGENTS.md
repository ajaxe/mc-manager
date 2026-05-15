# AGENTS.md: Minecraft Server Manager (mc-manager)

## 1. Project Mission
**MC-Manager** is a lightweight, web-based platform designed to simplify the lifecycle management of **Minecraft server instances**. It leverages **Docker** to provide a robust and isolated environment for creating, launching, and monitoring game worlds with ease.

## 2. Tech Stack & Core Dependencies
- **Language:** **Go (1.25+)** - Primary language for both Backend and Frontend.
- **Frontend Framework:** **go-app (v10)** - Compiles Go code into a **WebAssembly (WASM)** PWA.
- **Backend Framework:** **Echo (v4)** - Handles RESTful API routing and static file serving.
- **Database:** **MongoDB** - Used for persistent storage of world configurations and logs.
- **Infrastructure:** **Docker Engine API** - Orchestrates game server containers.
- **Configuration:** **Viper** - Manages application settings via `config.yaml`.
- **UI Styling:** **Bootstrap 5** - Provides responsive layout and pre-styled components.

## 3. Architecture & Design Patterns
- **Hybrid Codebase:** Shared domain models in `internal/models` are used by both WASM (frontend) and Native (backend) targets.
- **WASM PWA:** The frontend is a Single Page Application (SPA) that runs in the browser via WebAssembly, communicating with the backend via a JSON API.
- **Repository Pattern:** The `internal/db` package abstracts MongoDB interactions into a clean interface for handlers.
- **Component-Based UI:** The frontend uses the `go-app` component model, where UI elements are structured as Go structs with `Render()` methods.
- **Docker Integration:** The backend functions as a thin orchestration layer that translates user actions into Docker container operations.

## 4. Directory Mental Model
- **`cmd/webapp/`**: Main entry points. Coordinates `Frontend()` (WASM) and `Backend()` (Echo) initialization.
- **`internal/client/`**: Frontend-specific logic, including the API client, application context, and state management.
- **`internal/components/`**: Reusable UI components (e.g., `WorldItemCard`, `SidebarMenu`).
- **`internal/db/`**: MongoDB data access layer and collection-specific logic.
- **`internal/gameserver/`**: Core logic for Docker container management and Minecraft server property configuration.
- **`internal/handlers/`**: Echo HTTP handlers that implement the backend REST API.
- **`internal/models/`**: Shared Go structs for database entities, API payloads, and error types.
- **`internal/pages/`**: High-level PWA page components (e.g., `HomePage`, `WorldAddPage`).
- **`web/`**: Static assets (CSS, JS, Icons) and the destination for the compiled `app.wasm`.

## 5. Development Standards
- **Naming Conventions:**
    - Handlers: Grouped by domain in `internal/handlers/` (e.g., `worlds_handlers.go`).
    - UI Components: Use **PascalCase** for struct names (e.g., `LaunchItemCard`).
    - DB Methods: Use descriptive names like `InsertWorld` or `DeleteWorldByID`.
- **Error Handling:**
    - Use `models.AppError` to return consistent API responses. It encapsulates an HTTP status code, a developer-friendly message, and the original error.
- **State Management:**
    - Utilize `internal/client/AppContext` to manage global state and trigger asynchronous data loading via `SetState`.
- **API Communication:**
    - Always use the wrappers in `internal/client/api_client.go` (`httpGet`, `httpPost`, etc.) which handle WASM-specific `fetch` calls.

## 6. Hard Constraints & Anti-Patterns
- **Docker Requirement:** The backend **must** have access to a Docker daemon (via socket or TCP) to function.
- **Validation Rules:** World names are strictly validated (4-32 characters, regex: `^[a-zA-Z0-9 _-]{4,32}$`).
- **Deletion Safety:** Never allow deletion of worlds marked as `IsFavorite`.
- **Build Sync:** Ensure `make build` is run after frontend changes; otherwise, the server will serve an outdated `app.wasm`.
- **No Direct DOM Manipulation:** In the frontend, always use `go-app`'s declarative syntax; avoid direct JS calls unless encapsulated in a helper.

## 7. Operational Commands
- **`make build`**: Compiles the WASM frontend and the Go backend server into the `tmp/` directory.
- **`make run`**: Executes the full build and starts the server locally.
- **`make clean`**: Safely removes build artifacts and temporary files.
- **`make prod`**: Packages the application into a production-ready Docker image.
