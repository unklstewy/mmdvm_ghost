package dmr

import (
	"log"
	"testing"
)

func TestEncode(t *testing.T) {
	msg := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	parity, err := Encode(msg, 9)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if len(parity) != NPAR {
		t.Errorf("Expected %d parity bytes, got %d", NPAR, len(parity))
	}
}

func TestCheck(t *testing.T) {
	// Generate valid test data with parity bytes
	msg := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	parity, err := Encode(msg, len(msg))
	if err != nil {
		t.Fatalf("Failed to generate parity bytes: %v", err)
	}
	data := append(msg, parity...)

	log.Printf("TestCheck: Generated parity bytes: %v", parity)
	log.Printf("TestCheck: Data being checked: %v", data)

	valid, err := Check(data)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if !valid {
		t.Errorf("Expected data to be valid, but it was not")
	}
}
