package estructuras

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"os"
)

type Config struct {
	PortMemory     int    `json:"port_memory"`
	IpMemory       string `json:"ip_memory"`
	MemorySize     int    `json:"memory_size"`
	PagSize        int    `json:"pag_size"`
	EntriesPerPage int    `json:"entries_per_page"`
	NumberOfLevels int    `json:"number_of_levels"`
	MemoryDelay    int    `json:"memory_delay"`
	SwapfilePath   string `json:"swapfile_path"`
	SwapDelay      int    `json:"swap_delay"`
	LogLevel       string `json:"log_level"`
	DumpPath       string `json:"dump_path"`
	ScriptsPath    string `json:"scripts_path"`
}

func (cfg Config) Validate() error {
	if cfg.IpMemory == "" {
		return errors.New("falta el campo 'ip_memory'")
	}
	if cfg.PortMemory <= 0 {
		return errors.New("falta el campo 'port_memory' o es inválido")
	}
	if cfg.MemorySize <= 0 {
		return errors.New("falta el campo 'memory_size'")
	}
	if cfg.PagSize <= 0 {
		return errors.New("falta el campo 'pag_size'")
	}
	if cfg.EntriesPerPage <= 0 {
		return errors.New("falta el campo 'entries_per_page' o es inválido")
	}
	if cfg.NumberOfLevels <= 0 {
		return errors.New("falta el campo 'number_of_levels' o es inválido")
	}
	if cfg.MemoryDelay <= 0 {
		return errors.New("falta el campo 'memory_delay' o es inválido")
	}
	if cfg.SwapfilePath == "" {
		return errors.New("falta el campo 'swapfile_path'")
	}
	if cfg.SwapDelay <= 0 {
		return errors.New("falta el campo 'swap_delay' o es inválido")
	}
	if cfg.LogLevel == "" {
		return errors.New("falta el campo 'log_level'")
	}
	if cfg.DumpPath == "" {
		return errors.New("falta el campo 'dump_path'")
	}
	if cfg.ScriptsPath == "" {
		return errors.New("falta el campo 'scripts_path'")
	}
	return nil
}

func ConfigMemoria(archivoConfig string) *Config {
	path := fmt.Sprintf("../memoria/configs/%s.json", archivoConfig)

	var config Config

	configFile, err := os.Open(path)
	if err != nil {
		logger.Fatal("No se pudo abrir el archivo de configuración: %v", err)
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			logger.Error("%% Error al cerrar el config file: %v", err)
		}
	}(configFile)

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		logger.Fatal("Error al parsear el archivo de configuración: %v", err)
	}

	if err := config.Validate(); err != nil {
		logger.Fatal("Configuración inválida: %v", err)
	}

	errLogger := logger.ConfigureLogger("memoria.log", "INFO")
	if errLogger != nil {
		logger.Fatal("No se pudo crear el logger -%v", errLogger.Error())
		os.Exit(1)
	}

	errLoggerLevel := logger.SetLevel(config.LogLevel)
	if errLoggerLevel != nil {
		logger.Fatal("No se pudo establecer el log-level: %v", errLoggerLevel)
	}

	return &config
}
