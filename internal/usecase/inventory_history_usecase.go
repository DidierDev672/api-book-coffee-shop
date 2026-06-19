package usecase

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type HistoryService struct {
	db       *sql.DB
	makeRepo repository.HistoryRepoFactory
}

func NewHistoryService(db *sql.DB, makeRepo repository.HistoryRepoFactory) *HistoryService {
	return &HistoryService{db: db, makeRepo: makeRepo}
}

func (s *HistoryService) LogEvent(tx repository.DBTX, eventType domain.InventoryEventType, userID, companyID, documentID, documentType, description, ipAddress string, previousData, newData interface{}) error {
	now := time.Now()

	var prevStr, newStr *string
	if previousData != nil {
		b, err := json.Marshal(previousData)
		if err != nil {
			return err
		}
		p := string(b)
		prevStr = &p
	}
	if newData != nil {
		b, err := json.Marshal(newData)
		if err != nil {
			return err
		}
		n := string(b)
		newStr = &n
	}

	event := &domain.InventoryHistory{
		HistoryID:    generateID(),
		EventDate:    now,
		UserID:       userID,
		EventType:    eventType,
		CompanyID:    companyID,
		DocumentID:   documentID,
		DocumentType: documentType,
		PreviousData: prevStr,
		NewData:      newStr,
		Description:  description,
		IPAddress:    ipAddress,
		Result:       "SUCCESS",
		CreatedAt:    now,
	}

	repo := s.makeRepo(tx)
	return repo.Create(event)
}

func (s *HistoryService) LogRelation(tx repository.DBTX, orderID, shipmentID, userID, companyID, ipAddress string) error {
	now := time.Now()

	newData := map[string]string{
		"order_id":    orderID,
		"shipment_id": shipmentID,
	}
	b, _ := json.Marshal(newData)
	raw := string(b)

	repo := s.makeRepo(tx)

	event := &domain.InventoryHistory{
		HistoryID:    generateID(),
		EventDate:    now,
		UserID:       userID,
		EventType:    domain.EventTypeRELATION_CREATED,
		CompanyID:    companyID,
		DocumentID:   orderID,
		DocumentType: "order",
		NewData:      &raw,
		Description:  "Relation created between order " + orderID + " and shipment " + shipmentID,
		IPAddress:    ipAddress,
		Result:       "SUCCESS",
		CreatedAt:    now,
	}

	if err := repo.Create(event); err != nil {
		return err
	}

	event.HistoryID = generateID()
	event.DocumentID = shipmentID
	event.DocumentType = "shipment"
	return repo.Create(event)
}

func (s *HistoryService) LogStockUpdate(tx repository.DBTX, productCode string, previousStock, newStock float64, userID, companyID, ipAddress string) error {
	prevData := map[string]float64{"stock": previousStock}
	newData := map[string]float64{"stock": newStock}

	return s.LogEvent(tx,
		domain.EventTypeSTOCK_UPDATED,
		userID, companyID, productCode, "product",
		fmt.Sprintf("Stock updated from %.2f to %.2f", previousStock, newStock),
		ipAddress, prevData, newData,
	)
}

func (s *HistoryService) GetByDocument(documentType, documentID string) ([]*domain.InventoryHistory, error) {
	repo := s.makeRepo(s.db)
	return repo.GetByDocument(documentType, documentID)
}

func (s *HistoryService) GetAll() ([]*domain.InventoryHistory, error) {
	repo := s.makeRepo(s.db)
	return repo.GetAll()
}
