# ─── Builder ──────────────────────────────────────────────────────────────────
# Produces a fully static binary; CGO is disabled so the binary runs on any
# Linux base including distroless/static (no libc required).
FROM golang:1.26.3-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
        -trimpath -ldflags="-s -w" \
        -o /out/gompster . \
    && CGO_ENABLED=0 GOOS=linux go build \
        -trimpath -ldflags="-s -w" \
        -o /out/healthcheck ./cmd/healthcheck

# ─── distroless ───────────────────────────────────────────────────────────────
# Scratch-based image — contains only the two statically linked binaries.
# No shell, no package manager, no libc. Smallest possible attack surface.
# UID 65532 mirrors the conventional nonroot user ID used by distroless images.
FROM scratch AS distroless
COPY --from=builder /out/gompster    /gompster
COPY --from=builder /out/healthcheck /healthcheck

EXPOSE 80
USER 65532:65532
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["/healthcheck"]
ENTRYPOINT ["/gompster"]

# ─── alpine-rootless ──────────────────────────────────────────────────────────
# Alpine with a dedicated non-root service account.
# Larger than distroless but includes a shell and apk for ad-hoc debugging
# without requiring root access.
FROM alpine:3.23.4 AS alpine-rootless
RUN addgroup -S gompster \
    && adduser  -S -G gompster gompster

COPY --from=builder /out/gompster    /usr/local/bin/gompster
COPY --from=builder /out/healthcheck /usr/local/bin/healthcheck

EXPOSE 80
USER gompster
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["healthcheck"]
ENTRYPOINT ["gompster"]

# ─── alpine ───────────────────────────────────────────────────────────────────
# Alpine but running as root
FROM alpine:3.23.4 AS alpine

COPY --from=builder /out/gompster    /usr/local/bin/gompster
COPY --from=builder /out/healthcheck /usr/local/bin/healthcheck

EXPOSE 80
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["healthcheck"]
ENTRYPOINT ["gompster"]
