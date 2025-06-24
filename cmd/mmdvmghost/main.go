package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/unklstewy/mmdvm_ghost/pkg/ax25"
	"github.com/unklstewy/mmdvm_ghost/pkg/config"
	"github.com/unklstewy/mmdvm_ghost/pkg/dmr"
	"github.com/unklstewy/mmdvm_ghost/pkg/dstar"
	"github.com/unklstewy/mmdvm_ghost/pkg/m17"
	"github.com/unklstewy/mmdvm_ghost/pkg/nxdn"
	"github.com/unklstewy/mmdvm_ghost/pkg/pocsag"
	"github.com/unklstewy/mmdvm_ghost/pkg/ysf"
)

func main() {
	configPath := flag.String("config", "mmdvm_ghost.db", "Path to SQLite configuration database")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

	fmt.Printf("Starting mmdvm_ghost with config: %s\n", *configPath)
	if *verbose {
		fmt.Println("Verbose mode enabled")
	}

	// Initialize database schema and default values
	if err := config.InitializeDatabase(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	config, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize protocol handlers
	fmt.Println("Initializing protocol handlers...")
	dmr.Init(config.DMR)
	dstar.Init(config.DStar)
	m17.Init(config.M17)
	ax25.Init(config.AX25)
	nxdn.Init(config.NXDN)
	pocsag.Init(config.Pocsag)
	ysf.Init(config.YSF)

	// Main loop
	fmt.Println("Starting main loop...")
	for {
		// TODO: Process incoming packets and delegate to protocol handlers
		// Example: packet := modem.ReceivePacket()
		// Example: dmr.HandleDMRPacket(packet)
		// Example: dstar.HandleDStarPacket(packet)
		// Example: m17.HandleM17Packet(packet)

		// Prevent 100% CPU usage
		time.Sleep(100 * time.Millisecond)
	}
}
