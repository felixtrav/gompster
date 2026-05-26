# Dynamic Data

---

## `GET /bytes/{n}`

Returns `n` cryptographically random bytes as `application/octet-stream`. Capped at 100 MB — values above the cap are silently clamped.

### Path parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `n` | integer | yes | Number of bytes to return (0 – 104857600) |

```sh
# Save 1 KB of random bytes to a file
curl http://localhost/bytes/1024 -o random.bin

# Time how long it takes to receive 1 MB
time curl -s http://localhost/bytes/1048576 -o /dev/null
```

---

## `GET /delay/{seconds}`

Waits the given number of seconds before responding. Accepts fractional values. Capped at 300 seconds — values above the cap are clamped. If the client disconnects before the delay expires, the handler returns immediately.

### Path parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `seconds` | number | yes | Delay in seconds, float (0 – 300) |

The response body is the same echo payload as `/get`.

```sh
# 1.5 second delay
curl http://localhost/delay/1.5

# Test a client timeout shorter than the server delay
curl --max-time 2 http://localhost/delay/10
```

> **WriteTimeout interaction:** If the server's `--write-timeout` flag is shorter than the requested delay, the connection will be closed before the handler responds. The default is 30 seconds.
