## MEMORIA



## 🔌 1. Endpoint expuesto

La memoria se encarga de tener conexiones entrantes para los módulos:

`http://localhost:8083/memoria/kernel`
`http://localhost:8083/memoria/cpu`

## 📬 2. Formato del mensaje recibido

El cuerpo del mensaje (`body`) debe ser un JSON con una estructura dependiendo de cada Modulo:

```json
{
  "ip": "127.0.0.1",
  "puerto": 8000
}

{
  "nombre":"impresora",
  "ip": "127.0.0.1",
  "puerto": 8000
}
```

# Estructuras Globals
```go
var MemoryConfig *Config
var MemoriaPrincipal []byte         // MP simulada
var FramesLibres []bool             //los frames van a estar en True si están libres
var CantidadFramesLibres int        // simplemente recuenta la cantidad de frames
var ProcesosPorPID map[int]*Proceso // guardo procesos con los PID
```

# Estructuras con CPU
```go
var CPU DatosDeCPU

type DatosDeCPU struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type ContextoCPU struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type RespuestaInstruccion struct {
	Instruccion string `json:"instruccion"`
}

type RespuestaConfigMemoria struct {
	TamanioPagina    int `json:"tamanioPagina"`
	EntradasPorNivel int `json:"entradasPorNivel"`
	CantidadNiveles  int `json:"cantidadNiveles"`
}

type ConsultaMarco struct {
	PID            int   `json:"pid"`
	IndicesEntrada []int `json:"indices_entrada"`
}

type RespuestaMarco struct {
	NumeroMarco int `json:"numero_marco"`
}
```


# Estructuras con Kernel
```go
var RespuestaKernel InitProceso

// var Kernel DatosConsultaDeKernel

type DatosConsultaDeKernel struct {
	PID            int `json:"pid"`
	TamanioMemoria int `json:"tamanio_memoria"`
}

type InitProceso struct {
	PID            int    `json:"pid"`
	TamanioMemoria int    `json:"tamanio_memoria"`
	Pseudocodigo   string `json:"filename"`
}

type RespuestaMemoria struct {
	Exito   bool   `json:"exito"` // TODO: cambiar a error
	Mensaje string `json:"mensaje"`
}

type RespuestaEspacioLibre struct {
	EspacioLibre int `json:"espacio_libre"`
}

type LecturaProceso struct {
	PID              int `json:"pid"`
	DireccionFisica  int `json:"direccion_fisica"`
	TamanioARecorrer int `json:"tamanio_a_recorrer"`
}

type ExitoLecturaMemoria struct {
	Exito        error  `json:"exito"`
	DatosAEnviar string `json:"datos_a_enviar"`
}

type EscrituraProceso struct {
	PID              int 	`json:"pid"`
	DireccionFisica  int 	`json:"direccion_fisica"`
	TamanioARecorrer int 	`json:"tamanio_a_recorrer"`
	DatosAEscribir	 string `json:"datos_a_escribir"`
}

type FinalizacionProceso struct {
	PID int `json:"pid"`
}

type ExitoEdicionMemoria struct {
	Exito    error `json:"exito"`
	Booleano bool  `json:"booleano"`
}
```


# Estructuras de Paginación
```go
type EntradaPagina struct {
	NumeroPagina   int  `json:"numero_frame"`
	EstaPresente  bool `json:"esta_presente"`
	EstaEnUso     bool `json:"esta_en_uso"`
	FueModificado bool `json:"fue_modificado"`
}

type TablaPagina struct {
	Subtabla        map[int]*TablaPagina   `json:"subtabla"`
	EntradasPaginas map[int]*EntradaPagina `json:"entradas_pagina"`
}

type TablaPaginas map[int]*TablaPagina

type EscrituraPagina struct {
	PID                 int    `json:"pid"`
	DireccionFisica     int    `json:"direccion_fisica"`
	DatosASobreEscribir string `json:"datos_a_sobre_escribir"`
	TamanioNecesario    int    `json:"tamanio_necesario"`
}

type RespuestaEscritura struct {
	Exito           error  `json:"exito"`
	DireccionFisica int    `json:"direccion_fisica"`
	Mensaje         string `json:"mensaje"`
}

type LecturaPagina struct {
	PID             int `json:"pid"`
	DireccionFisica int `json:"direccion_fisica"`
}

type RespuestaLectura struct {
	Exito           error  `json:"exito"`
	PseudoCodigo    string `json:"pseudo_codigo"`
	DireccionFisica int    `json:"direccion_fisica"`
}
```


# Estructuras de Procesos
```go
type Proceso struct {
	PID       int             `json:"pid"`
	TablaRaiz TablaPaginas    `json:"tabla_paginas"`
	Metricas  MetricasProceso `json:"metricas_proceso"`
}

type Ocupante struct {
	PID          int `json:"pid"`
	NumeroPagina int `json:"numero_pagina"`
}

type MetricasProceso struct {
	AccesosTablasPaginas     int `json:"acceso_tablas_paginas"`
	InstruccionesSolicitadas int `json:"instrucciones_solicitadas"`
	BajadasSwap              int `json:"bajadas_swap"`
	SubidasMP                int `json:"subidas_mp"`
	LecturasDeMemoria        int `json:"lecturas_de_memoria"`
	EscriturasDeMemoria      int `json:"escrituras_de_memoria"`
}

type OperacionMetrica func(*MetricasProceso)

type ConsultaDump struct {
	PID       int    `json:"pid"`
	TimeStamp string `json:"timeStamp"`
}
type EntradaDump struct {
	DireccionFisica int `json:"direccion_fisica"`
	NumeroPagina     int `json:"numero_frame"`
}
```


# Estructuras de Semaforos
```go
var MutexProcesosPorPID sync.Mutex
var MutexMemoriaPrincipal sync.Mutex
var MutexCantidadFramesLibres sync.Mutex
var MutexEstructuraFramesLibres sync.Mutex
var MutexMetrica []sync.Mutex

func CambiarEstadoFrame(numeroFrame int) {
	MutexEstructuraFramesLibres.Lock()
	if FramesLibres[numeroFrame] == false {
		FramesLibres[numeroFrame] = true
	} else {
		FramesLibres[numeroFrame] = false
	}
	MutexEstructuraFramesLibres.Unlock()
}
```

# Estructuras para SWAP

```go
var ProcesosSuspendidos map[int]ProcesoSuspendido

type ProcesoSuspendido struct {
// 	PID 			int `json:"pid"` TODO: este voy a usar para mapear el vector
DireccionFisica string `json:"direccion_fisica"`
TamanioProceso  string `json:"tamanio_proceso"`
Base            int    `json:"base"`
Limite          int    `json:"limite"`
// TODO: ver que mas agrego
}

type SuspensionProceso struct {
PID    int   `json:"pid"`
Indice []int `json:"indice"`
}

type ExitoSuspensionProceso struct {
Exito           error `json:"exito"`
DireccionFisica int   `json:"direccion_fisica"`
TamanioProceso  int   `json:"tamanio_proceso"`
}

type DesuspensionProceso struct {
PID int `json:"pid"`
}

type ExitoDesuspensionProceso struct {
Exito           error `json:"exito"`
DireccionFisica int   `json:"direccion_fisica"`
TamanioProceso  int   `json:"tamanio_proceso"`
}
```