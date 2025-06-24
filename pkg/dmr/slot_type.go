package dmr

import (
	"errors"
)

// SlotType represents the DMR slot type data.
type SlotType struct {
	ColorCode uint8
	DataType  uint8
}

// PutData decodes the slot type data from the provided byte array.
func (s *SlotType) PutData(data []byte) error {
	if len(data) < 21 {
		return errors.New("data array too short")
	}

	DMRSlotType := make([]uint8, 3)
	DMRSlotType[0] = (data[12] << 2) & 0xFC
	DMRSlotType[0] |= (data[13] >> 6) & 0x03

	DMRSlotType[1] = (data[13] << 2) & 0xC0
	DMRSlotType[1] |= (data[19] << 2) & 0x3C
	DMRSlotType[1] |= (data[20] >> 6) & 0x03

	DMRSlotType[2] = (data[20] << 2) & 0xF0

	code, err := DecodeGolay2087(DMRSlotType)
	if err != nil {
		return err
	}

	s.ColorCode = (code >> 4) & 0x0F
	s.DataType = code & 0x0F
	return nil
}

// GetData encodes the slot type data into the provided byte array.
func (s *SlotType) GetData(data []byte) error {
	if len(data) < 21 {
		return errors.New("data array too short")
	}

	DMRSlotType := make([]uint8, 3)
	DMRSlotType[0] = (s.ColorCode << 4) & 0xF0
	DMRSlotType[0] |= s.DataType & 0x0F
	DMRSlotType[1] = 0x00
	DMRSlotType[2] = 0x00

	EncodeGolay2087(DMRSlotType)

	data[12] = (data[12] & 0xC0) | ((DMRSlotType[0] >> 2) & 0x3F)
	data[13] = (data[13] & 0x0F) | ((DMRSlotType[0] << 6) & 0xC0) | ((DMRSlotType[1] >> 2) & 0x30)
	data[19] = (data[19] & 0xF0) | ((DMRSlotType[1] >> 2) & 0x0F)
	data[20] = (data[20] & 0x03) | ((DMRSlotType[1] << 6) & 0xC0) | ((DMRSlotType[2] >> 2) & 0x3C)
	return nil
}

// DecodeGolay2087 decodes a Golay2087 code.
func DecodeGolay2087(data []uint8) (uint8, error) {
	// TODO: Implement Golay2087 decoding logic
	return 0, nil
}

// EncodeGolay2087 encodes a Golay2087 code.
func EncodeGolay2087(data []uint8) {
	// TODO: Implement Golay2087 encoding logic
}
