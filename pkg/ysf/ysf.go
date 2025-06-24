package ysf

import (
	"fmt"

	"github.com/unklstewy/mmdvm_ghost/pkg/config"
)

func HandleYSFPacket(packet []byte) {
	// TODO: Add YSF packet handling logic
}

// Init initializes the YSF protocol handler with the given configuration.
func Init(config config.YSFConfig) {
	fmt.Printf("YSF protocol handler initialized with Port: %s\n", config.Port)
}
