# Responses

---

## `GET /status/{code}`

Returns a response with the given HTTP status code. Useful for testing how your client handles specific codes.

### Path parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `code` | integer | yes | HTTP status code to return (100–599) |

Returns `400 Bad Request` if the code is outside the valid range.

```sh
curl -i http://localhost/status/418
curl -i http://localhost/status/503
```

### Response body

```json
{
  "code": 418,
  "description": "I'm a teapot"
}
```

---

## `GET /response-headers`

Promotes query parameters to response headers and also echoes them in the JSON body. Multiple pairs can be set in a single request.

### Query parameters

Any key=value pair is accepted. The key becomes the header name and the value becomes the header value. Only the first value is used if a key appears multiple times.

```sh
curl -i "http://localhost/response-headers?X-My-Header=hello&Cache-Control=no-store"
```

### Response

```
HTTP/1.1 200 OK
X-My-Header: hello
Cache-Control: no-store
Content-Type: application/json
```

```json
{
  "X-My-Header": "hello",
  "Cache-Control": "no-store"
}
```
