package repository

import "database/sql"

type DBTX interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type OrderRepoFactory func(tx DBTX) OrderRepository
type ShipmentRepoFactory func(tx DBTX) ShipmentRepository
type ProductEntryRepoFactory func(tx DBTX) ProductEntryRepository
type ProductRepoFactory func(tx DBTX) ProductRepository
type MovementRepoFactory func(tx DBTX) MovementRepository
type HistoryRepoFactory func(tx DBTX) InventoryHistoryRepository
type SaleRepoFactory func(tx DBTX) SaleRepository
