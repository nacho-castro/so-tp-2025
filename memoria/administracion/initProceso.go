package administracion

import (
	"bufio"
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"os"
	"strings"
	"sync"
)

func InicializarProceso(pid int, tamanioProceso int, nombreArchPseudocodigo string) (err error) {

	if !TieneTamanioNecesario(tamanioProceso) {
		logger.Error("No hay memoria suficiente para proceso PID <%d>", pid)
		return logger.ErrNoMemory
	}

	g.MutexProcesosPorPID.Lock()
	if g.ProcesosPorPID[pid] != nil {
		g.MutexProcesosPorPID.Unlock()
		logger.Error("El proceso PID <%d> ya existe", pid)
		return logger.ErrDuplicatePID
	} else {
		g.MutexProcesosPorPID.Unlock()
	}
	nuevoProceso := &g.Proceso{
		PID:                  pid,
		TablaRaiz:            InicializarTablaRaiz(),
		Metricas:             InicializarMetricas(),
		InstruccionesEnBytes: make(map[int][]byte),
	}

	g.MutexProcesosPorPID.Lock()
	g.ProcesosPorPID[pid] = nuevoProceso
	g.MutexProcesosPorPID.Unlock()
	g.MutexMetrica[pid] = new(sync.Mutex)

	if nuevoProceso.TablaRaiz == nil {
		logger.Error("TablaRaiz es nil para proceso PID <%d>", pid)
		return logger.ErrNoTabla
	}

	err = LecturaPseudocodigo(nuevoProceso, nombreArchPseudocodigo)
	if err != nil {
		logger.Error("error al leer pseudocódigo <%d> de error %v", nombreArchPseudocodigo, err)
		return logger.ErrBadRequest
	}

	err = AsignarPaginasParaPID(nuevoProceso, tamanioProceso)
	if err != nil {
		logger.Error("error al asignar espacio al pseudocódigo <%d> de error %v", nombreArchPseudocodigo, err)
		return logger.ErrBadRequest
	}

	logger.Info("## Proceso del PID <%d> instanciado exitosamente", pid)
	logger.Info("## Espacio de usuario del PID <%d> cargado exitosamente", pid)

	return nil
}

// ========== VERIFICO TAMAÑO PARA EL PROCESO ==========

func TieneTamanioNecesario(tamanioProceso int) (resultado bool) {
	framesNecesarios, err := g.CalcularCantidadEntradas(tamanioProceso)
	if err != nil {
		logger.Error("el tamanio pedido de espacio es 0: %v", logger.ErrBadRequest)
		resultado = false
	}

	g.MutexCantidadFramesLibres.Lock()
	resultado = framesNecesarios <= g.CantidadFramesLibres
	g.MutexCantidadFramesLibres.Unlock()

	return
}

// ========== LECTURA DE INSTRUCCIONES ==========

func LecturaPseudocodigo(proceso *g.Proceso, direccionPseudocodigo string) error {
	if direccionPseudocodigo == "" {
		logger.Error("El nombre de pseudocodigo dado es vacío")
		return logger.ErrNoInstance
	}
	ruta := "./scripts/" + direccionPseudocodigo
	file, err := os.Open(ruta)
	if err != nil {
		logger.Error("Error al abrir el archivo: %s\n", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file)

	pc := 0
	for scanner.Scan() {
		lineaEnString := scanner.Text()
		lineaEnBytes := []byte(lineaEnString)

		proceso.InstruccionesEnBytes[pc] = lineaEnBytes
		pc++

		if strings.TrimSpace(lineaEnString) == "EXIT" {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Error("Error al leer el archivo: %s", err)
	}

	IncrementarMetrica(proceso, pc, IncrementarInstruccionesSolicitadas)

	return nil
}

// ========== PAGINAS ==========

func InicializarTablaRaiz() g.TablaPaginas {
	cantidadEntradasPorTabla := g.MemoryConfig.EntriesPerPage
	tabla := make(g.TablaPaginas, cantidadEntradasPorTabla)
	for i := 0; i < cantidadEntradasPorTabla; i++ {
		tabla[i] = &g.TablaPagina{}
	}
	return tabla
}

func AsignarPaginasParaPID(proceso *g.Proceso, tamanio int) error {
	cantidadFrames, err := g.CalcularCantidadEntradas(tamanio)
	if err != nil {
		return err
	}
	for i := 0; i < cantidadFrames; i++ {
		numeroFrame, err := AsignarFrameLibre()
		if err != nil {
			logger.Error("No hay frames libres en el sistema %v", err)
			return err
		}
		entradaPagina := &g.EntradaPagina{
			NumeroFrame:  numeroFrame,
			EstaPresente: true,
		}
		InsertarEntradaPaginaEnTabla(proceso.TablaRaiz, i, entradaPagina)
		logger.Info("## PID <%d> - Entrada <%d> - Frame <%d>", i, proceso.PID, numeroFrame)
	}
	logger.Info("## Quedan <%d> frames libes", g.CantidadFramesLibres)
	return nil
}

func InsertarEntradaPaginaEnTabla(tablaRaiz g.TablaPaginas, numeroPagina int, entrada *g.EntradaPagina) {
	indices := CrearIndicePara(numeroPagina)
	actual := tablaRaiz[indices[0]]

	if actual == nil {
		actual = &g.TablaPagina{}
		tablaRaiz[indices[0]] = actual
	}

	for i := 1; i < len(indices)-1; i++ {
		if actual.Subtabla == nil {
			actual.Subtabla = make(map[int]*g.TablaPagina)
		}
		if actual.Subtabla[indices[i]] == nil {
			actual.Subtabla[indices[i]] = &g.TablaPagina{}
		}
		actual = actual.Subtabla[indices[i]]
	}

	if actual.EntradasPaginas == nil {
		actual.EntradasPaginas = make(map[int]*g.EntradaPagina)
	}
	actual.EntradasPaginas[indices[len(indices)-1]] = entrada

	// logger.Info("Insertada entrada para página <%d> (indices=%v)", numeroPagina, indices)
}

func AsignarFrameLibre() (int, error) {
	cantidadFramesTotales := g.MemoryConfig.MemorySize / g.MemoryConfig.PagSize

	g.MutexEstructuraFramesLibres.Lock()

	for numeroFrame := 0; numeroFrame < cantidadFramesTotales; numeroFrame++ {
		booleano := g.FramesLibres[numeroFrame]

		if booleano == true {
			MarcarOcupadoFrame(numeroFrame)
			g.MutexEstructuraFramesLibres.Unlock()
			return numeroFrame, nil
		}
	}
	g.MutexEstructuraFramesLibres.Unlock()
	return -10, logger.ErrNoMemory

}
