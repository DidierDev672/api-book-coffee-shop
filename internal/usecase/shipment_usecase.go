package usecase

import (
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
	Create(shipmentNumber, recordDate, movementType, status, companyID, warehouseID, responsibleID string, sourceDocument domain.SourceDocument, recipient domain.Recipient, details []domain.ShipmentDetail, financialSummary domain.ShipmentFinancialSummary, remarks string) (*domain.Shipment, error)
	GetByID(id string) (*domain.Shipment, error)
	GetAll() ([]*domain.Shipment, error)
	Update(id, shipmentNumber, recordDate, movementType, status, companyID, warehouseID, responsibleID string, sourceDocument domain.SourceDocument, recipient domain.Recipient, details []domain.ShipmentDetail, financialSummary domain.ShipmentFinancialSummary, remarks string) (*domain.Shipment, error)
	Delete(id string) error
}

type shipmentUseCase struct {
	repo repository.ShipmentRepository
}

func NewShipmentUseCase(repo repository.ShipmentRepository) ShipmentUseCase {
	return &shipmentUseCase{repo: repo}
}

func (uc *shipmentUseCase) Create(shipmentNumber, recordDate, movementType, status, companyID, warehouseID, responsibleID string, sourceDocument domain.SourceDocument, recipient domain.Recipient, details []domain.ShipmentDetail, financialSummary domain.ShipmentFinancialSummary, remarks string) (*domain.Shipment, error) {
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

	if err := uc.repo.Create(shipment); err != nil {
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

func (uc *shipmentUseCase) Update(id, shipmentNumber, recordDate, movementType, status, companyID, warehouseID, responsibleID string, sourceDocument domain.SourceDocument, recipient domain.Recipient, details []domain.ShipmentDetail, financialSummary domain.ShipmentFinancialSummary, remarks string) (*domain.Shipment, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateShipmentFields(shipmentNumber, recordDate, movementType, status, companyID, warehouseID, responsibleID, sourceDocument, recipient, details); err != nil {
		return nil, err
	}

	shipment, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	shipment.ShipmentNumber = shipmentNumber
	shipment.RecordDate = recordDate
	shipment.MovementType = movementType
	shipment.Status = status
	shipment.CompanyID = companyID
	shipment.WarehouseID = warehouseID
	shipment.ResponsibleID = responsibleID
	shipment.SourceDocument = sourceDocument
	shipment.Recipient = recipient
	shipment.Details = details
	shipment.FinancialSummary = financialSummary
	shipment.Remarks = remarks
	shipment.UpdatedAt = time.Now()

	if err := uc.repo.Update(shipment); err != nil {
		return nil, err
	}
	return shipment, nil
}

func (uc *shipmentUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
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
