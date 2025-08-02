package Utils

import (
	"sync"
)

var (
	// Mutex para coordinar creaciones concurrentes
	MutexPuedoCrearProceso *sync.Mutex

	// Mutex por cada cola
	MutexNuevo               sync.Mutex
	MutexReady               sync.Mutex
	MutexBloqueado           sync.Mutex
	MutexSalida              sync.Mutex
	MutexEjecutando          sync.Mutex
	MutexBloqueadoSuspendido sync.Mutex
	MutexSuspendidoReady     sync.Mutex
	MutexPedidosIO           sync.Mutex

	//Canales de señalización
	ChannelProcessArguments chan NewProcess
	InitProcess             chan struct{}
	ChannelFinishprocess    chan FinishProcess
	ChannelProcessBlocked   chan BlockProcess
	Desalojo                chan InterruptProcess

	//AVISAR AL DESPACHADOR CUANDO UN PROCESO CAMBIA SU ESTADO
	NotificarDespachador    chan int              //PASA A READY
	NotificarComienzoIO     chan MensajeIOChannel //PASA A BLOQUEADO
	NotificarIOLibre        chan IOEvent          //FIN DE IO
	NotificarDesconexion    chan IOEvent          //Desconexion DE IO
	ContextoInterrupcion    chan InterruptProcess //FIN DE EXECUTE
	NotificarTimeoutBlocked chan int
	FinIODesdeSuspBlocked   chan IOEvent
	IOWaiters               map[int]chan IOEvent
	MutexIOWaiters          sync.Mutex
	FinIOWaiters            map[int]chan IOEvent
	MutexIOFinishedWaiters  sync.Mutex

	InitSuspReady  chan struct{}
	InitNew        chan struct{}
	LiberarMemoria chan struct{}
)

// InicializarMutexes deja listas las variables de mutex.
// Solo MutexPuedoCrearProceso requiere puntero, el resto ya
// está listo con su valor cero.
func InicializarMutexes() {
	MutexPuedoCrearProceso = &sync.Mutex{}
	// MutexNuevo, MutexReady, ... ya funcionan sin más
}

// InicializarCanales crea y configura los canales con buffers adecuados.
func InicializarCanales() {
	ChannelProcessArguments = make(chan NewProcess, 20) // buffer para hasta 10 peticiones
	ChannelFinishprocess = make(chan FinishProcess, 20)
	InitProcess = make(chan struct{}) // sin buffer para sincronización exacta
	LiberarMemoria = make(chan struct{}, 1)
	Desalojo = make(chan InterruptProcess)

	InitSuspReady = make(chan struct{}) // sin buffer para sincronización exacta
	InitNew = make(chan struct{})       // sin buffer para sincronización exacta

	NotificarDespachador = make(chan int, 20) // buffer 10 procesos listos
	NotificarComienzoIO = make(chan MensajeIOChannel, 20)
	NotificarIOLibre = make(chan IOEvent, 20)
	NotificarDesconexion = make(chan IOEvent, 20)
	ContextoInterrupcion = make(chan InterruptProcess, 20)
	ChannelProcessBlocked = make(chan BlockProcess, 20)
	NotificarTimeoutBlocked = make(chan int)
	FinIODesdeSuspBlocked = make(chan IOEvent, 20)
	IOWaiters = make(map[int]chan IOEvent)
	FinIOWaiters = make(map[int]chan IOEvent)
}

type MensajeIOChannel struct {
	PID      int
	PC       int
	Nombre   string
	Duracion int
	CpuID    string
}
type FinishProcess struct {
	PID   int
	PC    int
	CpuID string
}
type InterruptProcess struct {
	PID   int
	PC    int
	CpuID string
}
type BlockProcess struct {
	PID      int
	PC       int
	Nombre   string
	Duracion int
	CpuID    string
}
type NewProcess struct {
	Filename string
	Tamanio  int
	PID      int
}
type IOEvent struct {
	PID    int
	Nombre string
	Puerto int
}
