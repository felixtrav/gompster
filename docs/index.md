# gompster — Endpoint Reference

Browse the full interactive spec at `http://localhost/` (Swagger UI) or `http://localhost/openapi.json`.

| Category | File | Endpoints |
|---|---|---|
| Health | [health.md](health.md) | `/health` |
| HTTP Methods | [http-methods.md](http-methods.md) | `/get` `/post` `/put` `/patch` `/delete` |
| Anything | [anything.md](anything.md) | `/anything` `/anything/*` |
| Request Inspection | [request-inspection.md](request-inspection.md) | `/headers` `/ip` `/ip/raw` `/user-agent` |
| Responses | [responses.md](responses.md) | `/status/{code}` `/response-headers` |
| Response Formats | [response-formats.md](response-formats.md) | `/gzip` `/deflate` `/brotli` `/utf8` `/html` `/json` `/xml` `/robots.txt` `/deny` |
| Dynamic Data | [dynamic.md](dynamic.md) | `/bytes/{n}` `/delay/{seconds}` |
| UUIDs | [uuid.md](uuid.md) | `/uuid` `/uuid/v1` `/uuid/v3` `/uuid/v4` `/uuid/v5` `/uuid/v7` |
| Cookies | [cookies.md](cookies.md) | `/cookies` `/cookies/set` `/cookies/delete` |
| Redirects | [redirects.md](redirects.md) | `/redirect/{n}` `/redirect-to` `/absolute-redirect/{n}` |
| Auth | [auth.md](auth.md) | `/basic-auth/{user}/{passwd}` `/bearer` |
| Streaming | [streaming.md](streaming.md) | `/stream/{n}` |
| Encoding | [encoding.md](encoding.md) | `/base64/{value}` |
| Images | [images.md](images.md) | `/image` `/image/attribution` `/image/{type}` |
