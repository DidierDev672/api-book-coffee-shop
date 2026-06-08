package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"book-coffee-shop/internal/config"
	"book-coffee-shop/internal/handler"
	"book-coffee-shop/internal/infrastructure"
	"book-coffee-shop/internal/usecase"

	_ "github.com/lib/pq"
)

func main() {
	pgCfg := config.DefaultPostgresConfig()

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

	orderRepo := infrastructure.NewPostgresOrderRepository(db)
	orderUC := usecase.NewOrderUseCase(orderRepo)
	orderH := handler.NewOrderHandler(orderUC)

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
	mux.HandleFunc("/orders", orderH.Handle)
	mux.HandleFunc("/orders/", orderH.Handle)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	fmt.Printf("Book Coffee Shop API running on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func runMigrations(db *sql.DB) error {
	db.Exec(`DROP TABLE IF EXISTS movements CASCADE`)

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
			product_code  VARCHAR(255) NOT NULL,
			categories    TEXT[] NOT NULL DEFAULT '{}',
			unit          VARCHAR(50) NOT NULL,
			minimum_stock DOUBLE PRECISION NOT NULL DEFAULT 0,
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
		`CREATE TABLE IF NOT EXISTS clients (
			id        VARCHAR(50) PRIMARY KEY,
			name_full VARCHAR(255) NOT NULL,
			phone     VARCHAR(50) NOT NULL,
			correo    VARCHAR(255),
			address   TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id             VARCHAR(50) PRIMARY KEY,
			order_numeric  VARCHAR(100) NOT NULL,
			date           DATE NOT NULL,
			hour           VARCHAR(20) NOT NULL,
			attended_by    VARCHAR(255) NOT NULL,
			client_id      VARCHAR(50) NOT NULL,
			details        JSONB NOT NULL DEFAULT '[]',
			payment_method VARCHAR(50) NOT NULL,
			status         VARCHAR(50) NOT NULL,
			observations   TEXT,
			created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
	}
	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return err
		}
	}
	return nil
}
