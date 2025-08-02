package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func EscribirEspacioEntrada(pid int, direccionFisica int, datosEscritura []byte) g.RespuestaEscritura {

	/*if len(datosEscritura) == 0 {
		logger.Debug("Los datos a escribir son vacios: %v", logger.ErrIsEmpty)
	}*/

	errEscritura := ModificarMemoria(direccionFisica, datosEscritura)
	if errEscritura != nil {
		return g.RespuestaEscritura{Exito: errEscritura, DireccionFisica: direccionFisica, Mensaje: errEscritura.Error()}
	}

	err := ModificarEstadoEntradaEscritura(pid)
	if err != nil {
		return g.RespuestaEscritura{Exito: err, DireccionFisica: direccionFisica, Mensaje: err.Error()}
	}

	exito := g.RespuestaEscritura{
		Exito:           nil,
		DireccionFisica: direccionFisica,
		Mensaje:         "Proceso fue modificado correctamente en memoria",
	}

	return exito
}

// =================== MODIFICO LOS VALORES EN LA MEMORIA PRINCIPAL ===================

func ModificarMemoria(direccionFisica int, datosEnBytes []byte) (err error) {
	tamanioPagina := g.MemoryConfig.PagSize
	numeroPagina := direccionFisica / tamanioPagina

	inicioFrame := numeroPagina * tamanioPagina
	finFrame := inicioFrame + tamanioPagina

	if direccionFisica+len(datosEnBytes) > finFrame {
		//logger.Error("Out of range - Escritura fuera del marco asignado")
		//return logger.ErrSegmentFault
	}

	g.MutexMemoriaPrincipal.Lock()
	copy(g.MemoriaPrincipal[direccionFisica:], datosEnBytes)
	g.MutexMemoriaPrincipal.Unlock()

	return nil
}

// =================== MODIFICO LOS VALORES GLOBALES PARA LA PAGINA ===================

func ModificarEstadoEntradaEscritura(pid int) error {
	g.MutexProcesosPorPID.Lock()
	procesoBuscado := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	if procesoBuscado == nil {
		logger.Error("Se intentó acceder a un proceso inexistente o nil para PID <%d>", pid)
		return logger.ErrProcessNil
	}

	IncrementarMetrica(procesoBuscado, 1, IncrementarEscrituraDeMemoria)
	// logger.Info("## Modificacion del estado entrada exitosa post ESCRITURA")

	return nil
}
