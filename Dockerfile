# ============================================================================
# STAGE 1: BUILDER
# ============================================================================
FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# ✅ IMPROVEMENT: Copy go.mod/go.sum dulu (cache optimization)
COPY go.mod go.sum .

RUN go mod tidy

# ✅ IMPROVEMENT: Copy source code setelah mod tidy
COPY . .

RUN go build -o backend-coffeeshop main.go

# ✅ Tidak perlu chmod (binary sudah executable)

# ============================================================================
# STAGE 2: RUNTIME
# ============================================================================
FROM alpine:latest

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /workspace/backend-coffeeshop /app/

EXPOSE 3002

# ✅ FIX: Path binary yang benar
ENTRYPOINT ["./backend-coffeeshop"]