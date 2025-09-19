# Stage 1: build
# Stage 1: build
FROM golang:1.22 AS builder

WORKDIR /app

# Download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

COPY .env .env

# Build your app (main.go is in root)
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Stage 2: run
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
