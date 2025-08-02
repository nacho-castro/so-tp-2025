package estructuras

// ================== PEDIDOS PARA UN PROCESO ==================

type InitProceso struct {
	PID            int    `json:"pid"`
	TamanioMemoria int    `json:"tamanio_memoria"`
	Pseudocodigo   string `json:"filename"`
}
type ConsultaProceso struct {
	PID int `json:"pid"`
}
type ConsultaDump struct {
	PID       int    `json:"pid"`
	TimeStamp string `json:"timeStamp"`
}

// ======================= CONSULTAS =======================

type RespuestaMemoria struct {
	Exito   bool   `json:"exito"` // TODO: cambiar a error
	Mensaje string `json:"mensaje"`
}

type RespuestaEspacioLibre struct {
	EspacioLibre int `json:"espacio_libre"`
}

// ================== PEDIDO INSTRUCCION ==================

type ConsultaContextCPU struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}
type RespuestaInstruccion struct {
	Exito       error  `json:"exito"`
	Instruccion string `json:"instruccion"`
}

// ================== RESPUESTA CONFIG DE MEMORIA ==================

type RespuestaConfigMemoria struct {
	TamanioPagina    int `json:"tamanioPagina"`
	EntradasPorNivel int `json:"entradasPorNivel"`
	CantidadNiveles  int `json:"cantidadNiveles"`
}

// ================== PEDIDO MARCO DE PAGINA ==================

type ConsultaMarco struct {
	PID            int   `json:"pid"`
	IndicesEntrada []int `json:"indices_entrada"`
}
type RespuestaMarco struct {
	NumeroMarco int `json:"numero_marco"`
}

// ================== ESCRITURA Y LECTURA ==================

type EscrituraProceso struct {
	PID             int    `json:"pid"`
	DireccionFisica int    `json:"direccion_fisica"`
	DatosAEscribir  string `json:"datos_a_escribir"`
}

type LecturaProceso struct {
	PID              int `json:"pid"`
	DireccionFisica  int `json:"direccion_fisica"`
	TamanioARecorrer int `json:"tamanio_a_recorrer"`
}

// ================== RESPUESTAS A ESCRITURA Y LECTURA ==================

type RespuestaEscritura struct {
	Exito           error  `json:"exito"`
	DireccionFisica int    `json:"direccion_fisica"`
	Mensaje         string `json:"mensaje"`
}

type RespuestaLectura struct {
	Exito error  `json:"exito"`
	Valor string `json:"valor_leido"`
}
