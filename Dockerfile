# syntax=docker/dockerfile:1
# GitPulse — https://github.com/sharanyashwant27-tech/GitPulse-CLI-GO
# Targets: runtime (CLI), docs (README site on :8098)

FROM golang:1.25-bookworm AS builder

LABEL org.opencontainers.image.source="https://github.com/sharanyashwant27-tech/GitPulse-CLI-GO" \
      org.opencontainers.image.title="GitPulse" \
      org.opencontainers.image.description="Interactive Git repository analytics CLI"

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/gitpulse .

FROM alpine:3.21 AS runtime

LABEL org.opencontainers.image.source="https://github.com/sharanyashwant27-tech/GitPulse-CLI-GO" \
      org.opencontainers.image.title="GitPulse CLI" \
      org.opencontainers.image.description="GitPulse runtime image"

RUN apk add --no-cache ca-certificates git tini \
    && adduser -D -u 1000 gitpulse

WORKDIR /app
COPY --from=builder /out/gitpulse /usr/local/bin/gitpulse
COPY templates /app/templates
COPY assets /app/assets
COPY README.md /app/README.md

USER gitpulse
ENTRYPOINT ["/sbin/tini", "--", "gitpulse"]
CMD ["--help"]

# Documentation / README static site (localhost:8098)
FROM nginx:1.27-alpine AS docs

LABEL org.opencontainers.image.source="https://github.com/sharanyashwant27-tech/GitPulse-CLI-GO" \
      org.opencontainers.image.title="GitPulse Docs" \
      org.opencontainers.image.description="GitPulse README docs site on port 8098"

RUN apk add --no-cache git \
    && rm -rf /usr/share/nginx/html/*

COPY README.md /usr/share/nginx/html/README.md
COPY assets/ /usr/share/nginx/html/assets/
COPY templates/ /usr/share/nginx/html/templates/
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY docker/index.html /usr/share/nginx/html/index.html

EXPOSE 8098
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8098/ >/dev/null || exit 1

CMD ["nginx", "-g", "daemon off;"]
