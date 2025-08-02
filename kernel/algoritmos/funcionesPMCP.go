package algoritmos

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
)

// ready_ingress_algorithm
// Criterio: mas chico el que menos memoria solicite (size).

func AddPMCPNew(p *pcb.PCB) {
	Utils.MutexNuevo.Lock()
	defer Utils.MutexNuevo.Unlock()

	insertado := false
	nuevaCola := make([]*pcb.PCB, 0, len(ColaNuevo.elements)+1)

	for _, actual := range ColaNuevo.elements {
		if !insertado && p.ProcessSize < actual.ProcessSize {
			nuevaCola = append(nuevaCola, p)
			insertado = true
		}
		nuevaCola = append(nuevaCola, actual)
	}

	if !insertado {
		nuevaCola = append(nuevaCola, p)
	}

	ColaNuevo.elements = nuevaCola
}

func AddPMCPSusp(p *pcb.PCB) {
	Utils.MutexSuspendidoReady.Lock()
	defer Utils.MutexSuspendidoReady.Unlock()

	insertado := false
	nuevaCola := make([]*pcb.PCB, 0, len(ColaSuspendidoReady.elements)+1)

	for _, actual := range ColaSuspendidoReady.elements {
		if !insertado && p.ProcessSize < actual.ProcessSize {
			nuevaCola = append(nuevaCola, p)
			insertado = true
		}
		nuevaCola = append(nuevaCola, actual)
	}

	if !insertado {
		nuevaCola = append(nuevaCola, p)
	}

	ColaSuspendidoReady.elements = nuevaCola
}

func AddPMCPReady(p *pcb.PCB) {
	Utils.MutexReady.Lock()
	defer Utils.MutexReady.Unlock()

	insertado := false
	nuevaCola := make([]*pcb.PCB, 0, len(ColaReady.elements)+1)

	for _, actual := range ColaReady.elements {
		if !insertado && p.ProcessSize < actual.ProcessSize {
			nuevaCola = append(nuevaCola, p)
			insertado = true
		}
		nuevaCola = append(nuevaCola, actual)
	}

	if !insertado {
		nuevaCola = append(nuevaCola, p)
	}

	ColaReady.elements = nuevaCola
}

func SeleccionarPMCPSusp() *pcb.PCB {
	Utils.MutexSuspendidoReady.Lock()
	defer Utils.MutexSuspendidoReady.Unlock()

	if len(ColaSuspendidoReady.elements) == 0 {
		return nil
	}

	masChico := ColaSuspendidoReady.elements[0]
	for _, p := range ColaSuspendidoReady.elements[1:] {
		if p.ProcessSize < masChico.ProcessSize {
			masChico = p
		}
	}
	return masChico
}

func SeleccionarPMCPNew() *pcb.PCB {
	Utils.MutexNuevo.Lock()
	defer Utils.MutexNuevo.Unlock()

	if len(ColaNuevo.elements) == 0 {
		return nil
	}

	masChico := ColaNuevo.elements[0]
	for _, p := range ColaNuevo.elements[1:] {
		if p.ProcessSize < masChico.ProcessSize {
			masChico = p
		}
	}
	return masChico
}
