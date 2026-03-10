# Product Catalog API

## Architecture

**Pattern:** Domain-Driven Design (DDD) with Hexagonal Architecture

```
internal/
├── domain/          # Business entities (Product, Category, Variant)
├── application/     # Use cases (ListCatalog, GetProductDetail, etc.)
├── infrastructure/  # Database repositories (GORM + PostgreSQL)
└── ports/          # HTTP handlers (adapters)
```

**Key Principles:**
- Domain layer has no external dependencies
- Use cases orchestrate domain operations
- Infrastructure implements domain interfaces
- Handlers translate HTTP to use case calls

## API Endpoints

### Products

#### List Catalog
```http
GET /catalog?offset=0&limit=10&category=CLOTHING&maxPrice=50.00
```

**Query Parameters:**
- `offset` (optional): Starting position, default 0
- `limit` (optional): Items per page, default 10, max 100, min 1
- `category` (optional): Filter by category code
- `maxPrice` (optional): Filter products with price less than value

**Response:**
```json
{
  "products": [
    {
      "id": 1,
      "code": "PROD001",
      "price": "29.99",
      "category": {
        "id": 1,
        "code": "CLOTHING",
        "name": "Clothing"
      },
      "variants": [...]
    }
  ],
  "total": 100
}
```

#### Get Product Detail
```http
GET /catalog/{code}
```

**Response:**
```json
{
  "code": "PROD001",
  "price": "29.99",
  "category": {
    "code": "CLOTHING",
    "name": "Clothing"
  },
  "variants": [
    {
      "name": "Small",
      "sku": "PROD001-S",
      "price": "25.99"
    }
  ]
}
```

**Note:** Variants without a price inherit the product price.

### Categories

#### List Categories
```http
GET /categories
```

**Response:**
```json
[
  {
    "id": 1,
    "code": "CLOTHING",
    "name": "Clothing"
  }
]
```

#### Create Category
```http
POST /categories
Content-Type: application/json

{
  "code": "ELECTRONICS",
  "name": "Electronics"
}
```

**Response:** `201 Created`
```json
{
  "id": 4,
  "code": "ELECTRONICS",
  "name": "Electronics"
}
```

## Error Responses

All errors return appropriate HTTP status codes with plain text messages:

- `400 Bad Request` - Invalid input
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Running the API

```bash
# Start infrastructure
make docker-up

# Seed database
make seed

# Run server
make run
```

Server runs on `http://localhost:8484`

## Testing

```bash
# Run all tests
make test

# Manual testing
# Use routes.http file with REST Client
```

## Database Schema

**Products**
- `id`, `code` (unique), `price`, `category_id`

**Categories**
- `id`, `code` (unique), `name`

**Variants**
- `id`, `product_id`, `name`, `sku` (unique), `price` (nullable)

## Design Decisions

1. **No "UseCase" suffix** - Cleaner naming (ListCatalog vs ListCatalogUseCase)
2. **Standard library for responses** - No custom wrappers, idiomatic Go
3. **Single filter parameter** - Not variadic, simpler and clearer
4. **Price inheritance** - Variants without price use product price
5. **Pagination required** - Always returns total count for client-side pagination
