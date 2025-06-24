package m17

import (
	"fmt"

	"github.com/unklstewy/mmdvm_ghost/pkg/config"
)

func HandleM17Packet(packet []byte) {
	// TODO: Add M17 packet handling logic
}

// Init initializes the M17 protocol handler with the given configuration.
func Init(cfg config.M17Config) {
	// TODO: Add initialization logic for M17
	fmt.Printf("M17 protocol handler initialized with CAN: %s\n", cfg.CAN)
}
