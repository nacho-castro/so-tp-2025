package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sisoputnfrba/tp-golang/io/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

var pid = -1

type MensajeAKernel struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
	Nombre string `json:"nombre"`
}

type MensajeDeKernel struct {
	PID      int `json:"pid"`
	Duracion int `json:"duracion"` // en milisegundos
}

type MensajeFin struct {
	PID         int    `json:"pid"`
	Desconexion bool   `json:"desconexion"`
	Nombre      string `json:"nombre"`
	Puerto      int    `json:"puerto"`
}

func main() {
	// ----------------------------------------------------
	// -------------------- CONFIG ------------------------
	// ----------------------------------------------------
	globals.IoConfig = globals.CargarConfig()

	// ----------------------------------------------------
	// --------------------- LOGGER -----------------------
	// ----------------------------------------------------
	logFileName := fmt.Sprintf("./logs/%s_%d.log", globals.IoConfig.Type, globals.IoConfig.PortIo)
	err := logger.ConfigureLogger(logFileName, "INFO")
	if err != nil {
		fmt.Println("No se pudo crear el logger -", err.Error())
		os.Exit(1)
	}
	logger.Debug("Logger creado")
	logger.Info("Comenzó ejecución del IO")
	logger.Info("Nombre de la interfaz IO: %s", globals.IoConfig.Type)
	logger.Info("Puerto asignado: %d", globals.IoConfig.PortIo)

	// ----------------------------------------------------
	// -------- ENVÍO IP/PUERTO/NOMBRE A KERNEL -----------
	// ----------------------------------------------------
	mensaje := MensajeAKernel{
		Ip:     globals.IoConfig.IpIo,
		Puerto: globals.IoConfig.PortIo,
		Nombre: globals.IoConfig.Type,
	}

	// ----------------------------------------------------
	// ------------- MANEJO DE SEÑALES --------------------
	// ----------------------------------------------------
	desconexion()

	// ----------------------------------------------------
	// ------------------ HANDLERS HTTP -------------------
	// ----------------------------------------------------
	mux := http.NewServeMux()
	mux.HandleFunc("/io/kernel", RecibirMensajeDeKernel)

	// ----------------------------------------------------
	// ------------------ SERVIDOR HTTP -------------------
	// ----------------------------------------------------

	//Lanzar goroutine para enviar al Kernel
	go func() {
		// Podés esperar brevemente o directamente intentar
		time.Sleep(300 * time.Millisecond)
		logger.Info("Enviando mensaje handshake a Kernel...")
		EnviarIpPuertoNombreAKernel(globals.IoConfig.IpKernel, globals.IoConfig.PortKernel, mensaje)
	}()

	direccion := fmt.Sprintf("%s:%d", globals.IoConfig.IpIo, globals.IoConfig.PortIo)
	logger.Info("Escuchando en %s", direccion)
	err = http.ListenAndServe(direccion, mux)
	if err != nil {
		logger.Fatal("Error al iniciar el servidor HTTP: %v", err)
	}
}

// Enviar IP y Puerto al Kernel
func EnviarIpPuertoNombreAKernel(ipDestino string, puertoDestino int, mensaje MensajeAKernel) {
	// Construye la URL del endpoint (url + path) a donde se va a enviar el mensaje
	url := fmt.Sprintf("http://%s:%d/kernel/io", ipDestino, puertoDestino)

	// Hace el POST al Kernel
	err := data.EnviarDatos(url, mensaje)
	// Verifico si hubo error y logueo si lo hubo
	if err != nil {
		logger.Error("Error enviando mensaje: %s", err.Error())
		return
	}
	// Si no hubo error, logueo que salió bien
	logger.Info("Mensaje enviado a Kernel")
}

// Al momento de recibir una petición del Kernel,
// el módulo deberá iniciar un usleep
// por el tiempo indicado en la request.
func RecibirMensajeDeKernel(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeDeKernel
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	pid = mensajeRecibido.PID
	//Realizo la operacion
	logger.Info("## PID: <%d> - Inicio de IO - Tiempo: %d", mensajeRecibido.PID, mensajeRecibido.Duracion)
	time.Sleep(time.Duration(mensajeRecibido.Duracion) * time.Millisecond)

	//IO Finalizada
	logger.Info("## PID: <%d> - Fin de IO", mensajeRecibido.PID)
	url := fmt.Sprintf("http://%s:%d/kernel/fin_io", globals.IoConfig.IpKernel, globals.IoConfig.PortKernel)

	mensaje := MensajeFin{
		PID:         mensajeRecibido.PID,
		Desconexion: false,
		Nombre:      globals.IoConfig.Type,
		Puerto:      globals.IoConfig.PortIo,
	}
	logger.Debug("Enviando PID <%d> a Kernel", mensaje.PID)

	err := data.EnviarDatos(url, mensaje)
	if err != nil {
		logger.Info("Error enviando PID y Nombre a Kernel: %s", err.Error())
		return
	}
	pid = -1
}

// Al finalizar deberá informar al Kernel que finalizó la solicitud de I/O
// quedará a la espera de la siguiente petición.
// ver el tema de FINALIZACION DE IO != FIN de timer de IO

//*El Módulo IO, deberá notificar al Kernel de su desconexion,
//para esto se deberá implementar el manejo de las señales SIGINT y SIGTERM,
//para enviar la notificación y finalizar de manera controlada.

// Mecanismo del SO para notificar a un proceso que debe hacer algo.
// En este caso se entera cuando IO muere o presionan ctrl c / kill
func desconexion() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Info("Se recibió una señal de finalización. Cerrando IO...")

		//Notificar al Kernel
		mensaje := MensajeFin{
			PID:         pid,
			Desconexion: true,
			Nombre:      globals.IoConfig.Type,
			Puerto:      globals.IoConfig.PortIo,
		}

		url := fmt.Sprintf("http://%s:%d/kernel/fin_io", globals.IoConfig.IpKernel, globals.IoConfig.PortKernel)
		const maxIntentos = 3
		const backoff = 200 * time.Millisecond

		var err error
		for intento := 1; intento <= maxIntentos; intento++ {
			err = data.EnviarDatos(url, mensaje)
			if err == nil {
				// Éxito: salimos
				logger.Info("======== Final de Ejecución IO ========")
				os.Exit(0)
				return
			}
			// Si no es connection refused, no tiene sentido reintentar
			if !strings.Contains(err.Error(), "connection refused") {
				logger.Warn("Error enviando a IO (no recoverable): %v", err)
				return
			}
			// es Connection refused: esperamos y reintentamos
			logger.Warn("Intento %d: connection refused al enviar a IO, reintentando en %s...", intento, backoff)
			time.Sleep(backoff)
		}
	}()
}
