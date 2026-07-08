# go-live-server

A live-reload HTTP server for static web development. Uses inotify to watch the filesystem and SSE (Server-Sent Events) to notify the browser of changes.

## Features

- Automatic browser reload on file changes
- Auto-injects the reload script into HTML pages
- Ignores `node_modules`, `.git`, `vendor`, `dist`, and hidden directories
- Event debouncing (100ms) to prevent multiple reloads
- Graceful shutdown

## Installation

```bash
git clone <repo-url>
cd go-live-server
go build -o main .
```

**Requires Go 1.26+ and Linux** (uses inotify, does not work on macOS/Windows).

## Usage

```bash
./main --dir ./my-site --port 8080
```

Or without building:

```bash
go run . --dir ./my-site --port 3000
```

### Flags

| Flag   | Default | Description            |
|--------|---------|------------------------|
| `--dir`  | `./`    | Directory to serve     |
| `--port` | `5000`  | HTTP server port       |

## How it works

1. A **watcher** uses inotify to monitor the entire directory tree
2. Changes are **debounced** (100ms) to coalesce nearby events
3. An **SSE hub** (`/__live`) notifies all connected clients
4. The server injects a `<script>` into HTML pages that opens an `EventSource` to `/__live`
5. On a `"reload"` event, the browser calls `location.reload()`

## Limitations

- Linux only (inotify)
- No selective CSS/JS reload (always full page reload)
- Hardcoded ignore patterns (not configurable via CLI)
- Port must be >= 1024
