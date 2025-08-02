package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func LiberarProceso(pid int) (g.MetricasProceso, error) {

	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	metricas := proceso.Metricas
	for _, tabla := range proceso.TablaRaiz {
		err := LiberarTablaPaginas(tabla)
		if err != nil {
			return g.MetricasProceso{}, err
		}
	}

	if !proceso.EstaEnSwap {
		_, errDesocupacion := DesocuparProcesoDeEstructurasGlobales(pid)
		if errDesocupacion != nil {
			return metricas, errDesocupacion
		}
	}
	if proceso.EstaEnSwap {
		errSwap := DesocuparProcesoDeSwap(pid)
		if errSwap != nil {
			return g.MetricasProceso{}, errSwap
		}
	}

	return metricas, nil
}

// ========== LIBERAR VECTORES GLOBALES ==========

func DesocuparProcesoDeEstructurasGlobales(pid int) (proceso *g.Proceso, err error) {
	err = nil
	g.MutexProcesosPorPID.Lock()
	proceso = g.ProcesosPorPID[pid]
	proceso.InstruccionesEnBytes = nil
	proceso.Metricas = g.MetricasProceso{}
	delete(g.ProcesosPorPID, pid)
	g.MutexProcesosPorPID.Unlock()

	if proceso == nil {
		logger.Debug("El proceso de PID <%d> no estaba presente en el vector para desocuparla", pid)
		return proceso, logger.ErrNoInstance
	}

	delete(g.MutexMetrica, pid)

	return
}

// ========== DEJAR NULO LOS PUNTEROS DE LA TABLA DE PAGINAS ==========

func LiberarTablaPaginas(tabla *g.TablaPagina) (err error) {
	err = nil

	if tabla.Subtabla != nil {
		for index, subtabla := range tabla.Subtabla {
			err := LiberarTablaPaginas(subtabla)
			if err != nil {
				logger.Error("Error al liberar la tabla de páginas: %v", err)
				return logger.ErrNoTabla
			}
			tabla.Subtabla[index] = nil
		}
		tabla.Subtabla = nil
	}
	if tabla.EntradasPaginas != nil {
		for _, entrada := range tabla.EntradasPaginas {

			tamanioPagina := g.MemoryConfig.PagSize
			direccionFisica := entrada.NumeroFrame * tamanioPagina
			err = RemoverEspacioMemoria(direccionFisica, direccionFisica+tamanioPagina)
			MarcarLibreFrame(entrada.NumeroFrame)
			if err != nil {
				logger.Error("Error al remover espacio de memoria del frame: %d ; %v", entrada.NumeroFrame, err)
				return err
			}

		}
		tabla.EntradasPaginas = nil
	}

	return
}

// ========== LIBERO EL ESPACIO EN MEMORIA ==========

func RemoverEspacioMemoria(inicio int, limite int) (err error) {
	espacioVacio := make([]byte, limite-inicio)
	if inicio < 0 || limite > len(g.MemoriaPrincipal) {
		logger.Error("El inicio es menor a cero o el limite excede el tamaño de la memoria principal")
		return logger.ErrBadRequest
	}

	g.MutexMemoriaPrincipal.Lock()
	copy(g.MemoriaPrincipal[inicio:limite], espacioVacio)
	g.MutexMemoriaPrincipal.Unlock()

	return nil
}
