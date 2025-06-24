package config

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitializeDatabase creates the database schema and populates default values using GORM.
func InitializeDatabase(dbPath string) error {
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_synchronous=NORMAL", dbPath)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	tables := []interface{}{
		&GeneralConfig{},
		&DMRConfig{},
		&DStarConfig{},
		&M17Config{},
		&AX25Config{},
		&NXDNConfig{},
		&PocsagConfig{},
		&YSFConfig{},
	}

	// Drop the GeneralConfig table if it exists to ensure schema consistency
	if err := db.Migrator().DropTable(&GeneralConfig{}); err != nil {
		return fmt.Errorf("failed to drop GeneralConfig table: %w", err)
	}

	for _, table := range tables {
		fmt.Printf("Migrating table: %T\n", table) // Log table being migrated
		if err := db.AutoMigrate(table); err != nil {
			return fmt.Errorf("failed to migrate table %T: %w", table, err)
		}
	}

	// Insert default values
	fmt.Println("Inserting default values...") // Log default value insertion
	defaults := map[string]interface{}{
		"GeneralConfig": GeneralConfig{Callsign: "NOCALL", Timeout: 60, Duplex: false},
		"DMRConfig":     DMRConfig{Enable: true, ColorCode: 1},
		"DStarConfig":   DStarConfig{Enable: true, Module: "C"},
		"M17Config":     M17Config{Enable: true, CAN: "A"},
		"AX25Config":    AX25Config{Enable: false, Port: ""},
		"NXDNConfig":    NXDNConfig{Enable: false, Port: ""},
		"PocsagConfig":  PocsagConfig{Enable: false, Frequency: 0},
		"YSFConfig":     YSFConfig{Enable: true, Port: ""},
	}

	for tableName, defaultValue := range defaults {
		fmt.Printf("Inserting default for table: %s\n", tableName) // Log each table default
		if err := db.Table(tableName).Create(defaultValue).Error; err != nil {
			return fmt.Errorf("failed to insert default value for table %s: %w", tableName, err)
		}
	}

	return nil
}
