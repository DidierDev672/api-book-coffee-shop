package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/usecase"
)

type SaleHandler struct {
	uc usecase.SaleUseCase
}

func NewSaleHandler(uc usecase.SaleUseCase) *SaleHandler {
	return &SaleHandler{uc: uc}
}

func (h *SaleHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/sales")
	path = strings.TrimSuffix(path, "/")
	id := strings.TrimPrefix(path, "/")

	if id != "" {
		if strings.HasSuffix(id, "/status") {
			saleID := strings.TrimSuffix(id, "/status")
			saleID = strings.TrimSuffix(saleID, "/")
			if r.Method == http.MethodPatch {
				h.updateStatus(w, r, saleID)
				return
			}
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if strings.HasSuffix(id, "/discount") {
			saleID := strings.TrimSuffix(id, "/discount")
			saleID = strings.TrimSuffix(saleID, "/")
			if r.Method == http.MethodPatch {
				h.updateDiscount(w, r, saleID)
				return
			}
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		if id != "" && !strings.Contains(id, "/") {
			h.getByID(w, r, id)
		} else {
			h.getAll(w, r)
		}
	case http.MethodPost:
		if id == "" {
			h.create(w, r)
		} else {
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodPut:
		if id != "" && !strings.Contains(id, "/") {
			h.update(w, r, id)
		} else {
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodDelete:
		if id != "" && !strings.Contains(id, "/") {
			h.delete(w, r, id)
		} else {
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *SaleHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SaleNumber    string                 `json:"sale_number"`
		OrderID       string                 `json:"order_id"`
		OrderType     string                 `json:"order_type"`
		ProviderID    string                 `json:"provider_id"`
		WarehouseID   string                 `json:"warehouse_id"`
		Products      []domain.SaleDetail    `json:"products"`
		Subtotal      float64                `json:"subtotal"`
		VAT           float64                `json:"vat"`
		Discount      float64                `json:"discount"`
		Total         float64                `json:"total"`
		PaymentMethod string                 `json:"payment_method"`
		CreatedBy     string                 `json:"created_by"`
		CompanyID     string                 `json:"company_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := usecase.CreateSaleInput{
		SaleNumber:    req.SaleNumber,
		OrderID:       req.OrderID,
		OrderType:     req.OrderType,
		ProviderID:    req.ProviderID,
		WarehouseID:   req.WarehouseID,
		Products:      req.Products,
		Subtotal:      req.Subtotal,
		VAT:           req.VAT,
		Discount:      req.Discount,
		Total:         req.Total,
		PaymentMethod: req.PaymentMethod,
		CreatedBy:     req.CreatedBy,
		CompanyID:     req.CompanyID,
	}

	sale, err := h.uc.Create(input, extractIP(r))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, sale, http.StatusCreated)
}

func (h *SaleHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		ProviderID    string              `json:"provider_id"`
		ClientID      string              `json:"client_id"`
		WarehouseID   string              `json:"warehouse_id"`
		OrderType     string              `json:"order_type"`
		Products      []domain.SaleDetail `json:"products"`
		Subtotal      float64             `json:"subtotal"`
		VAT           float64             `json:"vat"`
		Discount      float64             `json:"discount"`
		Total         float64             `json:"total"`
		PaymentMethod string              `json:"payment_method"`
		Status        string              `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	clientID := req.ClientID
	if clientID == "" {
		clientID = req.ProviderID
	}

	input := usecase.UpdateSaleInput{
		SaleID:        id,
		ClientID:      clientID,
		WarehouseID:   req.WarehouseID,
		OrderType:     req.OrderType,
		Products:      req.Products,
		Subtotal:      req.Subtotal,
		VAT:           req.VAT,
		Discount:      req.Discount,
		Total:         req.Total,
		PaymentMethod: req.PaymentMethod,
		Status:        req.Status,
	}

	sale, err := h.uc.Update(input, extractIP(r))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, sale, http.StatusOK)
}

func (h *SaleHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(id, extractIP(r)); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]string{"message": "sale deleted"}, http.StatusOK)
}

func (h *SaleHandler) getAll(w http.ResponseWriter, r *http.Request) {
	filters := make(map[string]string)
	if v := r.URL.Query().Get("client_id"); v != "" {
		filters["client_id"] = v
	}
	if v := r.URL.Query().Get("status"); v != "" {
		filters["status"] = strings.ToUpper(v)
	}
	if v := r.URL.Query().Get("payment_method"); v != "" {
		filters["payment_method"] = v
	}
	if v := r.URL.Query().Get("date_from"); v != "" {
		filters["date_from"] = v
	}
	if v := r.URL.Query().Get("date_to"); v != "" {
		filters["date_to"] = v
	}
	if v := r.URL.Query().Get("company_id"); v != "" {
		filters["company_id"] = v
	}

	// Pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	sales, err := h.uc.GetAll(filters)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Apply offset/limit manually for server-side pagination
	start := (page - 1) * limit
	if start >= len(sales) {
		sales = nil
	} else {
		end := start + limit
		if end > len(sales) {
			end = len(sales)
		}
		sales = sales[start:end]
	}

	response := map[string]interface{}{
		"data":  sales,
		"page":  page,
		"limit": limit,
	}
	writeJSON(w, response, http.StatusOK)
}

func (h *SaleHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	sale, err := h.uc.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, sale, http.StatusOK)
}

func (h *SaleHandler) updateStatus(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sale, err := h.uc.UpdateStatus(id, strings.ToUpper(req.Status), extractIP(r))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, sale, http.StatusOK)
}

func (h *SaleHandler) updateDiscount(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Discount float64 `json:"discount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sale, err := h.uc.UpdateDiscount(id, req.Discount, extractIP(r))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, sale, http.StatusOK)
}
