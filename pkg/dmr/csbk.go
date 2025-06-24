package dmr

import (
	"errors"
	"log"
)

// Define constants for CSBKO types
const (
	BSDWNACT = 0x01 // Example value, replace with actual
)

// CSBK represents a Control Signaling Block in the DMR protocol.
type CSBK struct {
	Data        []byte
	CSBKO       uint8
	FID         uint8
	GI          bool
	BsID        uint32
	SrcID       uint32
	DstID       uint32
	DataContent bool
	CBF         uint8
	OVCM        bool
}

// NewCSBK creates a new CSBK instance.
func NewCSBK() *CSBK {
	csbk := &CSBK{
		Data: make([]byte, 12),
	}
	log.Printf("NewCSBK: Initialized Data field: %v", csbk.Data)
	return csbk
}

// Put decodes the CSBK data from the provided byte array.
func (c *CSBK) Put(bytes []byte) error {
	if len(bytes) < 12 {
		log.Printf("Put: Data too short, length: %d", len(bytes))
		return errors.New("data too short")
	}

	log.Printf("Put: Received bytes: %v", bytes)
	log.Printf("Put: Data before copy: %v", c.Data)

	// Ensure the Data field is properly initialized before copying data
	if c.Data == nil || len(c.Data) != 12 {
		c.Data = make([]byte, 12)
		log.Printf("Put: Reinitialized Data field: %v", c.Data)
	}

	// Decode the data (placeholder for BPTC19696 decoding logic)
	copy(c.Data, bytes[:12])
	log.Printf("Put: Decoded data: %v", c.Data)

	// Validate CRC (placeholder for CRC check logic)
	if !checkCRC(c.Data) {
		log.Printf("Put: Invalid CRC for data: %v", c.Data)
		return errors.New("invalid CRC")
	}

	c.CSBKO = c.Data[0] & 0x3F
	c.FID = c.Data[1]
	log.Printf("Put: CSBKO: %d, FID: %d", c.CSBKO, c.FID)

	// Extract fields based on CSBKO type
	switch c.CSBKO {
	case BSDWNACT:
		c.GI = false
		c.BsID = uint32(c.Data[4])<<16 | uint32(c.Data[5])<<8 | uint32(c.Data[6])
		c.SrcID = uint32(c.Data[7])<<16 | uint32(c.Data[8])<<8 | uint32(c.Data[9])
		c.DataContent = false
		c.CBF = 0
	default:
		log.Printf("Put: Unsupported CSBKO type: %d", c.CSBKO)
		return errors.New("unsupported CSBKO type")
	}

	log.Printf("Put: Extracted CSBKO: %d, FID: %d", c.CSBKO, c.FID)
	log.Printf("Put: Extracted BsID: %d, SrcID: %d", c.BsID, c.SrcID)
	log.Printf("Put: Extracted CSBKO value: %d", c.CSBKO) // Debug log for CSBKO value
	log.Printf("CSBK: Data field state: %v", c.Data)

	// Append CRC to the data before further processing
	c.Data = addCCITT161(c.Data)
	log.Printf("Put: Data with appended CRC: %v", c.Data)

	return nil
}

// GetCSBKO returns the CSBKO field.
func (c *CSBK) GetCSBKO() uint8 {
	return c.CSBKO
}

// GetSrcID returns the source ID field.
func (c *CSBK) GetSrcID() uint32 {
	return c.SrcID
}

// addCCITT161 computes a CRC-16-CCITT value for the input data and appends it.
func addCCITT161(data []byte) []byte {
	crc := uint16(0xFFFF)
	log.Printf("addCCITT161: Initial data: %v", data)
	for i := 0; i < len(data)-2; i++ {
		crc = uint16(byte(crc>>8)) ^ CCITT16_TABLE1[byte(crc)^data[i]]
	}
	crc = ^crc
	log.Printf("addCCITT161: Appended CRC: 0x%X 0x%X", byte(crc>>8), byte(crc))
	log.Printf("addCCITT161: Final data: %v", data)
	return append(data[:len(data)-2], byte(crc>>8), byte(crc))
}

// checkCCITT161 validates the CRC-16-CCITT value appended to the input data.
func checkCCITT161(data []byte) bool {
	if len(data) < 2 {
		return false
	}
	crc := uint16(0xFFFF)
	log.Printf("checkCCITT161: Data to validate: %v", data)
	for i := 0; i < len(data)-2; i++ {
		crc = uint16(byte(crc>>8)) ^ CCITT16_TABLE1[byte(crc)^data[i]]
	}
	crc = ^crc
	log.Printf("checkCCITT161: Computed CRC: 0x%X 0x%X", byte(crc>>8), byte(crc))
	log.Printf("checkCCITT161: Validation result: %v", byte(crc>>8) == data[len(data)-2] && byte(crc) == data[len(data)-1])
	return byte(crc>>8) == data[len(data)-2] && byte(crc) == data[len(data)-1]
}

// checkCRC validates the CRC-16-CCITT value appended to the input data.
func checkCRC(data []byte) bool {
	return checkCCITT161(data)
}

// Define the CCITT16_TABLE1 lookup table.
var CCITT16_TABLE1 = [256]uint16{
	0x0000, 0x1021, 0x2042, 0x3063, 0x4084, 0x50A5, 0x60C6, 0x70E7,
	// ...remaining table values...
}

func HandleCSBK(data []byte) {
	csbk := NewCSBK()
	err := csbk.Put(data)
	if err != nil {
		log.Printf("Error processing CSBK: %v", err)
		return
	}

	log.Printf("Processed CSBK with CSBKO: %d, SrcID: %d", csbk.GetCSBKO(), csbk.GetSrcID())
}
