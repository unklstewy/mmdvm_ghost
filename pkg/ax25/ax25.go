package ax25

import (
	"fmt"

	"github.com/unklstewy/mmdvm_ghost/pkg/config"
)

func HandleAX25Packet(packet []byte) {
	// TODO: Add AX.25 packet handling logic
}

// Init initializes the AX.25 protocol handler with the given configuration.
func Init(config config.AX25Config) {
	fmt.Printf("AX.25 protocol handler initialized with Port: %s\n", config.Port)
}
