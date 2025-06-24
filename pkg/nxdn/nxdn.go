package nxdn

import (
	"fmt"

	"github.com/unklstewy/mmdvm_ghost/pkg/config"
)

func HandleNXDNPacket(packet []byte) {
	// TODO: Add NXDN packet handling logic
}

// Init initializes the NXDN protocol handler with the given configuration.
func Init(config config.NXDNConfig) {
	fmt.Printf("NXDN protocol handler initialized with Port: %s\n", config.Port)
}
