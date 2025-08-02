package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"os"
	"sync"
)

func InicializarMemoriaPrincipal() {
	InstanciarEstructurasGlobales(g.MemoryConfig.MemorySize / g.MemoryConfig.PagSize)

	logger.Info("Memoria Principal Inicializada con %d bytes de tamaño con %d frames de %d bytes.",
		g.MemoryConfig.MemorySize, g.MemoryConfig.MemorySize/g.MemoryConfig.PagSize, g.MemoryConfig.PagSize)
}

func InstanciarEstructurasGlobales(cantidadFrames int) {
	g.MemoriaPrincipal = make([]byte, g.MemoryConfig.MemorySize)
	g.ProcesosPorPID = make(map[int]*g.Proceso)
	g.SwapIndex = make(map[int]*g.SwapProcesoInfo)
	g.MutexMetrica = make(map[int]*sync.Mutex, g.MemoryConfig.MemorySize)

	ConfigurarFrames(cantidadFrames)
}

func ConfigurarFrames(cantidadFrames int) {
	g.FramesLibres = make([]bool, cantidadFrames)
	g.MutexEstructuraFramesLibres.Lock()
	for i := 0; i < cantidadFrames; i++ {
		g.FramesLibres[i] = true
	}
	g.MutexEstructuraFramesLibres.Unlock()
	g.CantidadFramesLibres = cantidadFrames
}

func LimpiarSwapFile() {
	file, err := os.Create(g.MemoryConfig.SwapfilePath)
	if err != nil {
		logger.Error("Error al crear/sobrescribir el swap file: %v", err)
		return
	}
	err = file.Close()
	if err != nil {
		logger.Error("Error al cerrar el swap file recién creado: %v", err)
	}
}
