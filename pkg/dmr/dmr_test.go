package dmr

import (
	"bytes"
	"io"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHandleDMRPacket(t *testing.T) {
	// Test case: Valid wakeup packet
	validPacket := []byte{TagData, DmrIdleRx | DmrSyncData | DtCSBK, 0x01, 0x02, 0x03}
	HandleDMRPacket(validPacket)

	// Test case: Invalid packet (too short)
	shortPacket := []byte{TagData}
	HandleDMRPacket(shortPacket)

	// Test case: Invalid wakeup packet
	invalidPacket := []byte{TagData, 0xFF, 0x01, 0x02, 0x03}
	HandleDMRPacket(invalidPacket)
}

func TestCSBK(t *testing.T) {
	csbk := &CSBK{}
	validData := []byte{BSDWNACT, 0x02, 0x03, 0x04, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	shortData := []byte{0x01, 0x02}

	t.Logf("TestCSBK: Valid data: %v", validData)
	t.Logf("TestCSBK: Short data: %v", shortData)

	// Test case: Valid CSBK data
	if err := csbk.Put(validData); err != nil {
		t.Errorf("Expected valid CSBK data to be processed successfully, got error: %v", err)
	}

	// Test case: Invalid CSBK data (too short)
	if err := csbk.Put(shortData); err == nil {
		t.Errorf("Expected error for short CSBK data, got nil")
	}

	// Test case: Retrieve CSBKO value
	csbko := csbk.GetCSBKO()
	if csbko != BSDWNACT {
		t.Errorf("Expected CSBKO to be %d, got %d", BSDWNACT, csbko)
	}
}

func TestCSBK_Put(t *testing.T) {
	csbk := &CSBK{}
	validData := []byte{BSDWNACT, 0x02, 0x03, 0x04, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	shortData := []byte{0x01, 0x02}

	// Check for errors instead of using the `!` operator.
	if err := csbk.Put(validData); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := csbk.Put(shortData); err == nil {
		t.Errorf("expected error for short data, got nil")
	}
}

func TestCSBK_GetCSBKO(t *testing.T) {
	// Use the Put method to initialize the Data field
	c := &CSBK{}
	data := []byte{BSDWNACT, 0x02, 0x03, 0x04, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	if err := c.Put(data); err != nil {
		t.Fatalf("unexpected error in Put: %v", err)
	}

	// Compare CSBKO to the expected value
	csbko := c.GetCSBKO()
	if csbko != BSDWNACT {
		t.Errorf("expected CSBKO to be %d, got %d", BSDWNACT, csbko)
	}
}

func TestHandleAccessControl(t *testing.T) {
	// Test case: Valid access control data
	validData := []byte("12345,admin")
	err := HandleAccessControl(validData)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Test case: Invalid access control data
	invalidData := []byte("invalid_access_data")
	err = HandleAccessControl(invalidData)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestHandleEmbeddedData(t *testing.T) {
	// Test case: Valid embedded data
	validData := []byte("valid_embedded_data")
	data, err := HandleEmbeddedData(validData)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if data == nil || data.Type == "" {
		t.Errorf("Expected valid embedded data, got: %v", data)
	}

	// Test case: Invalid embedded data
	invalidData := []byte("invalid_embedded_data")
	data, err = HandleEmbeddedData(invalidData)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if data != nil {
		t.Errorf("Expected nil data, got: %v", data)
	}
}

// Add unit tests for DMRSlot methods.

func TestDMRSlot_TransitionState(t *testing.T) {
	slot := &DMRSlot{ID: 1, State: "Idle"}
	slot.TransitionState("Active")

	if slot.State != "Active" {
		t.Errorf("expected state to be 'Active', got '%s'", slot.State)
	}
}

func TestDMRSlot_HandleEmbeddedData(t *testing.T) {
	slot := &DMRSlot{ID: 1}
	data := []byte{0x01, 0x02, 0x03}

	// Capture output for validation
	output := captureOutput(func() {
		slot.HandleEmbeddedData(data)
	})

	expectedOutput := "Slot 1 received embedded data: 010203\n"
	if output != expectedOutput {
		t.Errorf("expected output '%s', got '%s'", expectedOutput, output)
	}
}

// Helper function to capture output.
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

// Consolidated setupTestDatabase function to initialize the test database schema and insert test data.
func setupTestDatabase() *gorm.DB {
	// Reset the test database by deleting the existing file
	os.Remove("test.db")

	// Initialize the database connection
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	// Create the access_rules table if it doesn't exist
	db.AutoMigrate(&AccessRule{})

	// Updated test data to exclude the rule for `ID=54321` and `Role=user`
	testRules := []AccessRule{
		{ID: "12345", Role: "admin"},
	}
	for _, rule := range testRules {
		db.FirstOrCreate(&rule, AccessRule{ID: rule.ID})
	}

	return db
}

// Initialize the database with the `access_rules` table and test data
func initTestDatabase() {
	db := getDatabaseConnection()
	db.Exec("CREATE TABLE IF NOT EXISTS access_rules (id TEXT PRIMARY KEY, role TEXT);")
	db.Exec("INSERT OR IGNORE INTO access_rules (id, role) VALUES ('12345', 'admin');")
}

func TestIsValidAccess(t *testing.T) {
	dbInstance = setupTestDatabase() // Use the test database

	validData := &AccessControlData{ID: "12345", Role: "admin"}
	if !isValidAccess(validData) {
		t.Errorf("expected access to be valid for ID: %s, Role: %s", validData.ID, validData.Role)
	}

	invalidData := &AccessControlData{ID: "54321", Role: "user"}
	if isValidAccess(invalidData) {
		t.Errorf("expected access to be invalid for ID: %s, Role: %s", invalidData.ID, invalidData.Role)
	}
}

func TestUpdateBlackList(t *testing.T) {
	ac := &AccessControl{}
	newBlackList := []uint32{12345, 67890}
	ac.UpdateBlackList(newBlackList)

	if len(ac.BlackList) != len(newBlackList) {
		t.Errorf("expected blacklist length %d, got %d", len(newBlackList), len(ac.BlackList))
	}
}

func TestUpdateWhiteList(t *testing.T) {
	ac := &AccessControl{}
	newWhiteList := []uint32{11111, 22222}
	ac.UpdateWhiteList(newWhiteList)

	if len(ac.WhiteList) != len(newWhiteList) {
		t.Errorf("expected whitelist length %d, got %d", len(newWhiteList), len(ac.WhiteList))
	}
}

func TestReloadRules(t *testing.T) {
	dbInstance = setupTestDatabase() // Use the test database
	ac := &AccessControl{}

	if err := ac.ReloadRules(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// Call the initialization function in the test setup
func TestMain(m *testing.M) {
	initTestDatabase()
	os.Exit(m.Run())
}

func TestCSBKPut(t *testing.T) {
	csbk := NewCSBK()
	input := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C}
	err := csbk.Put(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(csbk.Data) != 12 {
		t.Errorf("Expected Data length 12, got: %d", len(csbk.Data))
	}

	for i, b := range input {
		if csbk.Data[i] != b {
			t.Errorf("Data mismatch at index %d: expected %v, got %v", i, b, csbk.Data[i])
		}
	}
}
