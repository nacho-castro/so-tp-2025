package estructuras

import "sync"

// ================== VARIABLES GLOBALES==================

var MemoryConfig *Config
var MemoriaPrincipal []byte         // MP simulada
var FramesLibres []bool             //los frames van a estar en True si están libres
var CantidadFramesLibres int        // simplemente recuenta la cantidad de frames
var ProcesosPorPID map[int]*Proceso // guardo procesos con los PID
var SwapIndex map[int]*SwapProcesoInfo
var PunteroSwap int

// ================== SEMAFOROS GLOBALES ==================

var MutexProcesosPorPID sync.Mutex
var MutexMemoriaPrincipal sync.Mutex
var MutexCantidadFramesLibres sync.Mutex
var MutexEstructuraFramesLibres sync.Mutex
var MutexMetrica map[int]*sync.Mutex
var MutexSwapIndex sync.Mutex
var MutexOperacionMemoria sync.Mutex

// ================== TABLA Y ENTRADA DE PÁGINAS ==================

type TablaPaginas map[int]*TablaPagina
type TablaPagina struct {
	Subtabla        map[int]*TablaPagina   `json:"subtabla"`
	EntradasPaginas map[int]*EntradaPagina `json:"entradas_pagina"`
}
type EntradaPagina struct {
	NumeroFrame  int  `json:"numero_frame"`
	EstaPresente bool `json:"esta_presente"`
	//EstaEnUso     bool `json:"esta_en_uso"`
	//FueModificado bool `json:"fue_modificado"`
}

// ================== PROCESOS ==================

type Proceso struct {
	PID                  int             `json:"pid"`
	TablaRaiz            TablaPaginas    `json:"tabla_paginas"`
	Metricas             MetricasProceso `json:"metricas_proceso"`
	InstruccionesEnBytes map[int][]byte  `json:"instrucciones_en_bytes"`
	EstaEnSwap           bool            `json:"esta_en_swap"`
}
type MetricasProceso struct {
	AccesosTablasPaginas     int `json:"acceso_tablas_paginas"`
	InstruccionesSolicitadas int `json:"instrucciones_solicitadas"`
	BajadasSwap              int `json:"bajadas_swap"`
	SubidasMP                int `json:"subidas_mp"`
	LecturasDeMemoria        int `json:"lecturas_de_memoria"`
	EscriturasDeMemoria      int `json:"escrituras_de_memoria"`
}
type OperacionMetrica func(*MetricasProceso, int)

// ================== SWAP ==================

type SwapProcesoInfo struct {
	Entradas             map[int]*EntradaSwapInfo `json:"entradas"`
	InstruccionesEnBytes map[int][]byte           `json:"instrucciones_en_bytes"`
	NumerosDePaginas     []int                    `json:"numeros_de_paginas"`
}
type EntradaSwapInfo struct {
	NumeroPagina   int `json:"numero_pagina"`
	PosicionInicio int `json:"posicion_inicio"`
	Tamanio        int `json:"tamanio"`
}
type EntradaSwap struct {
	NumeroPagina int    `json:"numero_pagina"`
	Datos        []byte `json:"datos"`
}
