# gompster (go-dumpster)

A lightweight, Go-based, open source HTTP debugging server in the spirit of httpbin — useful for testing clients, inspecting requests, and simulating a variety of HTTP behaviours.

---

## Quick start

```bash
docker run -p 80:80 felixtrav/gompster
```

Open **http://localhost/** in your browser to reach the Swagger UI. All available endpoints are documented there; individual endpoints are not covered in this README.

---

## Configuration

All parameters can be set via environment variable or CLI flag. Flags take precedence over environment variables, which take precedence over defaults.

| Flag | Env var | Default | Description |
|---|---|---|---|
| `-host` | `HOST` | `0.0.0.0` | Listener host address |
| `-port` | `PORT` | `80` | Listener port |
| `-read-timeout` | `READ_TIMEOUT` | `30s` | HTTP read timeout |
| `-write-timeout` | `WRITE_TIMEOUT` | `30s` | HTTP write timeout |
| `-idle-timeout` | `IDLE_TIMEOUT` | `120s` | HTTP idle (keep-alive) timeout |
| `-log` | `LOG_REQUESTS` | `true` | Log incoming requests to stdout |

Duration values use Go's duration syntax (`30s`, `1m`, `500ms`, etc.).

```bash
# Example: listen on port 8080, disable request logging
docker run -p 8080:80 -e PORT=8080 -e LOG_REQUESTS=false felixtrav/gompster
```

---

## Docker images

Three image variants are built from `docker/Dockerfile`:

| Variant | Base | Runs as | Use when |
|---|---|---|---|
| `distroless` | `scratch` | UID 65532 (nonroot) | Production — smallest possible image, no shell or package manager |
| `alpine-rootless` | `alpine:3.23.4` | Dedicated non-root service account | Debugging — has `apk` but still runs unprivileged |
| `alpine` | `alpine:3.23.4` | root | Building derived images that need to install extra packages in a later layer |

---

## Healthcheck

All three image variants embed a purpose-built healthcheck binary (`cmd/healthcheck/main.go`). It:

- Makes a `GET http://127.0.0.1:{PORT}/health` request with a 5 s timeout.
- Exits `0` if the response is HTTP 200; exits `1` on any connection error or non-200 status.
- Reads `PORT` from the environment (default: `80`) so it stays in sync with the server automatically.

The binary exists because distroless images ship no `curl` or `wget`. It is wired to the Docker `HEALTHCHECK` instruction in all three image variants.

---

## Tests

Tests live in `tests/` and use Go's `httptest` package — no external server or Docker image is needed.

```bash
# Run all tests
go test ./tests/

# Verbose output
go test -v ./tests/
```

Each test file spins up a real in-process HTTP server via `httptest.NewServer`, runs its assertions against it, then tears it down when the test binary exits. The server setup and shared `client`/`redir` helpers live in `tests/main_test.go`.
