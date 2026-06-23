package config

// AppConfig holds all runtime configuration for slide-forge.
type AppConfig struct {
	Transport string `yaml:"transport"` // stdio | http (default: stdio)
	Addr      string `yaml:"addr"`      // HTTP listen address (default: :8080)
	LogLevel  string `yaml:"log_level"` // debug | info | warn | error (default: info)
}
