package dmr

import (
	"log"
	"testing"
)

// MockLookup is a mock implementation of the Lookup interface.
type MockLookup struct{}

// Find returns a mock result for a given ID.
func (m *MockLookup) Find(id uint32) string {
	return "MockSource"
}

func TestControl_ProcessWakeup(t *testing.T) {
	lookup := &MockLookup{}
	accessControl := &AccessControl{}
	accessControl.Init(nil, nil, nil, nil, nil, false, 0)

	control := NewControl(1, 100, lookup, accessControl)

	// Test case: Valid wakeup packet
	validData := []byte{TagData, DmrIdleRx | DmrSyncData | DtCSBK, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C}
	log.Printf("TestControl_ProcessWakeup: Valid data: %v", validData)
	err := control.ProcessWakeup(validData)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Test case: Invalid wakeup packet
	invalidData := []byte{TagData, 0xFF, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C}
	log.Printf("TestControl_ProcessWakeup: Invalid data: %v", invalidData)
	err = control.ProcessWakeup(invalidData)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
