package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/usecase"
)

type ShipmentHandler struct {
	uc usecase.ShipmentUseCase
}

func NewShipmentHandler(uc usecase.ShipmentUseCase) *ShipmentHandler {
	return &ShipmentHandler{uc: uc}
}

func (h *ShipmentHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/shipments")
	id := strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		if id != "" {
			h.getByID(w, r, id)
		} else {
			h.getAll(w, r)
		}
	case http.MethodPost:
		h.create(w, r)
	case http.MethodPut:
		if id == "" {
			writeError(w, "id is required", http.StatusBadRequest)
			return
		}
		h.update(w, r, id)
	case http.MethodDelete:
		if id == "" {
			writeError(w, "id is required", http.StatusBadRequest)
			return
		}
		h.delete(w, r, id)
	default:
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ShipmentHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ShipmentNumber   string                        `json:"shipment_number"`
		RecordDate       string                        `json:"record_date"`
		MovementType     string                        `json:"movement_type"`
		Status           string                        `json:"status"`
		CompanyID        string                        `json:"company_id"`
		WarehouseID      string                        `json:"warehouse_id"`
		ResponsibleID    string                        `json:"responsible_id"`
		SourceDocument   domain.SourceDocument         `json:"source_document"`
		Recipient        domain.Recipient              `json:"recipient"`
		Details          []domain.ShipmentDetail       `json:"details"`
		FinancialSummary domain.ShipmentFinancialSummary `json:"financial_summary"`
		Remarks          string                        `json:"remarks"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	shipment, err := h.uc.Create(req.ShipmentNumber, req.RecordDate, req.MovementType, req.Status, req.CompanyID, req.WarehouseID, req.ResponsibleID, req.SourceDocument, req.Recipient, req.Details, req.FinancialSummary, req.Remarks, extractIP(r))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, shipment, http.StatusCreated)
}

func (h *ShipmentHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	shipment, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, shipment, http.StatusOK)
}

func (h *ShipmentHandler) getAll(w http.ResponseWriter, r *http.Request) {
	shipments, err := h.uc.GetAll()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, shipments, http.StatusOK)
}

func (h *ShipmentHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		ShipmentNumber   string                        `json:"shipment_number"`
		RecordDate       string                        `json:"record_date"`
		MovementType     string                        `json:"movement_type"`
		Status           string                        `json:"status"`
		CompanyID        string                        `json:"company_id"`
		WarehouseID      string                        `json:"warehouse_id"`
		ResponsibleID    string                        `json:"responsible_id"`
		SourceDocument   domain.SourceDocument         `json:"source_document"`
		Recipient        domain.Recipient              `json:"recipient"`
		Details          []domain.ShipmentDetail       `json:"details"`
		FinancialSummary domain.ShipmentFinancialSummary `json:"financial_summary"`
		Remarks          string                        `json:"remarks"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	shipment, err := h.uc.Update(id, req.ShipmentNumber, req.RecordDate, req.MovementType, req.Status, req.CompanyID, req.WarehouseID, req.ResponsibleID, req.SourceDocument, req.Recipient, req.Details, req.FinancialSummary, req.Remarks, extractIP(r))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, shipment, http.StatusOK)
}

func (h *ShipmentHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id, extractIP(r)); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
