package dmr

import (
	"testing"
)

func TestAccessControl_ValidateSrcID(t *testing.T) {
	ac := &AccessControl{}
	ac.Init([]uint32{12345}, []uint32{67890}, nil, nil, nil, false, 0)

	if ac.ValidateSrcID(12345) {
		t.Errorf("Expected ID 12345 to be invalid")
	}

	if !ac.ValidateSrcID(67890) {
		t.Errorf("Expected ID 67890 to be valid")
	}
}

func TestAccessControl_ValidateTGID(t *testing.T) {
	ac := &AccessControl{}
	ac.Init(nil, nil, []uint32{100, 200}, []uint32{300, 400}, nil, false, 0)

	if !ac.ValidateTGID(1, true, 100) {
		t.Errorf("Expected TGID 100 on slot 1 to be valid")
	}

	if ac.ValidateTGID(1, true, 300) {
		t.Errorf("Expected TGID 300 on slot 1 to be invalid")
	}

	if !ac.ValidateTGID(2, true, 300) {
		t.Errorf("Expected TGID 300 on slot 2 to be valid")
	}
}
