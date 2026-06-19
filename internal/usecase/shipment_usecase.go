package usecase

import (
	"database/sql"
	"errors"
	"slices"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

var validShipmentMovementTypes = []string{"SALE", "SUPPLIER_RETURN", "DONATION", "SHRINKAGE", "ADJUSTMENT", "TRANSFER"}
var validShipmentStatuses = []string{"DRAFT", "CONFIRMED", "CANCELED"}
var validRecipientTypes = []string{"CUSTOMER", "SUPPLIER", "WAREHOUSE", "INTERNAL"}

type ShipmentUseCase interface {
	Create(shipmentNumber, recordDate, movementType, status, companyID, warehouseID, responsibleID string, sourceDocument domain.SourceDocument, recipient domain.Recipient, details []domain.ShipmentDetail, financialSummary domain.ShipmentFinancialSummary, remarks, ipAddress string) (*domain.Shipment, error)
	GetByID(id string) (*domain.Shipment, error)
	GetAll() ([]*domain.Shipment, error)
	Update(id, shipmentNumber, recordDate, movementType, status, companyID, warehouseID, responsibleID string, sourceDocument domain.SourceDocument, recipient domain.Recipient, details []domain.ShipmentDetail, financialSummary domain.ShipmentFinancialSummary, remarks, ipAddress string) (*domain.Shipment, error)
	Delete(id, ipAddress string) error
}

type shipmentUseCase struct {
	db          *sql.DB
	repo        repository.ShipmentRepository
	repoFactory repository.ShipmentRepoFactory
	historySvc  *HistoryService
}

func NewShipmentUseCase(db *sql.DB, repo repository.ShipmentRepository, repoFactory repository.ShipmentRepoFactory, historySvc *HistoryService) ShipmentUseCase {
	return &shipmentUseCase{db: db, repo: repo, repoFactory: repoFactory, historySvc: historySvc}
}

func (uc *shipmentUseCase) Create(shipmentNumber, recordDate, movementType, status, companyID, warehouseID, responsibleID string, sourceDocument domain.SourceDocument, recipient domain.Recipient, details []domain.ShipmentDetail, financialSummary domain.ShipmentFinancialSummary, remarks, ipAddress string) (*domain.Shipment, error) {
	if err := validateShipmentFields(shipmentNumber, recordDate, movementType, status, companyID, warehouseID, responsibleID, sourceDocument, recipient, details); err != nil {
		return nil, err
	}

	shipment := &domain.Shipment{
		ID:               generateID(),
		ShipmentNumber:   shipmentNumber,
		RecordDate:       recordDate,
		MovementType:     movementType,
		Status:           status,
		CompanyID:        companyID,
		WarehouseID:      warehouseID,
		ResponsibleID:    responsibleID,
		SourceDocument:   sourceDocument,
		Recipient:        recipient,
		Details:          details,
		FinancialSummary: financialSummary,
		Remarks:          remarks,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	shipmentRepo := uc.repoFactory(tx)
	if err := shipmentRepo.Create(shipment); err != nil {
		return nil, err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeSHIPMENT_CREATED, responsibleID, companyID,
		shipment.ID, "shipment", "Shipment "+shipment.ShipmentNumber+" created",
		ipAddress, nil, shipment,
	); err != nil {
		return nil, err
	}

	if len(sourceDocument.EntryIDs) > 0 {
		for _, entryID := range sourceDocument.EntryIDs {
			if err := uc.historySvc.LogRelation(tx, entryID, shipment.ID, responsibleID, companyID, ipAddress); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return shipment, nil
}

func (uc *shipmentUseCase) GetByID(id string) (*domain.Shipment, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *shipmentUseCase) GetAll() ([]*domain.Shipment, error) {
	return uc.repo.GetAll()
}

func (uc *shipmentUseCase) Update(id, shipmentNumber, recordDate, movementType, status, companyID, warehouseID, responsibleID string, sourceDocument domain.SourceDocument, recipient domain.Recipient, details []domain.ShipmentDetail, financialSummary domain.ShipmentFinancialSummary, remarks, ipAddress string) (*domain.Shipment, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateShipmentFields(shipmentNumber, recordDate, movementType, status, companyID, warehouseID, responsibleID, sourceDocument, recipient, details); err != nil {
		return nil, err
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	shipmentRepo := uc.repoFactory(tx)
	existing, err := shipmentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	previousData := *existing

	existing.ShipmentNumber = shipmentNumber
	existing.RecordDate = recordDate
	existing.MovementType = movementType
	existing.Status = status
	existing.CompanyID = companyID
	existing.WarehouseID = warehouseID
	existing.ResponsibleID = responsibleID
	existing.SourceDocument = sourceDocument
	existing.Recipient = recipient
	existing.Details = details
	existing.FinancialSummary = financialSummary
	existing.Remarks = remarks
	existing.UpdatedAt = time.Now()

	if err := shipmentRepo.Update(existing); err != nil {
		return nil, err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeUPDATE, responsibleID, companyID,
		id, "shipment", "Shipment "+existing.ShipmentNumber+" updated",
		ipAddress, previousData, existing,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return existing, nil
}

func (uc *shipmentUseCase) Delete(id, ipAddress string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	shipmentRepo := uc.repoFactory(tx)
	shipment, err := shipmentRepo.GetByID(id)
	if err != nil {
		return err
	}

	previousData := *shipment
	shipment.Status = "CANCELED"
	shipment.UpdatedAt = time.Now()

	if err := shipmentRepo.Update(shipment); err != nil {
		return err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeSHIPMENT_CANCELLED, shipment.ResponsibleID, shipment.CompanyID,
		id, "shipment", "Shipment "+shipment.ShipmentNumber+" canceled",
		ipAddress, previousData, shipment,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func validateShipmentFields(shipmentNumber, recordDate, movementType, status, companyID, warehouseID, responsibleID string, sourceDocument domain.SourceDocument, recipient domain.Recipient, details []domain.ShipmentDetail) error {
	if shipmentNumber == "" {
		return errors.New("shipment_number cannot be empty")
	}
	if recordDate == "" {
		return errors.New("record_date cannot be empty")
	}
	if movementType == "" {
		return errors.New("movement_type cannot be empty")
	}
	if !slices.Contains(validShipmentMovementTypes, movementType) {
		return errors.New("movement_type must be one of: SALE, SUPPLIER_RETURN, DONATION, SHRINKAGE, ADJUSTMENT, TRANSFER")
	}
	if status == "" {
		return errors.New("status cannot be empty")
	}
	if !slices.Contains(validShipmentStatuses, status) {
		return errors.New("status must be one of: DRAFT, CONFIRMED, CANCELED")
	}
	if companyID == "" {
		return errors.New("company_id cannot be empty")
	}
	if warehouseID == "" {
		return errors.New("warehouse_id cannot be empty")
	}
	if responsibleID == "" {
		return errors.New("responsible_id cannot be empty")
	}
	if recipient.RecipientType == "" {
		return errors.New("recipient_type cannot be empty")
	}
	if !slices.Contains(validRecipientTypes, recipient.RecipientType) {
		return errors.New("recipient_type must be one of: CUSTOMER, SUPPLIER, WAREHOUSE, INTERNAL")
	}
	if recipient.RecipientID == "" {
		return errors.New("recipient_id cannot be empty")
	}
	if len(details) == 0 {
		return errors.New("details cannot be empty")
	}
	for _, d := range details {
		if d.Code == "" {
			return errors.New("detail code cannot be empty")
		}
		if d.Product == "" {
			return errors.New("detail product cannot be empty")
		}
		if d.Unit == "" {
			return errors.New("detail unit cannot be empty")
		}
		if d.Quantity <= 0 {
			return errors.New("detail quantity must be greater than 0")
		}
		if d.UnitCost < 0 {
			return errors.New("detail unit_cost cannot be negative")
		}
	}
	return nil
}
