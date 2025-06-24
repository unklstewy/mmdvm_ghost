package config

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Config holds settings from the SQLite database.
type Config struct {
	General   GeneralConfig
	Log       LogConfig
	Modem     ModemConfig
	DMR       DMRConfig
	DStar     DStarConfig
	M17       M17Config
	Network   NetworkConfig
	Display   DisplayConfig
	FilePaths FilePaths
	AX25      AX25Config
	NXDN      NXDNConfig
	Pocsag    PocsagConfig
	YSF       YSFConfig
}

// GeneralConfig stores general configuration parameters
// Add GORM tags for table and column mapping
type GeneralConfig struct {
	Callsign          string `gorm:"column:callsign"`
	Timeout           int    `gorm:"column:timeout"`
	Duplex            bool   `gorm:"column:duplex"`
	RFModeHang        int    `gorm:"column:rf_mode_hang"`
	NetModeHang       int    `gorm:"column:net_mode_hang"`
	DisplayLevel      int    `gorm:"column:display_level"`
	DisplayMode       string `gorm:"column:display_mode"`
	DisplayBrightness int    `gorm:"column:display_brightness"`
	DisplayInvert     bool   `gorm:"column:display_invert"`
}

// LogConfig stores logging configuration
type LogConfig struct {
	LogPath    string
	LogLevel   int
	DisplayLog bool
}

// ModemConfig stores modem connection settings
type ModemConfig struct {
	Port            string
	Protocol        string
	TXDelay         int
	RXLevel         int
	TXLevel         int
	DMRDelay        int
	RXOffset        int
	TXOffset        int
	RSSIMappingFile string
}

// DMRConfig stores DMR protocol configuration
// Add GORM tags for table and column mapping
type DMRConfig struct {
	Enable         bool `gorm:"column:enable"`
	Beacons        bool `gorm:"column:beacons"`
	ColorCode      int  `gorm:"column:color_code"`
	SelfOnly       bool `gorm:"column:self_only"`
	EmbeddedLCOnly bool `gorm:"column:embedded_lc_only"`
	DumpTAData     bool `gorm:"column:dump_ta_data"`
}

// DStarConfig stores D-Star protocol configuration
// Add GORM tags for table and column mapping
type DStarConfig struct {
	Enable bool   `gorm:"column:enable"`
	Module string `gorm:"column:module"`
}

// M17Config stores M17 protocol configuration
// Add GORM tags for table and column mapping
type M17Config struct {
	Enable bool   `gorm:"column:enable"`
	CAN    string `gorm:"column:can"`
}

// NetworkConfig stores network connection parameters
type NetworkConfig struct {
	Enable        bool
	Port          int
	HostsFile     string
	ReloadTime    int
	ParrotAddress string
	ParrotPort    int
	Startup       string
}

// DisplayConfig stores display device parameters
type DisplayConfig struct {
	Type       string
	Port       string
	Brightness int
}

// FilePaths stores paths to various auxiliary files
type FilePaths struct {
	DMRID     string
	NXDNID    string
	WhiteList string
	BlackList string
}

// AX25Config stores AX.25 protocol configuration
// Add GORM tags for table and column mapping
type AX25Config struct {
	Enable bool   `gorm:"column:enable"`
	Port   string `gorm:"column:port"`
}

// NXDNConfig stores NXDN protocol configuration
// Add GORM tags for table and column mapping
type NXDNConfig struct {
	Enable bool   `gorm:"column:enable"`
	Port   string `gorm:"column:port"`
}

// PocsagConfig stores POCSAG protocol configuration
// Add GORM tags for table and column mapping
type PocsagConfig struct {
	Enable    bool `gorm:"column:enable"`
	Frequency int  `gorm:"column:frequency"`
}

// YSFConfig stores YSF protocol configuration
// Add GORM tags for table and column mapping
type YSFConfig struct {
	Enable bool   `gorm:"column:enable"`
	Port   string `gorm:"column:port"`
}

// LoadConfig loads configuration from the specified SQLite database file.
func LoadConfig(path string) (*Config, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	config := &Config{}

	// Load GeneralConfig
	generalConfig, err := loadGeneralConfig(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load general config: %w", err)
	}
	config.General = generalConfig

	// Load LogConfig
	if err := loadLogConfig(db, &config.Log); err != nil {
		return nil, fmt.Errorf("failed to load log config: %w", err)
	}

	// Load ModemConfig
	if err := loadModemConfig(db, &config.Modem); err != nil {
		return nil, fmt.Errorf("failed to load modem config: %w", err)
	}

	// Load DMRConfig
	if err := loadDMRConfig(db, &config.DMR); err != nil {
		return nil, fmt.Errorf("failed to load DMR config: %w", err)
	}

	// Load DStarConfig
	if err := loadDStarConfig(db, &config.DStar); err != nil {
		return nil, fmt.Errorf("failed to load D-Star config: %w", err)
	}

	// Load M17Config
	if err := loadM17Config(db, &config.M17); err != nil {
		return nil, fmt.Errorf("failed to load M17 config: %w", err)
	}

	// Load NetworkConfig
	if err := loadNetworkConfig(db, &config.Network); err != nil {
		return nil, fmt.Errorf("failed to load network config: %w", err)
	}

	// Load DisplayConfig
	if err := loadDisplayConfig(db, &config.Display); err != nil {
		return nil, fmt.Errorf("failed to load display config: %w", err)
	}

	// Load FilePaths
	if err := loadFilePaths(db, &config.FilePaths); err != nil {
		return nil, fmt.Errorf("failed to load file paths: %w", err)
	}

	// Load AX25Config
	if err := loadAX25Config(db, &config.AX25); err != nil {
		return nil, fmt.Errorf("failed to load AX.25 config: %w", err)
	}

	// Load NXDNConfig
	if err := loadNXDNConfig(db, &config.NXDN); err != nil {
		return nil, fmt.Errorf("failed to load NXDN config: %w", err)
	}

	// Load PocsagConfig
	if err := loadPocsagConfig(db, &config.Pocsag); err != nil {
		return nil, fmt.Errorf("failed to load POCSAG config: %w", err)
	}

	// Load YSFConfig
	if err := loadYSFConfig(db, &config.YSF); err != nil {
		return nil, fmt.Errorf("failed to load YSF config: %w", err)
	}

	return config, nil
}

// createDefaultDatabase creates a new SQLite database with default configuration values.
func createDefaultDatabase(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	defer db.Close()

	queries := []string{
		`CREATE TABLE General (Callsign TEXT, Timeout INTEGER, Duplex BOOLEAN, RFModeHang INTEGER, NetModeHang INTEGER, DisplayLevel INTEGER, DisplayMode TEXT, DisplayBrightness INTEGER, DisplayInvert BOOLEAN);`,
		`INSERT INTO General VALUES ('NOCALL', 60, 0, 10, 10, 1, 'Normal', 100, 0);`,
		`CREATE TABLE Log (LogPath TEXT, LogLevel INTEGER, DisplayLog BOOLEAN);`,
		`INSERT INTO Log VALUES ('log/mmdvm_ghost.log', 1, 1);`,
		// Add more table creation and default insertion queries here...
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

// loadGeneralConfig loads the GeneralConfig section from the database.
func loadGeneralConfig(db *sql.DB) (GeneralConfig, error) {
	query := `SELECT callsign, timeout, duplex, rf_mode_hang, net_mode_hang, display_level, display_mode, display_brightness, display_invert FROM GeneralConfig LIMIT 1`
	fmt.Printf("Executing query: %s\n", query) // Log the query
	row := db.QueryRow(query)

	var general GeneralConfig
	err := row.Scan(
		&general.Callsign,
		&general.Timeout,
		&general.Duplex,
		&general.RFModeHang,
		&general.NetModeHang,
		&general.DisplayLevel,
		&general.DisplayMode,
		&general.DisplayBrightness,
		&general.DisplayInvert,
	)
	if err != nil {
		return GeneralConfig{}, fmt.Errorf("failed to scan general config: %w", err)
	}

	return general, nil
}

// loadLogConfig loads the Log configuration section from the database.
func loadLogConfig(db *sql.DB, log *LogConfig) error {
	row := db.QueryRow(`SELECT LogPath, LogLevel, DisplayLog FROM Log`)
	return row.Scan(&log.LogPath, &log.LogLevel, &log.DisplayLog)
}

// loadModemConfig loads the Modem configuration section from the database.
func loadModemConfig(db *sql.DB, modem *ModemConfig) error {
	row := db.QueryRow(`SELECT Port, Protocol, TXDelay, RXLevel, TXLevel, DMRDelay, RXOffset, TXOffset, RSSIMappingFile FROM Modem`)
	return row.Scan(&modem.Port, &modem.Protocol, &modem.TXDelay, &modem.RXLevel, &modem.TXLevel, &modem.DMRDelay, &modem.RXOffset, &modem.TXOffset, &modem.RSSIMappingFile)
}

// loadDMRConfig loads the DMR configuration section from the database.
func loadDMRConfig(db *sql.DB, dmr *DMRConfig) error {
	row := db.QueryRow(`SELECT Enable, Beacons, ColorCode, SelfOnly, EmbeddedLCOnly, DumpTAData FROM DMR`)
	return row.Scan(&dmr.Enable, &dmr.Beacons, &dmr.ColorCode, &dmr.SelfOnly, &dmr.EmbeddedLCOnly, &dmr.DumpTAData)
}

// loadDStarConfig loads the D-Star configuration section from the database.
func loadDStarConfig(db *sql.DB, dstar *DStarConfig) error {
	row := db.QueryRow(`SELECT Enable, Module FROM DStar`)
	return row.Scan(&dstar.Enable, &dstar.Module)
}

// loadM17Config loads the M17 configuration section from the database.
func loadM17Config(db *sql.DB, m17 *M17Config) error {
	row := db.QueryRow(`SELECT Enable, CAN FROM M17`)
	return row.Scan(&m17.Enable, &m17.CAN)
}

// loadNetworkConfig loads the Network configuration section from the database.
func loadNetworkConfig(db *sql.DB, network *NetworkConfig) error {
	row := db.QueryRow(`SELECT Enable, Port, HostsFile, ReloadTime, ParrotAddress, ParrotPort, Startup FROM Network`)
	return row.Scan(&network.Enable, &network.Port, &network.HostsFile, &network.ReloadTime, &network.ParrotAddress, &network.ParrotPort, &network.Startup)
}

// loadDisplayConfig loads the Display configuration section from the database.
func loadDisplayConfig(db *sql.DB, display *DisplayConfig) error {
	row := db.QueryRow(`SELECT Type, Port, Brightness FROM Display`)
	return row.Scan(&display.Type, &display.Port, &display.Brightness)
}

// loadFilePaths loads the FilePaths configuration section from the database.
func loadFilePaths(db *sql.DB, paths *FilePaths) error {
	row := db.QueryRow(`SELECT DMRID, NXDNID, WhiteList, BlackList FROM FilePaths`)
	return row.Scan(&paths.DMRID, &paths.NXDNID, &paths.WhiteList, &paths.BlackList)
}

// loadAX25Config loads the AX.25 configuration section from the database.
func loadAX25Config(db *sql.DB, ax25 *AX25Config) error {
	row := db.QueryRow(`SELECT Enable, Port FROM AX25`)
	return row.Scan(&ax25.Enable, &ax25.Port)
}

// loadNXDNConfig loads the NXDN configuration section from the database.
func loadNXDNConfig(db *sql.DB, nxdn *NXDNConfig) error {
	row := db.QueryRow(`SELECT Enable, Port FROM NXDN`)
	return row.Scan(&nxdn.Enable, &nxdn.Port)
}

// loadPocsagConfig loads the POCSAG configuration section from the database.
func loadPocsagConfig(db *sql.DB, pocsag *PocsagConfig) error {
	row := db.QueryRow(`SELECT Enable, Frequency FROM Pocsag`)
	return row.Scan(&pocsag.Enable, &pocsag.Frequency)
}

// loadYSFConfig loads the YSF configuration section from the database.
func loadYSFConfig(db *sql.DB, ysf *YSFConfig) error {
	row := db.QueryRow(`SELECT Enable, Port FROM YSF`)
	return row.Scan(&ysf.Enable, &ysf.Port)
}

func (GeneralConfig) TableName() string {
	return "GeneralConfig"
}

func (DMRConfig) TableName() string {
	return "DMRConfig"
}

func (DStarConfig) TableName() string {
	return "DStarConfig"
}

func (M17Config) TableName() string {
	return "M17Config"
}

func (AX25Config) TableName() string {
	return "AX25Config"
}

func (NXDNConfig) TableName() string {
	return "NXDNConfig"
}

func (PocsagConfig) TableName() string {
	return "PocsagConfig"
}

func (YSFConfig) TableName() string {
	return "YSFConfig"
}

// Ensure all required structs are present
// GeneralConfig, DMRConfig, DStarConfig, M17Config, AX25Config, NXDNConfig, PocsagConfig, YSFConfig are already defined.
// No additional structs are missing.
