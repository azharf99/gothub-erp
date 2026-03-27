# ==========================================
# TAHAP 1: BUILDER (Membangun Aplikasi)
# ==========================================
# Gunakan image Golang resmi versi ringan
FROM golang:1.26-alpine AS builder

# Set direktori kerja di dalam container
WORKDIR /app

# Salin file go.mod dan go.sum lebih dulu untuk caching dependency
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode sumber
COPY . .

# Build aplikasi Go menjadi file binary bernama 'main'
# CGO_ENABLED=0 memastikan binary bisa berjalan mandiri tanpa library C eksternal
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# ==========================================
# TAHAP 2: RUNNER (Menjalankan Aplikasi)
# ==========================================
FROM alpine:latest

# >>> TAMBAHKAN BARIS INI <<<
# Menginstal paket zona waktu dunia ke dalam Alpine
RUN apk add --no-cache tzdata

WORKDIR /app

# Salin HANYA file binary 'main' dari tahap builder sebelumnya
COPY --from=builder /app/main .

# Expose port yang digunakan aplikasi
EXPOSE 8080

# Perintah wajib saat container dijalankan
CMD ["./main"]