# Atreides — API REST

> *"Un comienzo es un tiempo muy delicado."* — Princesa Irulan, *Dune*

API REST escrita en Go siguiendo **Clean Architecture** con capas de dominio, casos de uso, repositorios e infraestructura. El nombre **Atreides** evoca la Casa Atreides de *Dune*: un linaje que crece, se adapta y se expande a nuevos territorios. Así como la Casa Atreides gobernó Caladan y luego Arrakis, esta API está diseñada para gobernar múltiples dominios de negocio — desde una librería-cafetería hasta proyectos más grandes que compartan el mismo núcleo fiscal, inventario y autenticación.

---

## Descripción general

Sistema backend monolitico que expone **29 endpoints HTTP** para administrar:

- **Gestión empresarial colombiana**: Compañías (NIT), direcciones, actividades económicas, información tributaria (regímenes, responsables de IVA, autorretenedores, grandes contribuyentes).
- **Inventario y producto**: Productos con categorías, unidades de medida, stock mínimo, cantidad, proveedores, bodegas, entradas de producto con detalle JSONB y resumen financiero, movimientos de inventario, resúmenes mensuales.
- **Bodegas**: Bodegas con área, unidades, fecha de registro, vinculadas a compañía.
- **Órdenes y ventas**: Órdenes con detalles JSONB, métodos de pago (`cash`, `transfer`, `debit-card`, `credit-card`), estados (`received`, `in-preparation`, `ready-for-delivery`, `delivered`, `cancelled`).
- **Autenticación**: Registro y login con JWT HS256 + bcrypt, middleware de autenticación para rutas protegidas.
- **Módulo editorial (legacy)**: Autores, libros, tópicos y notas.
- **Clientes**: CRUD de clientes.

### Stack técnico

| Capa | Tecnología |
|------|-----------|
| Lenguaje | **Go 1.26.3** |
| Base de datos | **PostgreSQL 15+** con JSONB, arrays (`TEXT[]`), `uuid-ossp` |
| Autenticación | **JWT HS256** (`github.com/golang-jwt/jwt/v5`) + **bcrypt** (`golang.org/x/crypto`) |
| Driver BD | `github.com/lib/pq` v1.12.3 |
| Validación | `github.com/go-playground/validator/v10` (solo en `models/`) |
| CORS | `github.com/rs/cors` v1.11.1 |
| Arquitectura | Clean Architecture (Domain → Repository → Usecase → Handler → Infrastructure) |

### Estructura del proyecto

```
cmd/api/
  main.go                  → Punto de entrada, wiring DI, migraciones inline
internal/
  config/config.go         → Config Postgres + JWT (con fallback hardcodeado)
  database/postgres.go     → EnsureDatabaseExists (crea DB si no existe)
  domain/                  → 19 structs de dominio (sin dependencias externas)
  repository/              → 21 interfaces + 2 servicios (TokenService, PasswordHasher)
  usecase/                 → 19 implementaciones de lógica de negocio
  handler/                 → 19 handlers HTTP activos
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
// Usado por los 18 usecases
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

18 structs sin dependencias externas. Solo `User` tiene método:

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

Sistema inline con 19 `CREATE TABLE IF NOT EXISTS` + 10 `ALTER TABLE ... ADD COLUMN IF NOT EXISTS`.

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
