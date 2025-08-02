package conexiones

import (
	"encoding/json"
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

func ObtenerEspacioLibreHandler(w http.ResponseWriter, r *http.Request) {

	g.MutexCantidadFramesLibres.Lock()
	espacioLibre := g.CantidadFramesLibres * g.MemoryConfig.PagSize
	g.MutexCantidadFramesLibres.Unlock()

	respuesta := g.RespuestaEspacioLibre{EspacioLibre: espacioLibre}

	logger.Info("## Espacio libre devuelto - Tamaño: <%d>", respuesta.EspacioLibre)

	err := json.NewEncoder(w).Encode(respuesta)
	if err != nil {
		logger.Error("Error al enviar la respuesta de esapcioLibre: %v", err)
		http.Error(w, "Error al enviar la respuesta de esapcioLibre", http.StatusInternalServerError)
		return
	}

}

func EnviarConfiguracionMemoriaHandler(w http.ResponseWriter, r *http.Request) {
	var pidData struct {
		PID int `json:"pid"`
	}

	err := data.LeerJson(w, r, &pidData)
	if err != nil {
		logger.Error("Error leyendo JSON del CPU\n", err)
		http.Error(w, "Error al decodear mensaje del JSON", http.StatusInternalServerError)
		return
	}
	// logger.Info("Recibí petición de configuración desde PID: %d", pidData.PID)

	mensaje := g.RespuestaConfigMemoria{
		TamanioPagina:    g.MemoryConfig.PagSize,
		EntradasPorNivel: g.MemoryConfig.EntriesPerPage,
		CantidadNiveles:  g.MemoryConfig.NumberOfLevels,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mensaje); err != nil {
		logger.Error("Error al codificar la respuesta JSON: %v", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
	}
}
