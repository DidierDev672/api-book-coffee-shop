package main

import (
	"book-coffee-shop/internal/config"
	"book-coffee-shop/internal/database"
	"book-coffee-shop/internal/handler"
	"book-coffee-shop/internal/infrastructure"
	"book-coffee-shop/internal/middleware"
	"book-coffee-shop/internal/usecase"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	pgCfg := config.DefaultPostgresConfig()

	if err := database.EnsureDatabaseExists(pgCfg); err != nil {
		log.Fatalf("failed to ensure database exists: %v", err)
	}

	db, err := sql.Open("postgres", pgCfg.DSN())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	fmt.Println("Connected to PostgreSQL")

	if err := runMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	fmt.Println("Migrations applied")

	authorRepo := infrastructure.NewPostgresAuthorRepository(db)
	authorUC := usecase.NewAuthorUseCase(authorRepo)
	authorH := handler.NewAuthorHandler(authorUC)

	bookRepo := infrastructure.NewPostgresBookRepository(db)
	bookUC := usecase.NewBookUseCase(bookRepo)
	bookH := handler.NewBookHandler(bookUC)

	topicRepo := infrastructure.NewPostgresTopicRepository(db)
	topicUC := usecase.NewTopicUseCase(topicRepo)
	topicH := handler.NewTopicHandler(topicUC)

	noteRepo := infrastructure.NewPostgresNoteRepository(db)
	noteUC := usecase.NewNoteUseCase(noteRepo)
	noteH := handler.NewNoteHandler(noteUC)

	estRepo := infrastructure.NewPostgresEstablishmentRepository(db)
	estUC := usecase.NewEstablishmentUseCase(estRepo)
	estH := handler.NewEstablishmentHandler(estUC)

	movTypeRepo := infrastructure.NewPostgresMovementTypeRepository(db)
	movTypeUC := usecase.NewMovementTypeUseCase(movTypeRepo)
	movTypeH := handler.NewMovementTypeHandler(movTypeUC)

	movRepo := infrastructure.NewPostgresMovementRepository(db)
	movUC := usecase.NewMovementUseCase(movRepo, movTypeRepo)
	movH := handler.NewMovementHandler(movUC)

	prodRepo := infrastructure.NewPostgresProductRepository(db)
	prodUC := usecase.NewProductUseCase(prodRepo)
	prodH := handler.NewProductHandler(prodUC)

	msRepo := infrastructure.NewPostgresMonthlySummaryRepository(db)
	msUC := usecase.NewMonthlySummaryUseCase(msRepo)
	msH := handler.NewMonthlySummaryHandler(msUC)

	clientRepo := infrastructure.NewPostgresClientRepository(db)
	clientUC := usecase.NewClientUseCase(clientRepo)
	clientH := handler.NewClientHandler(clientUC)

	companyRepo := infrastructure.NewPostgresCompanyRepository(db)
	companyUC := usecase.NewCompanyUseCase(companyRepo)
	companyH := handler.NewCompanyHandler(companyUC)

	mainAddressRepo := infrastructure.NewPostgresMainAddressRepository(db)
	mainAddressUC := usecase.NewMainAddressUseCase(mainAddressRepo)
	mainAddressH := handler.NewMainAddressHandler(mainAddressUC)

	taxInformationRepo := infrastructure.NewPostgresTaxInformationRepository(db)
	taxInformationUC := usecase.NewTaxInformationUseCase(taxInformationRepo)
	taxInformationH := handler.NewTaxInformationHandler(taxInformationUC)

	economicActivityRepo := infrastructure.NewPostgresEconomicActivityRepository(db)
	economicActivityUC := usecase.NewEconomicActivityUseCase(economicActivityRepo)
	economicActivityH := handler.NewEconomicActivityHandler(economicActivityUC)

	orderRepo := infrastructure.NewPostgresOrderRepository(db)
	orderUC := usecase.NewOrderUseCase(orderRepo)
	orderH := handler.NewOrderHandler(orderUC)

	providerRepo := infrastructure.NewPostgresProviderRepository(db)
	providerUC := usecase.NewProviderUseCase(providerRepo)
	providerH := handler.NewProviderHandler(providerUC)

	productEntryRepo := infrastructure.NewPostgresProductEntryRepository(db)
	productEntryUC := usecase.NewProductEntryUseCase(productEntryRepo)
	productEntryH := handler.NewProductEntryHandler(productEntryUC)

	userRepo := infrastructure.NewPostgresUserRepository(db)
	tokenService := infrastructure.NewJWTTokenService(config.JWTSecret())
	passwordHasher := infrastructure.NewBcryptPasswordHasher()
	authUC := usecase.NewAuthUseCase(userRepo, passwordHasher, tokenService)
	authH := handler.NewAuthHandler(authUC)

	wineryRepo := infrastructure.NewPostgresWineryRepository(db)
	wineryUC := usecase.NewWineryUseCase(wineryRepo)
	wineryH := handler.NewWineryHandler(wineryUC)

	shipmentRepo := infrastructure.NewPostgresShipmentRepository(db)
	shipmentUC := usecase.NewShipmentUseCase(shipmentRepo)
	shipmentH := handler.NewShipmentHandler(shipmentUC)

	mux := http.NewServeMux()
	mux.HandleFunc("/authors", authorH.Handle)
	mux.HandleFunc("/authors/", authorH.Handle)
	mux.HandleFunc("/books", bookH.Handle)
	mux.HandleFunc("/books/", bookH.Handle)
	mux.HandleFunc("/topics", topicH.Handle)
	mux.HandleFunc("/topics/", topicH.Handle)
	mux.HandleFunc("/notes", noteH.Handle)
	mux.HandleFunc("/notes/", noteH.Handle)
	mux.HandleFunc("/establishments", estH.Handle)
	mux.HandleFunc("/establishments/", estH.Handle)
	mux.HandleFunc("/movement-types", movTypeH.Handle)
	mux.HandleFunc("/movement-types/", movTypeH.Handle)
	mux.HandleFunc("/movements", movH.Handle)
	mux.HandleFunc("/movements/", movH.Handle)
	mux.HandleFunc("/products", prodH.Handle)
	mux.HandleFunc("/products/", prodH.Handle)
	mux.HandleFunc("/monthly-summaries", msH.Handle)
	mux.HandleFunc("/monthly-summaries/", msH.Handle)
	mux.HandleFunc("/clients", clientH.Handle)
	mux.HandleFunc("/clients/", clientH.Handle)
	mux.HandleFunc("/companies", companyH.Handle)
	mux.HandleFunc("/companies/", companyH.Handle)
	mux.HandleFunc("/companies/user/", companyH.Handle)
	mux.HandleFunc("/main-addresses", mainAddressH.Handle)
	mux.HandleFunc("/main-addresses/", mainAddressH.Handle)
	mux.HandleFunc("/tax-information", taxInformationH.Handle)
	mux.HandleFunc("/tax-information/", taxInformationH.Handle)
	mux.HandleFunc("/economic-activities", economicActivityH.Handle)
	mux.HandleFunc("/economic-activities/", economicActivityH.Handle)
	mux.HandleFunc("/providers", providerH.Handle)
	mux.HandleFunc("/providers/", providerH.Handle)
	mux.HandleFunc("/product-entries", productEntryH.Handle)
	mux.HandleFunc("/product-entries/", productEntryH.Handle)
	mux.HandleFunc("/orders", orderH.Handle)
	mux.HandleFunc("/orders/", orderH.Handle)
	mux.HandleFunc("/wineries", wineryH.Handle)
	mux.HandleFunc("/wineries/", wineryH.Handle)
	mux.HandleFunc("/shipments", shipmentH.Handle)
	mux.HandleFunc("/shipments/", shipmentH.Handle)
	mux.HandleFunc("/auth/register", authH.Register)
	mux.HandleFunc("/auth/login", authH.Login)
	mux.HandleFunc("/users", authH.ListUsers)
	mux.HandleFunc("/users/", authH.HandleUser)

	authMiddleware := middleware.NewAuthMiddleware(tokenService, userRepo, "/auth/register", "/auth/login")

	//! Configuracion de CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},                   // Dominios permitidos)
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, // Métodos permitidos
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug:          true, // Muestra logs en consola para ayudarte a depurar
	})

	handler := middleware.RecoveryMiddleware(c.Handler(authMiddleware(mux)))
	log.Println("Server listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS authors (
			id          VARCHAR(50) PRIMARY KEY,
			name        VARCHAR(255) NOT NULL,
			country     VARCHAR(255) NOT NULL,
			genres      TEXT[] NOT NULL DEFAULT '{}',
			birth_day   VARCHAR(20) NOT NULL,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS books (
			id               VARCHAR(50) PRIMARY KEY,
			title            VARCHAR(255) NOT NULL,
			description      TEXT NOT NULL,
			author           VARCHAR(255) NOT NULL,
			genres           TEXT[] NOT NULL DEFAULT '{}',
			photos           TEXT[] NOT NULL DEFAULT '{}',
			publication_date DATE,
			created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS topics (
			id          VARCHAR(50) PRIMARY KEY,
			name        VARCHAR(255) NOT NULL,
			type        VARCHAR(100) NOT NULL,
			description TEXT NOT NULL,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS notes (
			id         VARCHAR(50) PRIMARY KEY,
			name       VARCHAR(255) NOT NULL,
			content    TEXT NOT NULL,
			type       VARCHAR(100) NOT NULL,
			color      VARCHAR(50) NOT NULL,
			id_topic   VARCHAR(50) NOT NULL,
			id_book    VARCHAR(50),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS establishments (
			id                       VARCHAR(50) PRIMARY KEY,
			establishment_name       VARCHAR(255) NOT NULL,
			inventory_manager        VARCHAR(255) NOT NULL,
			warehouse_point_of_sale  VARCHAR(255) NOT NULL,
			created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS movement_types (
			id          VARCHAR(50) PRIMARY KEY,
			name        VARCHAR(255) NOT NULL,
			description TEXT,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id            VARCHAR(50) PRIMARY KEY,
			company_id    VARCHAR(50) NOT NULL DEFAULT '',
			supplier_id   VARCHAR(50) NOT NULL DEFAULT '',
			name          VARCHAR(255) NOT NULL DEFAULT '',
			product_code  VARCHAR(255) NOT NULL,
			categories    TEXT[] NOT NULL DEFAULT '{}',
			unit          VARCHAR(50) NOT NULL,
			quantity      DOUBLE PRECISION NOT NULL DEFAULT 0,
			minimum_stock DOUBLE PRECISION NOT NULL DEFAULT 0,
			winery_id     VARCHAR(50) NOT NULL DEFAULT '',
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS monthly_summaries (
			id               VARCHAR(50) PRIMARY KEY,
			product          VARCHAR(255) NOT NULL,
			beginning_stock  DOUBLE PRECISION NOT NULL DEFAULT 0,
			incoming_orders  DOUBLE PRECISION NOT NULL DEFAULT 0,
			outgoing_orders  DOUBLE PRECISION NOT NULL DEFAULT 0,
			ending_stock     DOUBLE PRECISION NOT NULL DEFAULT 0,
			created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id            VARCHAR(50) PRIMARY KEY,
			name_full     VARCHAR(255) NOT NULL,
			phone         VARCHAR(50) NOT NULL,
			id_number     VARCHAR(50) NOT NULL,
			date_of_birth DATE NOT NULL,
			email         VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			auth_token    TEXT,
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS clients (
			id        VARCHAR(50) PRIMARY KEY,
			name_full VARCHAR(255) NOT NULL,
			phone     VARCHAR(50) NOT NULL,
			correo    VARCHAR(255),
			address   TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS companies (
			id                 VARCHAR(50) PRIMARY KEY,
			user_id            VARCHAR(50) NOT NULL DEFAULT '' REFERENCES users(id),
			nit                VARCHAR(50) NOT NULL UNIQUE,
			social_reason      VARCHAR(255) NOT NULL,
			business_name      VARCHAR(255) NOT NULL,
			type_person        VARCHAR(100) NOT NULL,
			company_type       VARCHAR(100) NOT NULL,
			status             VARCHAR(50) NOT NULL,
			constitution_date  DATE NOT NULL,
			email              VARCHAR(255) NOT NULL DEFAULT '',
			phone              VARCHAR(50) NOT NULL DEFAULT '',
			cellphone          VARCHAR(50) NOT NULL DEFAULT '',
			created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS main_addresses (
			id          VARCHAR(50) PRIMARY KEY,
			user_id     VARCHAR(50) NOT NULL,
			company_id  VARCHAR(50) NOT NULL,
			country     VARCHAR(255) NOT NULL,
			department  VARCHAR(255) NOT NULL,
			address     TEXT NOT NULL,
			postcode    VARCHAR(50) NOT NULL,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS economic_activities (
			id          VARCHAR(50) PRIMARY KEY,
			user_id     VARCHAR(50) NOT NULL,
			company_id  VARCHAR(50) NOT NULL,
			code        VARCHAR(100) NOT NULL,
			description TEXT NOT NULL,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS tax_information (
			id                    VARCHAR(50) PRIMARY KEY,
			user_id               VARCHAR(50) NOT NULL,
			business_id           VARCHAR(50) NOT NULL,
			tax_regime            VARCHAR(100) NOT NULL,
			vat_responsible       BOOLEAN NOT NULL DEFAULT FALSE,
			withholding_taxpayer  BOOLEAN NOT NULL DEFAULT FALSE,
			large_taxpayer        BOOLEAN NOT NULL DEFAULT FALSE,
			created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS providers (
    id                  VARCHAR(50) PRIMARY KEY,
    code                VARCHAR(100) NOT NULL UNIQUE,
    type_person         VARCHAR(50) NOT NULL,
    document_type       VARCHAR(20) NOT NULL,
    document_number     VARCHAR(50) NOT NULL,
    verification_digit  VARCHAR(10) NOT NULL DEFAULT '',
    business_name       VARCHAR(255) NOT NULL,
    business_activity   TEXT NOT NULL DEFAULT '',
    status              BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
)`,
		`CREATE TABLE IF NOT EXISTS product_entries (
			id                 VARCHAR(50) PRIMARY KEY,
			entry_number       VARCHAR(100) NOT NULL,
			registered_date    DATE NOT NULL,
			movement_type      VARCHAR(50) NOT NULL,
			warehouse          VARCHAR(255) DEFAULT '',
			responsible_party  VARCHAR(50) NOT NULL,
			company_id         VARCHAR(50) NOT NULL DEFAULT '',
			details            JSONB NOT NULL DEFAULT '[]',
			financial_summary  JSONB NOT NULL DEFAULT '{}',
			observations       TEXT,
			created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id               VARCHAR(50) PRIMARY KEY,
			order_numeric    VARCHAR(100) NOT NULL,
			order_type       VARCHAR(50) NOT NULL DEFAULT 'PURCHASE',
			date             DATE NOT NULL,
			company_id       VARCHAR(50) NOT NULL DEFAULT '',
			user_id          VARCHAR(50) NOT NULL DEFAULT '',
			requested_by     VARCHAR(255) NOT NULL DEFAULT '',
			details          JSONB NOT NULL DEFAULT '[]',
			financial_summary JSONB NOT NULL DEFAULT '{}',
			status           VARCHAR(50) NOT NULL,
			reason_for_order TEXT,
			created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS movements (
			id                VARCHAR(50) PRIMARY KEY,
			date              DATE NOT NULL,
			code              VARCHAR(100) NOT NULL,
			product           VARCHAR(255) NOT NULL,
			unit              VARCHAR(50) NOT NULL,
			entrance          DOUBLE PRECISION NOT NULL DEFAULT 0,
			output            DOUBLE PRECISION NOT NULL DEFAULT 0,
			balance           DOUBLE PRECISION NOT NULL DEFAULT 0,
			unit_cost         DOUBLE PRECISION NOT NULL DEFAULT 0,
			valor_value       DOUBLE PRECISION NOT NULL DEFAULT 0,
			movement_type_id  VARCHAR(50) NOT NULL REFERENCES movement_types(id),
			observations      TEXT,
			created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS wineries (
			id               VARCHAR(50) PRIMARY KEY,
			registered_date  DATE NOT NULL,
			user_id          VARCHAR(50) NOT NULL DEFAULT '',
			company_id       VARCHAR(50) NOT NULL DEFAULT '',
			area             VARCHAR(50) NOT NULL,
			units            VARCHAR(50) NOT NULL,
			created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS shipments (
			id                VARCHAR(50) PRIMARY KEY,
			shipment_number   VARCHAR(100) NOT NULL,
			record_date       DATE NOT NULL,
			movement_type     VARCHAR(50) NOT NULL,
			status            VARCHAR(50) NOT NULL DEFAULT 'DRAFT',
			company_id        VARCHAR(50) NOT NULL DEFAULT '',
			warehouse_id      VARCHAR(50) NOT NULL DEFAULT '',
			responsible_id    VARCHAR(50) NOT NULL DEFAULT '',
			source_document   JSONB NOT NULL DEFAULT '{}',
			recipient         JSONB NOT NULL DEFAULT '{}',
			details           JSONB NOT NULL DEFAULT '[]',
			financial_summary JSONB NOT NULL DEFAULT '{}',
			remarks           TEXT,
			created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
	}
	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return err
		}
	}

	alterations := []string{
		`ALTER TABLE providers ADD COLUMN IF NOT EXISTS business_activity TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE main_addresses ADD COLUMN IF NOT EXISTS municipio VARCHAR(255) NOT NULL DEFAULT ''`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS company_id VARCHAR(50) NOT NULL DEFAULT ''`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS supplier_id VARCHAR(50) NOT NULL DEFAULT ''`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS name VARCHAR(255) NOT NULL DEFAULT ''`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS product_code VARCHAR(255) NOT NULL DEFAULT ''`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS categories TEXT[] NOT NULL DEFAULT '{}'`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS unit VARCHAR(50) NOT NULL DEFAULT ''`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS minimum_stock DOUBLE PRECISION NOT NULL DEFAULT 0`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS quantity DOUBLE PRECISION NOT NULL DEFAULT 0`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS winery_id VARCHAR(50) NOT NULL DEFAULT ''`,
		`ALTER TABLE orders ADD COLUMN IF NOT EXISTS order_type VARCHAR(50) NOT NULL DEFAULT 'PURCHASE'`,
		`ALTER TABLE orders ADD COLUMN IF NOT EXISTS company_id VARCHAR(50) NOT NULL DEFAULT ''`,
		`ALTER TABLE orders ADD COLUMN IF NOT EXISTS user_id VARCHAR(50) NOT NULL DEFAULT ''`,
		`ALTER TABLE orders ADD COLUMN IF NOT EXISTS requested_by VARCHAR(255) NOT NULL DEFAULT ''`,
		`ALTER TABLE orders ADD COLUMN IF NOT EXISTS reason_for_order TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE orders ADD COLUMN IF NOT EXISTS financial_summary JSONB NOT NULL DEFAULT '{}'`,
		`ALTER TABLE orders DROP COLUMN IF EXISTS hour`,
		`ALTER TABLE orders DROP COLUMN IF EXISTS attended_by`,
		`ALTER TABLE orders DROP COLUMN IF EXISTS client_id`,
		`ALTER TABLE orders DROP COLUMN IF EXISTS payment_method`,
		`ALTER TABLE orders DROP COLUMN IF EXISTS observations`,
		`ALTER TABLE orders DROP COLUMN IF EXISTS warehouse_id`,
		`ALTER TABLE orders DROP COLUMN IF EXISTS supplier_id`,
	}
	for _, a := range alterations {
		db.Exec(a)
	}
	return nil
}
