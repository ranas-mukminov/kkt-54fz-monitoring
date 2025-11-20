package domain

import "time"

// KKTStatus represents the operational status of a KKT device
type KKTStatus int

const (
	KKTStatusUnavailable KKTStatus = iota // 0 - device is unavailable
	KKTStatusRunning                       // 1 - device is running normally
	KKTStatusError                         // 2 - device has errors
)

// ShiftStatus represents the status of a shift
type ShiftStatus int

const (
	ShiftStatusClosed ShiftStatus = iota // 0 - shift is closed
	ShiftStatusOpen                       // 1 - shift is open
)

// OFDSyncStatus represents OFD synchronization status
type OFDSyncStatus int

const (
	OFDSyncStatusUnknown OFDSyncStatus = iota // 0 - status unknown
	OFDSyncStatusSynced                        // 1 - synchronized
	OFDSyncStatusPending                       // 2 - pending sync
	OFDSyncStatusError                         // 3 - sync error
)

// KKTDevice represents a cash register device
type KKTDevice struct {
	ID              string        `json:"id"`
	FactoryNumber   string        `json:"factory_number"`
	RegNumber       string        `json:"reg_number"`
	FiscalDriveNum  string        `json:"fiscal_drive_num"`
	Status          KKTStatus     `json:"status"`
	LastSeen        time.Time     `json:"last_seen"`
	ShiftStatus     ShiftStatus   `json:"shift_status"`
	OFDSyncStatus   OFDSyncStatus `json:"ofd_sync_status"`
	FiscalDriveInfo FiscalDrive   `json:"fiscal_drive_info"`
}

// FiscalDrive represents fiscal drive information
type FiscalDrive struct {
	Number       string    `json:"number"`
	ExpiryDate   time.Time `json:"expiry_date"`
	DocumentsMax int       `json:"documents_max"`
	DocumentsUsed int      `json:"documents_used"`
	MemoryUsage  float64   `json:"memory_usage"` // percentage
}

// FiscalDocument represents a fiscal document
type FiscalDocument struct {
	ID              string            `json:"id"`
	Type            DocumentType      `json:"type"`
	KKTID           string            `json:"kkt_id"`
	FiscalSign      string            `json:"fiscal_sign"`
	DocumentNumber  int               `json:"document_number"`
	ShiftNumber     int               `json:"shift_number"`
	DateTime        time.Time         `json:"date_time"`
	Amount          float64           `json:"amount"`
	OperationType   OperationType     `json:"operation_type"`
	TaxationSystem  TaxationSystem    `json:"taxation_system"`
	Items           []DocumentItem    `json:"items,omitempty"`
	RawData         map[string]interface{} `json:"raw_data,omitempty"`
}

// DocumentType represents the type of fiscal document
type DocumentType int

const (
	DocumentTypeReceipt DocumentType = iota + 1
	DocumentTypeReceiptReturn
	DocumentTypeReceiptCorrection
	DocumentTypeOpenShift
	DocumentTypeCloseShift
	DocumentTypeRegistration
	DocumentTypeReRegistration
	DocumentTypeCloseArchive
)

// OperationType represents the operation type
type OperationType int

const (
	OperationTypeSale OperationType = iota + 1
	OperationTypeSaleReturn
	OperationTypePurchase
	OperationTypePurchaseReturn
)

// TaxationSystem represents the taxation system
type TaxationSystem int

const (
	TaxationSystemCommon TaxationSystem = iota + 1
	TaxationSystemSimplified
	TaxationSystemSimplifiedMinusCosts
	TaxationSystemSingleTax
	TaxationSystemPatent
)

// DocumentItem represents an item in a fiscal document
type DocumentItem struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Amount   float64 `json:"amount"`
	VATRate  VATRate `json:"vat_rate"`
}

// VATRate represents VAT rate
type VATRate int

const (
	VATRateNone VATRate = iota
	VATRate0
	VATRate10
	VATRate20
)

// KKTError represents an error from KKT device
type KKTError struct {
	ID          string       `json:"id"`
	KKTID       string       `json:"kkt_id"`
	ErrorCode   string       `json:"error_code"`
	ErrorType   ErrorType    `json:"error_type"`
	Severity    ErrorSeverity `json:"severity"`
	Message     string       `json:"message"`
	Timestamp   time.Time    `json:"timestamp"`
	Resolved    bool         `json:"resolved"`
	ResolvedAt  *time.Time   `json:"resolved_at,omitempty"`
}

// ErrorType represents the type of error
type ErrorType int

const (
	ErrorTypeNetwork ErrorType = iota + 1
	ErrorTypeFiscalDrive
	ErrorTypeOFD
	ErrorTypePrinter
	ErrorTypeHardware
	ErrorTypeSoftware
	ErrorTypeConfiguration
)

// ErrorSeverity represents the severity of an error
type ErrorSeverity int

const (
	ErrorSeverityInfo ErrorSeverity = iota + 1
	ErrorSeverityWarning
	ErrorSeverityError
	ErrorSeverityCritical
)

// OFDTransaction represents a transaction with OFD
type OFDTransaction struct {
	ID             string    `json:"id"`
	KKTID          string    `json:"kkt_id"`
	DocumentID     string    `json:"document_id"`
	SentAt         time.Time `json:"sent_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	Status         OFDSyncStatus `json:"status"`
	RetryCount     int       `json:"retry_count"`
	LastError      string    `json:"last_error,omitempty"`
}

// Metrics represents aggregated metrics for KKT monitoring
type Metrics struct {
	KKTID              string              `json:"kkt_id"`
	Timestamp          time.Time           `json:"timestamp"`
	Status             KKTStatus           `json:"status"`
	DocumentsTotal     int64               `json:"documents_total"`
	ErrorsByType       map[ErrorType]int64 `json:"errors_by_type"`
	OFDSyncStatus      OFDSyncStatus       `json:"ofd_sync_status"`
	ShiftStatus        ShiftStatus         `json:"shift_status"`
	LastDocumentTime   time.Time           `json:"last_document_time"`
	FDMemoryUsage      float64             `json:"fd_memory_usage"`
	DocumentsPerHour   float64             `json:"documents_per_hour"`
	AverageSyncTime    float64             `json:"average_sync_time"` // seconds
}
