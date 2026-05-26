# Cookies

---

## `GET /cookies`

Returns all cookies currently sent by the client.

```sh
curl http://localhost/cookies --cookie "session=abc123; theme=dark"
```

```json
{
  "cookies": {
    "session": "abc123",
    "theme":   "dark"
  }
}
```

---

## `GET /cookies/set`

Sets one or more cookies via query parameters, then redirects to `/cookies` so you can see the result.

### Query parameters

Any key=value pair is accepted. Each pair becomes a `Set-Cookie` header in the `302` redirect response.

```sh
# Follow the redirect to confirm cookies were set
curl -c jar.txt -b jar.txt -L "http://localhost/cookies/set?session=abc123&theme=dark"
```

```json
{
  "cookies": {
    "session": "abc123",
    "theme":   "dark"
  }
}
```

---

## `GET /cookies/delete`

Expires one or more cookies by name via query parameters, then redirects to `/cookies`.

### Query parameters

Any key=value pair is accepted. Only the key matters — the value is ignored. Each named cookie is expired by setting `Max-Age=0` in a `Set-Cookie` header.

```sh
curl -c jar.txt -b jar.txt -L "http://localhost/cookies/delete?session&theme"
```

After the redirect the deleted cookies no longer appear in the `/cookies` response.
