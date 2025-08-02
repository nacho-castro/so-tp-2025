package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"io"
	"os"
)

func RecolectarEntradasParaSwap(pid int) (entradas map[int]g.EntradaSwap) {
	var contador *int
	contador = new(int)

	entradas = make(map[int]g.EntradaSwap)

	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	for i := 0; i < len(proceso.TablaRaiz); i++ {
		RecorrerTablaPaginasParaSwap(proceso.TablaRaiz[i], entradas, contador)
	}

	return
}

func RecorrerTablaPaginasParaSwap(tabla *g.TablaPagina, entradasSwap map[int]g.EntradaSwap, contador *int) {
	if tabla.Subtabla != nil {
		for i := 0; i < len(tabla.Subtabla); i++ {
			RecorrerTablaPaginasParaSwap(tabla.Subtabla[i], entradasSwap, contador)
		}
		return
	}
	if tabla.EntradasPaginas != nil {
		for i := 0; i < g.MemoryConfig.EntriesPerPage; i++ {
			tamanioPagina := g.MemoryConfig.PagSize

			entrada := tabla.EntradasPaginas[i]
			if entrada == nil {
				continue
			}

			inicio := entrada.NumeroFrame * tamanioPagina
			fin := inicio + tamanioPagina

			datos := make([]byte, tamanioPagina)

			g.MutexMemoriaPrincipal.Lock()
			copy(datos, g.MemoriaPrincipal[inicio:fin])
			copy(g.MemoriaPrincipal[inicio:fin], make([]byte, tamanioPagina))
			g.MutexMemoriaPrincipal.Unlock()

			datos = g.RecortarNulosFinales(datos)

			entraditaNueva := g.EntradaSwap{
				NumeroPagina: *contador,
				Datos:        datos,
			}

			MarcarLibreFrame(entrada.NumeroFrame)
			entrada.EstaPresente = false

			entradasSwap[entraditaNueva.NumeroPagina] = entraditaNueva
			*contador++
		}
	}

}

func CargarEntradasASwap(pid int, entradas map[int]g.EntradaSwap) error {

	archivo, err := os.OpenFile(g.MemoryConfig.SwapfilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func(archivo *os.File) {
		err := archivo.Close()
		if err != nil {
			logger.Error("Error al cerrar: %v", err)
		}
	}(archivo)

	_, errSeek := archivo.Seek(0, io.SeekEnd) // siempre se apunta al final del archivo! (pense que no, alto bobo)
	if errSeek != nil {
		logger.Error("Error al setear el puntero para SWAP: %v", errSeek)
		return errSeek
	}

	var info = &g.SwapProcesoInfo{
		Entradas:             make(map[int]*g.EntradaSwapInfo),
		NumerosDePaginas:     make([]int, 0),
		InstruccionesEnBytes: make(map[int][]byte),
	}

	g.MutexProcesosPorPID.Lock()
	info.InstruccionesEnBytes = g.ProcesosPorPID[pid].InstruccionesEnBytes
	g.ProcesosPorPID[pid].InstruccionesEnBytes = nil
	g.MutexProcesosPorPID.Unlock()

	punteroLocal := g.PunteroSwap

	for i := 0; i < len(entradas); i++ {
		entrada := entradas[i]

		posicionPunteroArchivo := punteroLocal
		longitudEscrito := len(entrada.Datos)

		_, errWrite := archivo.Write(entrada.Datos)
		if errWrite != nil {
			logger.Error("Error al escribir el archivo: %v", errWrite)
			return errWrite
		}

		punteroLocal += longitudEscrito

		info.Entradas[entrada.NumeroPagina] = &g.EntradaSwapInfo{
			NumeroPagina:   entrada.NumeroPagina,
			Tamanio:        longitudEscrito,
			PosicionInicio: posicionPunteroArchivo,
		}
		info.NumerosDePaginas = append(info.NumerosDePaginas, entrada.NumeroPagina)

		// VerificarLecturaDesdeSwap(archivo, posicionPunteroArchivo, longitudEscrito)

		logger.Info("## PID: <%d> - <MEMORIA A SWAP> - Entrada: <%d> - Posición en SWAP: <%d> - Tamaño: <%d>",
			pid,
			entrada.NumeroPagina,
			posicionPunteroArchivo,
			longitudEscrito,
		)
		g.PunteroSwap += longitudEscrito
	}
	g.MutexSwapIndex.Lock()
	g.SwapIndex[pid] = info
	g.MutexSwapIndex.Unlock()

	return nil
}
