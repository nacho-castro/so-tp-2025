package globals

import (
	"sync"
)

// Datos recibidos por el Kernel
type DatosIO struct {
	Tipo    string
	Ip      string
	Puerto  int
	Ocupada bool
	PID     int
}

func (d *DatosIO) Null() *DatosIO {
	return nil
}

func (d *DatosIO) Equal(other *DatosIO) bool {
	if d == nil && other == nil {
		return true
	}
	if d == nil || other == nil {
		return false
	}
	return d.Tipo == d.Tipo && d.Ip == other.Ip && d.Puerto == other.Puerto
}

type DatosCPU struct {
	ID      string
	Ip      string
	Puerto  int
	Ocupada bool
}

type EspacioLibreRTA struct {
	EspacioLibre int `json:"espacio_libre"`
}

var KConfig *KernelConfig

var CPU DatosCPU
var CPUs map[string]DatosCPU = make(map[string]DatosCPU) // clave: ID del CPU
var CPUMu sync.Mutex
var CPUCond = sync.NewCond(&CPUMu)

// MAP CLAVE: TIPO IO VALOR: ARREGLO DE INSTANCIAS DE TIPO DatosIO
var IO DatosIO
var IOs map[string][]DatosIO = make(map[string][]DatosIO)

var IOMu sync.Mutex

var EspacioLibreProceso EspacioLibreRTA
var UltimoPID int = -1
var PidMutex sync.Mutex

func GenerarNuevoPID() int {
	PidMutex.Lock()
	defer PidMutex.Unlock()

	UltimoPID++
	return UltimoPID
}

//crear nodos a punteros PCB, para instanciarlas en main (sera un puntero a pcb creado previamente)
