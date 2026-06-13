# API Atreides

> Versión 1.0 — Documentación oficial para desarrolladores externos.

---

## Tabla de contenidos

1. [Introducción](#introducción)
2. [Autenticación](#autenticación)
3. [Formato de datos](#formato-de-datos)
4. [Formato de respuestas](#formato-de-respuestas)
5. [Códigos de error globales](#códigos-de-error-globales)
6. [Endpoints](#endpoints)
   - [Autenticación y usuarios](#autenticación-y-usuarios)
   - [Recursos CRUD](#recursos-crud)
7. [Listado completo de recursos](#listado-completo-de-recursos)
   - [Authors](#1-authors)
   - [Books](#2-books)
   - [Topics](#3-topics)
   - [Notes](#4-notes)
   - [Establishments](#5-establishments)
   - [Movement Types](#6-movement-types)
   - [Movements](#7-movements)
   - [Products](#8-products)
   - [Monthly Summaries](#9-monthly-summaries)
   - [Clients](#10-clients)
   - [Companies](#11-companies)
   - [Main Addresses](#12-main-addresses)
   - [Tax Information](#13-tax-information)
   - [Economic Activities](#14-economic-activities)
   - [Orders](#15-orders)

---

## Introducción

**Atreides** es una API RESTful diseñada para la gestión integral de libros, inventarios, pedidos y datos fiscales de una cafetería-librería. Permite administrar catálogos de autores, libros, productos, inventario, clientes, empresas y documentos contables.

**Base URL:** `https://api.atreides.com/v1`

---

## Autenticación

La API utiliza **JSON Web Tokens (JWT)** firmados con HS256. El token se obtiene al iniciar sesión y debe incluirse en todas las peticiones que requieran autenticación.

| Campo | Valor |
|-------|-------|
| **Tipo** | `Bearer JWT` |
| **Header** | `Authorization: Bearer <token>` |
| **Expiración** | 24 horas |
| **Algoritmo** | HS256 |

**Flujo:**

1. El cliente envía `POST /auth/login` con email y contraseña.
2. El servidor valida las credenciales y devuelve un `token` JWT.
3. El cliente incluye el token en el header `Authorization` de las peticiones subsiguientes.

> Nota: El primer usuario del sistema se registra sin token. Los registros posteriores requieren un token válido con permisos de administrador.

---

## Formato de datos

Todas las solicitudes y respuestas utilizan **JSON** (`Content-Type: application/json`).

| Aspecto | Especificación |
|---------|----------------|
| **Formato** | JSON |
| **Charset** | UTF-8 |
| **Fechas** | ISO 8601 (`YYYY-MM-DD`) |
| **Fecha-hora** | RFC 3339 (`2025-06-13T15:04:05Z`) |
| **Números decimales** | `float64` (punto flotante) |

---

## Formato de respuestas

### Respuesta exitosa

El cuerpo contiene directamente el recurso solicitado (objeto o array):

```json
// GET /authors → 200 OK
{
    "id": "A001",
    "name": "Frank Herbert",
    "country": "United States",
    "genres": ["Science Fiction"],
    "birthDay": "1920-10-08",
    "createdAt": "2025-06-13T12:00:00Z",
    "updatedAt": "2025-06-13T12:00:00Z"
}
```

```json
// GET /authors → 200 OK (colección)
[{ ... }, { ... }]
```

### Respuesta de error

```json
{
    "error": "mensaje descriptivo",
    "code": 400
}
```

### Respuesta de error de validación

```json
{
    "error": "validation failed",
    "details": {
        "email": "Invalid email format",
        "password": "Minimum length is 8"
    },
    "code": 400
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `error` | `string` | Mensaje legible del error |
| `details` | `object` | (opcional) Mapa de campo → mensaje para errores de validación |
| `code` | `int` | Código de estado HTTP |

---

## Códigos de error globales

| Código | Significado | Descripción |
|--------|-------------|-------------|
| `200` | OK | La solicitud se completó exitosamente |
| `201` | Created | El recurso fue creado exitosamente |
| `204` | No Content | El recurso fue eliminado (sin cuerpo en la respuesta) |
| `400` | Bad Request | JSON inválido, validación fallida o falta el ID del recurso |
| `401` | Unauthorized | Token faltante, inválido o expirado |
| `403` | Forbidden | El usuario no tiene permisos para la operación |
| `404` | Not Found | El recurso solicitado no existe |
| `405` | Method Not Allowed | Método HTTP no soportado para la ruta |
| `500` | Internal Server Error | Error interno del servidor |

---

## Endpoints

---

## Autenticación y usuarios

---

### POST /auth/register

Registra un nuevo usuario en el sistema.

- **Descripción:** Crea una cuenta de usuario. El primer registro no requiere autenticación; los siguientes requieren un token Bearer válido.
- **Headers requeridos:**
  - `Content-Type: application/json`
  - `Authorization: Bearer <token>` (opcional solo para el primer registro)

**Cuerpo de la solicitud:**

```json
{
    "name_full": "string (requerido)",
    "phone": "string (requerido)",
    "id_number": "string (requerido, solo dígitos)",
    "date_of_birth": "string (requerido, formato YYYY-MM-DD)",
    "email": "string (requerido, email válido)",
    "password": "string (requerido, mínimo 8 caracteres)"
}
```

**Ejemplo de solicitud:**

```bash
curl -X POST https://api.atreides.com/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name_full": "Paul Atreides",
    "phone": "+56912345678",
    "id_number": "12345678",
    "date_of_birth": "1990-01-15",
    "email": "paul@atreides.com",
    "password": "MiClaveSegura123"
  }'
```

**Respuesta exitosa — `201 Created`:**

```json
{
    "id": "20250613120000A",
    "name_full": "Paul Atreides",
    "phone": "+56912345678",
    "id_number": "12345678",
    "date_of_birth": "1990-01-15",
    "email": "paul@atreides.com",
    "roles": null,
    "createdAt": "2025-06-13T12:00:00Z",
    "updatedAt": "2025-06-13T12:00:00Z"
}
```

> Los campos `password_hash` y `auth_token` nunca se exponen en las respuestas.

**Códigos de error:**

| Código | Condición |
|--------|-----------|
| `400` | Campos faltantes, formato inválido o email ya registrado |
| `401` | Token inválido o expirado (registros posteriores al primero) |

---

### POST /auth/login

Inicia sesión con credenciales y obtiene un token JWT.

- **Descripción:** Valida email y contraseña, y devuelve un token JWT con expiración de 24 horas junto con los datos del usuario.
- **Headers requeridos:**
  - `Content-Type: application/json`

**Cuerpo de la solicitud:**

```json
{
    "email": "string (requerido)",
    "password": "string (requerido)"
}
```

**Ejemplo de solicitud:**

```bash
curl -X POST https://api.atreides.com/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "paul@atreides.com",
    "password": "MiClaveSegura123"
  }'
```

**Respuesta exitosa — `200 OK`:**

```json
{
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
        "id": "20250613120000A",
        "name_full": "Paul Atreides",
        "phone": "+56912345678",
        "id_number": "12345678",
        "date_of_birth": "1990-01-15",
        "email": "paul@atreides.com",
        "roles": null,
        "createdAt": "2025-06-13T12:00:00Z",
        "updatedAt": "2025-06-13T12:00:00Z"
    }
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `token` | `string` | JWT para autenticación en peticiones posteriores |
| `user` | `object` | Datos del usuario autenticado |

**Códigos de error:**

| Código | Condición |
|--------|-----------|
| `400` | Email o contraseña vacíos |
| `401` | Email o contraseña incorrectos |

---

### GET /users

Obtiene el listado de todos los usuarios registrados.

- **Descripción:** Devuelve un array con todos los usuarios del sistema.
- **Headers requeridos:**
  - `Content-Type: application/json`

**Ejemplo de solicitud:**

```bash
curl -X GET https://api.atreides.com/v1/users \
  -H "Content-Type: application/json"
```

**Respuesta exitosa — `200 OK`:**

```json
[
    {
        "id": "20250613120000A",
        "name_full": "Paul Atreides",
        "phone": "+56912345678",
        "id_number": "12345678",
        "date_of_birth": "1990-01-15",
        "email": "paul@atreides.com",
        "roles": null,
        "createdAt": "2025-06-13T12:00:00Z",
        "updatedAt": "2025-06-13T12:00:00Z"
    }
]
```

**Códigos de error:**

| Código | Condición |
|--------|-----------|
| `500` | Error interno al consultar la base de datos |

---

## Recursos CRUD

Los siguientes 15 recursos comparten el **mismo patrón CRUD**. Cada recurso se maneja a través de una única ruta base que acepta los métodos `GET`, `POST`, `PUT` y `DELETE`.

### Patrón general

| Método | Ruta | Acción | Código respuesta |
|--------|------|--------|------------------|
| `GET` | `/{resource}` | Listar todos los registros | `200` |
| `GET` | `/{resource}/{id}` | Obtener un registro por ID | `200` |
| `POST` | `/{resource}` | Crear un nuevo registro | `201` |
| `PUT` | `/{resource}/{id}` | Actualizar un registro existente | `200` |
| `DELETE` | `/{resource}/{id}` | Eliminar un registro | `204` |

**Headers comunes requeridos:**
- `Content-Type: application/json`

**Parámetros de ruta:**

| Parámetro | Tipo | Obligatorio | Descripción |
|-----------|------|-------------|-------------|
| `id` | `string` | Sí (GET/PUT/DELETE) | Identificador único del recurso |

---

### 1. Authors

**Ruta base:** `/authors`

#### Estructura del recurso

```json
{
    "id": "string",
    "name": "string",
    "country": "string",
    "genres": ["string"],
    "birthDay": "string (YYYY-MM-DD)",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### POST /authors — Crear autor

**Cuerpo de la solicitud:**

```json
{
    "name": "string (requerido)",
    "country": "string (requerido)",
    "genres": ["string (requerido, arreglo no vacío)"],
    "birthDay": "string (requerido, YYYY-MM-DD)"
}
```

**Ejemplo:**

```bash
curl -X POST https://api.atreides.com/v1/authors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Frank Herbert",
    "country": "United States",
    "genres": ["Science Fiction"],
    "birthDay": "1920-10-08"
  }'
```

**Respuesta — `201 Created`:** Recurso creado (ver estructura arriba).

#### GET /authors — Listar autores

```bash
curl -X GET https://api.atreides.com/v1/authors
```

**Respuesta — `200 OK`:** Array de autores.

#### GET /authors/{id} — Obtener autor por ID

```bash
curl -X GET https://api.atreides.com/v1/authors/A001
```

**Respuesta — `200 OK`:** Objeto autor.

#### PUT /authors/{id} — Actualizar autor

**Cuerpo:** Misma estructura que POST.

```bash
curl -X PUT https://api.atreides.com/v1/authors/A001 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Frank Herbert",
    "country": "United States",
    "genres": ["Science Fiction", "Adventure"],
    "birthDay": "1920-10-08"
  }'
```

**Respuesta — `200 OK`:** Recurso actualizado.

#### DELETE /authors/{id} — Eliminar autor

```bash
curl -X DELETE https://api.atreides.com/v1/authors/A001
```

**Respuesta — `204 No Content`:** Sin cuerpo.

#### Códigos de error (aplican a todos los recursos CRUD)

| Código | Condición |
|--------|-----------|
| `400` | JSON inválido, campos faltantes o formato incorrecto |
| `404` | Recurso no encontrado (GET, PUT, DELETE con ID inexistente) |
| `405` | Método HTTP no permitido |
| `500` | Error interno del servidor |

---

### 2. Books

**Ruta base:** `/books`

#### Estructura del recurso

```json
{
    "id": "string",
    "title": "string",
    "description": "string",
    "author": "string",
    "genres": ["string"],
    "photos": ["string"],
    "publicationDate": "string (YYYY-MM-DD, opcional)",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "title": "string (requerido)",
    "description": "string (requerido)",
    "author": "string (requerido)",
    "genres": ["string (requerido)"],
    "photos": ["string (opcional)"],
    "publicationDate": "string (opcional, YYYY-MM-DD)"
}
```

---

### 3. Topics

**Ruta base:** `/topics`

#### Estructura del recurso

```json
{
    "id": "string",
    "name": "string",
    "type": "string",
    "description": "string",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "name": "string (requerido)",
    "type": "string (requerido)",
    "description": "string (requerido)"
}
```

---

### 4. Notes

**Ruta base:** `/notes`

#### Estructura del recurso

```json
{
    "id": "string",
    "name": "string",
    "content": "string",
    "type": "string",
    "color": "string",
    "id_topic": "string",
    "id_book": "string (opcional)",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "name": "string (requerido)",
    "content": "string (requerido)",
    "type": "string (requerido)",
    "color": "string (requerido)",
    "id_topic": "string (requerido)",
    "id_book": "string (opcional)"
}
```

---

### 5. Establishments

**Ruta base:** `/establishments`

#### Estructura del recurso

```json
{
    "id": "string",
    "establishment_name": "string",
    "inventory_manager": "string",
    "warehouse_point_of_sale": "string",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "establishment_name": "string (requerido)",
    "inventory_manager": "string (requerido)",
    "warehouse_point_of_sale": "string (requerido)"
}
```

---

### 6. Movement Types

**Ruta base:** `/movement-types`

#### Estructura del recurso

```json
{
    "id": "string",
    "name": "string",
    "description": "string (opcional)",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "name": "string (requerido)",
    "description": "string (opcional)"
}
```

---

### 7. Movements

**Ruta base:** `/movements`

#### Estructura del recurso

```json
{
    "id": "string",
    "date": "string (YYYY-MM-DD)",
    "code": "string",
    "product": "string",
    "unit": "string",
    "entrance": "number",
    "output": "number",
    "balance": "number",
    "unit_cost": "number",
    "valor_value": "number",
    "movement_type_id": "string",
    "observations": "string (opcional)",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "date": "string (requerido, YYYY-MM-DD)",
    "code": "string (requerido)",
    "product": "string (requerido)",
    "unit": "string (requerido)",
    "entrance": "number (opcional, default 0)",
    "output": "number (opcional, default 0)",
    "balance": "number (opcional, default 0)",
    "unit_cost": "number (opcional, default 0)",
    "valor_value": "number (opcional, default 0)",
    "movement_type_id": "string (requerido, FK a movement_types)",
    "observations": "string (opcional)"
}
```

---

### 8. Products

**Ruta base:** `/products`

#### Estructura del recurso

```json
{
    "id": "string",
    "product_code": "string",
    "categories": ["string"],
    "unit": "string",
    "minimum_stock": "number",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "product_code": "string (requerido)",
    "categories": ["string (requerido)"],
    "unit": "string (requerido, ej: Kg, Liter, Unit)",
    "minimum_stock": "number (opcional, default 0)"
}
```

---

### 9. Monthly Summaries

**Ruta base:** `/monthly-summaries`

#### Estructura del recurso

```json
{
    "id": "string",
    "product": "string",
    "beginning_stock": "number",
    "incoming_orders": "number",
    "outgoing_orders": "number",
    "ending_stock": "number",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "product": "string (requerido)",
    "beginning_stock": "number (opcional, default 0)",
    "incoming_orders": "number (opcional, default 0)",
    "outgoing_orders": "number (opcional, default 0)",
    "ending_stock": "number (opcional, default 0)"
}
```

---

### 10. Clients

**Ruta base:** `/clients`

#### Estructura del recurso

```json
{
    "id": "string",
    "name_full": "string",
    "phone": "string",
    "correo": "string (opcional)",
    "address": "string",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "name_full": "string (requerido)",
    "phone": "string (requerido)",
    "correo": "string (opcional, email)",
    "address": "string (requerido)"
}
```

---

### 11. Companies

**Ruta base:** `/companies`

#### Estructura del recurso

```json
{
    "id": "string",
    "nit": "string",
    "social_reason": "string",
    "business_name": "string",
    "type_person": "string",
    "company_type": "string",
    "status": "string",
    "constitution_date": "string (YYYY-MM-DD)",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "nit": "string (requerido, único)",
    "social_reason": "string (requerido)",
    "business_name": "string (requerido)",
    "type_person": "string (requerido)",
    "company_type": "string (requerido)",
    "status": "string (requerido)",
    "constitution_date": "string (requerido, YYYY-MM-DD)"
}
```

---

### 12. Main Addresses

**Ruta base:** `/main-addresses`

#### Estructura del recurso

```json
{
    "id": "string",
    "user_id": "string",
    "company_id": "string",
    "country": "string",
    "department": "string",
    "address": "string",
    "postcode": "string",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "user_id": "string (requerido)",
    "company_id": "string (requerido)",
    "country": "string (requerido)",
    "department": "string (requerido)",
    "address": "string (requerido)",
    "postcode": "string (requerido)"
}
```

---

### 13. Tax Information

**Ruta base:** `/tax-information`

#### Estructura del recurso

```json
{
    "id": "string",
    "user_id": "string",
    "business_id": "string",
    "tax_regime": "string",
    "vat_responsible": "boolean",
    "withholding_taxpayer": "boolean",
    "large_taxpayer": "boolean",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "user_id": "string (requerido)",
    "business_id": "string (requerido)",
    "tax_regime": "string (requerido)",
    "vat_responsible": "boolean (opcional, default false)",
    "withholding_taxpayer": "boolean (opcional, default false)",
    "large_taxpayer": "boolean (opcional, default false)"
}
```

---

### 14. Economic Activities

**Ruta base:** `/economic-activities`

#### Estructura del recurso

```json
{
    "id": "string",
    "user_id": "string",
    "company_id": "string",
    "code": "string",
    "description": "string",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "user_id": "string (requerido)",
    "company_id": "string (requerido)",
    "code": "string (requerido)",
    "description": "string (requerido)"
}
```

---

### 15. Orders

**Ruta base:** `/orders`

#### Estructura del recurso

```json
{
    "id": "string",
    "order_numeric": "string",
    "date": "string (YYYY-MM-DD)",
    "hour": "string (HH:MM)",
    "attended_by": "string",
    "client_id": "string",
    "details": [
        {
            "code": "string",
            "product": "string",
            "quantity": "number",
            "unit_price": "number",
            "subtotal": "number",
            "discount": "number",
            "taxes": "number",
            "total": "number"
        }
    ],
    "payment_method": "string",
    "status": "string",
    "observations": "string (opcional)",
    "createdAt": "string (RFC 3339)",
    "updatedAt": "string (RFC 3339)"
}
```

#### Cuerpo de solicitud (POST/PUT)

```json
{
    "order_numeric": "string (requerido)",
    "date": "string (requerido, YYYY-MM-DD)",
    "hour": "string (requerido, HH:MM)",
    "attended_by": "string (requerido)",
    "client_id": "string (requerido, FK a clients)",
    "details": [
        {
            "code": "string (requerido)",
            "product": "string (requerido)",
            "quantity": "number (requerido, > 0)",
            "unit_price": "number (requerido, > 0)",
            "subtotal": "number (opcional)",
            "discount": "number (opcional)",
            "taxes": "number (opcional)",
            "total": "number (opcional)"
        }
    ],
    "payment_method": "string (requerido, ej: cash, transfer, debit-card, credit-card)",
    "status": "string (requerido, ej: received, in-preparation, delivered, cancelled)",
    "observations": "string (opcional)"
}
```

---

## Resumen de rutas

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/auth/register` | Registrar nuevo usuario |
| `POST` | `/auth/login` | Iniciar sesión y obtener JWT |
| `GET` | `/users` | Listar usuarios |
| `GET` | `/{resource}` | Listar registros del recurso |
| `GET` | `/{resource}/{id}` | Obtener registro por ID |
| `POST` | `/{resource}` | Crear registro |
| `PUT` | `/{resource}/{id}` | Actualizar registro |
| `DELETE` | `/{resource}/{id}` | Eliminar registro |

**Recursos disponibles:** `authors`, `books`, `topics`, `notes`, `establishments`, `movement-types`, `movements`, `products`, `monthly-summaries`, `clients`, `companies`, `main-addresses`, `tax-information`, `economic-activities`, `orders`.
