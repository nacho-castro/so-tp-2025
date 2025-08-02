package administracion

import (
	"fmt"
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"os"
)

func RealizarDumpMemoria(pid int) (vector []string, err error) {
	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	if proceso == nil {
		logger.Error("No existe el proceso solicitado para DUMP")
		return vector, logger.ErrProcessNil
	}

	entradas := RecolectarEntradasProcesoDump(*proceso)
	tamanioPagina := g.MemoryConfig.PagSize

	vector = make([]string, len(g.FramesLibres))

	for i := 0; i < len(entradas); i++ {
		numeroFrame := entradas[i]
		inicio := numeroFrame * tamanioPagina
		fin := inicio + tamanioPagina

		if fin > len(g.MemoriaPrincipal) {
			logger.Error("Acceso fuera de rango al hacer dump del frame %d con PID: %d", numeroFrame, pid)
			fin = len(g.MemoriaPrincipal) - 1
			continue
		}
		copiado := make([]byte, fin-inicio)
		g.MutexMemoriaPrincipal.Lock()
		copy(copiado, g.MemoriaPrincipal[inicio:fin])
		g.MutexMemoriaPrincipal.Unlock()

		copiado = g.RecortarNulosFinales(copiado)

		resul := fmt.Sprintf("Direccion Fisica: %d | Frame: %d | Datos: %q\n", inicio, numeroFrame, string(copiado))

		vector[numeroFrame] = resul
	}

	return
}

//  ========== JUNTAR ENTRADAS DEL PID PARA EL DUMP ==========

func RecolectarEntradasProcesoDump(proceso g.Proceso) (resultados []int) {
	for _, subtabla := range proceso.TablaRaiz {
		RecorrerTablaPaginaDump(subtabla, &resultados)
	}
	return
}

// ========== RECORRER LA TABLA PARA JUNTAR ENTRADAS ==========

func RecorrerTablaPaginaDump(tabla *g.TablaPagina, resultados *[]int) {

	if tabla.Subtabla != nil {
		for _, subTabla := range tabla.Subtabla {
			RecorrerTablaPaginaDump(subTabla, resultados)
		}
		return
	}
	for i, entrada := range tabla.EntradasPaginas {
		if entrada.EstaPresente {
			*resultados = append(*resultados, entrada.NumeroFrame)
		} else {
			logger.Debug("%% Entrada <%d> en swap !", i)
		}
	}
}

// ========== COLOCO EL CONTENIDO EN EL .DMP ==========

func ParsearContenido(dumpFile *os.File, pid int, contenido []string) {
	comienzo := fmt.Sprintf("## Dump De Memoria Para PID: %d\n\n", pid)
	_, err := dumpFile.WriteString(comienzo)
	for i := 0; i < len(contenido); i++ {
		_, err = dumpFile.WriteString(contenido[i])
		if err != nil {
			logger.Error("Error al escribir contenido en el archivo dump: %v", err)
		}
	}
}
