package domain

import "time"

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
