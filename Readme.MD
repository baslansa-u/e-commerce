# E-Commerce Backend (Go)

REST API สำหรับระบบ e-commerce เขียนด้วย Go + Gin + GORM + PostgreSQL

## Tech Stack

- **Language**: Go 1.22
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Auth**: JWT

## Project Structure

```
ecommerce/
├── cmd/
│   └── main.go
├── internal/
│   ├── auth/
│   ├── product/
│   └── order/
├── pkg/
│   ├── middleware/
│   ├── database/
│   └── response/
├── docker-compose.yml
├── .env.example
└── go.mod
```

## Getting Started

### 1. Clone repo

```bash
git clone https://github.com/baslansa-u/e-commerce.git
cd ecommerce
```

### 2. Setup environment

```bash
cp .env.example .env
# แก้ไขค่าใน .env ให้ตรงกับเครื่องของคุณ
```

### 3. Start PostgreSQL

```bash
docker compose up -d
```

### 4. Run server

```bash
go mod tidy
go run cmd/main.go
```

Server จะเริ่มที่ `http://localhost:8080`

---

## API Endpoints

### Auth
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/auth/register` | ❌ | สมัครสมาชิก |
| POST | `/api/auth/login` | ❌ | เข้าสู่ระบบ |

### Products
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/products` | ❌ | ดูสินค้าทั้งหมด |
| GET | `/api/products/:id` | ❌ | ดูสินค้าชิ้นเดียว |
| POST | `/api/products` | ✅ | เพิ่มสินค้า |
| PUT | `/api/products/:id` | ✅ | แก้ไขสินค้า |
| DELETE | `/api/products/:id` | ✅ | ลบสินค้า |
| POST | `/api/categories` | ✅ | เพิ่ม category |
| GET | `/api/categories` | ✅ | ดู category ทั้งหมด |

### Orders
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/orders/checkout` | ✅ | สั่งซื้อสินค้า |
| GET | `/api/orders` | ✅ | ดูออเดอร์ของตัวเอง |
| GET | `/api/orders/:id` | ✅ | ดูออเดอร์ชิ้นเดียว |
| POST | `/api/orders/:id/pay` | ✅ | Mock ชำระเงิน |

---

## Example Requests

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@test.com","password":"123456"}'
```

### Checkout
```bash
curl -X POST http://localhost:8080/api/orders/checkout \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":1,"quantity":2}]}'
```