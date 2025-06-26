// Package dmr provides DMR protocol logic, including CSBK (Control Signaling Block) handling and related utilities.
package dmr

import (
	"errors" // For error handling
	"fmt"    // For formatted I/O
	"log"    // For logging debug/info messages

	"github.com/unklstewy/mmdvm_ghost/pkg/bptc" // For BPTC19696 decoding/correction
)

// Define constants for CSBKO types (Control Signaling Block Opcode)
const (
	BSDWNACT = 0x01 // Downlink Activate (example value, replace with actual if needed)
	UUVREQ   = 0x02 // Unit to Unit Service Request (example value, replace with actual if needed)
)

// CSBK represents a Control Signaling Block in the DMR protocol.
// It holds the decoded payload and extracted fields.
type CSBK struct {
	Data        []byte // Raw decoded CSBK payload (12 bytes)
	CSBKO       uint8  // CSBK Opcode (lower 6 bits of Data[0])
	FID         uint8  // Feature ID (Data[1])
	GI          bool   // Group/Individual flag
	BsID        uint32 // Base Station ID
	SrcID       uint32 // Source ID
	DstID       uint32 // Destination ID
	DataContent bool   // Indicates if data content is present
	CBF         uint8  // Control Block Flag
	OVCM        bool   // Over-the-air Voice Call Management flag
}

// NewCSBK creates a new CSBK instance with a zeroed 12-byte Data field.
func NewCSBK() *CSBK {
	csbk := &CSBK{
		Data: make([]byte, 12),
	}
	log.Printf("NewCSBK: Initialized Data field: %v", csbk.Data)
	return csbk
}

// DecodeBPTC19696 decodes the BPTC19696 data using the BPTC package's CorrectBPTCData function.
func DecodeBPTC19696(data []byte) ([]byte, error) {
	return bptc.CorrectBPTCData(data)
}

// Put decodes the CSBK data from the provided byte array.
// It expects 33 bytes, decodes/corrects them, extracts fields, and validates CRC.
func (c *CSBK) Put(bytes []byte) error {
	// Ensure input is exactly 33 bytes (pad or truncate as needed)
	if len(bytes) < 33 {
		log.Printf("Put: Data too short, padding to 33 bytes")
		paddedBytes := make([]byte, 33)
		copy(paddedBytes, bytes)
		bytes = paddedBytes
	} else if len(bytes) > 33 {
		log.Printf("Put: Data too long, truncating to 33 bytes")
		bytes = bytes[:33]
	}

	log.Printf("Put: Received bytes: %v", bytes)

	// Decode the data using BPTC19696 logic (calls CorrectBPTCData)
	log.Printf("Put: Passing data to DecodeBPTC19696: %v", bytes)
	decodedData, err := DecodeBPTC19696(bytes)
	if err != nil {
		log.Printf("Put: Error decoding BPTC19696 data: %v", err)
		return err
	}

	// Copy the decoded data to the CSBK structure (first 12 bytes)
	c.Data = make([]byte, 12)
	copy(c.Data, decodedData[:12])
	log.Printf("Put: Decoded data: %v", c.Data)

	// Append CRC to the data before validation (for checkCRC)
	c.Data = addCCITT161(c.Data)
	log.Printf("Put: Data with appended CRC: %v", c.Data)

	// Validate CRC
	if !checkCRC(c.Data) {
		log.Printf("Put: Invalid CRC for data: %v", c.Data)
		return errors.New("invalid CRC")
	}

	// Extract CSBKO (lower 6 bits of first byte)
	c.CSBKO = c.Data[0] & 0x3F
	c.FID = c.Data[1]
	log.Printf("Put: CSBKO: %d, FID: %d", c.CSBKO, c.FID)

	// Add detailed logging to trace CSBKO value and data processing
	log.Printf("Put: Initial bytes: %v", bytes)
	log.Printf("Put: Decoded data before CRC: %v", decodedData)
	log.Printf("Put: CSBKO extracted: %d", c.CSBKO)
	log.Printf("Put: Raw first byte of decoded data: %d", c.Data[0])
	log.Printf("Put: CSBKO after masking: %d", c.CSBKO)

	// Extract fields based on CSBKO type
	switch c.CSBKO {
	case BSDWNACT:
		// Downlink Activate: extract BsID and SrcID from payload
		c.GI = false
		c.BsID = uint32(c.Data[4])<<16 | uint32(c.Data[5])<<8 | uint32(c.Data[6])
		c.SrcID = uint32(c.Data[7])<<16 | uint32(c.Data[8])<<8 | uint32(c.Data[9])
		c.DataContent = false
		c.CBF = 0
		log.Printf("Downlink Activate CSBK: %+v", c)
	case UUVREQ:
		// Unit to Unit Service Request: extract DstID and SrcID
		c.GI = false
		c.DstID = uint32(c.Data[4])<<16 | uint32(c.Data[5])<<8 | uint32(c.Data[6])
		c.SrcID = uint32(c.Data[7])<<16 | uint32(c.Data[8])<<8 | uint32(c.Data[9])
		c.DataContent = false
		c.CBF = 0
		c.OVCM = (c.Data[2] & 0x04) == 0x04
		log.Printf("Unit to Unit Service Request CSBK: %+v", c)
	default:
		// Unknown/unsupported CSBKO
		log.Printf("Put: Unsupported CSBKO type: %d, Data: %v", c.CSBKO, c.Data)
		return fmt.Errorf("unsupported CSBKO type: %d", c.CSBKO)
	}

	return nil
}

// GetCSBKO returns the CSBKO field (opcode) from the CSBK structure.
func (c *CSBK) GetCSBKO() uint8 {
	return c.CSBKO
}

// GetSrcID returns the source ID field from the CSBK structure.
func (c *CSBK) GetSrcID() uint32 {
	return c.SrcID
}

// addCCITT161 computes a CRC-16-CCITT value for the input data and appends it (big-endian order).
func addCCITT161(data []byte) []byte {
	crc := uint16(0xFFFF)
	log.Printf("addCCITT161: Initial data: %v", data)
	for i := 0; i < len(data); i++ {
		log.Printf("addCCITT161: Before processing byte %d (0x%X), CRC: 0x%X", i, data[i], crc)
		crc = uint16(byte(crc>>8)) ^ CCITT16_TABLE1[byte(crc)^data[i]]
		log.Printf("addCCITT161: After processing byte %d (0x%X), CRC: 0x%X", i, data[i], crc)
	}
	crc = ^crc // Finalize CRC (invert bits)
	log.Printf("addCCITT161: Final computed CRC: 0x%X 0x%X", byte(crc>>8), byte(crc))
	return append(data, byte(crc>>8), byte(crc))
}

// checkCCITT161 validates the CRC-16-CCITT value appended to the input data.
func checkCCITT161(data []byte) bool {
	if len(data) < 2 {
		log.Printf("checkCCITT161: Data too short to validate CRC: %v", data)
		return false
	}
	crc := uint16(0xFFFF)
	log.Printf("checkCCITT161: Data to validate: %v", data)
	for i := 0; i < len(data)-2; i++ {
		log.Printf("checkCCITT161: Before processing byte %d (0x%X), CRC: 0x%X", i, data[i], crc)
		crc = uint16(byte(crc>>8)) ^ CCITT16_TABLE1[byte(crc)^data[i]]
		log.Printf("checkCCITT161: After processing byte %d (0x%X), CRC: 0x%X", i, data[i], crc)
	}
	crc = ^crc
	log.Printf("checkCCITT161: Final computed CRC: 0x%X 0x%X", byte(crc>>8), byte(crc))
	log.Printf("checkCCITT161: Validation result: %v", byte(crc>>8) == data[len(data)-2] && byte(crc) == data[len(data)-1])
	return byte(crc>>8) == data[len(data)-2] && byte(crc) == data[len(data)-1]
}

// checkCRC validates the CRC-16-CCITT value appended to the input data (wrapper for checkCCITT161).
func checkCRC(data []byte) bool {
	return checkCCITT161(data)
}

// CCITT16_TABLE1 is a lookup table for fast CRC-16-CCITT calculation.
var CCITT16_TABLE1 = [256]uint16{
	0x0000, 0x1021, 0x2042, 0x3063, 0x4084, 0x50A5, 0x60C6, 0x70E7,
	0x8108, 0x9129, 0xA14A, 0xB16B, 0xC18C, 0xD1AD, 0xE1CE, 0xF1EF,
	0x1231, 0x0210, 0x3273, 0x2252, 0x52B5, 0x4294, 0x72F7, 0x62D6,
	0x9339, 0x8318, 0xB37B, 0xA35A, 0xD3BD, 0xC39C, 0xF3FF, 0xE3DE,
	0x2462, 0x3443, 0x0420, 0x1401, 0x64E6, 0x74C7, 0x44A4, 0x5485,
	0xA56A, 0xB54B, 0x8528, 0x9509, 0xE5EE, 0xF5CF, 0xC5AC, 0xD58D,
	0x3653, 0x2672, 0x1611, 0x0630, 0x76D7, 0x66F6, 0x5695, 0x46B4,
	0xB75B, 0xA77A, 0x9719, 0x8738, 0xF7DF, 0xE7FE, 0xD79D, 0xC7BC,
	0x48C4, 0x58E5, 0x6886, 0x78A7, 0x0840, 0x1861, 0x2802, 0x3823,
	0xC9CC, 0xD9ED, 0xE98E, 0xF9AF, 0x8948, 0x9969, 0xA90A, 0xB92B,
	0x5AF5, 0x4AD4, 0x7AB7, 0x6A96, 0x1A71, 0x0A50, 0x3A33, 0x2A12,
	0xDBFD, 0xCBDC, 0xFBBF, 0xEB9E, 0x9B79, 0x8B58, 0xBB3B, 0xAB1A,
	0x6CA6, 0x7C87, 0x4CE4, 0x5CC5, 0x2C22, 0x3C03, 0x0C60, 0x1C41,
	0xEDAE, 0xFD8F, 0xCDEC, 0xDDCD, 0xAD2A, 0xBD0B, 0x8D68, 0x9D49,
	0x7E97, 0x6EB6, 0x5ED5, 0x4EF4, 0x3E13, 0x2E32, 0x1E51, 0x0E70,
	0xFF9F, 0xEFBE, 0xDFDD, 0xCFFC, 0xBF1B, 0xAF3A, 0x9F59, 0x8F78,
	0x9188, 0x81A9, 0xB1CA, 0xA1EB, 0xD10C, 0xC12D, 0xF14E, 0xE16F,
	0x1080, 0x00A1, 0x30C2, 0x20E3, 0x5004, 0x4025, 0x7046, 0x6067,
	0x83B9, 0x9398, 0xA3FB, 0xB3DA, 0xC33D, 0xD31C, 0xE37F, 0xF35E,
	0x02B1, 0x1290, 0x22F3, 0x32D2, 0x4235, 0x5214, 0x6277, 0x7256,
	0xB5EA, 0xA5CB, 0x95A8, 0x8589, 0xF56E, 0xE54F, 0xD52C, 0xC50D,
	0x34E2, 0x24C3, 0x14A0, 0x0481, 0x7466, 0x6447, 0x5424, 0x4405,
	0xA7DB, 0xB7FA, 0x8799, 0x97B8, 0xE75F, 0xF77E, 0xC71D, 0xD73C,
	0x26D3, 0x36F2, 0x0691, 0x16B0, 0x6657, 0x7676, 0x4615, 0x5634,
	0xD94C, 0xC96D, 0xF90E, 0xE92F, 0x99C8, 0x89E9, 0xB98A, 0xA9AB,
	0x5844, 0x4865, 0x7806, 0x6827, 0x18C0, 0x08E1, 0x3882, 0x28A3,
	0xCB7D, 0xDB5C, 0xEB3F, 0xFB1E, 0x8BF9, 0x9BD8, 0xABBB, 0xBB9A,
	0x4A75, 0x5A54, 0x6A37, 0x7A16, 0x0AF1, 0x1AD0, 0x2AB3, 0x3A92,
	0xFD2E, 0xED0F, 0xDD6C, 0xCD4D, 0xBDAA, 0xAD8B, 0x9DE8, 0x8DC9,
	0x7C26, 0x6C07, 0x5C64, 0x4C45, 0x3CA2, 0x2C83, 0x1CE0, 0x0CC1,
	0xEF1F, 0xFF3E, 0xCF5D, 0xDF7C, 0xAF9B, 0xBFBA, 0x8FD9, 0x9FF8,
	0x6E17, 0x7E36, 0x4E55, 0x5E74, 0x2E93, 0x3EB2, 0x0ED1, 0x1EF0,
}
