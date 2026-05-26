# Images

---

## `GET /image`

Returns a PNG image. Checks the `Accept` header to determine format — if the client advertises `image/jpeg`, `image/webp`, or `image/svg+xml` it returns that format instead. Falls back to PNG if no supported format is found.

Attribution details for the served image are included as response headers (`X-Attribution-Title`, `X-Attribution-Author`, `X-Attribution-Source`, `X-Attribution-License`, etc.). See [`/image/attribution`](#get-imageattribution) for the full field list.

```sh
# PNG (default)
curl http://localhost/image -o image.png

# Request a specific format via Accept
curl -H "Accept: image/jpeg" http://localhost/image -o image.jpg
curl -H "Accept: image/webp" http://localhost/image -o image.webp
curl -H "Accept: image/svg+xml" http://localhost/image -o image.svg
```

---

## `GET /image/attribution`

Returns a JSON object keyed by image type, where each value contains the full attribution details for that image. Fields marked optional are omitted when not applicable (e.g. `original_work` is only present for derivatives).

```sh
curl http://localhost/image/attribution
```

```json
{
  "png": {
    "title":      "gopher.png",
    "author":     "Takuya Ueda",
    "author_url": "https://twitter.com/tenntenn",
    "source":     "https://github.com/golang-samples/gopher-vector/...",
    "license":    "CC BY 3.0",
    "license_url": "https://creativecommons.org/licenses/by/3.0/",
    "original_work": {
      "title":      "Go gopher",
      "author":     "Renee French",
      "author_url": "https://reneefrench.blogspot.com/",
      "source":     "https://go.dev/blog/gopher",
      "license":    "CC BY 4.0",
      "license_url": "https://creativecommons.org/licenses/by/4.0/"
    }
  },
  "jpeg": { "..." },
  "webp": { "..." }
}
```

| Field | Type | Description |
|---|---|---|
| `title` | string | File or asset title |
| `author` | string | Creator of this file |
| `author_url` | string | (optional) Link to the author's profile |
| `source` | string | Canonical URL of the source file |
| `license` | string | License name |
| `license_url` | string | Link to the full license text |
| `modifications` | string | (optional) Description of changes made to the original |
| `original_work` | object | (optional) Attribution for the upstream work this file derives from; same fields as above minus `modifications` |

---

## `GET /image/{type}`

Returns an image in the explicitly requested format. Ignores the `Accept` header.

### Path parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `type` | string | yes | One of: `png`, `jpeg`, `webp`, `svg` |

Returns `404 Not Found` for unsupported types. Attribution headers are included on all responses (same headers as above).

```sh
curl http://localhost/image/png  -o image.png
curl http://localhost/image/jpeg -o image.jpg
curl http://localhost/image/webp -o image.webp
curl http://localhost/image/svg  -o image.svg
```
