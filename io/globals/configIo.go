package globals

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	IpKernel   string `json:"ip_kernel"`
	PortKernel int    `json:"port_kernel"`
	IpIo       string `json:"ip_io"`
	PortIo     int    `json:"port_io"`
	LogLevel   string `json:"log_level"`
	Type       string `json:"type"`
}

var IoConfig *Config

func (cfg Config) Validate() error {
	if cfg.IpKernel == "" {
		return errors.New("falta el campo 'ip_kernel'")
	}
	if cfg.PortKernel <= 0 {
		return errors.New("falta el campo 'port_kernel' o es inválido")
	}
	if cfg.PortIo <= 0 {
		return errors.New("falta el campo 'port_io' o es inválido")
	}
	if cfg.IpIo == "" {
		return errors.New("falta el campo 'ip_io'")
	}
	if cfg.LogLevel == "" {
		return errors.New("falta el campo 'log_level'")
	}
	if cfg.Type == "" {
		return errors.New("falta el campo 'type'")
	}
	return nil
}

func CargarConfig() *Config {
	configFile, err := os.Open("../io/configs/config.json")
	if err != nil {
		fmt.Printf("No se pudo abrir config: %v\n", err)
		os.Exit(1)
	}
	defer configFile.Close()

	var cfg Config
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&cfg); err != nil {
		fmt.Printf("Error parseando config: %v\n", err)
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Config inválida: %v\n", err)
		os.Exit(1)
	}
	return &cfg
}
