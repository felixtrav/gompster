# Response Formats

---

## `GET /gzip`

Returns a gzip-compressed JSON body.

```
Content-Encoding: gzip
Content-Type: application/json
```

```sh
curl --compressed http://localhost/gzip
```

---

## `GET /deflate`

Returns a deflate-compressed JSON body.

```
Content-Encoding: deflate
Content-Type: application/json
```

```sh
curl --compressed http://localhost/deflate
```

---

## `GET /brotli`

Returns a brotli-compressed JSON body.

```
Content-Encoding: br
Content-Type: application/json
```

```sh
curl --compressed http://localhost/brotli
```

The compressed endpoints all return the same payload shape:

```json
{
  "project":     "gompster",
  "description": "A lightweight, Go-based, open source HTTP debugging server.",
  "encoding":    "gzip",
  "origin":      "127.0.0.1",
  "features":    [ "..." ]
}
```

---

## `GET /utf8`

Returns a UTF-8 encoded plain text body containing a broad sample of Unicode characters, useful for testing character encoding handling.

```
Content-Type: text/plain; charset=utf-8
```

```sh
curl http://localhost/utf8
```

---

## `GET /html`

Returns a simple HTML document.

```
Content-Type: text/html; charset=utf-8
```

```sh
curl http://localhost/html
```

---

## `GET /json`

Returns a sample JSON document describing gompster.

```
Content-Type: application/json
```

```sh
curl http://localhost/json
```

---

## `GET /xml`

Returns a sample XML document describing gompster.

```
Content-Type: application/xml
```

```sh
curl http://localhost/xml
```

---

## `GET /robots.txt`

Returns a `robots.txt` that disallows `/deny`.

```
User-agent: *
Disallow: /deny
```

```sh
curl http://localhost/robots.txt
```

---

## `GET /deny`

Returns a plain text page that is disallowed by `/robots.txt`. The response is `200 OK` — the restriction exists only at the crawler protocol level.

```sh
curl http://localhost/deny
```
