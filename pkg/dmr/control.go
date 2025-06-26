// Package dmr provides DMR protocol logic, including control structures and wakeup packet processing.
package dmr

import (
	"errors" // For error handling
	"log"    // For logging debug/info messages
)

// Control manages DMR protocol control, including slots and lookup.
type Control struct {
	ColorCode uint32 // DMR color code
	Slot1     *Slot  // Slot 1 handler
	Slot2     *Slot  // Slot 2 handler
	Lookup    Lookup // Lookup interface for ID resolution
}

// NewControl initializes a new Control instance with slots and lookup.
func NewControl(colorCode uint32, timeout uint32, lookup Lookup, accessControl *AccessControl) *Control {
	if accessControl == nil || lookup == nil {
		log.Fatal("AccessControl and Lookup cannot be nil")
	}

	return &Control{
		ColorCode: colorCode,
		Slot1:     NewSlot(1, timeout),
		Slot2:     NewSlot(2, timeout),
		Lookup:    lookup,
	}
}

// ProcessWakeup processes a wakeup packet, validates it, and extracts CSBK info.
func (c *Control) ProcessWakeup(data []byte) error {
	if len(data) < 2 {
		return errors.New("data too short")
	}

	log.Printf("ProcessWakeup: Received data: %v", data)

	// Ensure data length is exactly 33 bytes
	if len(data) != 33 {
		return errors.New("invalid data length, expected 33 bytes")
	}

	// Validate wakeup packet header
	if data[0] != TagData || data[1] != (DmrIdleRx|DmrSyncData|DtCSBK) {
		return errors.New("invalid wakeup packet")
	}

	log.Printf("ProcessWakeup: Passing data to CSBK.Put: %v", data[2:])

	csbk := &CSBK{}
	if err := csbk.Put(data[2:]); err != nil {
		return err
	}

	if csbk.GetCSBKO() != BSDWNACT {
		return errors.New("invalid CSBKO")
	}

	srcID := csbk.GetSrcID()
	src := c.Lookup.Find(srcID)
	if !accessControl.ValidateSrcID(srcID) {
		log.Printf("Invalid Downlink Activate received from %s", src)
		return errors.New("access denied")
	}

	log.Printf("Downlink Activate received from %s", src)
	return nil
}

// Slot represents a DMR slot (time slot for communication).
type Slot struct {
	SlotNumber uint32 // Slot number (1 or 2)
	Timeout    uint32 // Timeout value
}

// NewSlot creates a new Slot instance with the given slot number and timeout.
func NewSlot(slotNumber uint32, timeout uint32) *Slot {
	return &Slot{
		SlotNumber: slotNumber,
		Timeout:    timeout,
	}
}

// Lookup defines an interface for looking up IDs (e.g., user or repeater IDs).
type Lookup interface {
	Find(id uint32) string
}

// Define and initialize accessControl as a global variable.
var accessControl = &AccessControl{}

func init() {
	// Initialize accessControl with default values or load from database.
	accessControl.Init(nil, nil, nil, nil, nil, false, 0)
}

// Expose ProcessWakeup as a package-level function for convenience.
func ProcessWakeup(data []byte) error {
	control := &Control{}
	return control.ProcessWakeup(data)
}
