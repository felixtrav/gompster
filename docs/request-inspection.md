# Request Inspection

---

## `GET /headers`

Returns a map of all HTTP headers sent by the client.

```sh
curl http://localhost/headers -H "X-Custom: value"
```

```json
{
  "headers": {
    "Accept":    "*/*",
    "X-Custom":  "value",
    "User-Agent": "curl/8.0"
  }
}
```

---

## `GET /ip`

Returns the client's IP address. Respects `X-Forwarded-For` and `X-Real-IP` proxy headers — the first address in `X-Forwarded-For` is returned if present.

```sh
curl http://localhost/ip
```

```json
{ "origin": "203.0.113.42" }
```

---

## `GET /ip/raw`

Returns `r.RemoteAddr` — the direct TCP peer address — with the port stripped.

> **Note:** The chi `RealIP` middleware rewrites `r.RemoteAddr` globally with the value from forwarding headers before any handler runs. In practice this endpoint returns the address of the last proxy in the chain, not the underlying TCP socket, unless no forwarding headers are present.

```sh
curl http://localhost/ip/raw
```

```json
{ "origin": "127.0.0.1" }
```

---

## `GET /user-agent`

Returns the `User-Agent` header sent by the client.

```sh
curl http://localhost/user-agent
```

```json
{ "user-agent": "curl/8.0" }
```
