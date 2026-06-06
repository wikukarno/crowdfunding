# Build stage
FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server .

# Run stage
FROM alpine:3.20

RUN adduser -D -g '' app
WORKDIR /app

COPY --from=build /out/server /app/server
COPY migrations ./migrations

# Uploaded avatars and campaign images land here, so it must be writable by the
# non-root user the app runs as.
RUN mkdir -p /app/images/campaigns && chown -R app /app/images

USER app
EXPOSE 8080

ENTRYPOINT ["/app/server"]
