# Auth

---

## `GET /basic-auth/{user}/{passwd}`

Challenges the client with HTTP Basic Authentication. Returns `401 Unauthorized` with a `WWW-Authenticate: Basic` challenge if credentials are absent or incorrect. Returns `200 OK` with the authenticated user if the credentials match the path parameters.

### Path parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `user` | string | yes | Expected username |
| `passwd` | string | yes | Expected password |

```sh
# Correct credentials
curl -u myuser:mypass http://localhost/basic-auth/myuser/mypass

# Missing credentials — returns 401
curl -i http://localhost/basic-auth/myuser/mypass
```

### Response `200 OK`

```json
{
  "authenticated": true,
  "user": "myuser"
}
```

---

## `GET /bearer`

Validates a Bearer token in the `Authorization` header. Returns `401 Unauthorized` if the header is absent or does not start with `Bearer `. Any non-empty token value is accepted — the token is not verified against a secret.

```sh
# Valid bearer token
curl -H "Authorization: Bearer mytoken123" http://localhost/bearer

# Missing header — returns 401
curl -i http://localhost/bearer
```

### Response `200 OK`

```json
{
  "authenticated": true,
  "token": "mytoken123"
}
```
