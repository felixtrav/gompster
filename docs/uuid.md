# UUIDs

All UUID endpoints return a JSON object with a single `uuid` field.

```json
{ "uuid": "550e8400-e29b-41d4-a716-446655440000" }
```

---

## `GET /uuid`

Returns a random version 4 UUID. Equivalent to `GET /uuid/v4`.

```sh
curl http://localhost/uuid
```

---

## `GET /uuid/v1`

Returns a version 1 UUID (time-based with MAC address). The node ID is randomly generated per process start.

```sh
curl http://localhost/uuid/v1
```

---

## `GET /uuid/v3`

Returns a version 3 UUID (MD5-hashed namespace+name).

### Query parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `namespace` | string | yes | One of: `dns`, `url`, `oid`, `x500` |
| `name` | string | yes | The name to hash |

```sh
curl "http://localhost/uuid/v3?namespace=dns&name=example.com"
```

---

## `GET /uuid/v4`

Returns a randomly generated version 4 UUID.

```sh
curl http://localhost/uuid/v4
```

---

## `GET /uuid/v5`

Returns a version 5 UUID (SHA-1-hashed namespace+name).

### Query parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `namespace` | string | yes | One of: `dns`, `url`, `oid`, `x500` |
| `name` | string | yes | The name to hash |

```sh
curl "http://localhost/uuid/v5?namespace=url&name=https://example.com"
```

---

## `GET /uuid/v7`

Returns a version 7 UUID (Unix timestamp + random bits). Monotonically increasing within a millisecond.

```sh
curl http://localhost/uuid/v7
```
