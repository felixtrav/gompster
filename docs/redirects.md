# Redirects

---

## `GET /redirect/{n}`

Issues `n` sequential relative redirects before returning `200 OK`. Each intermediate step redirects to `/redirect/{n-1}`. The final redirect (`n=1`) redirects to `/get`.

### Path parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `n` | integer | yes | Number of redirects to chain (1 – 100) |

Returns `400 Bad Request` if `n` is outside the allowed range.

```sh
# Follow 3 hops
curl -L http://localhost/redirect/3

# Observe each hop without following
curl -i http://localhost/redirect/3
```

---

## `ANY /redirect-to`

Redirects the client to an arbitrary URL using the status code of your choice. Accepts all HTTP methods.

### Query parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `url` | string | yes | Target URL to redirect to |
| `status_code` | integer | no | HTTP status code (default: `302`) |

Returns `400 Bad Request` if `url` is missing.

```sh
curl -L "http://localhost/redirect-to?url=https://example.com&status_code=301"
```

---

## `GET /absolute-redirect/{n}`

Issues `n` sequential absolute redirects before returning `200 OK`. Unlike `/redirect/{n}`, each intermediate step includes the full `http://host/...` URL in the `Location` header.

### Path parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `n` | integer | yes | Number of redirects to chain (1 – 100) |

Returns `400 Bad Request` if `n` is outside the allowed range.

```sh
curl -L http://localhost/absolute-redirect/2
```
