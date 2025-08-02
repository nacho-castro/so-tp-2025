package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
	"time"
)

//NO TIENE LOS DATOS PORQUE ESTA EN SWAP
//EL PROCESO SALE DE IO. INTERACTUA CON MEMORIA PERO EL PROCESO SIGUE EN SWAP.

//FIN IO -> SUSP. READY

// CREA PRIMER PROCESO (NEW)
func CrearPrimerProceso(fileName string, tamanio int) {
	// Paso 1: Crear el PCB
	pid := globals.GenerarNuevoPID()
	pcbNuevo := pcb.PCB{
		PID:            pid,
		PC:             0,
		ME:             make(map[string]int),
		MT:             make(map[string]float64),
		EstimadoRafaga: globals.KConfig.InitialEstimate,
		FileName:       fileName,
		ProcessSize:    tamanio,
		TiempoEstado:   time.Now(),
		CpuID:          "",
	}

	//Paso 2: Agregar el primero a la cola NEW
	algoritmos.ColaNuevo.Add(&pcbNuevo)
	pcb.CambiarEstado(&pcbNuevo, pcb.EstadoNew)
	logger.Info("## (<%d>) Se crea el Primer proceso - Estado: <%s>", pcbNuevo.PID, pcbNuevo.Estado)

	//PASO 3: Intentar crear en Memoria
	espacio := comunicacion.SolicitarEspacioEnMemoria(fileName, tamanio)
	if espacio < tamanio {
		logger.Warn("Memoria sin espacio. Abortando")
		return
	}

	//PASO 4: Mandar archivo pseudocodigo a Memoria
	comunicacion.EnviarArchivoMemoria(fileName, tamanio, pid)
}

// ARRANCAR LARGO PLAZO Y PRIMER PROCESO AL PRESIONAR ENTER
func PlanificadorLargoPlazo() {
	logger.Info("## Iniciando el planificador de largo plazo")

	//1. Obtener primer proceso de Cola NEW
	var primerProceso *pcb.PCB
	primerProceso = algoritmos.ColaNuevo.First()
	algoritmos.ColaNuevo.Remove(primerProceso)

	//2. Mandar a Ready
	algoritmos.ColaReady.Add(primerProceso)
	pcb.CambiarEstado(primerProceso, pcb.EstadoReady)

	Utils.NotificarDespachador <- primerProceso.PID //SIGNAL QUE PASO A READY. MANDO PID
	logger.Info("## (<%d>) Pasa del estado NEW al estado %s", primerProceso.PID, primerProceso.Estado)

	go ManejadorCreacionProcesos()
	go ManejadorInicializacionProcesos()
	go ManejadorFinalizacionProcesos()
}

func ManejadorInicializacionProcesos() {
	for {
		<-Utils.InitProcess //SIGNAL llega PROCESO a COLA NEW / SUSP.READY

		//Al llegar un nuevo proceso a NEW
		//y la misma esté vacía
		//y no se tengan procesos en la cola de SUSP READY,
		//se enviará un pedido a Memoria para inicializar el mismo.

		var p *pcb.PCB = nil

		logger.Debug("COLA NEW: %d, COLA SUSP.READY: %d", len(algoritmos.ColaNuevo.Values()), len(algoritmos.ColaSuspendidoReady.Values()))

		// 1. Prioridad: intentar obtener de SUSP.READY
		//(TOMADO POR ALGORITMO)
		if !algoritmos.ColaSuspendidoReady.IsEmpty() { //NO VACIA -> TIENE PRIORIDAD SUSP.READY
			switch globals.KConfig.ReadyIngressAlgorithm {
			case "FIFO":
				p = algoritmos.ColaSuspendidoReady.First()
			case "PMCP":
				p = algoritmos.SeleccionarPMCPSusp() //MAS CHICO
			default:
				logger.Error("Algoritmo de ingreso desconocido")
				return
			}
		} else {
			switch globals.KConfig.ReadyIngressAlgorithm { //BUSCAR EN NEW
			case "FIFO":
				p = algoritmos.ColaNuevo.First()
			case "PMCP":
				p = algoritmos.SeleccionarPMCPNew() //MAS CHICO
			default:
				logger.Error("Algoritmo de ingreso desconocido")
				return
			}
		}

		if p == nil {
			logger.Debug("No hay procesos para inicializar")
			continue
		}

		//logger.Debug("PID: %d, ESTADO: ", p.PID, p.Estado)
		filename := p.FileName
		size := p.ProcessSize
		estadoAnterior := p.Estado

		//Intentar crear en Memoria
		logger.Info("## (<%d>) Inicializa y pide espacio en memoria", p.PID)
		espacio := comunicacion.SolicitarEspacioEnMemoria(filename, size)
		if espacio < size {
			logger.Info("## Memoria sin espacio. PID <%d> queda pendiente", p.PID)

			// 1) Señal al planificador de corto plazo que continue
			Utils.NotificarDespachador <- p.PID
			continue
		}

		//DICE QUE SI, HAY ESPACIO
		//MANDAR PROCESO A READY
		logger.Info("## (<%d>) Pasa del estado %s al estado READY", p.PID, p.Estado)

		//Remover de su cola anterior
		if estadoAnterior == pcb.EstadoSuspReady {
			comunicacion.DesuspensionMemoria(p.PID)
			Utils.MutexSuspendidoReady.Lock()
			algoritmos.ColaSuspendidoReady.Remove(p)
			Utils.MutexSuspendidoReady.Unlock()

		} else if estadoAnterior == pcb.EstadoNew {
			comunicacion.EnviarArchivoMemoria(filename, size, p.PID)
			Utils.MutexNuevo.Lock()
			algoritmos.ColaNuevo.Remove(p)
			Utils.MutexNuevo.Unlock()
		}

		agregarProcesoAReady(p)
	}
}

// RECIBIR SYSCALLS DE CREAR PROCESO
func ManejadorCreacionProcesos() {
	logger.Debug("Esperando solicitudes de INIT_PROC para creación de procesos")
	for {
		//SIGNAL de SYSCALL INIT_PROC
		// Recibir filename, size, pid
		msg := <-Utils.ChannelProcessArguments
		fileName := msg.Filename
		size := msg.Tamanio
		pid := msg.PID
		logger.Debug("Solicitud INIT_PROC recibida: filename=%s, size=%d, pid=%d", fileName, size, pid)

		// AVISAR QUE SE CREO UN PROCESO AL LARGO PLAZO
		Utils.InitProcess <- struct{}{}
	}
}

func agregarProcesoAReady(proceso *pcb.PCB) {
	// 1) estado anterior
	estadoAnterior := proceso.Estado

	// 2) Agregar a READY
	Utils.MutexReady.Lock()
	pcb.CambiarEstado(proceso, pcb.EstadoReady)
	algoritmos.ColaReady.Add(proceso)
	Utils.MutexReady.Unlock()

	logger.Debug("## (<%d>) Pasa del estado %s al estado READY", proceso.PID, estadoAnterior)

	// 4) Señal al planificador de corto plazo
	Utils.NotificarDespachador <- proceso.PID //MANDO PID
}

// RECIBIR SYSCALLS DE EXIT
func ManejadorFinalizacionProcesos() {
	for {
		msg := <-Utils.ChannelFinishprocess
		pid := msg.PID
		pc := msg.PC
		cpuID := msg.CpuID

		// Avisar a Memoria para liberar recursos
		comunicacion.LiberarMemoria(pid)

		//Enviar a EXIT con metricas
		finalizarProceso(pid, pc, cpuID)
	}

}

func finalizarProceso(pid int, pc int, cpuID string) {
	var proceso *pcb.PCB = nil

	// 1. Buscar en EXECUTE y remover
	Utils.MutexEjecutando.Lock()
	for _, p := range algoritmos.ColaEjecutando.Values() {
		if p.PID == pid {
			algoritmos.ColaEjecutando.Remove(p)
			logger.Info("## (<%d>) Pasa del estado EXECUTE al estado EXIT", p.PID)
			p.PC = pc
			proceso = p
			break
		}
	}
	Utils.MutexEjecutando.Unlock()

	// 2. Si no está en EXECUTE, buscar en BLOCKED
	if proceso == nil {
		Utils.MutexBloqueado.Lock()
		for _, p := range algoritmos.ColaBloqueado.Values() {
			if p.PID == pid {
				algoritmos.ColaBloqueado.Remove(p)
				logger.Info("## (<%d>) Pasa del estado BLOCKED al estado EXIT", p.PID)
				proceso = p
				break
			}
		}
		Utils.MutexBloqueado.Unlock()
	}

	// 3. Si no está, buscar en SUSP.BLOCKED
	if proceso == nil {
		Utils.MutexBloqueadoSuspendido.Lock()
		for _, p := range algoritmos.ColaBloqueadoSuspendido.Values() {
			if p.PID == pid {
				algoritmos.ColaBloqueadoSuspendido.Remove(p)
				logger.Info("## (<%d>) Pasa del estado SUSP.BLOCKED al estado EXIT", p.PID)
				proceso = p
				break
			}
		}
		Utils.MutexBloqueadoSuspendido.Unlock()
	}

	// 4. Si no está en ninguna, loguear error
	if proceso == nil {
		//logger.Error("No se pudo finalizar PID=%d, no encontrado en ninguna cola", pid)
		return
	}

	// 5. Mover a EXIT
	pcb.CambiarEstado(proceso, pcb.EstadoExit)
	Utils.MutexSalida.Lock()
	algoritmos.ColaSalida.Add(proceso)
	Utils.MutexSalida.Unlock()

	// 6. Liberar CPU si corresponde
	if cpuID != "" {
		liberarCPU(cpuID)
	}

	// 7. Log y métricas
	logger.Info("## (<%d>) - Finaliza el proceso", proceso.PID)
	logger.Info(proceso.ImprimirMetricas())

	// 8. Señal al liberar memoria
	//Reintentos de creación pendientes
	Utils.InitProcess <- struct{}{}
}
