package dmr

import (
	"errors"
	"log"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// AccessControlData represents the structure for access control validation.
type AccessControlData struct {
	ID   string
	Role string
}

// AccessControl manages DMR access control logic.
type AccessControl struct {
	BlackList        []uint32
	WhiteList        []uint32
	Prefixes         []uint32
	Slot1TGWhiteList []uint32
	Slot2TGWhiteList []uint32
	SelfOnly         bool
	ID               uint32
}

// AccessRule represents the database model for access control rules.
type AccessRule struct {
	ID   string `gorm:"primaryKey"`
	Role string
}

var dbInstance *gorm.DB

// Init initializes the access control settings.
func (ac *AccessControl) Init(blacklist, whitelist, slot1TGWhitelist, slot2TGWhitelist, prefixes []uint32, selfOnly bool, id uint32) {
	ac.BlackList = blacklist
	ac.WhiteList = whitelist
	ac.Prefixes = prefixes
	ac.Slot1TGWhiteList = slot1TGWhitelist
	ac.Slot2TGWhiteList = slot2TGWhitelist
	ac.SelfOnly = selfOnly
	ac.ID = id
}

// HandleAccessControl validates and manages access based on provided data.
func HandleAccessControl(data []byte) error {
	// Parse the data into AccessControlData structure
	accessData, err := parseAccessControlData(data)
	if err != nil {
		log.Printf("Error parsing access control data: %v", err)
		return err
	}

	// Add debug logs to trace parsed data in HandleAccessControl
	log.Printf("HandleAccessControl: Parsed ID=%s, Role=%s", accessData.ID, accessData.Role)

	// Validate the access control data
	if !isValidAccess(accessData) {
		log.Printf("Access denied for ID: %s", accessData.ID)
		return errors.New("access denied")
	}

	log.Printf("Access granted for ID: %s with role: %s", accessData.ID, accessData.Role)
	return nil
}

// parseAccessControlData parses raw data into AccessControlData structure.
func parseAccessControlData(data []byte) (*AccessControlData, error) {
	parts := strings.Split(string(data), ",")
	if len(parts) != 2 {
		return nil, errors.New("invalid access control data format")
	}

	// Add debug logs to trace parsed values
	log.Printf("parseAccessControlData: Parsed ID: %s, Role: %s", parts[0], parts[1])
	return &AccessControlData{ID: parts[0], Role: parts[1]}, nil
}

// isValidAccess checks if the access control data is valid.
func isValidAccess(data *AccessControlData) bool {
	db := getDatabaseConnection() // Assume this function returns a *gorm.DB instance

	// Add debug logs to trace input data
	log.Printf("isValidAccess: Received ID: %s, Role: %s", data.ID, data.Role)

	// Add debug logs to trace database query
	log.Printf("isValidAccess: Querying database for ID: %s, Role: %s", data.ID, data.Role)
	// Add debug log to trace SQL query
	log.Printf("isValidAccess: Executing query: SELECT * FROM access_rules WHERE id = '%s' AND role = '%s'", data.ID, data.Role)

	// Use raw SQL query to validate access
	var count int64
	query := "SELECT COUNT(*) FROM access_rules WHERE id = ? AND role = ?"
	if err := db.Raw(query, data.ID, data.Role).Scan(&count).Error; err != nil {
		log.Printf("Database error: %v", err)
		return false
	}
	// Add debug log to trace query result
	log.Printf("isValidAccess: Query result count: %d", count)
	if count == 0 {
		log.Printf("Access denied: No matching rule for ID: %s, Role: %s", data.ID, data.Role)
		return false
	}
	log.Printf("Access granted for ID: %s with role: %s", data.ID, data.Role)
	return true
}

// ValidateSrcID validates a source ID against the access control rules.
func (ac *AccessControl) ValidateSrcID(id uint32) bool {
	if ac.SelfOnly {
		switch {
		case ac.ID > 99999999:
			return id == ac.ID/100
		case ac.ID > 9999999:
			return id == ac.ID/10
		default:
			return id == ac.ID
		}
	}

	for _, blacklisted := range ac.BlackList {
		if id == blacklisted {
			return false
		}
	}

	prefix := id / 10000
	if prefix == 0 || prefix > 999 {
		return false
	}

	if len(ac.Prefixes) > 0 {
		found := false
		for _, p := range ac.Prefixes {
			if prefix == p {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(ac.WhiteList) > 0 {
		for _, whitelisted := range ac.WhiteList {
			if id == whitelisted {
				return true
			}
		}
		return false
	}

	return true
}

// ValidateTGID validates a talk group ID for a given slot.
func (ac *AccessControl) ValidateTGID(slotNo uint32, group bool, id uint32) bool {
	if !group {
		return true
	}

	if id == 0 {
		return false
	}

	if slotNo == 1 {
		if len(ac.Slot1TGWhiteList) == 0 {
			return true
		}
		for _, whitelisted := range ac.Slot1TGWhiteList {
			if id == whitelisted {
				return true
			}
		}
		return false
	}

	if len(ac.Slot2TGWhiteList) == 0 {
		return true
	}
	for _, whitelisted := range ac.Slot2TGWhiteList {
		if id == whitelisted {
			return true
		}
	}
	return false
}

// getDatabaseConnection initializes and returns a GORM database connection.
func getDatabaseConnection() *gorm.DB {
	if dbInstance == nil {
		var err error
		dbInstance, err = gorm.Open(sqlite.Open("mmdvm_ghost.db"), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Add debug log to confirm database file path
		log.Printf("getDatabaseConnection: Using database file: mmdvm_ghost.db")
	}
	return dbInstance
}

// UpdateBlackList dynamically updates the blacklist.
func (ac *AccessControl) UpdateBlackList(newBlackList []uint32) {
	ac.BlackList = newBlackList
	log.Printf("Blacklist updated: %v", ac.BlackList)
}

// UpdateWhiteList dynamically updates the whitelist.
func (ac *AccessControl) UpdateWhiteList(newWhiteList []uint32) {
	ac.WhiteList = newWhiteList
	log.Printf("Whitelist updated: %v", ac.WhiteList)
}

// ReloadRules reloads access control rules from the database.
func (ac *AccessControl) ReloadRules() error {
	var rules []AccessRule
	db := getDatabaseConnection()
	if err := db.Find(&rules).Error; err != nil {
		log.Printf("Failed to reload rules from database: %v", err)
		return err
	}

	// Example: Update internal structures based on rules
	for _, rule := range rules {
		log.Printf("Loaded rule: ID=%s, Role=%s", rule.ID, rule.Role)
	}

	return nil
}
