package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	// General settings
	ConfigPath string `mapstructure:"config_path"`
	LogLevel   string `mapstructure:"log_level"`
	Network    string `mapstructure:"network"`

	// RPC settings
	RPC RPCConfig `mapstructure:"rpc"`

	// Keystore settings
	Keystore KeystoreConfig `mapstructure:"keystore"`

	// Node settings
	Node NodeConfig `mapstructure:"node"`

	// Validator settings
	Validator ValidatorConfig `mapstructure:"validator"`

	// API settings
	API APIConfig `mapstructure:"api"`

	// Logging settings
	Logging LoggingConfig `mapstructure:"logging"`
}

// RPCConfig contains RPC configuration
type RPCConfig struct {
	URL             string `mapstructure:"url"`
	Timeout         int    `mapstructure:"timeout"`
	MaxConnections  int    `mapstructure:"max_connections"`
	EnableWebSocket bool   `mapstructure:"enable_websocket"`
}

// KeystoreConfig contains keystore configuration
type KeystoreConfig struct {
	Path       string `mapstructure:"path"`
	Password   string `mapstructure:"password"`
	Encryption string `mapstructure:"encryption"`
}

// NodeConfig contains node configuration
type NodeConfig struct {
	DataDir     string `mapstructure:"data_dir"`
	GenesisFile string `mapstructure:"genesis_file"`
	HTTPPort    int    `mapstructure:"http_port"`
	WSPort      int    `mapstructure:"ws_port"`
	P2PPort     int    `mapstructure:"p2p_port"`
}

// ValidatorConfig contains validator configuration
type ValidatorConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	StakeAmount  string `mapstructure:"stake_amount"`
	Commission   string `mapstructure:"commission"`
	ValidatorKey string `mapstructure:"validator_key"`
	MinGasPrice  string `mapstructure:"min_gas_price"`
}

// APIConfig contains API configuration
type APIConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    int    `mapstructure:"port"`
	Host    string `mapstructure:"host"`
	CORS    string `mapstructure:"cors"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

// LoadConfig loads configuration from file and environment
func LoadConfig() (*Config, error) {
	// Set default values
	setDefaults()

	// Read configuration file
	configPath := getConfigPath()
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("diora")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.diora")
		viper.AddConfigPath("/etc/diora")
	}

	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvPrefix("DIORA")

	// Read configuration
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, use defaults
			fmt.Printf("Config file not found, using defaults\n")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal configuration
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Set config path
	config.ConfigPath = getConfigPath()

	return &config, nil
}

// Reload reloads configuration
func (c *Config) Reload() error {
	newConfig, err := LoadConfig()
	if err != nil {
		return err
	}
	*c = *newConfig
	return nil
}

// Save saves configuration to file
func (c *Config) Save() error {
	configPath := c.ConfigPath
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write configuration
	viper.Set("config_path", configPath)
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// General defaults
	viper.SetDefault("log_level", "info")
	viper.SetDefault("network", "mainnet")

	// RPC defaults
	viper.SetDefault("rpc.url", "http://localhost:8545")
	viper.SetDefault("rpc.timeout", 30)
	viper.SetDefault("rpc.max_connections", 100)
	viper.SetDefault("rpc.enable_websocket", true)

	// Keystore defaults
	viper.SetDefault("keystore.path", "$HOME/.diora/keystore")
	viper.SetDefault("keystore.encryption", "aes256")

	// Node defaults
	viper.SetDefault("node.data_dir", "$HOME/.diora/data")
	viper.SetDefault("node.genesis_file", "genesis.json")
	viper.SetDefault("node.http_port", 8545)
	viper.SetDefault("node.ws_port", 8546)
	viper.SetDefault("node.p2p_port", 30303)

	// Validator defaults
	viper.SetDefault("validator.enabled", false)
	viper.SetDefault("validator.stake_amount", "1000000000000000000") // 1 DIO
	viper.SetDefault("validator.commission", "0.1")                   // 10%
	viper.SetDefault("validator.min_gas_price", "1000000000")         // 1 Gwei

	// API defaults
	viper.SetDefault("api.enabled", true)
	viper.SetDefault("api.port", 8080)
	viper.SetDefault("api.host", "0.0.0.0")
	viper.SetDefault("api.cors", "*")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("logging.max_size", 100)
	viper.SetDefault("logging.max_backups", 3)
	viper.SetDefault("logging.max_age", 28)
}

// getConfigPath returns the configuration file path
func getConfigPath() string {
	if path := os.Getenv("DIORA_CONFIG"); path != "" {
		return path
	}
	if path := viper.GetString("config_path"); path != "" {
		return path
	}
	return ""
}

// getDefaultConfigPath returns the default configuration file path
func getDefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./diora.yaml"
	}
	return filepath.Join(home, ".diora", "diora.yaml")
}
