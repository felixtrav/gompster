# HTTP Methods

Each endpoint echoes back the full details of the incoming request as JSON. The method-specific routes enforce the correct HTTP verb — sending the wrong method returns `405 Method Not Allowed`.

---

## `GET /get`

Returns query parameters, headers, origin IP, and request URL.

```sh
curl "http://localhost/get?foo=bar" -H "X-My-Header: hello"
```

---

## `POST /post`

Returns the request body, parsed form fields, headers, and metadata. Supports `application/json`, `application/x-www-form-urlencoded`, and `multipart/form-data`.

```sh
curl -X POST http://localhost/post \
     -H "Content-Type: application/json" \
     -d '{"hello":"world"}'
```

---

## `PUT /put`

Returns the request body, headers, and metadata.

```sh
curl -X PUT http://localhost/put -d 'raw body'
```

---

## `PATCH /patch`

Returns the request body, headers, and metadata.

```sh
curl -X PATCH http://localhost/patch -d 'partial update'
```

---

## `DELETE /delete`

Returns query parameters, headers, and metadata.

```sh
curl -X DELETE "http://localhost/delete?id=123"
```

---

## Response body (all methods)

```json
{
  "args":    { "foo": "bar" },
  "data":    "",
  "files":   {},
  "form":    {},
  "headers": { "X-My-Header": "hello" },
  "json":    null,
  "method":  "GET",
  "origin":  "127.0.0.1",
  "url":     "http://localhost/get?foo=bar"
}
```

| Field | Description |
|---|---|
| `args` | Parsed query string parameters |
| `data` | Raw request body (non-form content types) |
| `files` | Uploaded file names keyed by field name (multipart only) |
| `form` | Parsed form fields (URL-encoded or multipart) |
| `headers` | All request headers |
| `json` | Parsed JSON body (if `Content-Type: application/json` and body is valid JSON) |
| `method` | HTTP method used |
| `origin` | Client IP, respecting `X-Forwarded-For` |
| `url` | Reconstructed full request URL |
