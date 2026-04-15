# Koda B6 Backend - Golang Gin API

REST API backend untuk Coffee Shop Management System menggunakan **Go Gin Framework**, PostgreSQL, dan Redis.

## 📋 Overview

Aplikasi ini menyediakan API lengkap untuk mengelola:
- **User Management**: Registrasi, login, forgot password, reset password
- **Product Management**: CRUD produk, kategori, variant, size, discount, images
- **Cart Management**: Tambah/hapus item, clear cart, lihat summary
- **Order Management**: Buat order, lihat history, update status, order details
- **Reviews & Ratings**: CRUD review, rating produk
- **Payments**: Proses pembayaran, update status
- **Admin Dashboard**: Daily sales statistics

## 🛠 Tech Stack

- **Runtime**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL 13+
- **Cache**: Redis 6+
- **Authentication**: JWT (JSON Web Tokens)
- **Containerization**: Docker & Docker Compose
- **Dependency Injection**: Custom DI Container

## 🏗 Architecture

```
koda-b6-backend/
├── internal/
│   ├── di/              # Dependency Injection Container
│   ├── middleware/      # Authentication & Authorization
│   ├── handler/         # HTTP Request Handlers
│   ├── repository/      # Database Operations
│   ├── service/         # Business Logic
│   └── model/           # Data Models
├── routes/              # Route definitions
├── config/              # Configuration
├── main.go              # Entry point
├── Dockerfile           # Docker build configuration
├── docker-compose.yml   # Docker Compose setup
└── .env                 # Environment variables
```

## 🚀 Setup & Instalasi

### Option 1: Setup Lokal (Development)

#### Prerequisites
- Go 1.21 atau lebih baru
- PostgreSQL 13+
- Redis 6+
- Git

#### Langkah-langkah

1. **Clone repository**
   ```bash
   git clone <repository-url>
   cd koda-b6-backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Setup environment variables**
   Buat file `.env` di root project:
   ```env
   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=coffee_shop
   DATABASE_URL=postgresql://postgres:password@localhost:5432/coffee_shop?sslmode=disable
   
   # Redis
   REDIS_URL=redis://localhost:6379
   REDIS_ADDR=localhost:6379
   
   # Server
   SERVER_PORT=8080
   JWT_SECRET=7x!A%D*F-JaNdRgUkXp2s5v8y/B?E(H+MbQeShVmYq3t6w9z$C&F)J@NcRfUjXn2r
   JWT_EXPIRATION=24h
   ```

4. **Jalankan aplikasi**
   ```bash
   go run main.go
   ```
   Aplikasi akan berjalan di `http://localhost:8080`

5. **Build untuk production**
   ```bash
   go build -o bin/koda-backend .
   ./bin/koda-backend
   ```

### Option 2: Setup dengan Docker Compose

#### Prasyarat
- Node.js backend harus sudah berjalan dengan Docker Compose (untuk sharing database dan redis)
- Network `coffee-network` sudah dibuat oleh Node.js setup

#### Langkah-langkah

1. **Pastikan Node.js backend sudah running**
   ```bash
   # Di folder Node.js project
   docker-compose up -d
   ```
   Ini akan membuat `coffee-network` dan services PostgreSQL + Redis

2. **Build dan jalankan Golang backend**
   ```bash
   # Di folder Golang project
   docker-compose up --build -d
   ```

3. **Verifikasi services**
   ```bash
   docker ps
   ```
   Anda seharusnya melihat:
   - `golang_backend` (port 8081:8080)
   - `backend-js-postgres-1` (port 5432)
   - `backend-js-redis-cache-1` (port 6379)

#### Docker Compose Configuration

Buat file `docker-compose.yml`:

```yaml
version: '3.8'

services:
  backend:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: golang_backend
    ports:
      - "8081:8080"
    environment:
      # Database configuration - UPDATED to match Node.js setup
      - DB_HOST=backend-js-postgres-1
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=password
      - DB_NAME=coffee_shop
      - DATABASE_URL=postgresql://postgres:password@backend-js-postgres-1:5432/coffeeshop?sslmode=disable
      
      # Redis configuration - UPDATED to match Node.js setup
      - REDIS_URL=redis://backend-js-redis-cache-1:6379
      - REDIS_ADDR=backend-js-redis-cache-1:6379
      
      # Server configuration
      - SERVER_PORT=8080
      - JWT_SECRET=7x!A%D*F-JaNdRgUkXp2s5v8y/B?E(H+MbQeShVmYq3t6w9z$C&F)J@NcRfUjXn2r
      - JWT_EXPIRATION=24h
    
    # Use external network dari Node.js setup
    networks:
      - coffee-network

# Use external network dari Node.js setup
networks:
  coffee-network:
    external: true
    name: backend-js_coffee-network
```

**Catatan Penting:**
- PostgreSQL & Redis sudah berjalan dari Node.js Docker Compose
- Network `coffee-network` harus external (dibuat oleh Node.js setup)
- Database hostname: `backend-js-postgres-1` (nama container Node.js)
- Redis hostname: `backend-js-redis-cache-1` (nama container Node.js)

#### Troubleshooting Docker

```bash
# Jika error "network not found"
# Pastikan Node.js container sudah running
docker network ls

# Jika error "cannot connect to database"
docker ps | grep backend-js

# Lihat logs Golang backend
docker logs -f golang_backend

# Stop semua services
docker-compose down
```

## 📊 Environment Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` | `backend-js-postgres-1` (Docker) |
| `DB_PORT` | PostgreSQL port | `5432` | `5432` |
| `DB_USER` | PostgreSQL username | - | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | - | `password` |
| `DB_NAME` | Database name | - | `coffee_shop` |
| `DATABASE_URL` | Full PostgreSQL URL | - | `postgresql://user:pass@host:5432/db` |
| `REDIS_URL` | Redis connection URL | - | `redis://localhost:6379` |
| `REDIS_ADDR` | Redis address (host:port) | - | `localhost:6379` |
| `SERVER_PORT` | Server port | `8080` | `8080` |
| `JWT_SECRET` | JWT signing secret | - | `7x!A%D*F-JaNdRgUkXp2s5v8y/...` |
| `JWT_EXPIRATION` | Token expiration | `24h` | `24h`, `72h` |

## 🔌 API Endpoints

### Authentication Routes
```
POST   /api/auth/register              - Register user baru
POST   /api/auth/login                 - Login & dapatkan JWT token
POST   /api/auth/forgot-password       - Request password reset
POST   /api/auth/reset-password        - Reset password dengan token
```

### User Management Routes
```
GET    /api/users                      - Get semua users
GET    /api/users/:id                  - Get user by ID
POST   /api/users                      - Create user baru
PUT    /api/users/:id                  - Update user
DELETE /api/users/:id                  - Delete user
```

### Product Routes
```
GET    /api/products                   - Get semua produk
GET    /api/products/:id               - Get produk by ID
GET    /api/products/recommended-products  - Get most reviewed products
GET    /api/products/top-products      - Get best selling products
POST   /api/products                   - Create produk baru
PUT    /api/products/:id               - Update produk

GET    /api/products/:id/variants      - Get variants by product
GET    /api/products/:id/sizes         - Get sizes by product
GET    /api/products/:id/discounts     - Get discounts by product
GET    /api/products/:id/images        - Get images by product
POST   /api/products/:id/images        - Upload single image
POST   /api/products/:id/images/multiple  - Upload multiple images
PATCH  /api/products/:id/images/:imageId/set-primary  - Set primary image
DELETE /api/products/:id/images/:imageId  - Delete image
```

### Product Category Routes
```
GET    /api/product-categories         - Get semua kategori
GET    /api/product-categories/:id     - Get kategori by ID
POST   /api/product-categories         - Create kategori baru
PUT    /api/product-categories/:id     - Update kategori
DELETE /api/product-categories/:id     - Delete kategori
```

### Cart Routes (Protected)
```
GET    /cart                           - Get user cart
GET    /cart/summary                   - Get cart summary
POST   /cart                           - Add item ke cart
PUT    /cart/:cart_id                  - Update cart item
DELETE /cart/:cart_id                  - Remove item dari cart
DELETE /cart                           - Clear semua cart items
```

### Order Routes (Protected)
```
POST   /api/orders                     - Create order baru
GET    /api/orders                     - Get user orders
GET    /api/orders/:id                 - Get order by ID
PUT    /api/orders/:id                 - Update order status
DELETE /api/orders/:id                 - Delete order

GET    /api/orders/:id/details         - Get order details
POST   /api/orders/:id/details         - Create order detail
PUT    /api/orders/:id/details/:detail_id  - Update order detail
DELETE /api/orders/:id/details/:detail_id  - Delete order detail
```

### Review Routes
```
GET    /api/reviews                    - Get semua reviews
GET    /api/reviews/:id                - Get review by ID
POST   /api/reviews                    - Create review baru
PUT    /api/reviews/:id                - Update review
```

### Payment Routes
```
POST   /api/payments                   - Create payment
GET    /api/payments/:id               - Get payment by ID
PUT    /api/payments/:id               - Update payment status
DELETE /api/payments/:id               - Delete payment
```

### Variant Routes
```
GET    /api/variants                   - Get semua variants
POST   /api/variants                   - Create variant baru
```

### Size Routes
```
GET    /api/sizes                      - Get semua sizes
POST   /api/sizes                      - Create size baru
```

### Public Routes
```
GET    /public/daily-sales             - Get daily sales statistics
```

## 🔐 Authentication

API menggunakan **JWT (JSON Web Tokens)** untuk authentication.

### Login dan dapatkan Token
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "Login successful"
}
```

### Gunakan Token di Request
Tambahkan header `Authorization: Bearer <token>` ke request:
```bash
curl -X GET http://localhost:8080/api/orders \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## 🏗 Dependency Injection

Project menggunakan custom DI Container di `internal/di` untuk manage dependencies:

```go
// Di main.go
container := di.NewContainer()

// Container auto-wire semua dependencies
userHandler := container.UserHandler()
productHandler := container.ProductHandler()
// dll...
```

## 📦 Routes Organization

Routes diorganisir dengan group di `routes/routes.go`:

```go
api := router.Group("/api")
{
    users := api.Group("/users")
    {
        users.GET("", userHandler.GetAllUsers)
        users.POST("", userHandler.CreateUser)
        // ...
    }
}
```

## 🔍 Troubleshooting

### Development Issues

#### Error: "database connection refused"
```bash
# Pastikan PostgreSQL running
psql -U postgres -h localhost

# Check DB credentials di .env
```

#### Error: "cannot connect to redis"
```bash
# Test Redis connection
redis-cli -h localhost ping
# Should return: PONG
```

#### Error: "port already in use"
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port di .env
SERVER_PORT=9080
```

### Docker Issues

#### Error: "network not found"
```bash
# Pastikan Node.js backend running
docker ps | grep backend-js

# Create network if needed
docker network create coffee-network
```

#### Error: "cannot connect to backend-js-postgres-1"
```bash
# Verify Node.js services
docker logs backend-js-postgres-1

# Test connectivity
docker exec golang_backend ping backend-js-postgres-1
```

#### View Docker Logs
```bash
# Golang backend logs
docker logs -f golang_backend

# PostgreSQL logs
docker logs backend-js-postgres-1

# Redis logs
docker logs backend-js-redis-cache-1
```

### JWT & Authentication Issues

#### Error: "JWT token invalid"
- Pastikan `JWT_SECRET` sama di semua instance
- Token sudah expired? Cek `JWT_EXPIRATION`
- Verify dengan https://jwt.io

#### Error: "Unauthorized"
- Format header: `Authorization: Bearer <token>`
- Check token masih valid
- Login lagi untuk generate token baru

## 📝 Development Tips

### Run with hot reload
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run dengan auto-reload
air
```

### Database migrations
```bash
# Seed database dengan data awal
# Update connection string dan jalankan seed scripts
```

## 📚 Dokumentasi Lengkap

Untuk dokumentasi API lengkap, gunakan Swagger/OpenAPI dokumentation yang tersedia di:
- `/docs/swagger` 

## 👥 Contributing

1. Create feature branch: `git checkout -b feature/amazing-feature`
2. Commit changes: `git commit -m 'Add amazing feature'`
3. Push ke branch: `git push origin feature/amazing-feature`
4. Open Pull Request


