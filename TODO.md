# TODO

## Epic: Port MMDVMHost-cpp to Go

### Chapter 1: Core Components

#### Story 1.1: Port Main Application Logic
**Priority:** High
- Reference: `MMDVMHost.cpp`, `MMDVMHost.h`
- **Task:** Create the main application loop and initialization logic.

```go
package main

func main() {
	// Initialize configuration
	config, err := config.LoadConfig("mmdvm_ghost.db")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Start main loop
	for {
		// TODO: Add main loop logic
	}
}
```

#### Story 1.2: Port Configuration Management
**Priority:** High
- Reference: `Conf.cpp`, `Conf.h`
- **Task:** Use SQLite for configuration management.

```go
package config

func LoadConfig(path string) (*Config, error) {
	// TODO: Implement SQLite-based configuration loading
	return nil, nil
}
```

### Chapter 2: Protocol Implementations

#### Story 2.1: Port DMR Protocol
**Priority:** High
- Reference: `DMRControl.cpp`, `DMRNetwork.cpp`, `DMRData.cpp`, etc.
- **Task:** Implement DMR protocol logic.

```go
package dmr

func HandleDMRPacket(packet []byte) {
	// TODO: Add DMR packet handling logic
}
```

#### Story 2.2: Port D-Star Protocol
**Priority:** Medium
- Reference: `DStarControl.cpp`, `DStarNetwork.cpp`, `DStarHeader.cpp`, etc.
- **Task:** Implement D-Star protocol logic.

```go
package dstar

func HandleDStarPacket(packet []byte) {
	// TODO: Add D-Star packet handling logic
}
```

#### Story 2.3: Port M17 Protocol
**Priority:** Medium
- Reference: `M17Control.cpp`, `M17Network.cpp`, `M17CRC.cpp`, etc.
- **Task:** Implement M17 protocol logic.

```go
package m17

func HandleM17Packet(packet []byte) {
	// TODO: Add M17 packet handling logic
}
```

### Chapter 3: Utilities

#### Story 3.1: Port Logging
**Priority:** High
- Reference: `Log.cpp`, `Log.h`
- **Task:** Implement structured logging.

```go
package log

func InitLogger() {
	// TODO: Initialize logger
}
```

#### Story 3.2: Port CRC and Error Correction
**Priority:** Medium
- Reference: `CRC.cpp`, `Golay2087.cpp`, `Hamming.cpp`, etc.
- **Task:** Implement CRC and error correction utilities.

```go
package utils

func CalculateCRC(data []byte) uint16 {
	// TODO: Add CRC calculation logic
	return 0
}
```

### Chapter 4: Hardware Interfaces

#### Story 4.1: Port Display Handling
**Priority:** Medium
- Reference: `HD44780.cpp`, `OLED.cpp`, etc.
- **Task:** Implement display handling logic.

```go
package display

func InitDisplay() {
	// TODO: Initialize display
}
```

#### Story 4.2: Port Modem Communication
**Priority:** Medium
- Reference: `Modem.cpp`, `ModemPort.cpp`, etc.
- **Task:** Implement modem communication logic.

```go
package modem

func InitModem() {
	// TODO: Initialize modem
}
```

### Chapter 5: Testing and Validation

#### Story 5.1: Write Unit Tests
**Priority:** High
- **Task:** Write unit tests for each package.

```go
package main

func TestMainLogic(t *testing.T) {
	// TODO: Add test cases
}
```

### Chapter 6: Documentation

#### Story 6.1: Create Documentation
**Priority:** Low
- **Task:** Document the Go implementation and its differences from the original C++ code.

```markdown
# Documentation

## Overview
- This project is a Go port of MMDVMHost-cpp.

## Features
- Feature-for-feature parity with the original C++ implementation.
```

### Chapter 7: Additional Protocols and Utilities

#### Story 7.1: Port AMBEFEC Logic
**Priority:** Medium
- Reference: `AMBEFEC.cpp`, `AMBEFEC.h`
- **Task:** Implement AMBE Forward Error Correction logic.

```go
package ambe

func CorrectAMBEData(data []byte) ([]byte, error) {
	// TODO: Add AMBE FEC logic
	return nil, nil
}
```

#### Story 7.2: Port AX.25 Protocol
**Priority:** Medium
- Reference: `AX25Control.cpp`, `AX25Control.h`, `AX25Defines.h`, `AX25Network.cpp`, `AX25Network.h`
- **Task:** Implement AX.25 protocol handling.

```go
package ax25

func HandleAX25Packet(packet []byte) {
	// TODO: Add AX.25 packet handling logic
}
```

#### Story 7.3: Port BCH Error Correction
**Priority:** Medium
- Reference: `BCH.cpp`, `BCH.h`
- **Task:** Implement BCH error correction.

```go
package bch

func CorrectBCHData(data []byte) ([]byte, error) {
	// TODO: Add BCH error correction logic
	return nil, nil
}
```

#### Story 7.4: Port DMR Slot Management
**Priority:** High
- Reference: `DMRSlot.cpp`, `DMRSlot.h`, `DMRSlotType.cpp`, `DMRSlotType.h`
- **Task:** Implement DMR slot management.

```go
package dmr

func ManageDMRSlot(slotID int, data []byte) {
	// TODO: Add DMR slot management logic
}
```

#### Story 7.5: Port NXDN Protocol
**Priority:** Medium
- Reference: `NXDNControl.cpp`, `NXDNNetwork.cpp`, `NXDNCRC.cpp`, etc.
- **Task:** Implement NXDN protocol handling.

```go
package nxdn

func HandleNXDNPacket(packet []byte) {
	// TODO: Add NXDN packet handling logic
}
```

#### Story 7.6: Port YSF Protocol
**Priority:** Medium
- Reference: `YSFControl.cpp`, `YSFNetwork.cpp`, `YSFFICH.cpp`, etc.
- **Task:** Implement YSF protocol handling.

```go
package ysf

func HandleYSFPacket(packet []byte) {
	// TODO: Add YSF packet handling logic
}
```

#### Story 7.7: Port P25 Protocol
**Priority:** Medium
- Reference: `P25Control.cpp`, `P25Network.cpp`, `P25Data.cpp`, etc.
- **Task:** Implement P25 protocol handling.

```go
package p25

func HandleP25Packet(packet []byte) {
	// TODO: Add P25 packet handling logic
}
```

#### Story 7.8: Port POCSAG Protocol
**Priority:** Medium
- Reference: `POCSAGControl.cpp`, `POCSAGNetwork.cpp`, etc.
- **Task:** Implement POCSAG protocol handling.

```go
package pocsag

func HandlePOCSAGPacket(packet []byte) {
	// TODO: Add POCSAG packet handling logic
}
```

#### Story 7.9: Port Display Utilities
**Priority:** Medium
- Reference: `Display.cpp`, `Display.h`, `HD44780.cpp`, `OLED.cpp`, etc.
- **Task:** Implement display utilities.

```go
package display

func RenderDisplay(data string) {
	// TODO: Add display rendering logic
}
```

#### Story 7.10: Port Miscellaneous Utilities
**Priority:** Low
- Reference: `Utils.cpp`, `Utils.h`, `CRC.cpp`, `Golay2087.cpp`, etc.
- **Task:** Implement miscellaneous utilities.

```go
package utils

func PerformUtilityTask(taskID int, data []byte) {
	// TODO: Add utility task logic
}
```

### Chapter 8: Unassigned Features

#### Story 8.1: Port BPTC19696 Logic
**Priority:** Low
- Reference: `BPTC19696.cpp`, `BPTC19696.h`
- **Task:** Implement BPTC19696 error correction logic.

```go
package bptc

func CorrectBPTCData(data []byte) ([]byte, error) {
	// TODO: Add BPTC19696 error correction logic
	return nil, nil
}
```

#### Story 8.2: Port CASTInfo Logic
**Priority:** Low
- Reference: `CASTInfo.cpp`, `CASTInfo.h`
- **Task:** Implement CAST information handling.

```go
package castinfo

func HandleCASTInfo(data []byte) {
	// TODO: Add CAST information handling logic
}
```

#### Story 8.3: Port Defines
**Priority:** Low
- Reference: `Defines.h`
- **Task:** Consolidate and implement global definitions.

```go
package defines

const (
	// TODO: Add global definitions
)
```

#### Story 8.4: Port DMR Access Control
**Priority:** Low
- Reference: `DMRAccessControl.cpp`, `DMRAccessControl.h`
- **Task:** Implement DMR access control logic.

```go
package dmr

func HandleAccessControl(data []byte) {
	// TODO: Add DMR access control logic
}
```

#### Story 8.5: Port DMR CSBK Logic
**Priority:** Low
- Reference: `DMRCSBK.cpp`, `DMRCSBK.h`
- **Task:** Implement DMR CSBK handling.

```go
package dmr

func HandleCSBK(data []byte) {
	// TODO: Add DMR CSBK handling logic
}
```

#### Story 8.6: Port DMR Data Header
**Priority:** Low
- Reference: `DMRDataHeader.cpp`, `DMRDataHeader.h`
- **Task:** Implement DMR data header parsing.

```go
package dmr

func ParseDataHeader(data []byte) {
	// TODO: Add DMR data header parsing logic
}
```

#### Story 8.7: Port DMR Defines
**Priority:** Low
- Reference: `DMRDefines.h`
- **Task:** Consolidate and implement DMR-specific definitions.

```go
package dmr

const (
	// TODO: Add DMR-specific definitions
)
```

#### Story 8.8: Port DMR Direct Network
**Priority:** Low
- Reference: `DMRDirectNetwork.cpp`, `DMRDirectNetwork.h`
- **Task:** Implement DMR direct network communication.

```go
package dmr

func HandleDirectNetwork(data []byte) {
	// TODO: Add DMR direct network communication logic
}
```

#### Story 8.9: Port DMR Embedded Data
**Priority:** Low
- Reference: `DMREmbeddedData.cpp`, `DMREmbeddedData.h`
- **Task:** Implement DMR embedded data handling.

```go
package dmr

func HandleEmbeddedData(data []byte) {
	// TODO: Add DMR embedded data handling logic
}
```

#### Story 8.10: Port DMR Full LC
**Priority:** Low
- Reference: `DMRFullLC.cpp`, `DMRFullLC.h`
- **Task:** Implement DMR full LC handling.

```go
package dmr

func HandleFullLC(data []byte) {
	// TODO: Add DMR full LC handling logic
}
```

#### Story 8.11: Port DMR Gateway Network
**Priority:** Low
- Reference: `DMRGatewayNetwork.cpp`, `DMRGatewayNetwork.h`
- **Task:** Implement DMR gateway network communication.

```go
package dmr

func HandleGatewayNetwork(data []byte) {
	// TODO: Add DMR gateway network communication logic
}
```

#### Story 8.12: Port DMR Ids
**Priority:** Low
- Reference: `DMRIds.dat`
- **Task:** Parse and manage DMR IDs.

```go
package dmr

func ParseDMRIds(data []byte) {
	// TODO: Add DMR ID parsing logic
}
```

#### Story 8.13: Port DMR LC
**Priority:** Low
- Reference: `DMRLC.cpp`, `DMRLC.h`
- **Task:** Implement DMR LC handling.

```go
package dmr

func HandleLC(data []byte) {
	// TODO: Add DMR LC handling logic
}
```

#### Story 8.14: Port DMR Lookup
**Priority:** Low
- Reference: `DMRLookup.cpp`, `DMRLookup.h`
- **Task:** Implement DMR lookup functionality.

```go
package dmr

func LookupDMRData(data []byte) {
	// TODO: Add DMR lookup functionality
}
```

#### Story 8.15: Port DMR Plus Startup Options
**Priority:** Low
- Reference: `DMRplus_startup_options.md`
- **Task:** Document DMR Plus startup options.

```markdown
# DMR Plus Startup Options

- TODO: Add documentation for DMR Plus startup options.
```

#### Story 8.16: Port DMR Short LC
**Priority:** Low
- Reference: `DMRShortLC.cpp`, `DMRShortLC.h`
- **Task:** Implement DMR short LC handling.

```go
package dmr

func HandleShortLC(data []byte) {
	// TODO: Add DMR short LC handling logic
}
```

#### Story 8.17: Port DMR TA
**Priority:** Low
- Reference: `DMRTA.cpp`, `DMRTA.h`
- **Task:** Implement DMR TA handling.

```go
package dmr

func HandleTA(data []byte) {
	// TODO: Add DMR TA handling logic
}
```

#### Story 8.18: Port DMR Trellis
**Priority:** Low
- Reference: `DMRTrellis.cpp`, `DMRTrellis.h`
- **Task:** Implement DMR trellis encoding/decoding.

```go
package dmr

func HandleTrellis(data []byte) {
	// TODO: Add DMR trellis encoding/decoding logic
}
```

#### Story 8.19: Port Dockerfile
**Priority:** Low
- Reference: `Dockerfile`
- **Task:** Create a Docker container for the application.

```dockerfile
# TODO: Add Dockerfile for containerization
```

#### Story 8.20: Port DStar Defines
**Priority:** Low
- Reference: `DStarDefines.h`
- **Task:** Consolidate and implement D-Star-specific definitions.

```go
package dstar

const (
	// TODO: Add D-Star-specific definitions
)
```

#### Story 8.21: Port DStar Slow Data
**Priority:** Low
- Reference: `DStarSlowData.cpp`, `DStarSlowData.h`
- **Task:** Implement D-Star slow data handling.

```go
package dstar

func HandleSlowData(data []byte) {
	// TODO: Add D-Star slow data handling logic
}
```

#### Story 8.22: Port FM Control
**Priority:** Low
- Reference: `FMControl.cpp`, `FMControl.h`
- **Task:** Implement FM control logic.

```go
package fm

func HandleFMControl(data []byte) {
	// TODO: Add FM control logic
}
```

#### Story 8.23: Port FM Network
**Priority:** Low
- Reference: `FMNetwork.cpp`, `FMNetwork.h`
- **Task:** Implement FM network communication.

```go
package fm

func HandleFMNetwork(data []byte) {
	// TODO: Add FM network communication logic
}
```

#### Story 8.24: Port Golay24128
**Priority:** Low
- Reference: `Golay24128.cpp`, `Golay24128.h`
- **Task:** Implement Golay24128 error correction.

```go
package golay

func CorrectGolay24128(data []byte) ([]byte, error) {
	// TODO: Add Golay24128 error correction logic
	return nil, nil
}
```

#### Story 8.25: Port HD44780 Layouts
**Priority:** Low
- Reference: `HD44780.layouts`
- **Task:** Implement HD44780 layout configurations.

```go
package hd44780

func ConfigureLayout(layout string) {
	// TODO: Add HD44780 layout configuration logic
}
```

### File-to-Story Mapping

```
.
├── AMBEFEC.cpp (Chapter 7: Story 7.1)
├── AMBEFEC.h (Chapter 7: Story 7.1)
├── AX25Control.cpp (Chapter 7: Story 7.2)
├── AX25Control.h (Chapter 7: Story 7.2)
├── AX25Defines.h (Chapter 7: Story 7.2)
├── AX25Network.cpp (Chapter 7: Story 7.2)
├── AX25Network.h (Chapter 7: Story 7.2)
├── BCH.cpp (Chapter 7: Story 7.3)
├── BCH.h (Chapter 7: Story 7.3)
├── CASTInfo.cpp (Chapter 8: Story 8.2)
├── CASTInfo.h (Chapter 8: Story 8.2)
├── Conf.cpp (Chapter 1: Story 1.2)
├── Conf.h (Chapter 1: Story 1.2)
├── CRC.cpp (Chapter 3: Story 3.2)
├── Dockerfile (Chapter 8: Story 8.19)
├── DMRAccessControl.cpp (Chapter 8: Story 8.4)
├── DMRAccessControl.h (Chapter 8: Story 8.4)
├── DMRCSBK.cpp (Chapter 8: Story 8.5)
├── DMRCSBK.h (Chapter 8: Story 8.5)
├── DMRControl.cpp (Chapter 2: Story 2.1)
├── DMRDataHeader.cpp (Chapter 8: Story 8.6)
├── DMRDataHeader.h (Chapter 8: Story 8.6)
├── DMRDefines.h (Chapter 8: Story 8.7)
├── DMRDirectNetwork.cpp (Chapter 8: Story 8.8)
├── DMRDirectNetwork.h (Chapter 8: Story 8.8)
├── DMREmbeddedData.cpp (Chapter 8: Story 8.9)
├── DMREmbeddedData.h (Chapter 8: Story 8.9)
├── DMRFullLC.cpp (Chapter 8: Story 8.10)
├── DMRFullLC.h (Chapter 8: Story 8.10)
├── DMRGatewayNetwork.cpp (Chapter 8: Story 8.11)
├── DMRGatewayNetwork.h (Chapter 8: Story 8.11)
├── DMRIds.dat (Chapter 8: Story 8.12)
├── DMRLC.cpp (Chapter 8: Story 8.13)
├── DMRLC.h (Chapter 8: Story 8.13)
├── DMRLookup.cpp (Chapter 8: Story 8.14)
├── DMRLookup.h (Chapter 8: Story 8.14)
├── DMRplus_startup_options.md (Chapter 8: Story 8.15)
├── DMRSlot.cpp (Chapter 7: Story 7.4)
├── DMRSlot.h (Chapter 7: Story 7.4)
├── DMRSlotType.cpp (Chapter 7: Story 7.4)
├── DMRSlotType.h (Chapter 7: Story 7.4)
├── DMRTA.cpp (Chapter 8: Story 8.17)
├── DMRTA.h (Chapter 8: Story 8.17)
├── DMRTrellis.cpp (Chapter 8: Story 8.18)
├── DMRTrellis.h (Chapter 8: Story 8.18)
├── DStarControl.cpp (Chapter 2: Story 2.2)
├── DStarDefines.h (Chapter 8: Story 8.20)
├── DStarHeader.cpp (Chapter 2: Story 2.2)
├── DStarNetwork.cpp (Chapter 2: Story 2.2)
├── DStarSlowData.cpp (Chapter 8: Story 8.21)
├── DStarSlowData.h (Chapter 8: Story 8.21)
├── Golay2087.cpp (Chapter 3: Story 3.2)
├── Golay2087.h (Chapter 3: Story 3.2)
├── Golay24128.cpp (Chapter 8: Story 8.24)
├── Golay24128.h (Chapter 8: Story 8.24)
├── HD44780.cpp (Chapter 4: Story 4.1)
├── HD44780.h (Chapter 4: Story 4.1)
├── HD44780.layouts (Chapter 8: Story 8.25)
├── OLED.cpp (Chapter 4: Story 4.1)
├── P25Control.cpp (Chapter 7: Story 7.7)
├── P25Network.cpp (Chapter 7: Story 7.7)
├── P25Data.cpp (Chapter 7: Story 7.7)
├── FMControl.cpp (Chapter 8: Story 8.22)
├── FMNetwork.cpp (Chapter 8: Story 8.23)
├── Log.cpp (Chapter 3: Story 3.1)
├── Log.h (Chapter 3: Story 3.1)
├── Utils.cpp (Chapter 8: Story 8.10)
├── Utils.h (Chapter 8: Story 8.10)
├── YSFControl.cpp (Chapter 7: Story 7.6)
├── YSFNetwork.cpp (Chapter 7: Story 7.6)
└── YSFFICH.cpp (Chapter 7: Story 7.6)
```
