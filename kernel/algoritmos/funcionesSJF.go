package algoritmos

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
)

// definido por archivo de configuración:
// - rafaga inicial estimada
// - alpha
// Criterio: Se elegirá el proceso que tenga la rafaga más corta.
func SeleccionarSJF() *pcb.PCB {
	Utils.MutexReady.Lock()
	defer Utils.MutexReady.Unlock()

	if len(ColaReady.elements) == 0 {
		return nil
	}

	//logger.Debug("Cola READY tiene %d procesos\n", len(ColaReady.elements))
	masChico := ColaReady.elements[0] //Tomo el primero y empiezo a comparar rafagas
	for _, p := range ColaReady.elements {
		//logger.Debug("<%d> %.0f |VS| <%d> %.0f", masChico.PID, masChico.EstimadoRafaga, p.PID, p.EstimadoRafaga)
		if p.EstimadoRafaga < masChico.EstimadoRafaga {
			masChico = p
		}
	}
	//logger.Debug("Seleccionado SJF: <%d> | Rafaga: %.0f", masChico.PID, masChico.EstimadoRafaga)
	return masChico
}

/*
SJF con Desalojo
Funciona igual que el anterior con la variante que,
al ingresar un proceso en la cola de Ready y no haber CPUs libres,
se debe evaluar si dicho proceso tiene una rafaga
más corta que los que se encuentran en ejecución.
*/
