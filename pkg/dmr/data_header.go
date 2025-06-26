// Package dmr provides DMR protocol logic, including the DataHeader structure for DMR data headers.
package dmr

import "errors" // For error handling

// DataHeader represents a DMR data header, holding decoded fields and payload.
type DataHeader struct {
	Data   []byte // Raw decoded data header (12 bytes)
	GI     bool   // Group/Individual flag
	A      bool   // Address type flag
	SrcID  uint32 // Source ID
	DstID  uint32 // Destination ID
	Blocks uint8  // Number of blocks
	F      bool   // F flag (purpose depends on spec)
	S      bool   // S flag (purpose depends on spec)
	Ns     uint8  // Ns field (purpose depends on spec)
}

// NewDataHeader creates a new DataHeader instance with a zeroed 12-byte Data field.
func NewDataHeader() *DataHeader {
	return &DataHeader{
		Data: make([]byte, 12),
	}
}

// Put decodes the data header from the provided byte array.
// It expects at least 12 bytes, copies them, and extracts fields.
func (d *DataHeader) Put(bytes []byte) error {
	if len(bytes) < 12 {
		return errors.New("data too short")
	}

	// Decode the data (placeholder for BPTC19696 decoding logic)
	copy(d.Data, bytes[:12])

	// Validate CRC (placeholder for CRC check logic)
	if !checkCRC(d.Data) {
		return errors.New("invalid CRC")
	}

	d.GI = (d.Data[0] & 0x80) == 0x80
	d.A = (d.Data[0] & 0x40) == 0x40

	d.DstID = uint32(d.Data[2])<<16 | uint32(d.Data[3])<<8 | uint32(d.Data[4])
	d.SrcID = uint32(d.Data[5])<<16 | uint32(d.Data[6])<<8 | uint32(d.Data[7])

	return nil
}
