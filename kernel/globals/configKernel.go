package globals

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"os"
)

type KernelConfig struct {
	MemoryAddress         string  `json:"ip_memory"`
	MemoryPort            int     `json:"port_memory"`
	KernelAddress         string  `json:"ip_kernel"`
	KernelPort            int     `json:"port_kernel"`
	SchedulerAlgorithm    string  `json:"scheduler_algorithm"`
	ReadyIngressAlgorithm string  `json:"ready_ingress_algorithm"`
	Alpha                 float64 `json:"alpha"`
	InitialEstimate       float64 `json:"initial_estimate"`
	SuspensionTime        int     `json:"suspension_time"`
	LogLevel              string  `json:"log_level"`
}

func (cfg KernelConfig) Validate() error {
	if cfg.MemoryAddress == "" {
		return errors.New("falta el campo 'ip_memoria'")
	}
	if cfg.MemoryPort == 0 {
		return errors.New("falta el campo 'puerto_memoria' o es inválido")
	}
	if cfg.KernelPort == 0 {
		return errors.New("falta el campo 'puerto_kernel' o es inválido")
	}
	if cfg.KernelAddress == "" {
		return errors.New("falta el campo 'ip_kernel'")
	}
	if cfg.SchedulerAlgorithm == "" {
		return errors.New("falta el campo 'scheduler_algorithm'")
	}
	if cfg.ReadyIngressAlgorithm == "" {
		return errors.New("falta el campo 'ready_ingress_algorithm'")
	}
	if cfg.Alpha <= 0 {
		return errors.New("falta el campo 'alpha' o es inválido")
	}
	if cfg.InitialEstimate <= 0 {
		return errors.New("falta el campo 'initial_estimate' o es inválido")
	}
	if cfg.SuspensionTime < 0 {
		return errors.New("falta el campo 'suspension_time' o es inválido")
	}
	if cfg.LogLevel == "" {
		return errors.New("falta el campo 'log_level'")
	}
	return nil
}

func CargarConfig(archivoConfig string) *KernelConfig {
	path := fmt.Sprintf("../kernel/configs/%s.json", archivoConfig)

	file, err := os.Open(path)
	if err != nil {
		logger.Fatal("No se pudo abrir el archivo de configuración: %v", err)
	}
	defer file.Close()

	var config KernelConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		logger.Fatal("Error al parsear el archivo de configuración: %v", err)
	}

	if err := config.Validate(); err != nil {
		logger.Fatal("Configuración inválida: %v", err)
	}

	return &config
}
