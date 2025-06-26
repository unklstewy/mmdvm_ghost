package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/unklstewy/mmdvm_ghost/pkg/ax25"
	"github.com/unklstewy/mmdvm_ghost/pkg/config"
	"github.com/unklstewy/mmdvm_ghost/pkg/dmr"
	"github.com/unklstewy/mmdvm_ghost/pkg/dstar"
	"github.com/unklstewy/mmdvm_ghost/pkg/log"
	"github.com/unklstewy/mmdvm_ghost/pkg/m17"
	"github.com/unklstewy/mmdvm_ghost/pkg/nxdn"
	"github.com/unklstewy/mmdvm_ghost/pkg/pocsag"
	"github.com/unklstewy/mmdvm_ghost/pkg/ysf"
)

func handleSignals() chan os.Signal {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	return signalChan
}

func main() {
	configPath := flag.String("config", "mmdvm_ghost.db", "Path to SQLite configuration database")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

	log.InitLogger("mmdvm_ghost.log", "info", 10, 3, 7) // Updated logger initialization with required arguments
	log.Info("Starting mmdvm_ghost with config:", *configPath)
	if *verbose {
		log.Info("Verbose mode enabled")
	}

	// Initialize database schema and default values
	if err := config.InitializeDatabase(*configPath); err != nil {
		log.Fatal("Error initializing database:", err)
	}

	// Load configuration
	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	signalChan := handleSignals()
	reload := false

	log.Info("Starting main loop...")
	// Main loop
	for {
		if reload {
			log.Info("Reloading configuration...")
			// TODO: Reload configuration
			reload = false
		}

		log.Info("Initializing protocol handlers...")
		// Initialize protocol handlers here
		dmr.Init(config.DMR)
		dstar.Init(config.DStar)
		m17.Init(config.M17)
		ax25.Init(config.AX25)
		nxdn.Init(config.NXDN)
		pocsag.Init(config.Pocsag)
		ysf.Init(config.YSF)

		// Example usage of ProcessWakeup in the main loop
		data := []byte{dmr.TAG_DATA, dmr.DMR_IDLE_RX | dmr.DMR_SYNC_DATA | dmr.DT_CSBK, 0x01, 0x02}
		if err := dmr.ProcessWakeup(data); err != nil {
			log.Error("ProcessWakeup failed:", err)
		}

		// Simplified signal handling logic
		for sig := range signalChan {
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				log.Info("Exiting on signal:", sig)
				os.Exit(0)
			case syscall.SIGHUP:
				log.Info("Reloading on signal:", sig)
				reload = true
			}
		}

		// Prevent 100% CPU usage
		time.Sleep(100 * time.Millisecond)
	}
}
