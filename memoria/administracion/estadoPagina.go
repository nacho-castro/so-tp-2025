package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
)

func MarcarOcupadoFrame(numeroFrame int) {
	g.FramesLibres[numeroFrame] = false

	g.MutexCantidadFramesLibres.Lock()
	g.CantidadFramesLibres--
	g.MutexCantidadFramesLibres.Unlock()
}

func MarcarLibreFrame(numeroFrameALiberar int) {
	g.MutexEstructuraFramesLibres.Lock()
	g.FramesLibres[numeroFrameALiberar] = true
	g.MutexEstructuraFramesLibres.Unlock()

	g.MutexCantidadFramesLibres.Lock()
	g.CantidadFramesLibres++
	g.MutexCantidadFramesLibres.Unlock()
}
