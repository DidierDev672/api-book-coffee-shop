<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Clean_Architecture-000?style=for-the-badge" alt="Clean Architecture">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
</p>

<h1 align="center">⚔️ Casa Atreides API</h1>

<p align="center">
  <em>«Un comienzo es un tiempo muy delicado.»</em> — Princesa Irulan, <em>Dune</em>
</p>

<p align="center">
  API REST escrita en Go que implementa el gobierno digital de tu negocio con la<br>
  misma nobleza y estrategia de la <strong>Casa Atreides</strong> de <em>Dune</em>.
</p>

---

## 📖 Descripción

**Casa Atreides API** es un backend RESTful construido con **Go** y **Clean Architecture** que gestiona un sistema integral de inventario, ventas, pedidos, empresas, clientes y autenticación JWT.

El proyecto está inspirado en la saga *Dune* de Frank Herbert: cada capa del sistema es un feudo con su propósito — dominio, casos de uso, repositorios e infraestructura.

### ✨ Características Principales

| Módulo | Descripción |
|--------|-------------|
| 🔐 **Autenticación JWT** | Login, registro, roles (admin/user), middleware de protección |
| 📚 **Libros & Autores** | CRUD completo con validación de dominio |
| 🏢 **Empresas** | Gestión de empresas con NIT, razón social, dirección fiscal |
| 📦 **Productos** | Catálogo de productos con código, categorías, stock mínimo |
| 🛒 **Pedidos (Orders)** | CRUD de órdenes con resumen financiero y estado |
| 💰 **Ventas (Sales)** | Gestión de ventas con consecutivos, descuentos, IVA |
| 🚚 **Envíos (Shipments)** | Registro de envíos con destinatario y documento fuente |
| 📊 **Movimientos** | Entradas/salidas de inventario con tipos y balance |
| 🏭 **Proveedores** | Directorio de proveedores con código y actividad económica |
| 📋 **Notas** | CRUD de notas clasificadas por tipo |
| 🏷️ **Temas (Topics)** | Clasificación de contenido por categoría |
| 📍 **Direcciones** | Direcciones principales por empresa |
| 💼 **Actividades Económicas** | Registro de actividades económicas por empresa |
| 🧾 **Información Tributaria** | Régimen tributario, IVA, retención |
| 🗺️ **Bodegas (Wineries)** | Gestión de bodegas por empresa |
| 📈 **Resúmenes Mensuales** | Resumen mensual de stock por producto |
| 📜 **Historial de Inventario** | Auditoría completa de cada evento del inventario |
| 👥 **Usuarios** | CRUD de usuarios con hash bcrypt y tokens JWT |
| 🧪 **Equipment Server** | Microservicio independiente para gestión de equipamiento |

---

## 🛠️ Stack Tecnológico

| Tecnología | Versión | Rol |
|------------|---------|-----|
| [Go](https://go.dev/) | 1.26 | Lenguaje de programación |
| [PostgreSQL](https://www.postgresql.org/) | 16 | Base de datos relacional |
| [JWT](https://github.com/golang-jwt/jwt) | v5 | Autenticación por tokens |
| [bcrypt](https://pkg.go.dev/golang.org/x/crypto) | — | Hash de contraseñas |
| [CORS](https://github.com/rs/cors) | — | Configuración de CORS |
| [lib/pq](https://github.com/lib/pq) | — | Driver PostgreSQL |

---

## 📋 Prerrequisitos

- **Go** ≥ 1.21
- **PostgreSQL** ≥ 14
- **Git**

---

## 🚀 Instalación y Ejecución

### 1. Clonar el repositorio

```bash
git clone <url-del-repositorio>
cd api-book-coffee-shop
```

### 2. Instalar dependencias

```bash
go mod tidy
```

### 3. Configurar base de datos

Crea la base de datos PostgreSQL (las migraciones se ejecutan automáticamente al iniciar):

```sql
CREATE DATABASE coffee_book;
```

> ⚠️ Las credenciales por defecto están en `internal/config/config.go`. Para personalizar, edita `DefaultPostgresConfig()` o configura variables de entorno.

### 4. Configurar JWT Secret (opcional)

```bash
export JWT_SECRET="tu-secreto-seguro"
```

Si no se define, se usa el valor por defecto: `book-coffee-shop-dev-secret`

### 5. Ejecutar el servidor

```bash
# Servidor API principal (puerto 8080)
go run ./cmd/api/

# Servidor Gin (puerto alternativo)
go run ./cmd/gin-server/

# Servidor Equipment (microservicio)
go run ./cmd/equipment/
```

### 📜 Scripts disponibles

| Comando | Descripción |
|---------|-------------|
| `go run ./cmd/api/` | Servidor API principal en `:8080` |
| `go run ./cmd/gin-server/` | Servidor con Gin framework |
| `go run ./cmd/equipment/` | Microservicio de equipamiento |
| `go build -o bin/api.exe ./cmd/api/` | Compilar binario |
| `go test ./...` | Ejecutar tests |
| `go vet ./...` | Verificar código |

---

## 📡 Endpoints Principales

### 🔐 Autenticación

| Método | Ruta | Descripción | Auth |
|--------|------|-------------|------|
| `POST` | `/auth/register` | Registrar usuario | ❌ |
| `POST` | `/auth/login` | Iniciar sesión | ❌ |
| `GET` | `/users` | Listar usuarios | ✅ Admin |
| `GET` | `/users/{id}` | Obtener usuario | ✅ |
| `PUT` | `/users/{id}` | Actualizar usuario | ✅ |

### 📚 Contenido

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET/POST` | `/authors`, `/authors/` | CRUD Autores |
| `GET/POST` | `/books`, `/books/` | CRUD Libros |
| `GET/POST` | `/topics`, `/topics/` | CRUD Temas |
| `GET/POST` | `/notes`, `/notes/` | CRUD Notas |

### 🏢 Negocio

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET/POST` | `/companies`, `/companies/` | CRUD Empresas |
| `GET/POST` | `/clients`, `/clients/` | CRUD Clientes |
| `GET/POST` | `/providers`, `/providers/` | CRUD Proveedores |
| `GET/POST` | `/establishments`, `/establishments/` | CRUD Establecimientos |

### 📦 Inventario

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET/POST` | `/products`, `/products/` | CRUD Productos |
| `GET/POST` | `/product-entries`, `/product-entries/` | Entradas de producto |
| `GET/POST` | `/movements`, `/movements/` | Movimientos de inventario |
| `GET/POST` | `/movement-types`, `/movement-types/` | Tipos de movimiento |
| `GET/POST` | `/wineries`, `/wineries/` | Bodegas |
| `GET/POST` | `/history`, `/history/` | Historial de inventario |

### 🛒 Ventas y Pedidos

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET/POST` | `/orders`, `/orders/` | CRUD Pedidos |
| `GET/POST` | `/sales`, `/sales/` | CRUD Ventas |
| `GET/POST` | `/shipments`, `/shipments/` | CRUD Envíos |

### 📊 Datos Fiscales

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET/POST` | `/tax-information`, `/tax-information/` | Información tributaria |
| `GET/POST` | `/economic-activities`, `/economic-activities/` | Actividades económicas |
| `GET/POST` | `/main-addresses`, `/main-addresses/` | Direcciones principales |
| `GET/POST` | `/monthly-summaries`, `/monthly-summaries/` | Resúmenes mensuales |

---

## 📁 Estructura de Carpetas

```
api-book-coffee-shop/
├── cmd/                            # 🚀 Puntos de entrada
│   ├── api/                        #   Servidor API principal (net/http)
│   │   └── main.go                 #     Entry point + migraciones + DI
│   ├── gin-server/                 #   Servidor alternativo con Gin
│   │   └── main.go
│   └── equipment/                  #   Microservicio de equipamiento
│       └── main.go
│
├── internal/                       # 🏛️ Lógica interna (Clean Architecture)
│   ├── config/                     #   ⚙️  Configuración (Postgres, JWT)
│   │   └── config.go
│   ├── database/                   #   🗄️  Utilidades de base de datos
│   │   └── postgres.go
│   ├── domain/                     #   📦 Entidades de dominio (22 entidades)
│   │   ├── author.go               #     Author, Book, Client, Company...
│   │   ├── book.go                 #     Order, Sale, Product, Movement...
│   │   ├── user.go                 #     User, Provider, Shipment...
│   │   └── ...
│   ├── repository/                 #   📋 Interfaces de repositorio (25 interfaces)
│   │   ├── author_repository.go
│   │   ├── book_repository.go
│   │   ├── user_repository.go
│   │   └── ...
│   ├── usecase/                    #   ⚙️  Casos de uso (22 usecases)
│   │   ├── auth_usecase.go         #     Auth, Author, Book, Order...
│   │   ├── author_usecase.go
│   │   ├── order_usecase.go
│   │   └── ...
│   ├── handler/                    #   🖥️  Handlers HTTP (20 handlers)
│   │   ├── auth_handler.go
│   │   ├── book_handler.go
│   │   ├── order_handler.go
│   │   └── ...
│   ├── infrastructure/             #   🔌 Implementaciones (repositorios concretos)
│   │   ├── postgres_author_repository.go
│   │   ├── postgres_book_repository.go
│   │   ├── jwt_token_service.go
│   │   ├── bcrypt_password_hasher.go
│   │   └── ...
│   ├── middleware/                  #   🛡️  Middleware (Auth, Recovery, Validation)
│   │   ├── auth.go
│   │   ├── authorization.go
│   │   ├── RecoveryMiddleware.go
│   │   └── validation.go
│   ├── models/                     #   📨 Modelos de request/response
│   │   └── requests.go
│   ├── gin_auth/                   #   🔐 Módulo de auth (Gin)
│   │   ├── handler.go
│   │   ├── jwt.go
│   │   ├── middleware.go
│   │   └── models.go
│   └── utils/                      #   🧰 Utilidades
│       ├── errors.go
│       ├── response.go
│       ├── validation.go
│       └── TryExecut.go
│
├── docs/                           # 📚 Documentación
├── graphify-out/                   # 🕸️  Grafo de conocimiento
├── go.mod                          # Dependencias Go
├── go.sum                          # Lock de dependencias
└── README.md                       # Este archivo
```

### 🏗️ Arquitectura (Clean Architecture)

```
Handler (HTTP) → UseCase (Lógica) → Repository (Interface) → Infrastructure (PostgreSQL)
      ↑                                                              ↓
   Request                                                    Response JSON
```

**Capas:**
- **Domain** — Entidades puras, sin dependencias externas
- **Repository** — Interfaces que definen contratos de acceso a datos
- **UseCase** — Lógica de negocio y validaciones
- **Handler** — Endpoints HTTP, parseo de requests, respuestas
- **Infrastructure** — Implementaciones concretas (PostgreSQL, JWT, bcrypt)
- **Middleware** — Autenticación, recuperación de panics, validación

---

## 🧪 Ejemplos de Uso

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "abc123",
    "name_full": "Paul Atreides",
    "email": "user@example.com"
  }
}
```

### Crear Autor

```bash
curl -X POST http://localhost:8080/authors \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"Frank Herbert","country":"USA","genres":["Sci-Fi"],"birth_day":"1920-10-08"}'
```

### Listar Productos

```bash
curl http://localhost:8080/products \
  -H "Authorization: Bearer <token>"
```

### Crear Pedido

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "order_numeric": "ORD-001",
    "date": "2026-06-24",
    "status": "PENDING",
    "details": [{"product_id":"p1","quantity":10,"unit_cost":5.0}]
  }'
```

---

## 🔑 Variables de Entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `JWT_SECRET` | Secreto para firmar tokens JWT | `book-coffee-shop-dev-secret` |

> Las credenciales de PostgreSQL se configuran en `internal/config/config.go`.

---

## 🤝 Contribuir

1. **Fork** el repositorio
2. **Crear** una rama para tu feature:
   ```bash
   git checkout -b feature/nueva-feature
   ```
3. **Confirmar** tus cambios:
   ```bash
   git commit -m "Add: nueva feature"
   ```
4. **Push** a la rama:
   ```bash
   git push origin feature/nueva-feature
   ```
5. **Abrir** un Pull Request

### Convenciones de Commits

| Prefijo | Descripción |
|---------|-------------|
| `Add:` | Nueva funcionalidad |
| `Fix:` | Corrección de bug |
| `Update:` | Mejora o refactor |
| `Docs:` | Documentación |
| `Test:` | Tests |

---

## 📜 Licencia

Este proyecto está bajo la licencia **MIT**. Ver el archivo [LICENSE](LICENSE) para más detalles.

---

<p align="center">
  <em>«El poder de gobernar está en los pequeños detalles.»</em> — Thufir Hawat
</p>

<p align="center">
  Desarrollado con ⚔️ por la Casa Atreides
</p>
