# ---- Build stage ----
FROM golang:1.22 AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod ./
RUN go mod download

# Copy source
COPY . .

# Build (static binary, stripped)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./src

# ---- Run stage ----
FROM gcr.io/distroless/base-debian12 AS final
# (alternative: FROM scratch or FROM alpine:latest)

WORKDIR /app

# Copy binary only
COPY --from=builder /app/main .

# Create non-root user
USER nonroot:nonroot

EXPOSE 8080
ENTRYPOINT ["./main"]
