package traducciones

import (
	"container/list"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

type EntradaTLB struct {
	NroPagina int
	Marco     int
}

type TLB struct {
	entradas    map[int]*list.Element
	orden       *list.List
	maxEntradas int
	algoritmo   string
}

func NuevaTLB() *TLB {
	algoritmo := globals.ClientConfig.TlbReplacement
	maxEntradas := globals.ClientConfig.TlbEntries

	if algoritmo != "FIFO" && algoritmo != "LRU" {
		logger.Warn("Algoritmo TLB inválido: %s", algoritmo)
	}

	return &TLB{
		entradas:    make(map[int]*list.Element),
		orden:       list.New(),
		maxEntradas: maxEntradas,
		algoritmo:   algoritmo,
	}
}

func (tlb *TLB) Buscar(nroPagina int) (int, bool) {
	if tlb.maxEntradas == 0 {
		//logger.Info("TLB deshabilitada (0 entradas)")
		return -1, false
	}

	if elem, ok := tlb.entradas[nroPagina]; ok {
		if tlb.algoritmo == "LRU" {
			tlb.orden.MoveToFront(elem)
		}
		logger.Info("PID: %d - TLB HIT - Página: %d", globals.PIDActual, nroPagina)
		return elem.Value.(EntradaTLB).Marco, true
	}

	logger.Info("PID: %d - TLB MISS - Página: %d", globals.PIDActual, nroPagina)
	return -1, false
}

func (tlb *TLB) AgregarEntrada(nroPagina int, marco int) {
	if tlb.maxEntradas == 0 {
		logger.Debug("TLB deshabilitada: no se agrega entrada")
		return
	}

	if elem, ok := tlb.entradas[nroPagina]; ok {
		elem.Value = EntradaTLB{NroPagina: nroPagina, Marco: marco}
		if tlb.algoritmo == "LRU" {
			tlb.orden.MoveToFront(elem)
		}
		return
	}

	if len(tlb.entradas) >= tlb.maxEntradas {
		victima := tlb.orden.Back()
		if victima != nil {
			entrada := victima.Value.(EntradaTLB)
			delete(tlb.entradas, entrada.NroPagina)
			tlb.orden.Remove(victima)
		}
	}

	nuevaEntrada := EntradaTLB{NroPagina: nroPagina, Marco: marco}
	elem := tlb.orden.PushFront(nuevaEntrada)
	tlb.entradas[nroPagina] = elem
}

func (tlb *TLB) Limpiar() {
	tlb.entradas = make(map[int]*list.Element)
	tlb.orden.Init()
	//logger.Info("PID: %d - TLB limpiada", globals.PIDActual)
}
