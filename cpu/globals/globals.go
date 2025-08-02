package globals

import (
	"errors"
	"sync"
)

// No se si es correcto crear una carpeta globals
type Config struct {
	IpSelf           string `json:"ip_self"`
	PortSelf         int    `json:"port_self"`
	IpMemory         string `json:"ip_memory"`
	PortMemory       int    `json:"port_memory"`
	IpKernel         string `json:"ip_kernel"`
	PortKernel       int    `json:"port_kernel"`
	TlbEntries       int    `json:"tlb_entries"`
	TlbReplacement   string `json:"tlb_replacement"`
	CacheEntries     int    `json:"cache_entries"`
	CacheDelay       int    `json:"cache_delay"`
	CacheReplacement string `json:"cache_replacement"`
	LogLevel         string `json:"log_level"`
}

// En el paquete globals
var MutexPID sync.Mutex
var MutexPC sync.Mutex

var PIDActual int
var PCActual int
var ClientConfig *Config
var InterrupcionPendiente bool
var PIDInterrumpido int
var ErrSyscallBloqueante = errors.New("proceso bloqueado por syscall IO")
var SaltarIncrementoPC bool
var ID string
var TamanioPagina int
var EntradasPorNivel int
var CantidadNiveles int
var MutexInterrupcion sync.Mutex
