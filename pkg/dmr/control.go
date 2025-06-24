package dmr

import (
	"errors"
	"log"
)

// Control manages DMR protocol control.
type Control struct {
	ColorCode uint32
	Slot1     *Slot
	Slot2     *Slot
	Lookup    Lookup
}

// NewControl initializes a new Control instance.
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

// ProcessWakeup processes a wakeup packet.
func (c *Control) ProcessWakeup(data []byte) error {
	if len(data) < 2 {
		return errors.New("data too short")
	}

	log.Printf("ProcessWakeup: Received data: %v", data)

	// Validate wakeup packet
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

// Slot represents a DMR slot.
type Slot struct {
	SlotNumber uint32
	Timeout    uint32
}

// NewSlot creates a new Slot instance.
func NewSlot(slotNumber uint32, timeout uint32) *Slot {
	return &Slot{
		SlotNumber: slotNumber,
		Timeout:    timeout,
	}
}

// Lookup defines an interface for looking up IDs.
type Lookup interface {
	Find(id uint32) string
}

// Define and initialize accessControl as a global variable.
var accessControl = &AccessControl{}

func init() {
	// Initialize accessControl with default values or load from database.
	accessControl.Init(nil, nil, nil, nil, nil, false, 0)
}
