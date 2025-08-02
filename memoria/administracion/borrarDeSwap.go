package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"io"
	"os"
)

// ========== DESOCUPO ESTRUCTURAS SWAP ==========

func DesocuparProcesoDeSwap(pid int) error {
	g.MutexSwapIndex.Lock()
	defer g.MutexSwapIndex.Unlock()
	for i := 0; i < len(g.SwapIndex[pid].Entradas); i++ {
		entrada := g.SwapIndex[pid].Entradas[i]
		if entrada.Tamanio != 0 {
			errBorrar := BorrarSeccionSwap(int64(entrada.PosicionInicio), entrada.Tamanio)
			if errBorrar != nil {
				return errBorrar
			}
		}
	}

	delete(g.SwapIndex, pid)

	return nil
}

// ========== LIBERO ESPACIO EN SWAP ==========

func BorrarSeccionSwap(posicionInicial int64, tamanio int) error {
	archivo, errAbrir := os.OpenFile(g.MemoryConfig.SwapfilePath, os.O_WRONLY, 0666)
	if errAbrir != nil {
		return errAbrir
	}
	defer func(archivo *os.File) {
		errCerrar := archivo.Close()
		if errCerrar != nil {
		}
	}(archivo)

	_, errSeek := archivo.Seek(posicionInicial, io.SeekStart)
	if errSeek != nil {
		return errSeek
	}

	relleno := make([]byte, tamanio)

	_, errWrite := archivo.Write(relleno)
	if errWrite != nil {
		logger.Error("Error al borrar espacio en SWAP: %v", errWrite)
	}

	return nil
}

func VerificarLecturaDesdeSwap(file *os.File, posicionInicio int, tamanio int) {

	_, errSeek := file.Seek(int64(posicionInicio), io.SeekStart)
	if errSeek != nil {
		logger.Error("Error al hacer seek para verificar: %v", errSeek)
		return
	}

	buffer := make([]byte, tamanio)
	_, errRead := file.Read(buffer)
	if errRead != nil {
		logger.Error("Error al leer para verificar: %v", errRead)
		return
	}

	logger.Debug("## Posición: %d | Tamaño: %d | Contenido: %q", posicionInicio, tamanio, buffer)
}
