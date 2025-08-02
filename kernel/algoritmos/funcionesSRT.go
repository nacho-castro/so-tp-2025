package algoritmos

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"time"
)

func Desalojar(procesoEntrante *pcb.PCB) (cpuID string, desalojado *pcb.PCB) {
	tiempoEntrante := procesoEntrante.EstimadoRafaga
	logger.Debug("SRT: Evaluando posible desalojo por llegada de <%d> con ráfaga estimada %.0f ms", procesoEntrante.PID, tiempoEntrante)

	Utils.MutexEjecutando.Lock()
	defer Utils.MutexEjecutando.Unlock()

	if len(ColaEjecutando.Values()) == 0 {
		logger.Error("SRT: No hay procesos en ejecución")
		return "", nil
	}

	var procesoAInterrumpir *pcb.PCB
	var cpuAInterrumpir string
	var mayorTiempoRestante float64 = 0

	for _, p := range ColaEjecutando.Values() {
		//Queremos interrumpir al proceso con mayor tiempo restante,
		//si el nuevo tiene menor
		tiempoEjecutado := float64(time.Since(p.TiempoEstado).Milliseconds())
		tiempoRestante := p.EstimadoRafaga - tiempoEjecutado
		//logger.Debug("SRT: PID <%d> - Ejecutado: %.0f ms - Restante: %.0f ms", p.PID, tiempoEjecutado, tiempoRestante)

		if tiempoRestante > mayorTiempoRestante {
			mayorTiempoRestante = tiempoRestante
			procesoAInterrumpir = p
			cpuAInterrumpir = p.CpuID
		}
	}

	if procesoAInterrumpir == nil || tiempoEntrante >= mayorTiempoRestante {
		//logger.Debug("SRT: Proceso <%d> NO tiene menor tiempo restante que los procesos ejecutando (%.0f > %.0f)",
		//	procesoEntrante.PID, tiempoEntrante, mayorTiempoRestante)
		return "", nil
	}

	return cpuAInterrumpir, procesoAInterrumpir
}
