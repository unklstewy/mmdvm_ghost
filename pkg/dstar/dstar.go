package dstar

import (
	"fmt"

	"github.com/unklstewy/mmdvm_ghost/pkg/config"
)

func HandleDStarPacket(packet []byte) {
	// TODO: Add D-Star packet handling logic
}

// Init initializes the D-Star protocol handler with the given configuration.
func Init(cfg config.DStarConfig) {
	// TODO: Add initialization logic for D-Star
	fmt.Printf("D-Star protocol handler initialized with Module: %s\n", cfg.Module)
}
