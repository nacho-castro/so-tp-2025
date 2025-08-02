package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"time"
)

func PlanificadorMedianoPlazo() {
	logger.Info("Iniciando el planificador de Mediano Plazo")
	go ManejadorMedianoPlazo()
	//ESTA FUNCION VA A ATENDER EL REQUERIMIENTO QUE PASE DE SUSP.BLOCKED A SUSP.READY.
	go AtenderSuspBlockedAFinIO()
	go DespacharIO()
}

// ManejadorMedianoPlazo se quedará escuchando
// bloqueos e iniciará un timer por cada PID que llegue.

func ManejadorMedianoPlazo() {
	for bp := range Utils.ChannelProcessBlocked {
		// arrancá un timer en paralelo para CADA proceso bloqueado
		//y lleno el MAP con un mutex entonces cada pid tiene un mutex individual
		/*ch := make(chan Utils.IOEvent, 1)
		Utils.MutexIOWaiters.Lock()
		Utils.IOWaiters[bp.PID] = ch
		Utils.MutexIOWaiters.Unlock()*/
		ch1 := make(chan Utils.IOEvent, 1)

		// 2) Guárdalo en el map protegido con mutex
		Utils.MutexIOFinishedWaiters.Lock()
		Utils.FinIOWaiters[bp.PID] = ch1
		//logger.Error("INICIALICE CANAL DE ")
		Utils.MutexIOFinishedWaiters.Unlock()
		go monitorBloqueado(bp)
	}
}

// moverDeBlockedAReady quita de BLOCKED y encola en READY
func moverDeBlockedAReady(ioLibre Utils.IOEvent) bool {
	// busca en ColaBloqueado
	Utils.MutexBloqueado.Lock()

	// Remover de BLOCKED
	var proceso *pcb.PCB
	for _, p := range algoritmos.ColaBloqueado.Values() {
		if p.PID == ioLibre.PID {
			proceso = p
			algoritmos.ColaBloqueado.Remove(p)
			break
		}
	}
	Utils.MutexBloqueado.Unlock()

	if proceso == nil {
		// No se encontró el proceso en BLOCKED
		return false
	}

	globals.IOMu.Lock()
	instancias, ok := globals.IOs[ioLibre.Nombre]

	if ok {
		for i := range instancias {
			if instancias[i].PID == ioLibre.PID && instancias[i].Puerto == ioLibre.Puerto {
				globals.IOs[ioLibre.Nombre][i].Ocupada = false
				break
			}
		}
	}
	globals.IOMu.Unlock()

	// Agregar a READY
	logger.Info("## (<%d>) finalizó IO y pasa a READY", ioLibre.PID)
	logger.Info("## (<%d>) Pasa del estado %s al estado READY", proceso.PID, proceso.Estado)
	agregarProcesoAReady(proceso) //Señal al corto plazo para despachar

	return true
}

// moverDeBlockedASuspBlocked quita de BLOCKED y encola en SUSP.BLOCKED
func moverDeBlockedASuspBlocked(pid int) bool {
	Utils.MutexBloqueado.Lock()

	// busca en ColaBloqueado
	var proceso *pcb.PCB
	for _, p := range algoritmos.ColaBloqueado.Values() {
		if p.PID == pid {
			proceso = p
			break
		}
	}

	if proceso == nil {
		// No se encontró el proceso en BLOCKED
		//logger.Warn("MedianoPlazo: No se encontró en Blocked")
		//Utils.MutexBloqueado.Unlock()
		return false
	}

	algoritmos.ColaBloqueado.Remove(proceso)
	Utils.MutexBloqueado.Unlock()

	Utils.MutexBloqueadoSuspendido.Lock()
	pcb.CambiarEstado(proceso, pcb.EstadoSuspBlocked)
	algoritmos.ColaBloqueadoSuspendido.Add(proceso)
	Utils.MutexBloqueadoSuspendido.Unlock()

	return true
}

func monitorBloqueado(bp Utils.BlockProcess) {
	pid := bp.PID

	// Obtenemos ambos canales bajo sus mutex
	/*Utils.MutexIOWaiters.Lock()
	blockedCh, okB := Utils.IOWaiters[pid]
	Utils.MutexIOWaiters.Unlock()
	if !okB {
		logger.Warn("monitorBloqueado: no existe blockedCh para PID %d", pid)
		return
	}*/

	Utils.MutexIOFinishedWaiters.Lock()
	finishedCh, okF := Utils.FinIOWaiters[pid]
	Utils.MutexIOFinishedWaiters.Unlock()
	if !okF {
		logger.Warn("monitorBloqueado: no existe finishedCh para PID %d", pid)
		return
	}

	// 1) Esperamos señal de que realmente se bloqueó (del corto plazo)
	/*<-blockedCh*/

	logger.Debug("Arrancó TIMER para PID <%d>", pid)
	suspensión := time.Duration(globals.KConfig.SuspensionTime) * time.Millisecond
	timer := time.NewTimer(suspensión)
	defer timer.Stop()

	// encolamos el pedido de I/O
	Utils.MutexPedidosIO.Lock()
	algoritmos.PedidosIO.Add(&algoritmos.PedidoIO{
		Nombre:   bp.Nombre,
		PID:      bp.PID,
		Duracion: bp.Duracion,
	})
	Utils.MutexPedidosIO.Unlock()
	Utils.NotificarIOLibre <- Utils.IOEvent{Nombre: bp.Nombre, PID: bp.PID}

	select {
	case ioEvt := <-finishedCh:
		// fin de I/O antes del timeout → READY
		moverDeBlockedAReady(ioEvt)

	case <-timer.C:
		// timeout → SUSP.BLOCKED
		if moverDeBlockedASuspBlocked(pid) {
			logger.Info("## (<%d>) Pasa del estado BLOCKED al estado SUSP.BLOCKED (timeout)", pid)
			if err := comunicacion.SolicitarSuspensionEnMemoria(pid); err == nil {
				//Señal al liberar memoria
				//Reintentos de creación pendientes
				Utils.InitProcess <- struct{}{}
			}
			//Avisamos al subplanificador para pasarlo LUEGO a SUSP.READY
			Utils.FinIODesdeSuspBlocked <- Utils.IOEvent{PID: pid, Nombre: bp.Nombre}
		}
	}
}

// pasa de SUSP BLOCKED A SUSP READY (con FIN DE IO)
func AtenderSuspBlockedAFinIO() {
	for ev := range Utils.FinIODesdeSuspBlocked {
		// Obtener el canal dedicado a este PID
		Utils.MutexIOWaiters.Lock()
		finIOChan, ok := Utils.FinIOWaiters[ev.PID]
		Utils.MutexIOWaiters.Unlock()

		if !ok {
			//logger.Warn("AtenderSuspBlockedAFinIO: no hay canal para PID %d", ev.PID)
			continue
		}

		// Arrancar un goroutine que espere el verdadero fin de I/O
		go func(pid int, ch chan Utils.IOEvent) {
			//logger.Info("Esperando fin de I/O para PID <%d> en SUSP.BLOCKED...", pid)
			ioEvt := <-ch
			logger.Debug("FIN DE IO PARA PROCESO PID: <%d>", ioEvt.PID)

			// Pasar de SUSP.BLOCKED a SUSP.READY
			Utils.MutexBloqueadoSuspendido.Lock()
			var proc *pcb.PCB
			for _, p := range algoritmos.ColaBloqueadoSuspendido.Values() {
				if p.PID == pid {
					proc = p
					algoritmos.ColaBloqueadoSuspendido.Remove(p)
					break
				}
			}
			Utils.MutexBloqueadoSuspendido.Unlock()

			if proc != nil {
				Utils.MutexSuspendidoReady.Lock()
				pcb.CambiarEstado(proc, pcb.EstadoSuspReady)
				algoritmos.ColaSuspendidoReady.Add(proc)
				Utils.MutexSuspendidoReady.Unlock()

				logger.Info("## (<%d>) Pasa del estado SUSP.BLOCKED al estado SUSP.READY", pid)
				// Notificar al planificador de largo plazo
				Utils.InitProcess <- struct{}{}
			} else {
				logger.Warn("AtenderSuspBlockedAFinIO: PID %d no estaba en SUSP.BLOCKED", pid)
			}

			// Limpiar el canal dedicado a este PID
			Utils.MutexIOWaiters.Lock()
			delete(Utils.IOWaiters, pid)
			Utils.MutexIOWaiters.Unlock()

			//muestro cola suspready cada vez que llega un fin IO DEBERIA ESTAR CARGADA CON LOS PROCESOS
			//MostrarColasSUSPREADY()

		}(ev.PID, finIOChan)
	}
}

func DespacharIO() {
	for {
		<-Utils.NotificarIOLibre // Esperar señal de IO libre

		Utils.MutexPedidosIO.Lock()
		pedidos := algoritmos.PedidosIO.Values()

		//Si no hay pedidos continua...
		if len(pedidos) == 0 {
			Utils.MutexPedidosIO.Unlock()
			continue
		}

		//BUSCAR PEDIDOS DE IO pendientes
		pedido := algoritmos.PedidosIO.First() // FIFO
		Utils.MutexPedidosIO.Unlock()

		// Buscar una instancia de IO LIBRE del tipo Nombre
		globals.IOMu.Lock()
		var ioAsignada *globals.DatosIO
		for i := range globals.IOs[pedido.Nombre] {
			if !globals.IOs[pedido.Nombre][i].Ocupada {
				globals.IOs[pedido.Nombre][i].Ocupada = true
				globals.IOs[pedido.Nombre][i].PID = pedido.PID
				ioAsignada = &globals.IOs[pedido.Nombre][i]
				break
			}
		}
		globals.IOMu.Unlock()

		if ioAsignada == nil {
			// No se encontró una IO libre
			logger.Debug("No se encontró IO libre: %s. PID <%d> Debe esperar", pedido.Nombre, pedido.PID)
			continue
		}

		logger.Debug("Asignada IO <%s> (puerto %d) a proceso <%d>", ioAsignada.Tipo, ioAsignada.Puerto, pedido.PID)
		algoritmos.PedidosIO.Remove(pedido)
		go comunicacion.EnviarContextoIO(*ioAsignada, pedido.PID, pedido.Duracion)
	}
}
