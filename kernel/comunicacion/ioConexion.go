package comunicacion

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
	"strings"
	"time"
)

// Body JSON a recibir
type MensajeDeIO struct {
	Nombre string `json:"nombre"`
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
}

type MensajeAIO struct {
	Pid      int `json:"pid"`
	Duracion int `json:"duracion"` //en segundos
}

type MensajeFin struct {
	PID         int    `json:"pid"`
	Desconexion bool   `json:"desconexion"`
	Nombre      string `json:"nombre"`
	Puerto      int    `json:"puerto"`
}

// w http.ResponseWriter. Se usa para escribir la respuesta al Cliente
// r *http.Request es la peticion que se recibio
func RecibirMensajeDeIO(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeDeIO
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	tipo := mensajeRecibido.Nombre

	globals.IOMu.Lock()
	instancia := globals.DatosIO{
		Tipo:    mensajeRecibido.Nombre,
		Ip:      mensajeRecibido.Ip,
		Puerto:  mensajeRecibido.Puerto,
		Ocupada: false,
	}

	// Agrega a la lista correspondiente
	globals.IOs[tipo] = append(globals.IOs[tipo], instancia)
	globals.IOMu.Unlock()

	logger.Info("## Se ha recibido IO: Nombre: %s Ip: %s Puerto: %d",
		instancia.Tipo, instancia.Ip, instancia.Puerto)

	Utils.NotificarIOLibre <- Utils.IOEvent{
		PID:    -1,
		Nombre: mensajeRecibido.Nombre,
		Puerto: mensajeRecibido.Puerto,
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("STATUS OK"))
}

func EnviarContextoIO(instanciaIO globals.DatosIO, pid int, duracion int) {
	url := fmt.Sprintf("http://%s:%d/io/kernel", instanciaIO.Ip, instanciaIO.Puerto)
	mensaje := MensajeAIO{
		Pid:      pid,
		Duracion: duracion,
	}

	logger.Info("## (<%d>) - Bloqueado por IO: %s", pid, instanciaIO.Tipo)

	const maxIntentos = 3
	const backoff = 200 * time.Millisecond

	var err error
	for intento := 1; intento <= maxIntentos; intento++ {
		err = data.EnviarDatos(url, mensaje)
		if err == nil {
			// Envío exitoso
			return
		}

		// Si el error es "connection refused" o "connection reset" (IO desconectada)
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "forcibly closed") {
			// Log leve para debug, NO error fatal
			logger.Warn("No se pudo conectar a IO %s (puerto %d), está desconectada o no responde (intento %d/%d)", instanciaIO.Tipo, instanciaIO.Puerto, intento, maxIntentos)
			time.Sleep(backoff)
			continue
		}

		// Otros errores menos comunes los logueamos también pero no fatales
		logger.Warn("Error enviando a IO: %v", err)
		return
	}

	// Luego de reintentos no exitosos, solo log leve
	logger.Warn("No se pudo conectar a IO %s tras %d intentos, se ignorará temporalmente", instanciaIO.Tipo, maxIntentos)
	return
}

// Al momento de recibir un mensaje de una IO se deberá verificar
// que el mismo sea una confirmación de fin de IO, en caso afirmativo,
// se deberá validar si hay más procesos esperando realizar dicha IO.
// En caso de que el mensaje corresponda a una desconexión de la IO,
// el proceso que estaba ejecutando en dicha IO, se deberá pasar al estado EXIT.

func RecibirFinDeIO(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeFin
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		logger.Error("Error al recibir mensaje de fin de IO: %s", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	evt := Utils.IOEvent{
		PID:    mensajeRecibido.PID,
		Nombre: mensajeRecibido.Nombre,
		Puerto: mensajeRecibido.Puerto,
	}

	//  Si es desconexión
	if mensajeRecibido.Desconexion {
		logger.Info("## Desconexión de IO: %s - Puerto %d - PID %d", evt.Nombre, evt.Puerto, evt.PID)
		Utils.NotificarDesconexion <- evt
	} else {
		//Fin de IO normal
		logger.Debug("FIN de IO: %s - PID %d", evt.Nombre, evt.PID)
		globals.IOMu.Lock()
		for i := range globals.IOs[evt.Nombre] {
			if evt.Puerto == globals.IOs[evt.Nombre][i].Puerto {
				globals.IOs[evt.Nombre][i].Ocupada = false
				globals.IOs[evt.Nombre][i].PID = -1
				break
			}
		}
		globals.IOMu.Unlock()
	}

	//1. Aviso a DespacharIO (Mediano Plazo)
	//logger.Debug("AVISO A DESPACHADOR, PID: <%d>", evt.PID)
	Utils.NotificarIOLibre <- evt

	//2. Reenvío al canal individual de este PID (si existe)
	Utils.MutexIOFinishedWaiters.Lock()
	ch, exists := Utils.FinIOWaiters[evt.PID]
	Utils.MutexIOFinishedWaiters.Unlock()

	if !exists {
		//logger.Warn("RecibirFinDeIO: no hay canal FinIOWaiters para PID %d", evt.PID)
	} else {
		select {
		case ch <- evt:
			//logger.Debug("ENVIE FIN DE IO A CANAL DE %d", evt.PID)
		default:
			logger.Debug("RecibirFinDeIO: canal FinIOWaiters[%d] lleno, descartando evento", evt.PID)
		}
	}
	//Utils.FinIODesdeSuspBlocked <- Utils.IOEvent{PID: evt.PID, Nombre: evt.Nombre, Puerto: evt.Puerto}

	w.WriteHeader(http.StatusOK)
}
