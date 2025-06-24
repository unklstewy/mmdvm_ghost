package dmr

import "errors"

// DataHeader represents a DMR data header.
type DataHeader struct {
	Data   []byte
	GI     bool
	A      bool
	SrcID  uint32
	DstID  uint32
	Blocks uint8
	F      bool
	S      bool
	Ns     uint8
}

// NewDataHeader creates a new DataHeader instance.
func NewDataHeader() *DataHeader {
	return &DataHeader{
		Data: make([]byte, 12),
	}
}

// Put decodes the data header from the provided byte array.
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
