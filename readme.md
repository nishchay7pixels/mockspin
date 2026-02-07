# MockSpin 🚀

<img width="1536" height="1024" alt="image" src="https://github.com/user-attachments/assets/94bbdabb-b4f9-4e55-9ae6-e846b8a21333" />


**OpenAPI → Real Dummy Backend (Offline, CRUD, Zero Boilerplate)**

MockSpin is a CLI tool that turns an **OpenAPI 3.x specification** into a **fully working local dummy backend** with real CRUD-like behavior — no manual coding, no cloud dependency, no heavy frameworks.

Unlike traditional mock servers, MockSpin is built to:
- run fully **offline**
- behave like a **real CRUD backend**
- be **OpenAPI-driven**
- require **zero handwritten backend code**

---

## Why MockSpin?

Most existing mock servers validate contracts but fail to behave like a real backend:

| Capability | Traditional Mocks | MockSpin |
|-----------|------------------|----------|
| POST echoes input | ❌ | ✅ |
| GET reflects POST | ❌ | ✅ |
| In-memory persistence | ❌ | ✅ |
| Offline-first | ⚠️ | ✅ |
| OpenAPI as source of truth | ✅ | ✅ |

MockSpin exists to close this gap.

---

## Core Philosophy

- **OpenAPI is the source of truth**
- **Runtime-first (no code generation required)**
- **Offline-first and lightweight**
- **Predictable, debuggable behavior**
- **Frontend-friendly dummy backend**

---

## Current Status

🚧 **Active development**

- **Phase 1**: ✅ Implemented  
- **Phase 2A**: 🔨 Under development  
- **Phase 3+**: 🧪 Planned  

---

## Tech Stack 

| Area | Choice | Why | 
| -------- | ------------- | --------------------------- | 
| Language | Go | Single binary, fast startup | 
| Router | chi | Lightweight, idiomatic | 
| OpenAPI | kin-openapi | Mature, schema validation | 
| Storage | In-memory map | Fast, predictable | 
| UUIDs | google/uuid | Standard |

--- 

## What MockSpin does (MVP)

### Input
- OpenAPI 3.x (`openapi.yaml` / `openapi.json`)

### Output
- A running local HTTP server that:
  - exposes all paths from the spec
  - validates requests against schemas
  - persists data in memory
  - returns schema-compliant responses

### Example
```bash
mockspin up --spec openapi.yaml --port 4010
POST /users        → stores user
GET  /users        → returns stored users
GET  /users/{id}   → returns stored user
```
## Prerequisites

MockSpin is designed to be lightweight and offline-first, but it assumes a few system-level tools are already installed.

### Required
- **Docker**
  - Used for runtime isolation and packaging
  - Docker Desktop (macOS / Windows) or Docker Engine (Linux)
- **Free local port**
  - Default: `4010` (configurable via `--port`)

### Supported Platforms
- macOS (amd64 / arm64)
- Linux (amd64 / arm64)
- Windows via WSL2

## Install

MockSpin is distributed as a Go CLI binary.

### Option 1: Install using Go (recommended)

Requires **Go 1.22+**.

```bash
go install github.com/nishchay7pixels/mockspin@latest
```
This installs the mockspin binary into:
```bash
$HOME/go/bin
```
Add mockspin to your **PATH**
If mockspin is not found after installation, add the Go bin directory to your PATH.

macOS / Linux (bash / zsh)

Add this line to ~/.bashrc, ~/.zshrc, or ~/.profile:
```bash
export PATH="$HOME/go/bin:$PATH"
```
Reload your shell:
```bash
source ~/.zshrc
```
Verify:
```bash
mockspin --help
```
### Option 2: Download prebuilt binary (coming soon)

Prebuilt binaries for macOS and Linux will be available via GitHub Releases.

**Until then, installing via Go is the recommended method.**

## Build locally
```bash
go mod tidy
go build -o mockspin .
./mockspin doctor
./mockspin up --spec ./openapi.yaml --port 4010
# Ctrl+C stops container and removes session
./mockspin status
./mockspin stop
```

# Roadmap

## Phase 1 — CLI & Runtime Infrastructure (Implemented)
### Goal

Build a reliable CLI and runtime shell that can:
- start a local server from an OpenAPI spec - manage its lifecycle 
- surface errors clearly 
- behave consistently across platforms

> Phase 1 builds the engine, not the backend logic.
### Phase 1 Features
#### CLI Commands
```bash
mockspin up --spec openapi.yaml --port 4010
mockspin stop
mockspin status
mockspin doctor
```

**Note - Validate prerequisites using**
```bash
mockspin doctor
```

What it does:
- check Docker is installed
- check Docker daemon is running
- check architecture compatibility
- check required ports are free
- print actionable errors, not stack traces

Example output:

```bash
❌ Docker not found
→ Install Docker Desktop from https://docker.com

⚠️ Port 4010 already in use
→ Use --port to specify a different port
```
### Responsibilities

- CLI argument parsing

- Spec file existence & basic validation

- Runtime startup and shutdown

- Foreground and detached execution

- Session tracking

- Signal handling (Ctrl+C cleanup)

- Crash detection and log surfacing

- Cross-platform compatibility (amd64 / arm64)

### What Phase 1 does NOT do

❌ CRUD persistence

❌ Request/response schema enforcement

❌ Dummy data generation

❌ Code generation

---

## Phase 2 — Runtime Dummy Backend (Under Development)
### Goal

Replace external mock engines with a native Go HTTP server that behaves like a real backend.

#### Planned Features
**✅ CRUD behavior**
```
POST /resource → create + persist

GET /resource → list

GET /resource/{id} → fetch

PUT /resource/{id} → update

DELETE /resource/{id} → delete
```

**✅ Schema-based validation**
```
Required field enforcement

Enum checks

Format checks (email, uuid, date-time)

Invalid requests return 400 / 422
```

**✅ Automatic server-generated fields**
```
If missing from request:

id (UUID)

createdAt

updatedAt
```
**✅ In-memory store**
```
Fast

Isolated per run

Reset on restart (by design)
```

## Phase 3 (Planned)

- Optional persistence (file / SQLite)

- Relationship awareness

- Seed data generation (offline AI)

## Phase 4 (Planned)

- Code generation mode

- Docker image output

- WireMock compatibility

# Example workflow
```bash
# Start backend
# refer openapi.yaml in ./example
mockspin up --spec ./example/openapi.yaml 

# Create data
curl -X POST http://localhost:4010/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Nishchya","email":"nishchya@example.com"}'

# Fetch data
curl http://localhost:4010/users
```
# Feedback Wanted

MockSpin is in early development.

If you try it and have thoughts on:
- CLI UX
- error messages
- missing commands
- developer experience

please open an issue or start a discussion.

Honest feedback is more valuable than stars.

# License 
###  MIT 

# Status 🚧 
**Active development** 
Phase 2 in progress.
