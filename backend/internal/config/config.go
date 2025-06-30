package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Nuclei   NucleiConfig   `yaml:"nuclei"`
	GRPC     GRPCConfig     `yaml:"grpc"`
	Logging  LoggingConfig  `yaml:"logging"`
	RF       RFConfig       `yaml:"rf"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

type NucleiConfig struct {
	TemplatesPath string `yaml:"templates_path"`
	Concurrency   int    `yaml:"concurrency"`
	RateLimit     int    `yaml:"rate_limit"`
}

type GRPCConfig struct {
	Port       int  `yaml:"port"`
	TLSEnabled bool `yaml:"tls_enabled"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

type RFConfig struct {
	DevicePort    string  `yaml:"device_port"`
	SampleRate    int     `yaml:"sample_rate"`
	FrequencyMin  float64 `yaml:"frequency_min"`
	FrequencyMax  float64 `yaml:"frequency_max"`
	GainMode      string  `yaml:"gain_mode"`
	Gain          float64 `yaml:"gain"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
			Mode: getEnv("SERVER_MODE", "development"),
		},
		Database: DatabaseConfig{
			Type: getEnv("DB_TYPE", "sqlite"),
			Path: getEnv("DB_PATH", "./data/scanner.db"),
		},
		Nuclei: NucleiConfig{
			TemplatesPath: getEnv("NUCLEI_TEMPLATES_PATH", "./nuclei-templates"),
			Concurrency:   getEnvAsInt("NUCLEI_CONCURRENCY", 10),
			RateLimit:     getEnvAsInt("NUCLEI_RATE_LIMIT", 150),
		},
		GRPC: GRPCConfig{
			Port:       getEnvAsInt("GRPC_PORT", 9090),
			TLSEnabled: getEnvAsBool("GRPC_TLS_ENABLED", false),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", "./logs/scanner.log"),
		},
		RF: RFConfig{
			DevicePort:   getEnv("RF_DEVICE_PORT", "/dev/ttyUSB0"),
			SampleRate:   getEnvAsInt("RF_SAMPLE_RATE", 2048000),
			FrequencyMin: getEnvAsFloat("RF_FREQUENCY_MIN", 2400000000),
			FrequencyMax: getEnvAsFloat("RF_FREQUENCY_MAX", 5800000000),
			GainMode:     getEnv("RF_GAIN_MODE", "auto"),
			Gain:         getEnvAsFloat("RF_GAIN", 0),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}
