# Build stage
FROM golang:1.27.3-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build Flux
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /flux ./cmd/flux

# Runtime stage
FROM alpine:3.24

WORKDIR /app

# Run as a non-root user
RUN addgroup -S flux && adduser -S flux -G flux

COPY --from=builder /flux /app/flux

USER flux

ENTRYPOINT ["/app/flux"]
