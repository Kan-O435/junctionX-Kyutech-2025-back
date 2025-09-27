# 公式の最新Goイメージ（1.25以上）を指定
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN go build -o server ./cmd/server

# 実行用イメージ（軽量化のためalpineなど）
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
