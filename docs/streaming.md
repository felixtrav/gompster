# Streaming

---

## `GET /stream/{n}`

Streams `n` newline-delimited JSON objects (NDJSON), one per line, written and flushed individually. Useful for testing clients that consume chunked/streaming responses.

### Path parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `n` | integer | yes | Number of JSON lines to stream (1 – 100) |

Returns `400 Bad Request` if `n` is outside the allowed range.

### Response

```
Content-Type: application/x-ndjson
Transfer-Encoding: chunked
```

Each flushed line is a JSON object:

```json
{"id":0,"origin":"127.0.0.1","headers":{"Accept":"*/*"},"url":"http://localhost/stream/5"}
{"id":1,"origin":"127.0.0.1","headers":{"Accept":"*/*"},"url":"http://localhost/stream/5"}
{"id":2,"origin":"127.0.0.1","headers":{"Accept":"*/*"},"url":"http://localhost/stream/5"}
{"id":3,"origin":"127.0.0.1","headers":{"Accept":"*/*"},"url":"http://localhost/stream/5"}
{"id":4,"origin":"127.0.0.1","headers":{"Accept":"*/*"},"url":"http://localhost/stream/5"}
```

If the client disconnects before all lines are written, the handler stops immediately rather than continuing to write.

```sh
# Stream 10 lines and display them as they arrive
curl -N http://localhost/stream/10

# Simulate a slow consumer
curl -N http://localhost/stream/5 | while IFS= read -r line; do
  echo "$line"
  sleep 0.5
done
```
