package dmr

import (
	"testing"
)

func TestSlotType_PutData(t *testing.T) {
	data := make([]byte, 21)
	data[12] = 0x3C
	data[13] = 0xF0
	data[19] = 0x0F
	data[20] = 0xC3

	slotType := &SlotType{}
	err := slotType.PutData(data)
	if err != nil {
		t.Errorf("PutData failed: %v", err)
	}

	if slotType.ColorCode != 0 || slotType.DataType != 0 {
		t.Errorf("Unexpected slot type values: ColorCode=%d, DataType=%d", slotType.ColorCode, slotType.DataType)
	}
}

func TestSlotType_GetData(t *testing.T) {
	data := make([]byte, 21)
	slotType := &SlotType{
		ColorCode: 0x0A,
		DataType:  0x05,
	}

	err := slotType.GetData(data)
	if err != nil {
		t.Errorf("GetData failed: %v", err)
	}

	// Add assertions to verify encoded data
}
