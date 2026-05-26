# Anything

Catch-all endpoints that echo the full request for **any HTTP method**. Useful when you want to send an arbitrary method without needing a method-specific route.

---

## `ANY /anything`

Accepts GET, POST, PUT, PATCH, DELETE (and any other method). Returns the same echo payload as the method-specific endpoints.

```sh
curl -X POST http://localhost/anything -d 'hello'
curl -X DELETE http://localhost/anything
```

---

## `ANY /anything/{path}`

Accepts any method at any sub-path. The path itself is not included in the response body — it is available via the `url` field.

```sh
curl http://localhost/anything/some/nested/path?debug=true
```

---

## Response body

Same structure as the [HTTP Methods](http-methods.md) echo response.

```json
{
  "args":    {},
  "data":    "hello",
  "files":   {},
  "form":    {},
  "headers": { "Content-Length": "5" },
  "json":    null,
  "method":  "POST",
  "origin":  "127.0.0.1",
  "url":     "http://localhost/anything"
}
```
