# Book Coffee Shop — REST API Reference

## Application Context

The **Book Coffee Shop API** is a backend service for managing a combined bookstore and café operation. It exposes REST endpoints for catalog management (authors, books, topics, notes), inventory (products, establishments, movement types, movements, monthly summaries), sales (clients, orders), business registration (companies, main addresses, tax information, economic activities), and user authentication.

Built with **Go** and **PostgreSQL**.

| Property | Value |
|----------|-------|
| **Base URL** | `http://localhost:8080` |
| **Content-Type** | `application/json` |
| **Port override** | Set environment variable `PORT` |

---

## Table of Contents

1. [Global Conventions](#global-conventions)
2. [Authentication](#authentication)
3. [Users](#users)
4. [Authors](#authors)
5. [Books](#books)
6. [Topics](#topics)
7. [Notes](#notes)
8. [Establishments](#establishments)
9. [Movement Types](#movement-types)
10. [Movements](#movements)
11. [Products](#products)
12. [Monthly Summaries](#monthly-summaries)
13. [Clients](#clients)
14. [Companies](#companies)
15. [Main Addresses](#main-addresses)
16. [Tax Information](#tax-information)
17. [Economic Activities](#economic-activities)
18. [Orders](#orders)
19. [CRUD Quick Reference](#crud-quick-reference)

---

## Global Conventions

### Resource identifiers

IDs are auto-generated on create using the format `YYYYMMDDHHMMSSX`, where `X` is an uppercase letter (A–Z).

### Dates and timestamps

| Format | Usage | Example |
|--------|-------|---------|
| `YYYY-MM-DD` | Date fields | `"2026-06-08"` |
| RFC 3339 | `createdAt`, `updatedAt` | `"2026-06-08T22:41:07.862350309Z"` |

### Error response

All errors return a JSON object:

```json
{
  "error": "Human-readable error message"
}
```

### Standard HTTP status codes

| Code | Meaning |
|------|---------|
| `200 OK` | Successful GET or PUT |
| `201 Created` | Successful POST |
| `204 No Content` | Successful DELETE (empty body) |
| `400 Bad Request` | Validation error, malformed JSON, or missing path ID on PUT/DELETE |
| `401 Unauthorized` | Invalid or missing authentication token |
| `404 Not Found` | Resource not found |
| `405 Method Not Allowed` | HTTP method not supported on route |
| `500 Internal Server Error` | Unexpected server or database error |

### CRUD routing pattern

All resource collections follow the same URL pattern:

| Operation | Method | Path |
|-----------|--------|------|
| List | `GET` | `/{resource}` |
| Get by ID | `GET` | `/{resource}/{id}` |
| Create | `POST` | `/{resource}` |
| Update | `PUT` | `/{resource}/{id}` |
| Delete | `DELETE` | `/{resource}/{id}` |

**Path variables**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | `string` | Yes (GET-one, PUT, DELETE) | Unique resource identifier |

No query parameters are used on any endpoint.

---

## Authentication

### Register User

**`POST /auth/register`**

Creates a new user account in the `users` table. The first user in the system can register without a token. Subsequent registrations require a valid Bearer token from an existing authenticated user.

#### Functional description

Validates user profile fields, hashes the password with bcrypt, and persists the record. Returns the created user without sensitive fields (`password_hash`, `auth_token`).

#### Headers

| Header | Type | Required | Description |
|--------|------|----------|-------------|
| `Content-Type` | `string` | Yes | `application/json` |
| `Authorization` | `string` | Conditional | `Bearer <token>` — required when at least one user already exists |

#### Request body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name_full` | `string` | Yes | Full name |
| `phone` | `string` | Yes | Phone number |
| `id_number` | `string` | Yes | Numeric identification number |
| `date_of_birth` | `string` | Yes | Date in `YYYY-MM-DD` format |
| `email` | `string` | Yes | Valid email address (stored lowercase) |
| `password` | `string` | Yes | Minimum 8 characters |

#### Response body (`201 Created`)

```json
{
  "id": "20260608143022A",
  "name_full": "Ana García",
  "phone": "+505 8888-0001",
  "id_number": "001234567890",
  "date_of_birth": "1990-05-15",
  "email": "ana@example.com",
  "createdAt": "2026-06-08T14:30:22.123456789Z",
  "updatedAt": "2026-06-08T14:30:22.123456789Z"
}
```

#### Status codes and errors

| Code | Condition |
|------|-----------|
| `201` | User created successfully |
| `400` | Validation failure, duplicate email, or password processing error |
| `401` | Missing or invalid Bearer token (when users already exist) |
| `405` | Method other than POST |

**Common error messages:** `name_full cannot be empty`, `id_number must be numeric`, `date_of_birth must be in YYYY-MM-DD format`, `email format is invalid`, `password must be at least 8 characters`, `email already registered`, `authorization token is required`, `invalid or expired token`

#### Example (cURL)

```bash
curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name_full": "Ana García",
    "phone": "+505 8888-0001",
    "id_number": "001234567890",
    "date_of_birth": "1990-05-15",
    "email": "ana@example.com",
    "password": "securePass123"
  }'
```

#### Example (Go)

```go
payload := map[string]string{
    "name_full":     "Ana García",
    "phone":         "+505 8888-0001",
    "id_number":     "001234567890",
    "date_of_birth": "1990-05-15",
    "email":         "ana@example.com",
    "password":      "securePass123",
}
body, _ := json.Marshal(payload)
req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/auth/register", bytes.NewReader(body))
req.Header.Set("Content-Type", "application/json")
resp, err := http.DefaultClient.Do(req)
```

---

### Login

**`POST /auth/login`**

Authenticates a user by email and password. On first login, a JWT is generated and stored. On subsequent logins, the client must send the existing Bearer token in the `Authorization` header.

#### Functional description

Verifies credentials against the `users` table. Returns a JWT token (24-hour expiry) and the user profile on success.

#### Headers

| Header | Type | Required | Description |
|--------|------|----------|-------------|
| `Content-Type` | `string` | Yes | `application/json` |
| `Authorization` | `string` | Conditional | `Bearer <token>` — required when user already has a stored auth token |

#### Request body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | `string` | Yes | Registered email |
| `password` | `string` | Yes | Account password |

#### Response body (`200 OK`)

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "20260608143022A",
    "name_full": "Ana García",
    "phone": "+505 8888-0001",
    "id_number": "001234567890",
    "date_of_birth": "1990-05-15",
    "email": "ana@example.com",
    "createdAt": "2026-06-08T14:30:22.123456789Z",
    "updatedAt": "2026-06-08T14:30:22.123456789Z"
  }
}
```

#### Status codes and errors

| Code | Condition |
|------|-----------|
| `200` | Login successful |
| `400` | Empty email or password |
| `401` | Invalid credentials or token issues |
| `405` | Method other than POST |

**Common error messages:** `invalid email or password`, `authorization token is required`, `token does not belong to this user`, `invalid or expired token`

#### Example (cURL)

```bash
curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ana@example.com","password":"securePass123"}'
```

---

## Users

### List Users

**`GET /users`**

Returns all registered users. Password hashes and auth tokens are never exposed.

#### Functional description

Retrieves the full list of user records from the `users` table. No authentication middleware is applied in the current implementation.

#### Parameters

None.

#### Response body (`200 OK`)

```json
[
  {
    "id": "20260608143022A",
    "name_full": "Ana García",
    "phone": "+505 8888-0001",
    "id_number": "001234567890",
    "date_of_birth": "1990-05-15",
    "email": "ana@example.com",
    "createdAt": "2026-06-08T14:30:22.123456789Z",
    "updatedAt": "2026-06-08T14:30:22.123456789Z"
  }
]
```

#### Status codes and errors

| Code | Condition |
|------|-----------|
| `200` | Success |
| `405` | Method other than GET |
| `500` | Database error |

#### Example (cURL)

```bash
curl -s http://localhost:8080/users
```

---

## Authors

Resource path: `/authors` · Database table: `authors`

### List Authors — `GET /authors`

Returns all authors ordered by creation date.

**Response (`200 OK`):** Array of `Author` objects.

```json
[
  {
    "id": "20260608224107A",
    "name": "Gabriel García Márquez",
    "country": "Colombia",
    "genres": ["Realismo mágico", "Novela"],
    "birthDay": "1927-03-06",
    "createdAt": "2026-06-08T22:41:07.862350309Z",
    "updatedAt": "2026-06-08T22:41:07.862350309Z"
  }
]
```

| Code | Condition |
|------|-----------|
| `200` | Success |
| `500` | Database error |

```bash
curl -s http://localhost:8080/authors
```

---

### Get Author by ID — `GET /authors/{id}`

Returns a single author record.

**Path variables:** `id` (`string`, required)

**Response (`200 OK`):** Single `Author` object (same schema as above).

| Code | Condition |
|------|-----------|
| `200` | Found |
| `404` | `author not found` |

```bash
curl -s http://localhost:8080/authors/20260608224107A
```

---

### Create Author — `POST /authors`

Creates a new author with auto-generated ID.

#### Request body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | `string` | Yes | Author full name |
| `country` | `string` | Yes | Country of origin |
| `genres` | `string[]` | Yes | Non-empty array of literary genres |
| `birthDay` | `string` | Yes | Birth date |

#### Response (`201 Created`)

```json
{
  "id": "20260608224107A",
  "name": "Gabriel García Márquez",
  "country": "Colombia",
  "genres": ["Realismo mágico", "Novela"],
  "birthDay": "1927-03-06",
  "createdAt": "2026-06-08T22:41:07.862350309Z",
  "updatedAt": "2026-06-08T22:41:07.862350309Z"
}
```

| Code | Condition |
|------|-----------|
| `201` | Created |
| `400` | Validation error or invalid JSON |

**Validation errors:** `name cannot be empty`, `country cannot be empty`, `genres cannot be empty`, `birthDay cannot be empty`

```bash
curl -s -X POST http://localhost:8080/authors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gabriel García Márquez",
    "country": "Colombia",
    "genres": ["Realismo mágico", "Novela"],
    "birthDay": "1927-03-06"
  }'
```

---

### Update Author — `PUT /authors/{id}`

Replaces all mutable fields of an existing author.

**Path variables:** `id` (`string`, required)

**Request body:** Same schema as Create.

**Response (`200 OK`):** Updated `Author` object.

| Code | Condition |
|------|-----------|
| `200` | Updated |
| `400` | Validation error, missing ID, or invalid JSON |
| `404` | `author not found` |

```bash
curl -s -X PUT http://localhost:8080/authors/20260608224107A \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gabriel García Márquez",
    "country": "Colombia",
    "genres": ["Novela"],
    "birthDay": "1927-03-06"
  }'
```

---

### Delete Author — `DELETE /authors/{id}`

Permanently removes an author record.

**Path variables:** `id` (`string`, required)

**Response:** Empty body.

| Code | Condition |
|------|-----------|
| `204` | Deleted |
| `400` | Missing ID in path |
| `404` | `author not found` |

```bash
curl -s -X DELETE http://localhost:8080/authors/20260608224107A
```

---

## Books

Resource path: `/books` · Database table: `books`

### Request body schema (POST / PUT)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | `string` | Yes | Book title |
| `description` | `string` | Yes | Book description |
| `author` | `string` | Yes | Author name |
| `genres` | `string[]` | Yes | Non-empty genre array |
| `photos` | `string[]` | No | Image URLs (defaults to `[]`) |
| `publicationDate` | `string` | No | Publication date |

### Response schema

```json
{
  "id": "string",
  "title": "string",
  "description": "string",
  "author": "string",
  "genres": ["string"],
  "photos": ["string"],
  "publicationDate": "string",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
```

### Endpoints

| Method | Path | Success | Key errors |
|--------|------|---------|------------|
| `GET` | `/books` | `200` array | `500` |
| `GET` | `/books/{id}` | `200` object | `404` book not found |
| `POST` | `/books` | `201` object | `400` validation |
| `PUT` | `/books/{id}` | `200` object | `400`, `404` |
| `DELETE` | `/books/{id}` | `204` empty | `404` |

**Validation errors:** `title cannot be empty`, `description cannot be empty`, `author cannot be empty`, `genres cannot be empty`

```bash
curl -s -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Cien años de soledad",
    "description": "Historia de la familia Buendía",
    "author": "Gabriel García Márquez",
    "genres": ["Realismo mágico"],
    "photos": ["https://example.com/cover.jpg"],
    "publicationDate": "1967-06-05"
  }'
```

---

## Topics

Resource path: `/topics` · Database table: `topics`

### Request body schema (POST / PUT)

| Field | Type | Required |
|-------|------|----------|
| `name` | `string` | Yes |
| `type` | `string` | Yes |
| `description` | `string` | Yes |

### Response schema

```json
{
  "id": "string",
  "name": "string",
  "type": "string",
  "description": "string",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
```

### Endpoints

| Method | Path | Success | Key errors |
|--------|------|---------|------------|
| `GET` | `/topics` | `200` | `500` |
| `GET` | `/topics/{id}` | `200` | `404` topic not found |
| `POST` | `/topics` | `201` | `400` |
| `PUT` | `/topics/{id}` | `200` | `400`, `404` |
| `DELETE` | `/topics/{id}` | `204` | `404` |

```bash
curl -s -X POST http://localhost:8080/topics \
  -H "Content-Type: application/json" \
  -d '{"name":"Literatura","type":"general","description":"Temas literarios"}'
```

---

## Notes

Resource path: `/notes` · Database table: `notes`

### Request body schema (POST / PUT)

| Field | Type | Required |
|-------|------|----------|
| `name` | `string` | Yes |
| `content` | `string` | Yes |
| `type` | `string` | Yes |
| `color` | `string` | Yes |
| `id_topic` | `string` | Yes |
| `id_book` | `string` | No |

### Response schema

```json
{
  "id": "string",
  "name": "string",
  "content": "string",
  "type": "string",
  "color": "string",
  "id_topic": "string",
  "id_book": "string",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
```

### Endpoints

| Method | Path | Success | Key errors |
|--------|------|---------|------------|
| `GET` | `/notes` | `200` | `500` |
| `GET` | `/notes/{id}` | `200` | `404` note not found |
| `POST` | `/notes` | `201` | `400` |
| `PUT` | `/notes/{id}` | `200` | `400`, `404` |
| `DELETE` | `/notes/{id}` | `204` | `404` |

```bash
curl -s -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Idea de negocio",
    "content": "Catas de café con maridaje de libros",
    "type": "idea",
    "color": "yellow",
    "id_topic": "TOPIC_ID",
    "id_book": "BOOK_ID"
  }'
```

---

## Establishments

Resource path: `/establishments` · Database table: `establishments`

### Request body schema (POST / PUT)

| Field | Type | Required |
|-------|------|----------|
| `establishment_name` | `string` | Yes |
| `inventory_manager` | `string` | Yes |
| `warehouse_point_of_sale` | `string` | Yes |

### Response schema

```json
{
  "id": "string",
  "establishment_name": "string",
  "inventory_manager": "string",
  "warehouse_point_of_sale": "string",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
```

### Endpoints

| Method | Path | Success | Key errors |
|--------|------|---------|------------|
| `GET` | `/establishments` | `200` | `500` |
| `GET` | `/establishments/{id}` | `200` | `404` establishment not found |
| `POST` | `/establishments` | `201` | `400` |
| `PUT` | `/establishments/{id}` | `200` | `400`, `404` |
| `DELETE` | `/establishments/{id}` | `204` | `404` |

```bash
curl -s -X POST http://localhost:8080/establishments \
  -H "Content-Type: application/json" \
  -d '{
    "establishment_name": "Sucursal Central",
    "inventory_manager": "Carlos Ruiz",
    "warehouse_point_of_sale": "Almacén A"
  }'
```

---

## Movement Types

Resource path: `/movement-types` · Database table: `movement_types`

### Request body schema (POST / PUT)

| Field | Type | Required |
|-------|------|----------|
| `name` | `string` | Yes |
| `description` | `string` | No |

### Response schema

```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
```

### Endpoints

| Method | Path | Success | Key errors |
|--------|------|---------|------------|
| `GET` | `/movement-types` | `200` | `500` |
| `GET` | `/movement-types/{id}` | `200` | `404` movement type not found |
| `POST` | `/movement-types` | `201` | `400` |
| `PUT` | `/movement-types/{id}` | `200` | `400`, `404` |
| `DELETE` | `/movement-types/{id}` | `204` | `404` |

```bash
curl -s -X POST http://localhost:8080/movement-types \
  -H "Content-Type: application/json" \
  -d '{"name":"Compra proveedor","description":"Entrada por compra"}'
```

---

## Movements

Resource path: `/movements` · Database table: `movements`

### Request body schema (POST / PUT)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `date` | `string` | Yes | `YYYY-MM-DD` |
| `code` | `string` | Yes | Movement code |
| `product` | `string` | Yes | Product name |
| `unit` | `string` | Yes | Unit of measure |
| `entrance` | `number` | No | Default `0` |
| `output` | `number` | No | Default `0` |
| `balance` | `number` | No | Default `0` |
| `unit_cost` | `number` | No | Default `0` |
| `valor_value` | `number` | No | Default `0` |
| `movement_type_id` | `string` | Yes | FK to `movement_types.id` |
| `observations` | `string` | No | Free text |

### Response schema

```json
{
  "id": "string",
  "date": "2026-06-08",
  "code": "MOV-001",
  "product": "Café Guatemalteco",
  "unit": "Kg",
  "entrance": 50,
  "output": 0,
  "balance": 50,
  "unit_cost": 12.5,
  "valor_value": 625.0,
  "movement_type_id": "MOVEMENT_TYPE_ID",
  "observations": "Compra mensual",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
```

### Endpoints

| Method | Path | Success | Key errors |
|--------|------|---------|------------|
| `GET` | `/movements` | `200` | `500` |
| `GET` | `/movements/{id}` | `200` | `404` movement not found |
| `POST` | `/movements` | `201` | `400` |
| `PUT` | `/movements/{id}` | `200` | `400`, `404` |
| `DELETE` | `/movements/{id}` | `204` | `404` |

**Validation errors:** `movement_type_id cannot be empty`, `movement type not found`

```bash
curl -s -X POST http://localhost:8080/movements \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2026-06-08",
    "code": "MOV-001",
    "product": "Café Guatemalteco",
    "unit": "Kg",
    "entrance": 50,
    "output": 0,
    "balance": 50,
    "unit_cost": 12.50,
    "valor_value": 625.00,
    "movement_type_id": "MOVEMENT_TYPE_ID",
    "observations": "Compra mensual"
  }'
```

---

## Products

Resource path: `/products` · Database table: `products`

### Request body schema (POST / PUT)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `product_code` | `string` | Yes | Unique product code |
| `categories` | `string[]` | Yes | Non-empty category array |
| `unit` | `string` | Yes | One of: `Kg`, `Liter`, `Pound`, `Grams`, `Unit` |
| `minimum_stock` | `number` | No | Default `0` |

### Response schema

```json
{
  "id": "string",
  "product_code": "CAFE-001",
  "categories": ["Bebidas", "Calientes"],
  "unit": "Unit",
  "minimum_stock": 10,
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
```

### Endpoints

| Method | Path | Success | Key errors |
|--------|------|---------|------------|
| `GET` | `/products` | `200` | `500` |
| `GET` | `/products/{id}` | `200` | `404` product not found |
| `POST` | `/products` | `201` | `400` |
| `PUT` | `/products/{id}` | `200` | `400`, `404` |
| `DELETE` | `/products/{id}` | `204` | `404` |

**Validation errors:** `unit must be one of: Kg, Liter, Pound, Grams, Unit`

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "product_code": "CAFE-001",
    "categories": ["Bebidas", "Calientes"],
    "unit": "Unit",
    "minimum_stock": 10
  }'
```

---

## Monthly Summaries

Resource path: `/monthly-summaries` · Database table: `monthly_summaries`

### Request body schema (POST / PUT)

| Field | Type | Required |
|-------|------|----------|
| `product` | `string` | Yes |
| `beginning_stock` | `number` | No (default `0`) |
| `incoming_orders` | `number` | No (default `0`) |
| `outgoing_orders` | `number` | No (default `0`) |
| `ending_stock` | `number` | No (default `0`) |

### Response schema

```json
{
  "id": "string",
  "product": "Café Guatemalteco",
  "beginning_stock": 100,
  "incoming_orders": 50,
  "outgoing_orders": 30,
  "ending_stock": 120,
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
```

### Endpoints

| Method | Path | Success | Key errors |
|--------|------|---------|------------|
| `GET` | `/monthly-summaries` | `200` | `500` |
| `GET` | `/monthly-summaries/{id}` | `200` | `404` monthly summary not found |
| `POST` | `/monthly-summaries` | `201` | `400` |
| `PUT` | `/monthly-summaries/{id}` | `200` | `400`, `404` |
| `DELETE` | `/monthly-summaries/{id}` | `204` | `404` |

---

## Clients

Resource path: `/clients` · Database table: `clients`

### Request body schema (POST / PUT)

| Field | Type | Required |
|-------|------|----------|
| `name_full` | `string` | Yes |
| `phone` | `string` | Yes |
| `correo` | `string` | No |
| `address` | `string` | Yes |

### Response schema

```json
{
  "id": "string",
  "name_full": "Juan Pérez",
  "phone": "+505 8888-9999",
  "correo": "juan@example.com",
  "address": "Managua, Nicaragua",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
```

### Endpoints

| Method | Path | Success | Key errors |
|--------|------|---------|------------|
| `GET` | `/clients` | `200` | `500` |
| `GET` | `/clients/{id}` | `200` | `404` client not found |
| `POST` | `/clients` | `201` | `400` |
| `PUT` | `/clients/{id}` | `200` | `400`, `404` |
| `DELETE` | `/clients/{id}` | `204` | `404` |

```bash
curl -s -X POST http://localhost:8080/clients \
  -H "Content-Type: application/json" \
  -d '{
    "name_full": "Juan Pérez",
    "phone": "+505 8888-9999",
    "correo": "juan@example.com",
    "address": "Managua, Nicaragua"
  }'
```

---

## Companies

Resource path: `/companies` · Database table: `companies`

### Request body schema (POST / PUT)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `nit` | `string` | Yes | Unique tax ID |
| `social_reason` | `string` | Yes | Legal name |
| `business_name` | `string` | Yes | Trade name |
| `type_person` | `string` | Yes | Person type |
| `company_type` | `string` | Yes | Company classification |
| `status` | `string` | Yes | Current status |
| `constitution_date` | `string` | Yes | Incorporation date |

### Response schema

```json
{
  "id": "string",
  "nit": "0614-280891-001-0",
  "social_reason": "Café y Libros S.A.",
  "business_name": "Book Coffee Shop",
  "type_person": "Jurídica",
  "company_type": "Sociedad Anónima",
  "status": "Activa",
  "constitution_date": "2020-01-15",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
```

### Endpoints

| Method | Path | Success | Key errors |
|--------|------|---------|------------|
| `GET` | `/companies` | `200` | `500` |
| `GET` | `/companies/{id}` | `200` | `404` company not found |
| `POST` | `/companies` | `201` | `400` |
| `PUT` | `/companies/{id}` | `200` | `400`, `404` |
| `DELETE` | `/companies/{id}` | `204` | `404` |

**Validation errors:** `a company with this nit already exists`

```bash
curl -s -X POST http://localhost:8080/companies \
  -H "Content-Type: application/json" \
  -d '{
    "nit": "0614-280891-001-0",
    "social_reason": "Café y Libros S.A.",
    "business_name": "Book Coffee Shop",
    "type_person": "Jurídica",
    "company_type": "Sociedad Anónima",
    "status": "Activa",
    "constitution_date": "2020-01-15"
  }'
```

---

## Main Addresses

Resource path: `/main-addresses` · Database table: `main_addresses`

### Request body schema (POST / PUT)

| Field | Type | Required |
|-------|------|----------|
| `user_id` | `string` | Yes |
| `company_id` | `string` | Yes |
| `country` | `string` | Yes |
| `department` | `string` | Yes |
| `address` | `string` | Yes |
| `postcode` | `string` | Yes |

### Response schema

```json
{
  "id": "string",
  "user_id": "USER_ID",
  "company_id": "COMPANY_ID",
  "country": "Nicaragua",
  "department": "Managua",
  "address": "Km 8 Carretera Sur",
  "postcode": "14000",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
```

### Endpoints

| Method | Path | Success | Key errors |
|--------|------|---------|------------|
| `GET` | `/main-addresses` | `200` | `500` |
| `GET` | `/main-addresses/{id}` | `200` | `404` main address not found |
| `POST` | `/main-addresses` | `201` | `400` |
| `PUT` | `/main-addresses/{id}` | `200` | `400`, `404` |
| `DELETE` | `/main-addresses/{id}` | `204` | `404` |

---

## Tax Information

Resource path: `/tax-information` · Database table: `tax_information`

### Request body schema (POST / PUT)

| Field | Type | Required | Default |
|-------|------|----------|---------|
| `user_id` | `string` | Yes | — |
| `business_id` | `string` | Yes | — |
| `tax_regime` | `string` | Yes | — |
| `vat_responsible` | `boolean` | No | `false` |
| `withholding_taxpayer` | `boolean` | No | `false` |
| `large_taxpayer` | `boolean` | No | `false` |

### Response schema

```json
{
  "id": "string",
  "user_id": "USER_ID",
  "business_id": "COMPANY_ID",
  "tax_regime": "General",
  "vat_responsible": true,
  "withholding_taxpayer": false,
  "large_taxpayer": false,
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
```

### Endpoints

| Method | Path | Success | Key errors |
|--------|------|---------|------------|
| `GET` | `/tax-information` | `200` | `500` |
| `GET` | `/tax-information/{id}` | `200` | `404` tax information not found |
| `POST` | `/tax-information` | `201` | `400` |
| `PUT` | `/tax-information/{id}` | `200` | `400`, `404` |
| `DELETE` | `/tax-information/{id}` | `204` | `404` |

```bash
curl -s -X POST http://localhost:8080/tax-information \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "USER_ID",
    "business_id": "COMPANY_ID",
    "tax_regime": "General",
    "vat_responsible": true,
    "withholding_taxpayer": false,
    "large_taxpayer": false
  }'
```

---

## Economic Activities

Resource path: `/economic-activities` · Database table: `economic_activities`

### Request body schema (POST / PUT)

| Field | Type | Required |
|-------|------|----------|
| `user_id` | `string` | Yes |
| `company_id` | `string` | Yes |
| `code` | `string` | Yes |
| `description` | `string` | Yes |

### Response schema

```json
{
  "id": "string",
  "user_id": "USER_ID",
  "company_id": "COMPANY_ID",
  "code": "4711",
  "description": "Retail sale in non-specialized stores",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
```

### Endpoints

| Method | Path | Success | Key errors |
|--------|------|---------|------------|
| `GET` | `/economic-activities` | `200` | `500` |
| `GET` | `/economic-activities/{id}` | `200` | `404` economic activity not found |
| `POST` | `/economic-activities` | `201` | `400` |
| `PUT` | `/economic-activities/{id}` | `200` | `400`, `404` |
| `DELETE` | `/economic-activities/{id}` | `204` | `404` |

```bash
curl -s -X POST http://localhost:8080/economic-activities \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "USER_ID",
    "company_id": "COMPANY_ID",
    "code": "4711",
    "description": "Retail sale in non-specialized stores"
  }'
```

---

## Orders

Resource path: `/orders` · Database table: `orders`

### Create Order — `POST /orders`

Creates a sales order with line-item details stored as JSONB.

#### Request body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `order_numeric` | `string` | Yes | Human-readable order number |
| `date` | `string` | Yes | `YYYY-MM-DD` |
| `hour` | `string` | Yes | Time in `HH:MM` format |
| `attended_by` | `string` | Yes | Staff member name |
| `client_id` | `string` | Yes | FK to `clients.id` |
| `details` | `array` | Yes | Minimum 1 line item (see below) |
| `payment_method` | `string` | Yes | See allowed values |
| `status` | `string` | Yes | See allowed values |
| `observations` | `string` | No | Free text notes |

**`details[]` item schema**

| Field | Type | Required |
|-------|------|----------|
| `code` | `string` | Yes |
| `product` | `string` | Yes |
| `quantity` | `number` | Yes (> 0) |
| `unit_price` | `number` | Yes (> 0) |
| `subtotal` | `number` | No |
| `discount` | `number` | No |
| `taxes` | `number` | No |
| `total` | `number` | No |

**Allowed `payment_method` values:** `cash`, `transfer`, `debit-card`, `credit-card`

**Allowed `status` values:** `received`, `in-preparation`, `ready-for-delivery`, `delivered`, `cancelled`

#### Response body (`201 Created`)

```json
{
  "id": "20260608143022B",
  "order_numeric": "ORD-001",
  "date": "2026-06-08",
  "hour": "14:30",
  "attended_by": "Carlos",
  "client_id": "CLIENT_ID",
  "details": [
    {
      "code": "P001",
      "product": "Café Latte",
      "quantity": 2,
      "unit_price": 3.5,
      "subtotal": 7.0,
      "discount": 0.5,
      "taxes": 0.91,
      "total": 7.41
    }
  ],
  "payment_method": "cash",
  "status": "received",
  "observations": "Sin azúcar",
  "createdAt": "2026-06-08T14:30:22.123456789Z",
  "updatedAt": "2026-06-08T14:30:22.123456789Z"
}
```

| Code | Condition |
|------|-----------|
| `201` | Order created |
| `400` | Validation error |

**Validation errors:** `details cannot be empty`, `details quantity must be greater than 0`, `payment_method must be one of: cash, transfer, debit-card, credit-card`, `status must be one of: received, in-preparation, ready-for-delivery, delivered, cancelled`

```bash
curl -s -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "order_numeric": "ORD-001",
    "date": "2026-06-08",
    "hour": "14:30",
    "attended_by": "Carlos",
    "client_id": "CLIENT_ID",
    "details": [{
      "code": "P001",
      "product": "Café Latte",
      "quantity": 2,
      "unit_price": 3.50,
      "subtotal": 7.00,
      "discount": 0.50,
      "taxes": 0.91,
      "total": 7.41
    }],
    "payment_method": "cash",
    "status": "received",
    "observations": "Sin azúcar"
  }'
```

---

### List Orders — `GET /orders`

Returns all orders.

**Response (`200 OK`):** Array of `Order` objects.

| Code | Condition |
|------|-----------|
| `200` | Success |
| `500` | Database error |

```bash
curl -s http://localhost:8080/orders
```

---

### Get Order by ID — `GET /orders/{id}`

**Path variables:** `id` (`string`, required)

**Response (`200 OK`):** Single `Order` object.

| Code | Condition |
|------|-----------|
| `200` | Found |
| `404` | `order not found` |

```bash
curl -s http://localhost:8080/orders/ORDER_ID
```

---

### Update Order — `PUT /orders/{id}`

Full replacement of order fields. Commonly used to update order status (e.g., `received` → `in-preparation`).

**Path variables:** `id` (`string`, required)

**Request body:** Same schema as Create.

**Response (`200 OK`):** Updated `Order` object.

| Code | Condition |
|------|-----------|
| `200` | Updated |
| `400` | Validation error |
| `404` | `order not found` |

```bash
curl -s -X PUT http://localhost:8080/orders/ORDER_ID \
  -H "Content-Type: application/json" \
  -d '{
    "order_numeric": "ORD-001",
    "date": "2026-06-08",
    "hour": "14:30",
    "attended_by": "Carlos",
    "client_id": "CLIENT_ID",
    "details": [{"code":"P001","product":"Café Latte","quantity":2,"unit_price":3.5}],
    "payment_method": "cash",
    "status": "in-preparation",
    "observations": "En preparación"
  }'
```

---

### Delete Order — `DELETE /orders/{id}`

**Path variables:** `id` (`string`, required)

**Response:** Empty body.

| Code | Condition |
|------|-----------|
| `204` | Deleted |
| `404` | `order not found` |

```bash
curl -s -X DELETE http://localhost:8080/orders/ORDER_ID
```

---

## CRUD Quick Reference

All resources below support the standard CRUD pattern described in [Global Conventions](#global-conventions).

| Resource | Base path | DB table |
|----------|-----------|----------|
| Authors | `/authors` | `authors` |
| Books | `/books` | `books` |
| Topics | `/topics` | `topics` |
| Notes | `/notes` | `notes` |
| Establishments | `/establishments` | `establishments` |
| Movement Types | `/movement-types` | `movement_types` |
| Movements | `/movements` | `movements` |
| Products | `/products` | `products` |
| Monthly Summaries | `/monthly-summaries` | `monthly_summaries` |
| Clients | `/clients` | `clients` |
| Companies | `/companies` | `companies` |
| Main Addresses | `/main-addresses` | `main_addresses` |
| Tax Information | `/tax-information` | `tax_information` |
| Economic Activities | `/economic-activities` | `economic_activities` |
| Orders | `/orders` | `orders` |

### Shared CRUD status matrix

| Operation | Success | Client errors | Not found / server |
|-----------|---------|---------------|---------------------|
| List (`GET`) | `200` | — | `500` |
| Get by ID (`GET`) | `200` | — | `404` |
| Create (`POST`) | `201` | `400` | — |
| Update (`PUT`) | `200` | `400` | `404` |
| Delete (`DELETE`) | `204` | `400` (missing id) | `404` |
| Wrong method | — | `405` | — |

**Shared handler errors:** `invalid request body`, `id is required`, `method not allowed`

---

## End-to-End Example Script

```bash
#!/bin/bash
BASE="http://localhost:8080"

# 1. Register first user
curl -s -X POST $BASE/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name_full":"Admin User","phone":"+505 8888-0000","id_number":"1234567890","date_of_birth":"1990-01-01","email":"admin@example.com","password":"adminPass123"}'

# 2. Login
TOKEN=$(curl -s -X POST $BASE/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"adminPass123"}' | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# 3. Create client
CLIENT=$(curl -s -X POST $BASE/clients \
  -H "Content-Type: application/json" \
  -d '{"name_full":"Juan Pérez","phone":"+505 8888-9999","address":"Managua"}')
CLIENT_ID=$(echo $CLIENT | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

# 4. Create order
curl -s -X POST $BASE/orders \
  -H "Content-Type: application/json" \
  -d "{\"order_numeric\":\"ORD-001\",\"date\":\"2026-06-08\",\"hour\":\"14:30\",\"attended_by\":\"Carlos\",\"client_id\":\"$CLIENT_ID\",\"details\":[{\"code\":\"P1\",\"product\":\"Café Latte\",\"quantity\":2,\"unit_price\":3.5}],\"payment_method\":\"cash\",\"status\":\"received\"}"
```

---

## Database Entity Relationship Overview

```
users ──────────────┬── main_addresses ── companies
                    ├── tax_information
                    └── economic_activities ── companies

authors ── (referenced by name in books.author)
clients ── orders (client_id)
movement_types ── movements (movement_type_id)
```
