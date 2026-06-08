# API REST - Book Coffee Shop

API de gestión para una cafetería-librería construida con Go y PostgreSQL.

**Base URL:** `http://localhost:8080`

---

## Índice de entidades

| Entidad | Endpoint | Descripción |
|---------|----------|-------------|
| Autores | `/authors` | CRUD de autores de libros |
| Libros | `/books` | CRUD de libros |
| Temas de interés | `/topics` | CRUD de temas |
| Notas | `/notes` | CRUD de notas |
| Establecimientos | `/establishments` | CRUD de sucursales/almacenes |
| Tipos de movimiento | `/movement-types` | CRUD de tipos de movimiento |
| Movimientos | `/movements` | CRUD de movimientos de inventario |
| Productos | `/products` | CRUD de productos (catálogo) |
| Resúmenes mensuales | `/monthly-summaries` | CRUD de resúmenes mensuales |
| Clientes | `/clients` | CRUD de clientes |
| Pedidos | `/orders` | CRUD de pedidos |

---

## Convenciones generales

- **Content-Type:** `application/json`
- **Errores:** `{"error": "mensaje"}` con código HTTP 400 o 404
- **IDs autogenerados:** formato `YYYYMMDDHHMMSSX` donde X es una letra mayúscula
- **Fechas:** ISO 8601 (`YYYY-MM-DD`)
- **Timestamps:** RFC 3339 (`2026-06-08T22:41:07.862350309Z`)

---

## 1. Autores (`/authors`)

### Crear autor

```bash
curl -s -X POST http://localhost:8080/authors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gabriel García Márquez",
    "biography": "Escritor colombiano, premio Nobel de literatura"
  }'
```

**Campos:**
| Campo | Tipo | Obligatorio |
|-------|------|-------------|
| `name` | string | Sí |
| `biography` | string | Sí |

### Listar autores

```bash
curl -s http://localhost:8080/authors
```

### Obtener autor por ID

```bash
curl -s http://localhost:8080/authors/20260608224107A
```

### Actualizar autor

```bash
curl -s -X PUT http://localhost:8080/authors/20260608224107A \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gabriel García Márquez",
    "biography": "Nobel de literatura 1982"
  }'
```

### Eliminar autor

```bash
curl -s -X DELETE http://localhost:8080/authors/20260608224107A
```

---

## 2. Libros (`/books`)

### Crear libro

```bash
curl -s -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Cien años de soledad",
    "isbn": "978-84-376-0494-7",
    "publication_date": "1967-06-05",
    "genre": ["Realismo mágico", "Novela"],
    "synopsis": "Historia de la familia Buendía en Macondo",
    "photos": ["https://ejemplo.com/portada.jpg"],
    "id_author": "ID_AQUI"
  }'
```

**Campos:**
| Campo | Tipo | Obligatorio |
|-------|------|-------------|
| `title` | string | Sí |
| `isbn` | string | Sí |
| `publication_date` | string | No (pero si se envía, no vacío) |
| `genre` | array de strings | Sí |
| `synopsis` | string | Sí |
| `photos` | array de strings | No |
| `id_author` | string | Sí |

**Nota:** `genre` se pasa como arreglo JSON y se almacena como `TEXT[]` en PostgreSQL.

```bash
# Ejemplo con solo obligatorios
curl -s -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "El principito",
    "isbn": "978-84-7880-952-7",
    "genre": ["Infantil", "Fábula"],
    "synopsis": "Un piloto perdido en el desierto conoce a un pequeño príncipe",
    "id_author": "ID_AQUI"
  }'
```

### Listar / Obtener / Actualizar / Eliminar

Misma estructura que autores:

```bash
curl -s http://localhost:8080/books
curl -s http://localhost:8080/books/{id}
curl -s -X PUT http://localhost:8080/books/{id} -H "Content-Type: application/json" -d '{...}'
curl -s -X DELETE http://localhost:8080/books/{id}
```

---

## 3. Temas de interés (`/topics`)

### Crear tema

```bash
curl -s -X POST http://localhost:8080/topics \
  -H "Content-Type: application/json" \
  -d '{
    "name": ["Literatura", "Poesía"],
    "description": "Temas relacionados con la literatura y poesía"
  }'
```

**Campos:**
| Campo | Tipo | Obligatorio |
|-------|------|-------------|
| `name` | array de strings | Sí |
| `description` | string | No |

**Nota:** `name` se almacena como `TEXT[]`.

---

## 4. Notas (`/notes`)

### Crear nota

```bash
curl -s -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Idea de negocio",
    "content": "Ofrecer catas de café con maridaje de libros",
    "id_book": "ID_LIBRO"
  }`
```

**Campos:**
| Campo | Tipo | Obligatorio |
|-------|------|-------------|
| `title` | string | Sí |
| `content` | string | Sí |
| `id_book` | string | No |

---

## 5. Establecimientos (`/establishments`)

### Crear establecimiento

```bash
curl -s -X POST http://localhost:8080/establishments \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sucursal Central",
    "address": "Managua, Nicaragua",
    "phone": "+505 8888-0000"
  }'
```

**Campos:**
| Campo | Tipo | Obligatorio |
|-------|------|-------------|
| `name` | string | Sí |
| `address` | string | Sí |
| `phone` | string | Sí |

---

## 6. Tipos de movimiento (`/movement-types`)

### Crear tipo de movimiento

```bash
curl -s -X POST http://localhost:8080/movement-types \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Compra proveedor",
    "description": "Entrada por compra a proveedores"
  }'
```

**Campos:**
| Campo | Tipo | Obligatorio |
|-------|------|-------------|
| `name` | string | Sí |
| `description` | string | No |

---

## 7. Movimientos (`/movements`)

### Crear movimiento

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
    "movement_type_id": "ID_TIPO_MOVIMIENTO",
    "observations": "Compra mensual"
  }'
```

**Campos:**
| Campo | Tipo | Obligatorio |
|-------|------|-------------|
| `date` | string (`YYYY-MM-DD`) | Sí |
| `code` | string | Sí |
| `product` | string | Sí |
| `unit` | string | Sí |
| `entrance` | number | Sí |
| `output` | number | Sí |
| `balance` | number | Sí |
| `unit_cost` | number | Sí |
| `valor_value` | number | Sí |
| `movement_type_id` | string (FK) | Sí |
| `observations` | string | No |

**Unidades válidas:** `Kg, Liter, Pound, Grams, Unit`

---

## 8. Productos (`/products`)

### Crear producto

```bash
curl -s -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "code": "CAFE-001",
    "description": "Café latte",
    "category": ["Bebidas", "Calientes"],
    "unit": "Unit"
  }'
```

**Campos:**
| Campo | Tipo | Obligatorio |
|-------|------|-------------|
| `code` | string | Sí |
| `description` | string | Sí |
| `category` | array de strings | Sí |
| `unit` | string | Sí (uno de: `Kg, Liter, Pound, Grams, Unit`) |

**Nota:** `category` se almacena como `TEXT[]`.

---

## 9. Resúmenes mensuales (`/monthly-summaries`)

### Crear resumen mensual

```bash
curl -s -X POST http://localhost:8080/monthly-summaries \
  -H "Content-Type: application/json" \
  -d '{
    "year": 2026,
    "month": 6,
    "total_income": 15000.50,
    "total_expenses": 8200.00,
    "net_balance": 6800.50,
    "notes": "Buen mes, aumento de ventas"
  }'
```

**Campos:**
| Campo | Tipo | Obligatorio |
|-------|------|-------------|
| `year` | number | Sí |
| `month` | number (1-12) | Sí |
| `total_income` | number | Sí |
| `total_expenses` | number | Sí |
| `net_balance` | number | Sí |
| `notes` | string | No |

---

## 10. Clientes (`/clients`)

### Crear cliente

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

**Campos:**
| Campo | Tipo | Obligatorio |
|-------|------|-------------|
| `name_full` | string | Sí |
| `phone` | string | Sí |
| `correo` | string | No |
| `address` | string | Sí |

```bash
# Sin correo (opcional)
curl -s -X POST http://localhost:8080/clients \
  -H "Content-Type: application/json" \
  -d '{
    "name_full": "María López",
    "phone": "+505 7777-8888",
    "address": "León, Nicaragua"
  }'
```

---

## 11. Pedidos (`/orders`)

### Crear pedido

```bash
curl -s -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "order_numeric": "ORD-001",
    "date": "2026-06-08",
    "hour": "14:30",
    "attended_by": "Carlos",
    "client_id": "ID_CLIENTE",
    "details": [
      {
        "code": "P001",
        "product": "Café Latte",
        "quantity": 2,
        "unit_price": 3.50,
        "subtotal": 7.00,
        "discount": 0.50,
        "taxes": 0.91,
        "total": 7.41
      },
      {
        "code": "P002",
        "product": "Croissant",
        "quantity": 1,
        "unit_price": 2.00,
        "subtotal": 2.00,
        "discount": 0,
        "taxes": 0.26,
        "total": 2.26
      }
    ],
    "payment_method": "cash",
    "status": "received",
    "observations": "Sin azúcar"
  }'
```

**Campos del pedido:**
| Campo | Tipo | Obligatorio |
|-------|------|-------------|
| `order_numeric` | string | Sí |
| `date` | string (`YYYY-MM-DD`) | Sí |
| `hour` | string (`HH:MM`) | Sí |
| `attended_by` | string | Sí |
| `client_id` | string (FK) | Sí |
| `details` | array de objetos | Sí (mínimo 1) |
| `payment_method` | string | Sí (ver valores válidos) |
| `status` | string | Sí (ver valores válidos) |
| `observations` | string | No |

**Payment methods válidos:** `cash`, `transfer`, `debit-card`, `credit-card`

**Status válidos:** `received`, `in-preparation`, `ready-for-delivery`, `delivered`, `cancelled`

**Campos de cada `detail`:**
| Campo | Tipo | Obligatorio |
|-------|------|-------------|
| `code` | string | Sí |
| `product` | string | Sí |
| `quantity` | number (> 0) | Sí |
| `unit_price` | number (> 0) | Sí |
| `subtotal` | number | No |
| `discount` | number | No |
| `taxes` | number | No |
| `total` | number | No |

### Actualizar estado del pedido

```bash
curl -s -X PUT "http://localhost:8080/orders/ID_PEDIDO" \
  -H "Content-Type: application/json" \
  -d '{
    "order_numeric": "ORD-001",
    "date": "2026-06-08",
    "hour": "14:30",
    "attended_by": "Carlos",
    "client_id": "ID_CLIENTE",
    "details": [...],
    "payment_method": "cash",
    "status": "in-preparation",
    "observations": "En preparación"
  }'
```

### Listar todos los pedidos

```bash
curl -s http://localhost:8080/orders
```

### Obtener pedido por ID

```bash
curl -s http://localhost:8080/orders/ID_PEDIDO
```

### Eliminar pedido

```bash
curl -s -X DELETE http://localhost:8080/orders/ID_PEDIDO
# Respuesta: 204 No Content
```

---

## Script de ejemplo completo

```bash
#!/bin/bash
BASE="http://localhost:8080"

echo "=== 1. Crear autor ==="
AUTHOR=$(curl -s -X POST $BASE/authors \
  -H "Content-Type: application/json" \
  -d '{"name":"Gabriel García Márquez","biography":"Nobel 1982"}')
AUTHOR_ID=$(echo $AUTHOR | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "Autor ID: $AUTHOR_ID"

echo ""
echo "=== 2. Crear cliente ==="
CLIENT=$(curl -s -X POST $BASE/clients \
  -H "Content-Type: application/json" \
  -d '{"name_full":"Juan Pérez","phone":"+505 8888-9999","address":"Managua"}')
CLIENT_ID=$(echo $CLIENT | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "Cliente ID: $CLIENT_ID"

echo ""
echo "=== 3. Crear pedido ==="
curl -s -X POST $BASE/orders \
  -H "Content-Type: application/json" \
  -d "{\"order_numeric\":\"ORD-001\",\"date\":\"2026-06-08\",\"hour\":\"14:30\",\"attended_by\":\"Carlos\",\"client_id\":\"$CLIENT_ID\",\"details\":[{\"code\":\"P1\",\"product\":\"Café Latte\",\"quantity\":2,\"unit_price\":3.5,\"subtotal\":7,\"discount\":0.5,\"taxes\":0.91,\"total\":7.41}],\"payment_method\":\"cash\",\"status\":\"received\"}" | python3 -m json.tool

echo ""
echo "=== 4. Listar pedidos ==="
curl -s $BASE/orders | python3 -c "import sys,json; d=json.load(sys.stdin); print(f'Total pedidos: {len(d)}')"
```

---

## Validaciones por entidad

| Entidad | Validaciones clave |
|---------|-------------------|
| Autores | `name` y `biography` obligatorios |
| Libros | `title`, `isbn`, `genre` (array), `synopsis`, `id_author` obligatorios |
| Temas | `name` (array) obligatorio |
| Notas | `title`, `content` obligatorios |
| Establecimientos | `name`, `address`, `phone` obligatorios |
| Tipos de movimiento | `name` obligatorio |
| Movimientos | `date`, `code`, `product`, `unit`, `entrance`, `output`, `balance`, `unit_cost`, `valor_value`, `movement_type_id` obligatorios; `unit` debe ser uno de: `Kg, Liter, Pound, Grams, Unit` |
| Productos | `code`, `description`, `category` (array), `unit` obligatorios; `unit` debe ser uno de: `Kg, Liter, Pound, Grams, Unit` |
| Resúmenes mensuales | `year`, `month` (1-12), `total_income`, `total_expenses`, `net_balance` obligatorios |
| Clientes | `name_full`, `phone`, `address` obligatorios; `correo` opcional |
| Pedidos | `order_numeric`, `date`, `hour`, `attended_by`, `client_id`, `details` (mínimo 1), `payment_method`, `status` obligatorios; `payment_method` validado contra lista; `status` validado contra lista |

---

## Códigos de respuesta HTTP

| Código | Significado |
|--------|-------------|
| `200 OK` | GET/PUT exitoso |
| `201 Created` | POST exitoso |
| `204 No Content` | DELETE exitoso |
| `400 Bad Request` | Error de validación (campo faltante, formato inválido, etc.) |
| `404 Not Found` | Recurso no encontrado |
| `405 Method Not Allowed` | Método HTTP no soportado |
