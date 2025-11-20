package domain

import (
	"testing"
	"time"
)

func TestKKTDevice(t *testing.T) {
	device := KKTDevice{
		ID:            "kkt-001",
		FactoryNumber: "12345678",
		RegNumber:     "0001234567890123",
		Status:        KKTStatusRunning,
		LastSeen:      time.Now(),
		ShiftStatus:   ShiftStatusOpen,
		OFDSyncStatus: OFDSyncStatusSynced,
	}

	if device.ID != "kkt-001" {
		t.Errorf("Expected ID kkt-001, got %s", device.ID)
	}

	if device.Status != KKTStatusRunning {
		t.Errorf("Expected status Running, got %d", device.Status)
	}
}

func TestFiscalDocument(t *testing.T) {
	doc := FiscalDocument{
		ID:             "doc-001",
		Type:           DocumentTypeReceipt,
		KKTID:          "kkt-001",
		DocumentNumber: 12345,
		ShiftNumber:    5,
		DateTime:       time.Now(),
		Amount:         1000.50,
		OperationType:  OperationTypeSale,
		TaxationSystem: TaxationSystemCommon,
		Items: []DocumentItem{
			{
				Name:     "Test Item",
				Quantity: 2,
				Price:    500.25,
				Amount:   1000.50,
				VATRate:  VATRate20,
			},
		},
	}

	if doc.Type != DocumentTypeReceipt {
		t.Errorf("Expected type Receipt, got %d", doc.Type)
	}

	if len(doc.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(doc.Items))
	}

	if doc.Items[0].Name != "Test Item" {
		t.Errorf("Expected item name 'Test Item', got %s", doc.Items[0].Name)
	}
}

func TestKKTError(t *testing.T) {
	err := KKTError{
		ID:        "err-001",
		KKTID:     "kkt-001",
		ErrorCode: "E001",
		ErrorType: ErrorTypeFiscalDrive,
		Severity:  ErrorSeverityCritical,
		Message:   "Fiscal drive error",
		Timestamp: time.Now(),
		Resolved:  false,
	}

	if err.ErrorType != ErrorTypeFiscalDrive {
		t.Errorf("Expected error type FiscalDrive, got %d", err.ErrorType)
	}

	if err.Severity != ErrorSeverityCritical {
		t.Errorf("Expected severity Critical, got %d", err.Severity)
	}

	if err.Resolved {
		t.Error("Expected error to be unresolved")
	}
}

func TestMetrics(t *testing.T) {
	metrics := Metrics{
		KKTID:            "kkt-001",
		Timestamp:        time.Now(),
		Status:           KKTStatusRunning,
		DocumentsTotal:   1000,
		ErrorsByType:     map[ErrorType]int64{ErrorTypeNetwork: 5},
		OFDSyncStatus:    OFDSyncStatusSynced,
		ShiftStatus:      ShiftStatusOpen,
		LastDocumentTime: time.Now(),
		FDMemoryUsage:    45.5,
		DocumentsPerHour: 50.0,
		AverageSyncTime:  2.5,
	}

	if metrics.KKTID != "kkt-001" {
		t.Errorf("Expected KKTID kkt-001, got %s", metrics.KKTID)
	}

	if metrics.DocumentsTotal != 1000 {
		t.Errorf("Expected DocumentsTotal 1000, got %d", metrics.DocumentsTotal)
	}

	if metrics.ErrorsByType[ErrorTypeNetwork] != 5 {
		t.Errorf("Expected 5 network errors, got %d", metrics.ErrorsByType[ErrorTypeNetwork])
	}
}
