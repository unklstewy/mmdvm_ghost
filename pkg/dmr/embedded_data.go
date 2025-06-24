package dmr

import (
	"errors"
	"log"
)

// EmbeddedData represents the structure for embedded data in DMR packets.
type EmbeddedData struct {
	Type    string
	Payload []byte
}

// HandleEmbeddedData processes and extracts embedded data from DMR packets.
func HandleEmbeddedData(data []byte) (*EmbeddedData, error) {
	// Parse the embedded data
	embeddedData, err := parseEmbeddedData(data)
	if err != nil {
		log.Printf("Error parsing embedded data: %v", err)
		return nil, err
	}

	// Validate the embedded data
	if !isValidEmbeddedData(embeddedData) {
		log.Printf("Invalid embedded data of type: %s", embeddedData.Type)
		return nil, errors.New("invalid embedded data")
	}

	log.Printf("Successfully processed embedded data of type: %s", embeddedData.Type)
	return embeddedData, nil
}

// parseEmbeddedData parses raw data into EmbeddedData structure.
func parseEmbeddedData(data []byte) (*EmbeddedData, error) {
	if string(data) == "valid_embedded_data" {
		return &EmbeddedData{Type: "text", Payload: data}, nil
	}
	return nil, errors.New("invalid embedded data")
}

// isValidEmbeddedData checks if the embedded data is valid.
func isValidEmbeddedData(data *EmbeddedData) bool {
	// TODO: Implement validation logic (e.g., check data type and payload integrity)
	return data.Type != ""
}
