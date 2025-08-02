package administracion

import g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"

// ========== METRICAS ==========

func InicializarMetricas() (metricas g.MetricasProceso) {
	metricas = g.MetricasProceso{
		AccesosTablasPaginas:     0,
		InstruccionesSolicitadas: 0,
		BajadasSwap:              0,
		SubidasMP:                0,
		LecturasDeMemoria:        0,
		EscriturasDeMemoria:      0,
	}
	return
}

func IncrementarMetrica(proceso *g.Proceso, cantidad int, funcMetrica g.OperacionMetrica) {
	g.MutexMetrica[proceso.PID].Lock()
	funcMetrica(&proceso.Metricas, cantidad)
	g.MutexMetrica[proceso.PID].Unlock()
}
func IncrementarAccesosTablasPaginas(metrica *g.MetricasProceso, cantidad int) {
	metrica.AccesosTablasPaginas += cantidad
}
func IncrementarInstruccionesSolicitadas(metrica *g.MetricasProceso, cantidad int) {
	metrica.InstruccionesSolicitadas += cantidad
}
func IncrementarBajadasSwap(metrica *g.MetricasProceso, cantidad int) {
	metrica.BajadasSwap += cantidad
}
func IncrementarSubidasMP(metrica *g.MetricasProceso, cantidad int) {
	metrica.SubidasMP += cantidad
}
func IncrementarLecturaDeMemoria(metrica *g.MetricasProceso, cantidad int) {
	metrica.LecturasDeMemoria += cantidad
}
func IncrementarEscrituraDeMemoria(metrica *g.MetricasProceso, cantidad int) {
	metrica.EscriturasDeMemoria += cantidad
}
