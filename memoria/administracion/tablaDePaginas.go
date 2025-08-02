package administracion

import (
	"fmt"
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

// ============= SE DEBE USAR EXCLUSIVAMENTE PARA CREAR INDICES CON NUMEROS DE PAGINAS LOGICAS =============

func CrearIndicePara(numeroPaginaLogica int) (indicesParaTabla []int) {
	cantidadNiveles := g.MemoryConfig.NumberOfLevels
	cantidadEntradasPorTabla := g.MemoryConfig.EntriesPerPage

	indicesParaTabla = make([]int, cantidadNiveles)
	divisor := 1

	for i := cantidadNiveles - 1; i >= 0; i-- {
		indicesParaTabla[i] = (numeroPaginaLogica / divisor) % cantidadEntradasPorTabla
		divisor *= cantidadEntradasPorTabla
	}
	return
}

func ObtenerEntradaPagina(pid int, indices []int) (int, error) {
	g.MutexProcesosPorPID.Lock()
	proceso, errPro := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()
	if !errPro {
		logger.Error("Processo Buscado no existe")
		return -1, fmt.Errorf("el proceso no existe o nunca fue inicializada: %w", logger.ErrNoInstance)
	}
	entradaPagina, errPag := BuscarEntradaPagina(proceso, indices)
	if errPag != nil {
		logger.Error("Error al buscar la entrada de página")
		return -1, fmt.Errorf("la entrada no existe o nunca fue inicializada: %w", logger.ErrNoInstance)
	}
	return entradaPagina.NumeroFrame, nil
}

func BuscarEntradaPagina(proceso *g.Proceso, indicesParaTabla []int) (entradaDeseada *g.EntradaPagina, err error) {
	if proceso == nil {
		logger.Error("Proceso es nil en BuscarEntradaPagina")
		return nil, logger.ErrProcessNil
	}

	if proceso.TablaRaiz == nil {
		logger.Error("TablaRaiz es nil en BuscarEntradaPagina")
		return nil, logger.ErrNoTabla
	}

	if len(indicesParaTabla) == 0 {
		logger.Error("Indices vacíos en BuscarEntradaPagina")
		return nil, logger.ErrNoIndices
	}

	tablaApuntada := proceso.TablaRaiz[indicesParaTabla[0]]
	if tablaApuntada == nil {
		logger.Error("La tabla no existe o nunca fue inicializada")
		return nil, logger.ErrNoTabla
	}

	for i := 1; i < len(indicesParaTabla)-1; i++ {
		if tablaApuntada.Subtabla == nil {
			logger.Error("La subtabla no existe o nunca fue inicializada")
			return nil, logger.ErrNoTabla
		}
		tablaApuntada = tablaApuntada.Subtabla[indicesParaTabla[i]]
		if tablaApuntada == nil {
			logger.Error("La subtabla no existe en el índice <%d>", indicesParaTabla[i])
			return nil, logger.ErrNoInstance
		}
	}

	if tablaApuntada.EntradasPaginas == nil {
		logger.Error("Las EntradasPaginas era nil para el índice <%v> para el PID <%d>", indicesParaTabla, proceso.PID)
		return nil, logger.ErrNoInstance
	}

	entradaDeseada = tablaApuntada.EntradasPaginas[indicesParaTabla[len(indicesParaTabla)-1]]
	if entradaDeseada == nil {
		logger.Error("La entrada buscada no existe del PID <%d>", proceso.PID)
		return nil, logger.ErrNoInstance
	}

	if entradaDeseada.EstaPresente == false {
		//logger.Info("Se pasa a presente entrada <%d> de PID <%d>", entradaDeseada, proceso.PID)
		return entradaDeseada, nil
	}

	IncrementarMetrica(proceso, 1, IncrementarAccesosTablasPaginas)
	return entradaDeseada, nil
}
