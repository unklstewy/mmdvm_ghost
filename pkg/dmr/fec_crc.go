package dmr

import (
	"errors"
	"log"
)

const (
	NPAR = 3 // Number of parity bytes
)

var (
	POLY = []byte{64, 56, 14, 1, 0, 0, 0, 0, 0, 0, 0, 0}
)

// gmult performs multiplication in Galois Field (2^8).
func gmult(a, b byte) byte {
	var p byte
	for i := 0; i < 8; i++ {
		if b&1 != 0 {
			p ^= a
		}
		carry := a & 0x80
		a <<= 1
		if carry != 0 {
			a ^= 0x1D // Primitive polynomial for GF(2^8)
		}
		b >>= 1
	}
	return p
}

// Encode generates parity bytes for the given message.
func Encode(msg []byte, nbytes int) ([]byte, error) {
	if len(msg) < nbytes {
		return nil, errors.New("message too short")
	}

	parity := make([]byte, NPAR)
	log.Printf("Encode: Initial parity state: %v", parity)
	for i := 0; i < nbytes; i++ {
		dbyte := msg[i] ^ parity[NPAR-1]
		log.Printf("Encode: Processing byte %d: 0x%X, dbyte: 0x%X", i, msg[i], dbyte)
		for j := NPAR - 1; j > 0; j-- {
			parity[j] = parity[j-1] ^ gmult(POLY[j], dbyte)
			log.Printf("Encode: Updated parity[%d]: 0x%X", j, parity[j])
		}
		parity[0] = gmult(POLY[0], dbyte)
		log.Printf("Encode: Updated parity[0]: 0x%X", parity[0])
	}
	log.Printf("Encode: Final parity state: %v", parity)

	// Add debug log to trace generated parity bytes
	log.Printf("Encode: Generated parity bytes: %v", parity)
	log.Printf("Encode: Parity byte order: %v", parity)

	return parity, nil
}

// Check validates the input data using Reed-Solomon error correction.
func Check(data []byte) (bool, error) {
	if len(data) < 12 {
		return false, errors.New("data too short")
	}

	parity, err := Encode(data[:9], 9)
	if err != nil {
		return false, err
	}

	// Add debug log to trace comparison in Check function
	log.Printf("Check: Data parity bytes: %v, Generated parity bytes: %v", data[9:], parity)
	log.Printf("Check: Expected parity bytes: %v, Actual parity bytes: %v", data[9:], parity)
	log.Printf("Check: Comparison results: %v, %v, %v", data[9] == parity[2], data[10] == parity[1], data[11] == parity[0])
	log.Printf("Check: Comparing data parity bytes: %v with generated parity bytes: %v", data[9:], parity)

	return data[9] == parity[0] && data[10] == parity[1] && data[11] == parity[2], nil
}
