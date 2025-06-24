package pocsag

import (
	"fmt"

	"github.com/unklstewy/mmdvm_ghost/pkg/config"
)

func HandlePOCSAGPacket(packet []byte) {
	// TODO: Add POCSAG packet handling logic
}

// Init initializes the POCSAG protocol handler with the given configuration.
func Init(config config.PocsagConfig) {
	fmt.Printf("POCSAG protocol handler initialized with Frequency: %d\n", config.Frequency)
}
