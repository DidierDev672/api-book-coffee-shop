# Atreides — API REST de la Casa Atreides

> *"Un comienzo es un tiempo muy delicado."* — Princesa Irulan, *Dune*
>
> *"El poder de gobernar está en los pequeños detalles."* — Thufir Hawat, Mentat de la Casa Atreides

API REST escrita en Go que implementa el gobierno digital de tu negocio con la
misma nobleza y estrategia de la **Casa Atreides** de *Dune*.

Así como la Casa Atreides gobernó Caladan (el dominio de los datos puros) y
luego conquistó Arrakis (la base de datos donde fluye la especia del inventario),
esta API está construida con **Clean Architecture** — cada capa es un feudo con
su propósito: dominio, casos de uso, repositorios e infraestructura.

| Concepto Dune | Componente Atreides API |
|---|---|
| **Casa Atreides** | El sistema completo — 31 endpoints, 20 handlers |
| **Caladan** (mundo oceánico) | `internal/domain/` — entidades puras sin dependencias externas |
| **Arrakis** (planeta desierto) | `internal/infrastructure/` — repositorios PostgreSQL donde se extraen los datos |
| **Mentat** (Thufir Hawat) | `internal/config/` + `internal/usecase/` — la lógica que calcula cada movimiento |
| **Maestro de Armas** (Gurney Halleck) | `internal/middleware/` — AuthMiddleware, RecoveryMiddleware, validación |
| **Consejo de Guerra** | `internal/handler/` — cada handler decide cómo responder al Imperio (el cliente HTTP) |
| **Especia Melange** | Los datos — el recurso más valioso que debe fluir sin interrupción |
| **Sardaukars** | JSONB + PostgreSQL — implacables, eficientes, almacenan todo |
| **Archivos de la Casa** (Crónicas de Irulan) | `internal/domain/InventoryHistory` + `internal/infrastructure/PostgresInventoryHistoryRepository` + `internal/usecase/HistoryService` — cada evento del inventario queda registrado como un pergamino en la Biblioteca de Caladan |

> Inspirado en el linaje Atreides: un linaje que crece, se adapta y se expande
> a nuevos territorios. Desde una librería-cafetería hasta cualquier dominio
> de negocio que comparta el mismo núcleo fiscal, inventario y autenticación.

---

## Descripción general — *Los Dominios de la Casa*

> *"La Casa Atreides gobierna donde otros apenas sobreviven."*

Sistema backend monolitico que expone **29 endpoints HTTP** para administrar
los dominios de tu negocio, como la Casa Atreides administraba los recursos
de Caladan y luego las riquezas de Arrakis:

| Dominio Atreides | Recurso API | Descripción |
|---|---|---|
| **El Gran Almacén** | Productos, Bodegas, Entradas | Control de inventario con JSONB y resúmenes financieros |
| **El Feudo Fiscal** | Compañías, Direcciones, Actividades Económicas, Info Tributaria | Gestión empresarial colombiana (NIT, regímenes, IVA) |
| **El Consejo de Órdenes** | Órdenes, Despachos, Clientes | Ventas con métodos de pago y estados tipo Dune |
| **La Guardia del Palacio** | Auth (JWT + bcrypt) | Registro y login, middleware de autenticación para rutas protegidas |
| **La Biblioteca de Caladan** | Autores, Libros, Tópicos, Notas | Módulo editorial (legacy, como los archivos históricos de la Casa) |

> Detalle de los endpoints y operaciones en la sección [Catálogo completo de rutas](#catálogo-completo-de-rutas).

### Stack técnico — *El Arsenal de la Casa*

> Como el Mentat prepara sus cálculos y el Maestro de Armas afila sus espadas,
> cada tecnología en este stack fue elegida con precisión.

| Capa | Tecnología | Rol en la Casa Atreides |
|------|-----------|------------------------|
| Lenguaje | **Go 1.26.3** | El carácter del Duque — rápido, confiable, sin concesiones |
| Base de datos | **PostgreSQL 15+** (JSONB, arrays, `uuid-ossp`) | Arrakis — almacena la especia (datos) en su forma más pura |
| Autenticación | **JWT HS256** + **bcrypt** | El Sello de la Casa — solo quienes portan el token correcto cruzan el puente |
| Driver BD | `github.com/lib/pq` v1.12.3 | Los gusanos de arena — transportan datos desde las profundidades |
| Validación | `go-playground/validator/v10` | El Catecismo Fremen — cada campo debe cumplir su ley |
| CORS | `github.com/rs/cors` v1.11.1 | El Muro de Shield — solo los orígenes permitidos traspasan |
| Arquitectura | Clean Architecture | La Estructura del Landsraad — cada capa conoce su lugar |

### Estructura del proyecto

```
cmd/api/
  main.go                  → Punto de entrada, wiring DI, migraciones inline
internal/
  config/config.go         → Config Postgres + JWT (con fallback hardcodeado)
  database/postgres.go     → EnsureDatabaseExists (crea DB si no existe)
  domain/                  → 20 structs de dominio (sin dependencias externas)
  repository/              → 21 interfaces + 2 servicios (TokenService, PasswordHasher)
  usecase/                 → 20 implementaciones de lógica de negocio
  handler/                 → 20 handlers HTTP activos
  handlers/                → 1 handler legacy (AuthHandler con middleware roles)
  infrastructure/          → 20 impls PostgreSQL + bcrypt + JWT
  middleware/              → Recovery, Auth (JWT), Authorization (roles), Validation
  models/                  → DTOs con tags de validación (LoginRequest, RegisterRequest)
  utils/                   → response.go, context.go, validation.go, TryExecut.go
```

---

## Requisitos previos

| Herramienta | Versión | Propósito |
|-------------|---------|-----------|
| Go | **1.26.3+** | Compilación y ejecución |
| PostgreSQL | **15+** | Base de datos |
| Git | Cualquiera | Control de versiones |

### Dependencias Go

```
github.com/golang-jwt/jwt/v5       v5.3.1
github.com/lib/pq                   v1.12.3
github.com/rs/cors                  v1.11.1
github.com/go-playground/validator  v10.30.3
golang.org/x/crypto                v0.53.0
```

---

## Instalación

### 1. Clonar e instalar dependencias

```powershell
git clone <repo-url> book-coffee-shop
cd book-coffee-shop
go mod tidy
go mod download
```

### 2. Configurar PostgreSQL

El sistema usa credenciales hardcodeadas por defecto en `internal/config/config.go`:

```go
func DefaultPostgresConfig() PostgresConfig {
    return PostgresConfig{
        Host:     "localhost",
        Port:     5432,
        User:     "postgres",
        Password: "1234567890",
        DBName:   "coffee_book",
    }
}
```

Asegúrate de que PostgreSQL esté corriendo en `localhost:5432` con usuario `postgres` y contraseña `1234567890`. La base de datos `coffee_book` se crea automáticamente al iniciar.

### 3. Configurar JWT Secret (opcional)

Por defecto usa `"book-coffee-shop-dev-secret"`. Para sobrescribir:

```powershell
$env:JWT_SECRET="mi-secreto-personalizado"
```

### 4. Compilar y ejecutar

```powershell
go build -o server.exe ./cmd/api
./server.exe
```

O directamente:

```powershell
go run ./cmd/api
```

El servidor inicia en `http://localhost:8080` con logs en consola.

### 5. Verificar compilación

```powershell
go build ./...
go vet ./...
```

---

## Guía de uso

### Flujo de autenticación

#### 1. Registrar usuario

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name_full": "Juan Pérez",
    "phone": "3001234567",
    "id_number": "1234567890",
    "date_of_birth": "1990-01-15",
    "email": "juan@example.com",
    "password": "miPassword123"
  }'
```

Respuesta:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "1749372100_001",
    "name_full": "Juan Pérez",
    "email": "juan@example.com",
    "roles": null
  }
}
```

> **Nota:** El registro también genera un token JWT automáticamente, igual que el login. No es necesario iniciar sesión después de registrar.

#### 2. Iniciar sesión

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "juan@example.com", "password": "miPassword123"}'
```

Respuesta:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "1749372100_001",
    "name_full": "Juan Pérez",
    "email": "juan@example.com",
    "roles": null
  }
}
```

#### 3. Usar el token en rutas protegidas

```bash
curl http://localhost:8080/products \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

### CRUD de productos

```bash
# Crear producto (requiere auth)
curl -X POST http://localhost:8080/products \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "id-de-compañía",
    "supplier_id": "id-de-proveedor",
    "name": "Café Especial",
    "product_code": "CAFE-001",
    "categories": ["Bebidas", "Café"],
    "unit": "Kg",
    "quantity": 50,
    "minimum_stock": 10,
    "winery_id": "id-de-bodega"
  }'

# Listar productos (por compañía opcional)
curl "http://localhost:8080/products?company_id=abc123" \
  -H "Authorization: Bearer <token>"

# Obtener por ID
curl http://localhost:8080/products/abc123 \
  -H "Authorization: Bearer <token>"

# Actualizar
curl -X PUT http://localhost:8080/products/abc123 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Café Premium", "product_code": "CAFE-002", "unit": "Kg", "quantity": 50, "minimum_stock": 10, "winery_id": "id-de-bodega", "categories": ["Bebidas", "Café"]}'

# Eliminar
curl -X DELETE http://localhost:8080/products/abc123 \
  -H "Authorization: Bearer <token>"
```

### Catálogo completo de rutas

| Método | Ruta | Protegida | Descripción |
|--------|------|-----------|-------------|
| GET/POST | `/authors` | Sí | CRUD autores |
| GET/PUT/DEL | `/authors/{id}` | Sí |  |
| GET/POST | `/books` | Sí | CRUD libros |
| GET/PUT/DEL | `/books/{id}` | Sí |  |
| GET/POST | `/topics` | Sí | CRUD tópicos |
| GET/PUT/DEL | `/topics/{id}` | Sí |  |
| GET/POST | `/notes` | Sí | CRUD notas |
| GET/PUT/DEL | `/notes/{id}` | Sí |  |
| GET/POST | `/establishments` | Sí | CRUD establecimientos |
| GET/PUT/DEL | `/establishments/{id}` | Sí |  |
| GET/POST | `/movement-types` | Sí | CRUD tipos de movimiento |
| GET/PUT/DEL | `/movement-types/{id}` | Sí |  |
| GET/POST | `/movements` | Sí | CRUD movimientos (valida FK movement_type) |
| GET/PUT/DEL | `/movements/{id}` | Sí |  |
| GET/POST | `/products` | Sí | CRUD productos (+ filtro `?company_id=`) |
| GET/PUT/DEL | `/products/{id}` | Sí |  |
| GET/POST | `/monthly-summaries` | Sí | CRUD resúmenes mensuales |
| GET/PUT/DEL | `/monthly-summaries/{id}` | Sí |  |
| GET/POST | `/clients` | Sí | CRUD clientes |
| GET/PUT/DEL | `/clients/{id}` | Sí |  |
| GET/POST | `/companies` | Sí | CRUD compañías (+ `GET /companies/user/{id}`) |
| GET/PUT/DEL | `/companies/{id}` | Sí |  |
| GET/POST | `/main-addresses` | Sí | CRUD direcciones principales |
| GET/PUT/DEL | `/main-addresses/{id}` | Sí |  |
| GET/POST | `/tax-information` | Sí | CRUD información tributaria |
| GET/PUT/DEL | `/tax-information/{id}` | Sí |  |
| GET/POST | `/economic-activities` | Sí | CRUD actividades económicas |
| GET/PUT/DEL | `/economic-activities/{id}` | Sí |  |
| GET/POST | `/providers` | Sí | CRUD proveedores |
| GET/PUT/DEL | `/providers/{id}` | Sí |  |
| GET/POST | `/product-entries` | Sí | CRUD entradas de producto (JSONB) |
| GET/PUT/DEL | `/product-entries/{id}` | Sí |  |
| GET/POST | `/orders` | Sí | CRUD órdenes (JSONB) |
| GET/PUT/DEL | `/orders/{id}` | Sí |  |
| GET/POST | `/wineries` | Sí | CRUD bodegas (+ filtro `?company_id=`) |
| GET/PUT/DEL | `/wineries/{id}` | Sí |  |
| POST | `/auth/register` | **No** | Registro de usuario |
| POST | `/auth/login` | **No** | Inicio de sesión |
| GET | `/history` | Sí | Listar todo el historial de inventario (Crónicas de la Casa) |
| GET | `/history/{document_type}/{document_id}` | Sí | Historial filtrado por documento (Archivos de un feudo específico) |
| GET | `/users` | Sí | Listar todos los usuarios |
| GET/PUT | `/users/{id}` | Sí | Obtener/actualizar perfil |

---

## Documentación de funciones

### Capa de infraestructura — `internal/infrastructure/`

#### `bcrypt_password_hasher.go`

```go
func NewBcryptPasswordHasher() *BcryptPasswordHasher
```

Crea un hasher con costo `bcrypt.DefaultCost`.

```go
func (h *BcryptPasswordHasher) Hash(password string) (string, error)
```

Retorna el hash bcrypt del password. Error si el password es muy largo (>72 bytes).

```go
func (h *BcryptPasswordHasher) Compare(hash, password string) error
```

Compara un hash contra un password. Retorna `nil` si coinciden, error en caso contrario.

#### `jwt_token_service.go`

```go
func NewJWTTokenService(secret string) *JWTTokenService
```

Crea servicio JWT con clave HMAC-SHA256.

```go
func (s *JWTTokenService) Generate(userID string) (string, error)
```

Genera token JWT con `sub: userID`, expiración 24h, firmado HS256.

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `userID` | `string` | ID del usuario para el claim `sub` |

Retorna el token JWT string o error si la firma falla.

```go
func (s *JWTTokenService) Validate(tokenString string) (string, error)
```

Valida un token JWT. Retorna el `userID` del claim `sub` o error si el token es inválido/expirado/firmado con método distinto a HMAC.

#### `PostgresXxxRepository` (patrón común, 19 impls)

Cada repositorio sigue el mismo patrón:

```go
func NewPostgresXxxRepository(db *sql.DB) *XxxRepository

func (r *XxxRepository) Create(entity *domain.Xxx) error
func (r *XxxRepository) GetByID(id string) (*domain.Xxx, error)
func (r *XxxRepository) GetAll() ([]*domain.Xxx, error)
func (r *XxxRepository) Update(entity *domain.Xxx) error
func (r *XxxRepository) Delete(id string) error
```

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `db` | `*sql.DB` | Conexión a PostgreSQL inyectada |

**Métodos adicionales específicos:**

| Repositorio | Método extra |
|-------------|-------------|
| `PostgresCompanyRepository` | `GetByNIT(nit)`, `GetByUserID(userID)` |
| `PostgresProductRepository` | `GetByCompanyID(companyID)` |
| `PostgresProviderRepository` | `GetByCode(code)` |
| `PostgresUserRepository` | `GetByEmail(email)`, `GetByAuthToken(token)`, `UpdateAuthToken(id, token)`, `Count()` |

**Manejo de datos JSONB** (`PostgresOrderRepository`, `PostgresProductEntryRepository`):

```go
detailsJSON, err := json.Marshal(order.Details)    // Al guardar
json.Unmarshal(detailsJSON, &order.Details)         // Al leer
```

**Manejo de arrays PostgreSQL** (`PostgresAuthorRepository`, `PostgresProductRepository`):

```go
pq.Array(author.Genres)      // Al guardar
pq.Array(&genres)            // Al leer
```

**Repositorio de Historial — `PostgresInventoryHistoryRepository`:**

> *"Lo que no se registra, no ha ocurrido."*
> — Proverbio Fremen adaptado por la Casa Atreides

```go
func NewPostgresInventoryHistoryRepository(db repository.DBTX) *PostgresInventoryHistoryRepository
```

Crea un repositorio que lee y escribe en la tabla `inventory_history`. Cada
pergamino se almacena como una fila en Arrakis (PostgreSQL).

| Método | SQL | Descripción |
|--------|-----|-------------|
| `Create(event)` | `INSERT INTO inventory_history (...)` | Graba un nuevo pergamino en los Archivos |
| `GetByDocument(docType, docID)` | `SELECT ... WHERE document_type=$1 AND document_id=$2 ORDER BY event_date DESC` | Busca eventos por feudo y documento |
| `GetAll()` | `SELECT ... ORDER BY event_date DESC` | Abre el Libro Mayor de la Casa |

Los campos `previous_data` y `new_data` viajan como `JSONB` en PostgreSQL — la
Especia se almacena en su forma más pura dentro del desierto digital.

**Helper compartido:**

```go
func nullIfEmpty(s string) *string
// Retorna nil si s == "", puntero a s en caso contrario
// Ubicación: postgres_book_repository.go:143
```

---

### Capa de middleware — `internal/middleware/`

#### `RecoveryMiddleware`

```go
func RecoveryMiddleware(next http.Handler) http.Handler
```

Recupera panics, loggea el stack trace con `log.Printf`, retorna 500.

#### `NewAuthMiddleware`

```go
func NewAuthMiddleware(
    tokenService repository.TokenService,
    userRepo repository.UserRepository,
    publicPaths ...string,
) func(http.Handler) http.Handler
```

Middleware de autenticación JWT. Flujo:

1. OPTIONS → pasa (CORS preflight)
2. Path en `publicPaths` → pasa (sin auth)
3. Extrae Bearer token del header `Authorization`
4. Valida JWT → obtiene `userID`
5. Busca usuario por `GetByAuthToken(token)`
6. Verifica que `user.ID == userID`
7. Inyecta `*domain.User` en `context.Context`
8. Llama al siguiente handler

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `tokenService` | `repository.TokenService` | Servicio de validación JWT |
| `userRepo` | `repository.UserRepository` | Repositorio para lookup por token |
| `publicPaths` | `...string` | Rutas que omiten auth (ej: `/auth/register`) |

**Errores:** Retorna 401 con JSON `{"error":"...","code":401}` si:
- Token ausente
- Token inválido/expirado
- Usuario no encontrado por token
- Discrepancia entre userID del token y usuario en BD

#### `RequireRoles`

```go
func RequireRoles(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc
```

Middleware de autorización por roles. Lee usuario del contexto (inyectado por `AuthMiddleware`). Retorna 401 si no hay usuario en contexto, 403 si el usuario no tiene ningún rol permitido.

#### `RequirePermission`

```go
func RequirePermission(permission string) func(http.HandlerFunc) http.HandlerFunc
```

Middleware de autorización granular. Usa `user.HasPermission()`. Retorna 401/403.

#### `ValidatePayload`

```go
func ValidatePayload(target any) func(http.Handler) http.Handler
```

Decodifica JSON del body en `target` y valida usando `go-playground/validator`. Retorna 400 con detalles de validación si falla.

---

### Capa de casos de uso — `internal/usecase/`

Patrón general (18 módulos):

```go
type XxxUseCase interface {
    // Métodos CRUD específicos del dominio
}

type xxxUseCase struct {
    repo repository.XxxRepository
}

func NewXxxUseCase(repo repository.XxxRepository) XxxUseCase
```

Cada usecase implementa las interfaces definidas y contiene su propia función `validateXxxFields()` privada con reglas de negocio.

#### `authUseCase` (el más complejo)

```go
func NewAuthUseCase(
    repo repository.UserRepository,
    hasher repository.PasswordHasher,
    tokens repository.TokenService,
) AuthUseCase
```

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `repo` | `UserRepository` | CRUD de usuarios |
| `hasher` | `PasswordHasher` | bcrypt hashing |
| `tokens` | `TokenService` | JWT generación/validación |

```go
func (uc *authUseCase) Register(token, nameFull, phone, idNumber, dateOfBirth, email, password string) (*domain.User, string, error)
```

Valida campos (incluyendo regex email, formato fecha `YYYY-MM-DD`, password ≥8 chars, id numérico), verifica email único, hashea password con bcrypt, crea usuario, genera nuevo JWT y lo persiste en BD, retorna `(user, token, error)`.

```go
func (uc *authUseCase) Login(_ string, email, password string) (*domain.User, string, error)
```

Busca por email, verifica password con bcrypt, genera nuevo JWT, almacena en BD, retorna `(user, token, error)`.

```go
func (uc *authUseCase) verifyToken(token string) error
```

Método privado que valida JWT + verifica que el token exista en BD (`GetByAuthToken`).

#### `HistoryService` — *Las Crónicas de la Casa*

> *"Un pueblo sin historia es como un hombre sin memoria."*
> — Leto Atreides II, *God Emperor of Dune*

El `HistoryService` es el **Archivero de la Casa Atreides**. No es un CRUD
tradicional — solo escribe y consulta. Así como la Princesa Irulan documentó
cada movimiento del imperio, este servicio registra cada evento del inventario
para que el Mentat (y el Dueño del negocio) puedan reconstruir el pasado.

```go
func NewHistoryService(db *sql.DB, makeRepo repository.HistoryRepoFactory) *HistoryService
```

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `db` | `*sql.DB` | Conexión a PostgreSQL — el ojo del Archivero |
| `makeRepo` | `HistoryRepoFactory` | Fábrica que crea repositorios — como un Escriba que prepara sus pergaminos |

**Métodos públicos:**

```go
func (s *HistoryService) LogEvent(tx repository.DBTX, eventType domain.InventoryEventType,
    userID, companyID, documentID, documentType, description, ipAddress string,
    previousData, newData interface{}) error
```

Registra un evento en el **Gran Pergamino**. El Mentat serializa `previousData`
y `newData` a JSON automáticamente (como una visión del antes y después de un
cálculo). Retorna `nil` si el evento quedó grabado en los Archivos.

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `tx` | `DBTX` | Transacción activa — el Escriba escribe con tinta indeleble |
| `eventType` | `InventoryEventType` | Tipo de decreto imperial (ver tabla abajo) |
| `userID` | `string` | ID del Noble que ejecutó la acción |
| `companyID` | `string` | Feudo al que pertenece el evento |
| `documentID` | `string` | ID del documento afectado (orden, entrada, etc.) |
| `documentType` | `string` | Tipo de documento ("order", "product-entry", "shipment", etc.) |
| `description` | `string` | Narración del evento — lo que el Cronista escribe |
| `ipAddress` | `string` | Dirección desde donde se emitió la orden |
| `previousData` | `interface{}` | Estado anterior del recurso (puede ser `nil`) |
| `newData` | `interface{}` | Estado posterior del recurso (puede ser `nil`) |

```go
func (s *HistoryService) LogRelation(tx repository.DBTX, orderID, shipmentID,
    userID, companyID, ipAddress string) error
```

Registra la **Alianza entre dos Casas** — la relación entre una orden y su
despacho (`shipment`). Crea dos eventos: uno asociado a la orden y otro al
despacho, como un tratado bilateral en el Landsraad.

```go
func (s *HistoryService) LogStockUpdate(tx repository.DBTX, productCode string,
    previousStock, newStock float64, userID, companyID, ipAddress string) error
```

Registra un **ajuste de especia** en los silos. Método de conveniencia que
envuelve `LogEvent` con `EventType = STOCK_UPDATED`. Útil cuando el Mentat
recalcula el inventario después de una tormenta de arena (un error humano).

```go
func (s *HistoryService) GetByDocument(documentType, documentID string) ([]*domain.InventoryHistory, error)
```

Consulta los Archivos por **feudo y pergamino**. Retorna todos los eventos
asociados a un documento específico, ordenados del más reciente al más antiguo
(como los anales de Irulan, del presente hacia el pasado).

```go
func (s *HistoryService) GetAll() ([]*domain.InventoryHistory, error)
```

Abre el **Gran Libro de la Casa** y devuelve todos los eventos registrados.
Ordenados por fecha descendente. Útil para el Dashboard del Duque.

#### Casos de uso con lógica adicional

| Usecase | Complejidad |
|---------|-------------|
| `productUseCase` | Valida unidad contra `validUnits = ["Kg","Litro","Libra","Gramos","Unidad"]` |
| `wineryUseCase` | Valida área contra `["Tienda","Almacén","Cafetería","Otro"]` y unidades contra `["Unidades","Cajas","Litros","Kilogramos"]` (ambos en español) |
| `productEntryUseCase` | Valida `movementType` contra 5 valores, cada detail requiere `code,product,unit,quantity>0,unitCost>=0` |
| `orderUseCase` | Valida `paymentMethod` (4 valores) y `status` (5 valores), cada detail requiere `code,product,quantity>0,unitPrice>0` |
| `providerUseCase` | Unicidad de `code` (tanto en Create como Update) |
| `companyUseCase` | Unicidad de `NIT` (tanto en Create como Update) |
| `movementUseCase` | Doble dependencia: valida FK `movementTypeID` contra `MovementTypeRepository` |
| `movementTypeUseCase` | Solo valida `name` no vacío |
| `monthlySummaryUseCase` | Solo valida `product` no vacío |

**Helper compartido `generateID()`:**

```go
// Ubicación: author_usecase.go:109 (primer usecase del proyecto)
func generateID() string
// Formato: timestamp_unix_milisegundos_númeroSecuencial
// Ejemplo: "1749372100_001"
// Usado por los 19 usecases
```

---

### Capa de handlers HTTP — `internal/handler/`

Patrón general (18 módulos):

```go
type XxxHandler struct {
    uc usecase.XxxUseCase
}

func NewXxxHandler(uc usecase.XxxUseCase) *XxxHandler

func (h *XxxHandler) Handle(w http.ResponseWriter, r *http.Request)
```

Cada `Handle` parsea la ruta, hace routing por método HTTP y delega a métodos privados `create/getByID/getAll/update/delete`. Formato de respuesta consistente:

```go
func writeJSON(w http.ResponseWriter, data any, status int)
func writeError(w http.ResponseWriter, msg string, status int)
```

#### `AuthHandler` (excepción — 4 métodos públicos)

```go
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request)
// POST /auth/register — registra usuario, genera JWT, usa TryExecute con goroutine, retorna {token, user}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request)
// POST /auth/login — login, retorna {token, user}
```

#### `CompanyHandler` (excepción — ruta adicional)

```go
func (h *CompanyHandler) Handle(w http.ResponseWriter, r *http.Request)
// Maneja: GET /companies, GET /companies/{id}, GET /companies/user/{userID}
//         POST /companies, PUT /companies/{id}, DELETE /companies/{id}
```

#### `ProductHandler` (excepción — query param)

```go
// GET /products?company_id=xxx filtra por compañía
// GET /products lista todos
```

---

### Capa de utilidades — `internal/utils/`

#### `response.go`

```go
type ErrorResponse struct {
    Error   string            `json:"error"`
    Details map[string]string `json:"details,omitempty"`
    Code    int               `json:"code"`
}

func WriteJSON(w http.ResponseWriter, data interface{}, status int)
func WriteError(w http.ResponseWriter, message string, status int)
func WriteValidationError(w http.ResponseWriter, errors map[string]string, status int)
```

#### `context.go`

```go
const UserContextKey contextKey = "user"

func SetUserContext(ctx context.Context, user *domain.User) context.Context
// Inyecta usuario en context

func GetUserFromContext(ctx context.Context) (*domain.User, bool)
// Recupera usuario del context
```

#### `validation.go`

```go
func ValidateRegisterFields(nameFull, phone, idNumber, dateOfBirth, email, password string) error
// Valida: name no vacío, password ≥8 chars, email regex
// (phone, idNumber, dateOfBirth aceptados pero no validados aún)
```

#### `TryExecut.go`

```go
func TryExecute(ctx context.Context, fn func() error) (err error)
```

Ejecuta `fn()` en una goroutine con:
- Recuperación de panic (`recover`)
- Respeto de `ctx.Done()` (cancelación de contexto)
- Retorna error si panic o si contexto cancelado

**Uso exclusivo:** `auth_handler.go:45` (`Register`).

---

### Capa de dominio — `internal/domain/`

19 structs sin dependencias externas. El dominio incluye ahora el **Archivo
Histórico de la Casa** — `InventoryHistory` — que registra cada decreto,
movimiento y transacción como un pergamino en la Biblioteca de Caladan.

**`InventoryHistory` — *El Pergamino Digital*:**

```go
type InventoryEventType string

const (
    EventTypeCREATE            InventoryEventType = "CREATE"
    EventTypeUPDATE            InventoryEventType = "UPDATE"
    EventTypeCANCEL            InventoryEventType = "CANCEL"
    EventTypeORDER_CREATED     InventoryEventType = "ORDER_CREATED"
    EventTypeORDER_UPDATED     InventoryEventType = "ORDER_UPDATED"
    EventTypeORDER_APPROVED    InventoryEventType = "ORDER_APPROVED"
    EventTypeSHIPMENT_CREATED  InventoryEventType = "SHIPMENT_CREATED"
    EventTypeSHIPMENT_CANCELLED InventoryEventType = "SHIPMENT_CANCELLED"
    EventTypeENTRY_CREATED     InventoryEventType = "ENTRY_CREATED"
    EventTypeENTRY_DELETED     InventoryEventType = "ENTRY_DELETED"
    EventTypeSTOCK_UPDATED     InventoryEventType = "STOCK_UPDATED"
    EventTypeINVOICE_LINKED    InventoryEventType = "INVOICE_LINKED"
    EventTypeRELATION_CREATED  InventoryEventType = "RELATION_CREATED"
)
```

Cada constante es un **Decreto Imperial** — un tipo de evento que la Casa
reconoce. El Archivero (`HistoryService`) los usa para clasificar cada entrada.

| Decreto (EventType) | Significado en el Imperio |
|---------------------|--------------------------|
| `CREATE` | Se fundó un nuevo feudo (recurso creado) |
| `UPDATE` | Se modificó un trato (recurso actualizado) |
| `CANCEL` | Se anuló una orden del Duque |
| `ORDER_CREATED` | Se emitió una orden de compra |
| `ORDER_UPDATED` | Se modificó una orden existente |
| `ORDER_APPROVED` | El Consejo aprobó la orden |
| `SHIPMENT_CREATED` | Un cargamento de especia partió de Arrakis |
| `SHIPMENT_CANCELLED` | El cargamento fue abortado |
| `ENTRY_CREATED` | Llegó especia a los almacenes (entrada de producto) |
| `ENTRY_DELETED` | Se eliminó un registro de entrada |
| `STOCK_UPDATED` | El Mentat ajustó el inventario |
| `INVOICE_LINKED` | Se vinculó una factura al feudo |
| `RELATION_CREATED` | Se firmó un tratado (relación orden↔despacho) |

```go
type InventoryHistory struct {
    HistoryID             string             `json:"history_id"`
    EventDate             time.Time          `json:"event_date"`
    UserID                string             `json:"user_id"`
    EventType             InventoryEventType `json:"event_type"`
    CompanyID             string             `json:"company_id"`
    DocumentID            string             `json:"document_id"`
    DocumentType          string             `json:"document_type"`
    ProviderDestinationID *string            `json:"provider_destination_id,omitempty"`
    PreviousData          *string            `json:"previous_data,omitempty"`
    NewData               *string            `json:"new_data,omitempty"`
    Description           string             `json:"description"`
    IPAddress             string             `json:"ip_address"`
    Result                string             `json:"result"`
    CreatedAt             time.Time          `json:"created_at"`
}
```

| Campo | Tipo | Propósito en la Casa |
|-------|------|---------------------|
| `HistoryID` | `string` | Identificador único del pergamino (generado por `generateID()`) |
| `EventDate` | `time.Time` | Momento exacto en que ocurrió el evento — el sello temporal del Archivero |
| `UserID` | `string` | El Noble que ejecutó la acción |
| `EventType` | `InventoryEventType` | Tipo de decreto imperial |
| `CompanyID` | `string` | Feudo al que pertenece el registro |
| `DocumentID` | `string` | ID del documento asociado (orden, entrada, despacho) |
| `DocumentType` | `string` | Tipo de documento (`"order"`, `"product-entry"`, `"shipment"`, `"product"`) |
| `ProviderDestinationID` | `*string` | Destino del cargamento (opcional, como una ruta en el mapa de Arrakis) |
| `PreviousData` | `*string` | JSON del estado anterior — la foto del feudo antes del cambio |
| `NewData` | `*string` | JSON del estado posterior — la foto del feudo después del cambio |
| `Description` | `string` | Narración del evento en lenguaje del Landsraad |
| `IPAddress` | `string` | Coordenadas del originante (dirección IP) |
| `Result` | `string` | Verdict final: `"SUCCESS"` o `"FAILURE"` |
| `CreatedAt` | `time.Time` | Cuándo se selló el pergamino en los Archivos |

Solo `User` tiene método:

```go
func (u *User) HasPermission(permission string) bool
// Retorna true si user tiene el rol "admin" o el permiso específico
```

Relaciones entre entidades:

```
User (1) ──→ Company (N)         vía UserID
Company (1) ──→ MainAddress (N)  vía CompanyID
Company (1) ──→ EconomicActivity (N)
Company (1) ──→ TaxInformation (N)
Company (1) ──→ Product (N)      vía CompanyID
Company (1) ──→ Winery (N)       vía CompanyID
Product (N) ──→ Provider (1)     vía SupplierID
Product (N) ──→ Winery (1)       vía WineryID
ProductEntry (N) ──→ Product (N) vía Details[].Product
ProductEntry (N) ──→ Provider (1) vía SupplierID
Movement (N) ──→ MovementType (1) vía MovementTypeID (FK)
Order (N) ──→ Client (1)         vía ClientID
InventoryHistory (N) ──→ User (1)  vía UserID
InventoryHistory (N) ──→ Company (1) vía CompanyID
```

### Capa de configuración — `internal/config/`

```go
type PostgresConfig struct {
    Host, User, Password, DBName string
    Port int
}

func DefaultPostgresConfig() PostgresConfig
// Valores por defecto: localhost:5432/postgres/1234567890/coffee_book

func (c PostgresConfig) DSN() string
// Formatea DSN para lib/pq: "host=... port=... user=... password=... dbname=... sslmode=disable"

func JWTSecret() string
// Retorna $JWT_SECRET o "book-coffee-shop-dev-secret"
```

### Capa de base de datos — `internal/database/`

```go
func EnsureDatabaseExists(cfg config.PostgresConfig) error
```

Conecta a `postgres` database, verifica si `cfg.DBName` existe, la crea si no. No migra esquemas — eso lo hace `runMigrations()` en `main.go`.

### Migraciones — `cmd/api/main.go:runMigrations()`

Sistema inline con 20 `CREATE TABLE IF NOT EXISTS` + 10 `ALTER TABLE ... ADD COLUMN IF NOT EXISTS`.

**Advertencia:** Las primeras líneas ejecutan `DROP TABLE IF EXISTS movements CASCADE` y `DROP TABLE IF EXISTS companies CASCADE` — **peligroso en producción** porque elimina datos existentes.

---

### WineryUseCase

```go
func NewWineryUseCase(repo repository.WineryRepository) WineryUseCase
```

Valida que `area` esté en `["Tienda","Almacén","Cafetería","Otro"]` y que `units` esté en `["Unidades","Cajas","Litros","Kilogramos"]`. Ambos conjuntos usan español para coincidir con los formularios frontend.

---

### Capa de handlers HTTP — winery

```go
type WineryHandler struct {
    uc usecase.WineryUseCase
}
```

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/wineries` | Lista todas las bodegas (filtro `?company_id=`) |
| POST | `/wineries` | Crea una bodega |
| GET | `/wineries/{id}` | Obtiene bodega por ID |
| PUT | `/wineries/{id}` | Actualiza bodega |
| DELETE | `/wineries/{id}` | Elimina bodega |

---

### Capa de handlers HTTP — inventory history (Las Crónicas)

> *"El pasado es un libro que siempre está abierto."*
> — Princesa Irulan, *Dune*

```go
type InventoryHistoryHandler struct {
    svc *usecase.HistoryService
}

func NewInventoryHistoryHandler(svc *usecase.HistoryService) *InventoryHistoryHandler
```

El handler de historial es el **Bibliotecario de Caladan**. No crea, no actualiza,
no elimina — solo **consulta los Archivos**. Sus métodos son de solo lectura,
como los anales de Irulan que narran lo ocurrido sin poder cambiarlo.

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/history` | Abre el Gran Libro de la Casa (todos los eventos) |
| GET | `/history/{document_type}/{document_id}` | Busca pergaminos por feudo y documento específico |

El handler parsea la ruta extrayendo `document_type` y `document_id` del path,
y delega al `HistoryService`:

- `getAll()` → `svc.GetAll()` — retorna todos los eventos, ordenados por fecha descendente.
- `getByDocument(documentType, documentID)` → `svc.GetByDocument()` — retorna solo los eventos de un documento concreto.

**Respuesta típica:**
```json
[
  {
    "history_id": "1749372100_042",
    "event_date": "2026-06-19T10:30:00Z",
    "user_id": "1749372100_001",
    "event_type": "ORDER_CREATED",
    "company_id": "1749372100_015",
    "document_id": "1749372100_020",
    "document_type": "order",
    "provider_destination_id": null,
    "previous_data": null,
    "new_data": "{\"order_numeric\":\"OC-001\",\"status\":\"received\",\"details\":[...]}",
    "description": "Orden de compra OC-001 creada por Juan Pérez",
    "ip_address": "192.168.1.10",
    "result": "SUCCESS",
    "created_at": "2026-06-19T10:30:00Z"
  }
]
```

---

## Preguntas Frecuentes (FAQ)

### ¿Cómo agrego un nuevo módulo CRUD?

1. Crear struct en `internal/domain/`
2. Crear interface en `internal/repository/` con métodos CRUD
3. Crear implementación en `internal/infrastructure/` (PostgreSQL)
4. Crear caso de uso en `internal/usecase/` con validaciones
5. Crear handler en `internal/handler/` siguiendo el patrón `Handle()` + métodos privados
6. Agregar migración en `main.go:runMigrations()`
7. Wirear en `main.go:` `NewRepo → NewUC → NewHandler → mux.HandleFunc`

### ¿Cómo desactivo la autenticación para desarrollo?

En `main.go:157`, modifica `NewAuthMiddleware` para incluir todas las rutas como públicas o simplemente no aplicar el middleware:

```go
// handler := middleware.RecoveryMiddleware(c.Handler(authMiddleware(mux)))
handler := middleware.RecoveryMiddleware(c.Handler(mux)) // sin auth
```

### ¿Por qué hay dos packages `handler/` y `handlers/`?

`internal/handler/` (18 archivos) es el activo, conectado en `main.go`. `internal/handlers/` (1 archivo) es legacy, usa `middleware.RequireRoles` + `models` DTOs con `ValidatePayload`, pero **no está conectado al router**. Es deuda técnica.

### El token JWT dice "invalid or expired token" constantemente

El middleware valida:
1. JWT (firma + expiración) → `tokenService.Validate()`
2. Token existe en BD → `userRepo.GetByAuthToken()`

Si el token fue regenerado en otro login, el anterior queda invalidado en BD. Vuelve a hacer login para obtener un token fresco.

### ¿Cómo cambio el tiempo de expiración del JWT?

En `jwt_token_service.go:23`:

```go
"exp": time.Now().Add(24 * time.Hour).Unix(), // Cambia 24h a lo que necesites
```

### La migración falla porque una tabla ya existe con columnas faltantes

El sistema maneja esto con `ALTER TABLE ... ADD COLUMN IF NOT EXISTS` en el segundo bloque de migraciones (líneas 383-396). Si una columna específica no se agrega, puede ser que ya exista con otro tipo. Revisa los logs de PostgreSQL.

### ¿Cómo ejecuto consultas personalizadas?

Usa `db.Query()` o `db.QueryRow()` directamente en los repositorios de infraestructura. No hay capa ORM.

### El campo `roles` del usuario siempre está vacío

La tabla `users` en la migración no tiene columna `roles`, y `PostgresUserRepository.scanUser()` no la escanea. Para usar `RequireRoles`, necesitas agregar la columna y modificar el repositorio.

### Tipos de movimiento válidos para entradas de producto

```
"Purchase", "Return", "Donation", "Inventory Adjustment", "Internal Production"
```

### Unidades de medida válidas para productos

```
"Kg", "Litro", "Libra", "Gramos", "Unidad"
```

### Áreas y unidades válidas para bodegas

```
Áreas:   "Tienda", "Almacén", "Cafetería", "Otro"
Unidades: "Unidades", "Cajas", "Litros", "Kilogramos"
```

Nota: Los valores están en español para coincidir con los formularios frontend. El `wineryUseCase` rechazará valores en inglés.

### El campo `winery_id` en productos es obligatorio?

No, es opcional. Si no se envía en la petición (o se envía como `""`), el producto se crea sin bodega asociada. El frontend permite al usuario decidir si vincular el producto a una bodega durante el registro o hacerlo manualmente después.

### Métodos de pago válidos para órdenes

```
"cash", "transfer", "debit-card", "credit-card"
```

### Estados válidos para órdenes

```
"received", "in-preparation", "ready-for-delivery", "delivered", "cancelled"
```

### ¿El historial de inventario soporta escritura desde los endpoints?

No. El historial es solo de lectura vía API (GET). La escritura ocurre
**automáticamente** desde los casos de uso de producto, orden, entrada y despacho
cuando ejecutan acciones. El `HistoryService` se inyecta como dependencia y cada
operación relevante registra su evento. Esto es como los **Archivos Imperiales**:
solo los Escribas autorizados (los casos de uso) pueden añadir pergaminos; el
resto del Imperio solo puede consultarlos.

### ¿Qué módulos registran eventos en el historial?

Actualmente 5 módulos alimentan las Crónicas:

| Módulo (Casa) | Eventos que registra |
|---------------|---------------------|
| `productUseCase` | `CREATE`, `UPDATE`, `STOCK_UPDATED` |
| `orderUseCase` | `ORDER_CREATED`, `ORDER_UPDATED`, `ORDER_APPROVED` |
| `productEntryUseCase` | `ENTRY_CREATED`, `ENTRY_DELETED` |
| `shipmentUseCase` | `SHIPMENT_CREATED`, `SHIPMENT_CANCELLED`, `RELATION_CREATED` |
| `movementUseCase` | `CREATE`, `UPDATE`, `CANCEL` |

Cada uno inyecta `HistoryService` en su constructor y llama a `LogEvent()`
(o `LogRelation()`, `LogStockUpdate()`) dentro de la misma transacción de base
de datos. Si el evento falla, toda la operación se revierte — la especia no
se mueve si el Escriba no puede escribir.

---

## Sobre la Casa Atreides

> *"Un hombre debe ver antes de poder actuar."*
> — Leto Atreides, *Dune*

Esta API no es solo un conjunto de endpoints. Es el **gobierno digital** de tu
negocio, construido con el honor y la estrategia de la más noble Casa del
Landsraad.

### Principios de la Casa

1. **Honor en los datos** — Cada transacción es ACID, cada respuesta es JSON,
   cada error tiene un código. No mentimos al Imperio.
2. **Estrategia en la arquitectura** — Clean Architecture no es un lujo, es una
   necesidad. Como el Mentat calcula rutas de especia, cada capa tiene un
   propósito definido.
3. **Adaptación constante** — Así como los Atreides pasaron de Caladan a Arrakis,
   esta API puede gobernar cualquier dominio de negocio. El mismo núcleo fiscal,
   inventario y autenticación funciona para una cafetería, una bodega o una
   cadena de suministro completa.
4. **Legado mantenible** — El código que escribes hoy será la base del imperio
   de mañana. Sin deuda técnica evitable, sin magia, sin atajos.

> *"El poder de gobernar está en los pequeños detalles."*
> — Thufir Hawat, Mentat de la Casa Atreides

---

*Atreides API v2.0 — Documentación de la Casa*
*Inspirado en la saga *Dune* de Frank Herbert*
