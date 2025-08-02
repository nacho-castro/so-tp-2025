package syscalls

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"log"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
)

// Body JSON a recibir
type MensajeInit struct {
	PID      int    `json:"pid"`
	PC       int    `json:"pc"`
	Filename string `json:"filename"`
	Tamanio  int    `json:"tamanio_memoria"`
}

type MensajeIo struct {
	PID      int    `json:"pid"`
	PC       int    `json:"pc"`
	Duracion int    `json:"tiempo"`
	Nombre   string `json:"nombre"`
	ID       string `json:"id"`
}

type MensajeSyscall struct {
	PID int    `json:"pid"`
	PC  int    `json:"pc"`
	ID  string `json:"id"`
}

type MensajeInterrupt struct {
	PID int    `json:"pid"`
	PC  int    `json:"pc"`
	ID  string `json:"id"`
}

type MensajeDUMP struct {
	PID       int    `json:"pid"`
	Timestamp string `json:"timestamp"`
}

// RECIBIR RESPUESTA DE CPU
func ContextoInterrumpido(w http.ResponseWriter, r *http.Request) {
	// 1) Leer y parsear el JSON entrante
	//se recibe el PID
	//y el Program Counter (PC) actualizado
	//con motivo de la interrupción.
	var msg MensajeInterrupt
	if err := data.LeerJson(w, r, &msg); err != nil {
		return
	}

	logger.Debug("<%d> Se ha recibido contexto de Interrupcion. CPU <%s>", msg.PID, msg.ID)
	//SIGNAL A Planif. CORTO PLAZO QUE SE INTERRUMPIO
	go func(p int) {
		Utils.ContextoInterrupcion <- Utils.InterruptProcess{
			PID:   msg.PID,
			PC:    msg.PC,
			CpuID: msg.ID,
		}
	}(msg.PID)

	//Responder de inmediato
	w.WriteHeader(http.StatusOK)
}

func InitProcess(w http.ResponseWriter, r *http.Request) {
	// 1) Leer y parsear el JSON entrante (sin usar PID desde la CPU)
	var msg MensajeInit
	if err := data.LeerJson(w, r, &msg); err != nil {
		return
	}

	// 2) Generar el PID dentro del kernel
	pid := globals.GenerarNuevoPID()
	fileName := msg.Filename
	tamanio := msg.Tamanio

	estimado := globals.KConfig.InitialEstimate

	logger.Info("## (<%d>) - Solicitó syscall: <INIT_PROC>", msg.PID)
	//logger.Info("Se ha recibido: Filename: %s Tamaño Memoria: %d", fileName, tamanio)

	// 3) Despachar la creación en segundo plano
	go func() {
		// Construir el PCB con el PID generado
		pcbNuevo := pcb.PCB{
			PID:            pid,
			PC:             0,
			FileName:       fileName,
			ProcessSize:    tamanio,
			ME:             make(map[string]int),
			MT:             make(map[string]float64),
			EstimadoRafaga: estimado,
			TiempoEstado:   time.Now(),
			CpuID:          "",
		}
		logger.Info("## (<%d>) Se crea el proceso - Estado: NEW", pid)

		// Encolar en NEW
		Utils.MutexNuevo.Lock()
		algoritmos.ColaNuevo.Add(&pcbNuevo)
		pcb.CambiarEstado(&pcbNuevo, pcb.EstadoNew)
		Utils.MutexNuevo.Unlock()

		//planificadores.MostrarColaNew()
		//logger.Debug("Notificado planificador de largo plazo para PID <%d>", pid)

		// Notificar al Planificador de Largo Plazo
		Utils.ChannelProcessArguments <- Utils.NewProcess{
			PID:      pid,
			Filename: fileName,
			Tamanio:  tamanio,
		}
	}()

	//<-Utils.SemProcessCreateOK

	// 4) Responder de inmediato
	w.WriteHeader(http.StatusOK)
}

func Exit(w http.ResponseWriter, r *http.Request) {
	var msg MensajeSyscall
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	pid := msg.PID
	pc := msg.PC
	cpuid := msg.ID
	logger.Info("## (<%d>) - Solicitó syscall: <EXIT>", pid)

	// Despachamos la señal en segundo plano para no bloquear el handler HTTP
	go func(p int) {
		Utils.ChannelFinishprocess <- Utils.FinishProcess{
			PID:   pid,
			PC:    pc,
			CpuID: cpuid,
		}
	}(pid)

	// Respondemos de inmediato
	w.WriteHeader(http.StatusOK)
}

func DumpMemory(w http.ResponseWriter, r *http.Request) {
	var mensaje MensajeDUMP
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return // El error ya fue respondido por LeerJson
	}

	// Generar el timestamp en formato yyyyMMddTHHmmss
	timestamp := time.Now().Format("20060102T150405")

	// Crear mensaje para enviar a Memoria
	req := MensajeDUMP{
		PID:       mensaje.PID,
		Timestamp: timestamp,
	}

	// Armar URL del módulo Memoria
	url := fmt.Sprintf("http://%s:%d/memoria/dump", globals.KConfig.MemoryAddress, globals.KConfig.MemoryPort)

	// Usar helper para enviar datos
	if err := data.EnviarDatos(url, req); err != nil {
		log.Printf("Error enviando dump a Memoria: %v", err)
		http.Error(w, "Error comunicando con Memoria", http.StatusInternalServerError)
		return
	}

	log.Printf("Se envió correctamente el pedido de dump del PID %d a Memoria", req.PID)
	w.WriteHeader(http.StatusOK)
}

func Io(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeIo
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return
	}
	pid := mensajeRecibido.PID
	pc := mensajeRecibido.PC
	tipoIO := mensajeRecibido.Nombre

	logger.Info("## (<%d>) - Solicitó syscall: <IO>”", pid)
	logger.Debug("Nombre IO: %s Duracion: %d", tipoIO, mensajeRecibido.Duracion)

	//SIGNAL A Planif. CORTO PLAZO QUE LLEGO I/O
	go func(p int) {
		Utils.NotificarComienzoIO <- Utils.MensajeIOChannel{
			PID:      pid,
			PC:       pc,
			Nombre:   mensajeRecibido.Nombre,
			Duracion: mensajeRecibido.Duracion,
			CpuID:    mensajeRecibido.ID,
		}
	}(pid)

	w.WriteHeader(http.StatusOK)
}

// ESTE ES EL ENDPOINT PARA QUE MEMORIA AVISE QUE PUDO HACER T O D O EL LABURO DE SWAP
func ConfirmarSuspensionHandler(w http.ResponseWriter, r *http.Request) {
	var msg comunicacion.PedidoKernel
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	logger.Info("## (<%d>) - Confirmación de suspensión recibida desde Memoria", msg.PID)

	// Signal al Planificador de Mediano Plazo: ya puede reintentar NEW/SUSP.READY
	go func(pid int) {
		Utils.NotificarTimeoutBlocked <- pid
	}(msg.PID)

	w.WriteHeader(http.StatusOK)
}
