# Health

## `GET /health`

Returns the server's liveness status and uptime. Intended for use with Docker and Kubernetes health probes.

### Response `200 OK`

```json
{
  "status": "ok",
  "uptime_seconds": 42.7
}
```

| Field | Type | Description |
|---|---|---|
| `status` | string | Always `"ok"` while the server is running |
| `uptime_seconds` | number | Seconds since the process started |

### Example

```sh
curl http://localhost/health
```
