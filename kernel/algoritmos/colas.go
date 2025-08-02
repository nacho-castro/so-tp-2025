package algoritmos

import (
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"sync"
)

type Cola[T Nulleable[T]] struct {
	elements []T
	mutex    sync.Mutex
	Priority int
}

type PedidoIO struct {
	Nombre   string
	PID      int
	Duracion int
}

// Saber si son nulos
func (a *PedidoIO) Null() *PedidoIO {
	return nil
}

// Comparar pedidos
func (a *PedidoIO) Equal(b *PedidoIO) bool {
	return a.PID == b.PID && a.Duracion == b.Duracion && a.Nombre == b.Nombre
}

// ESTAS SON VARIABLES GLOBALES OJO¡¡¡¡
var ColaNuevo Cola[*pcb.PCB]
var ColaBloqueado Cola[*pcb.PCB]
var ColaSalida Cola[*pcb.PCB]
var ColaEjecutando Cola[*pcb.PCB]
var ColaReady Cola[*pcb.PCB]
var ColaBloqueadoSuspendido Cola[*pcb.PCB]
var ColaSuspendidoReady Cola[*pcb.PCB]
var PedidosIO Cola[*PedidoIO]

func MostrarCOLABLOQUEADO() {
	lista := ColaBloqueado.Values()

	if len(lista) == 0 {
		logger.Info("Cola NEW vacía")
		return
	}

	logger.Info("Contenido de la cola New:")
	for _, proceso := range lista {
		logger.Info(" - PCB EN COLA New con PID: %d, TAMAÑO: %d", proceso.PID, proceso.ProcessSize)
	}
}

func MostrarColaReady() {
	lista := ColaReady.Values()

	if len(lista) == 0 {
		logger.Info("Cola READY vacía")
		return
	}

	logger.Info("Contenido de la cola READY:")
	for _, proceso := range lista {
		logger.Info(" - PCB EN COLA READY con PID: %d", proceso.PID)
	}
}

func MostrarColaNew() {
	lista := ColaNuevo.Values()

	if len(lista) == 0 {
		logger.Info("Cola NEW vacía")
		return
	}

	logger.Info("Contenido de la cola New:")
	for _, proceso := range lista {
		logger.Info(" - PCB EN COLA New con PID: %d, TAMAÑO: %d", proceso.PID, proceso.ProcessSize)
	}
}

func MostrarColasSUSPREADY() {
	lista := ColaSuspendidoReady.Values()

	if len(lista) == 0 {
		logger.Info("Cola SUSPENDIDO READY vacía")
		return
	}

	logger.Info("Contenido de la cola SUSPENDIDO READY:")
	for _, proceso := range lista {
		logger.Info(" - PCB EN COLA SUSPENDIDO READY con PID: %d, TAMAÑO: %d", proceso.PID, proceso.ProcessSize)
	}
}
