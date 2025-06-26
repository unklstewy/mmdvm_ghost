// Package bptc provides functions for decoding and correcting BPTC19696 data blocks.
package bptc

import (
	"errors" // Provides error handling utilities
	"fmt"
	"log" // Provides logging utilities for debugging
	"testing"
)

// CorrectBPTCData decodes and corrects BPTC19696 data.
// It expects a 33-byte input, processes it through several steps (bit extraction, deinterleaving, error correction, payload extraction),
// and returns the final 12-byte payload or an error if the input is invalid.
func CorrectBPTCData(data []byte) ([]byte, error) {
	// Check that the input data is exactly 33 bytes (as required by BPTC19696)
	if len(data) != 33 {
		log.Printf("CorrectBPTCData: Invalid data length: %d, expected 33 bytes", len(data))
		return nil, errors.New("invalid data length, expected 33 bytes")
	}

	log.Printf("CorrectBPTCData: Starting with input data: %v", data)

	// Step 1: Convert the 33 bytes into a slice of 196 bits (bools)
	rawData := extractBinary(data)
	log.Printf("CorrectBPTCData: Extracted raw binary data: %v", rawData)

	// Step 2: Deinterleave the bits to restore their original order
	deInterData := deInterleave(rawData)
	log.Printf("CorrectBPTCData: Deinterleaved data: %v", deInterData)

	// Step 3: Perform error correction (Hamming code) on the deinterleaved bits
	correctedData := errorCheck(deInterData)
	log.Printf("CorrectBPTCData: Corrected data after errorCheck: %v", correctedData)

	// Step 4: Extract the 12-byte payload from the corrected bits
	payload := extractPayload(correctedData)
	log.Printf("CorrectBPTCData: Final extracted payload: %v", payload)

	return payload, nil
}

// extractBinary converts a 33-byte input into a 196-bit slice (bools),
// extracting each bit in big-endian order (MSB first for each byte).
func extractBinary(data []byte) []bool {
	rawData := make([]bool, 196) // Output slice for 196 bits
	bitIndex := 0                // Tracks the current bit position

	log.Printf("extractBinary: Input data: %v", data)

	// Iterate over each byte in the input
	for byteIndex, byteVal := range data {
		log.Printf("extractBinary: Processing byte %d (0x%X)", byteIndex, byteVal)
		// For each byte, extract bits from MSB (bit 7) to LSB (bit 0)
		for i := 7; i >= 0; i-- {
			// Stop if we've filled all 196 bits
			if bitIndex >= len(rawData) {
				log.Printf("extractBinary: Warning - Excess bits ignored after index %d", bitIndex)
				return rawData
			}
			// Extract the i-th bit and store as a bool
			rawData[bitIndex] = (byteVal & (1 << i)) != 0
			log.Printf("extractBinary: Bit %d: %v", bitIndex, rawData[bitIndex])
			bitIndex++
		}
	}

	log.Printf("extractBinary: Extracted binary data: %v", rawData)
	return rawData
}

// deInterleave reverses the interleaving process applied to the 196 bits.
// It uses the formula (i * 181) % 196 to map each output bit to its original position.
func deInterleave(rawData []bool) []bool {
	deInterData := make([]bool, 196) // Output slice for deinterleaved bits
	log.Printf("deInterleave: Input raw data: %v", rawData)
	// For each output bit position, compute the original (interleaved) index
	for i := 0; i < 196; i++ {
		interleaveSequence := (i * 181) % 196        // Calculate the original bit position
		deInterData[i] = rawData[interleaveSequence] // Assign the bit to its deinterleaved position
	}
	log.Printf("deInterleave: Deinterleaved data: %v", deInterData)
	return deInterData
}

// errorCheck applies Hamming(8,4) error correction to the deinterleaved bits.
// It processes each group of 8 bits, checks/corrects parity bits, and flips a bit if a single-bit error is detected.
func errorCheck(deInterData []bool) []bool {
	correctedData := make([]bool, len(deInterData)) // Output slice for corrected bits
	copy(correctedData, deInterData)                // Start with a copy of the input

	log.Printf("errorCheck: Input deinterleaved data: %v", deInterData)

	// Process each 8-bit block (Hamming(8,4) code)
	for i := 0; i < len(correctedData); i += 8 {
		// If less than 8 bits remain, skip (incomplete block)
		if i+7 >= len(correctedData) {
			log.Printf("errorCheck: Warning - Incomplete byte encountered at index %d", i)
			break
		}

		// Parity bit positions for Hamming(8,4): 0, 1, 3, 7
		parityBits := []int{0, 1, 3, 7}
		parityError := 0 // Accumulates the error position (if any)
		// For each parity bit, check if the parity matches
		for _, p := range parityBits {
			parity := false
			for j := 0; j < 8; j++ {
				// Only include bits that are not the parity bit itself and are covered by this parity
				if j != p && ((j+1)&(p+1)) != 0 {
					parity = parity != correctedData[i+j]
				}
			}
			// If the calculated parity does not match the stored parity, increment error position
			if parity != correctedData[i+p] {
				parityError += (p + 1)
			}
		}

		// If a single-bit error is detected (error position in 1..8), flip the bit
		if parityError > 0 && parityError <= 8 {
			log.Printf("errorCheck: Correcting bit at index %d", i+parityError-1)
			correctedData[i+parityError-1] = !correctedData[i+parityError-1]
		}
	}

	log.Printf("errorCheck: Corrected data: %v", correctedData)
	return correctedData
}

// extractPayload converts the first 96 bits (12 bytes) of the corrected data into a byte slice.
// It packs each group of 8 bits into a byte, with the most significant bit first.
func extractPayload(correctedData []bool) []byte {
	payload := make([]byte, 12) // Output slice for 12 bytes
	// For each byte in the payload
	for i := 0; i < 12; i++ {
		// For each bit in the byte (MSB to LSB)
		for j := 0; j < 8; j++ {
			bitIndex := i*8 + j // Compute the bit index in the corrected data
			// If the bit is set, set the corresponding bit in the output byte
			if bitIndex < len(correctedData) && correctedData[bitIndex] {
				payload[i] |= (1 << (7 - j)) // Set bit (MSB first)
			}
		}
	}
	log.Printf("extractPayload: Extracted payload: %v", payload)
	return payload
}

func TestDMRFrameFromCapture(t *testing.T) {
	// Example: First three frames' UDP payloads from the capture (as extracted from first10frames.hex)
	frames := [][]byte{
		// Frame 1
		{
			0x4a, 0x3d, 0xa8, 0x48, 0x00, 0x00, 0x5b, 0xa0,
			0x59, 0x90, 0xf7, 0x05, 0x2b, 0x00, 0x00, 0x00,
			0x99, 0xe9, 0xa3, 0x72, 0x20, 0x64, 0x54, 0x4f,
			0x7e, 0x99, 0xe9, 0xa3, 0x72, 0x21, 0x10, 0x00,
			0x00, 0x00, 0x0e, 0x20, 0x64, 0x54, 0x4f, 0x7e,
			0x99, 0xe9, 0xa3, 0x72, 0x20, 0x64, 0x54, 0x4f,
			0x7e, 0x00, 0x00,
		},
		// Frame 2
		{
			0x4b, 0x3d, 0xa8, 0x48, 0x00, 0x00, 0x5b, 0xa0,
			0x59, 0x90, 0xf7, 0x10, 0x2b, 0x00, 0x00, 0x00,
			0x99, 0xe9, 0xa3, 0x72, 0x20, 0x64, 0x54, 0x4f,
			0x7e, 0x99, 0xe9, 0xa3, 0x72, 0x27, 0x55, 0xfd,
			0x7d, 0xf7, 0x5f, 0x70, 0x64, 0x54, 0x4f, 0x7e,
			0x99, 0xe9, 0xa3, 0x72, 0x20, 0x64, 0x54, 0x4f,
			0x7e, 0x00, 0x00,
		},
		// Frame 3
		{
			0x4c, 0x3d, 0xa8, 0x48, 0x00, 0x00, 0x5b, 0xa0,
			0x59, 0x90, 0xf7, 0x01, 0x2b, 0x00, 0x00, 0x00,
			0x99, 0xe9, 0xa3, 0x72, 0x20, 0x64, 0x54, 0x4f,
			0x7e, 0x99, 0xe9, 0xa3, 0x72, 0x21, 0x30, 0x30,
			0xa0, 0x02, 0x19, 0x10, 0x64, 0x54, 0x4f, 0x7e,
			0x99, 0xe9, 0xa3, 0x72, 0x20, 0x64, 0x54, 0x4f,
			0x7e, 0x00, 0x00,
		},
	}

	for i, payload := range frames {
		t.Run(fmt.Sprintf("Frame_%d", i+1), func(t *testing.T) {
			result, err := CorrectBPTCData(payload)
			if err != nil {
				t.Errorf("Failed to decode DMR frame %d: %v", i+1, err)
			} else {
				t.Logf("Frame %d decoded payload: %x", i+1, result)
			}
		})
	}
}
