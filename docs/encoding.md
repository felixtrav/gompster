# Encoding

---

## `GET /base64/{value}`

Decodes a Base64-encoded path segment and returns the raw decoded bytes. Supports both standard and URL-safe Base64 alphabets, with or without padding.

### Path parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `value` | string | yes | Base64-encoded string (URL-safe or standard) |

Returns `400 Bad Request` if the value cannot be decoded.

### Response

```
Content-Type: application/octet-stream
```

The response body is the raw decoded bytes. If the decoded content is valid text you can pipe or redirect it normally.

```sh
# Encode a string and decode it back
echo -n 'Hello, gompster!' | base64
# SGVsbG8sIGdvbXBzdGVyIQ==

curl http://localhost/base64/SGVsbG8sIGdvbXBzdGVyIQ==
# Hello, gompster!

# URL-safe Base64 (- instead of +, _ instead of /)
curl "http://localhost/base64/SGVsbG8sIGdvbXBzdGVyIQ"
```
