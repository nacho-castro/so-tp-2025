package comunicacion

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

// Body JSON a recibir
type MensajeAMemoria struct {
	Pseudocodigo   string `json:"filename"` //filename
	TamanioMemoria int    `json:"tamanio_memoria"`
	PID            int    `json:"pid"`
}

type ConsultaAMemoria struct {
	Hilo      Hilo        `json:"hilo"`
	Tipo      string      `json:"tipo"`
	Arguments interface{} `json:"argumentos"` // <-- puede ser cualquier tipo ahora (map, struct, etc.)
}

type FinProceso struct {
	PID int `json:"pid"`
}

type Pid int

type Hilo struct {
	PID Pid `json:"pid"`
}

type RespuestaMemoria struct {
	EspacioLibre int `json:"espacio_libre"`
}

// ENVIAR ARCHIVO DE PSEUDOCODIGO Y TAMAÑO
func SolicitarEspacioEnMemoria(fileName string, tamanio int) int {
	url := fmt.Sprintf("http://%s:%d/memoria/espaciolibre", globals.KConfig.MemoryAddress, globals.KConfig.MemoryPort)

	mensaje := MensajeAMemoria{
		Pseudocodigo:   fileName,
		TamanioMemoria: tamanio,
	}

	resp, err := data.EnviarDatosConRespuesta(url, mensaje)
	if err != nil {
		logger.Error("Error enviando pseudocódigo a Memoria: %s", err.Error())
	}
	defer resp.Body.Close()

	var rta RespuestaMemoria
	err = json.NewDecoder(resp.Body).Decode(&rta)
	if err != nil {
		logger.Error("Error al decodificar respuesta de Memoria: %s", err.Error())
	}

	logger.Info("## Memoria dice => Espacio libre: %d", rta.EspacioLibre)
	return rta.EspacioLibre
}

// ENVIAR ARCHIVO DE PSEUDOCODIGO Y TAMAÑO
func EnviarArchivoMemoria(fileName string, tamanio int, pid int) {
	url := fmt.Sprintf("http://%s:%d/memoria/inicializacionProceso", globals.KConfig.MemoryAddress, globals.KConfig.MemoryPort)

	mensaje := MensajeAMemoria{
		Pseudocodigo:   fileName,
		TamanioMemoria: tamanio,
		PID:            pid,
	}

	err := data.EnviarDatos(url, mensaje)
	if err != nil {
		logger.Error("Error enviando pseudocódigo a Memoria: %s", err.Error())
	}
}

func LiberarMemoria(pid int) {
	url := fmt.Sprintf("http://%s:%d/memoria/finalizacionProceso", globals.KConfig.MemoryAddress, globals.KConfig.MemoryPort)

	mensaje := FinProceso{
		PID: pid,
	}

	err := data.EnviarDatos(url, mensaje)
	if err != nil {
		logger.Error("Error enviando EXIT a Memoria: %s", err.Error())
	}
}

// PARA MANJERAR LOS MENSAJES DEL ENDPOINT QUE ESTAN EN MEMORIA
// por ejemplo: http.HandleFunc("/kernel/createProcess", createProcess)
const (
	CreateProcess = "createProcess"
	FinishProcess = "finishProcess"
	MemoryDump    = "memoryDump"
)

var ErrorRequestType = map[string]error{
	CreateProcess: errors.New("memoria: No hay espacio disponible en memoria "),
	FinishProcess: errors.New("memoria: No se puedo finalizar el proceso"),
}

type inicializacionRequest struct {
	Filename string `json:"filename"`
	Tamanio  int    `json:"tamanio_memoria"`
	PID      int    `json:"pid"`
}

type inicializacionResponse struct {
	Exito   bool   `json:"exito"`
	Mensaje string `json:"mensaje"`
}
type PedidoKernel struct {
	PID int `json:"pid"`
}
type RespuestaMemoriaSWAP struct {
	Exito   bool   `json:"exito"`
	Mensaje string `json:"mensaje"`
}

func SolicitarSuspensionEnMemoria(pid int) error {
	url := fmt.Sprintf("http://%s:%d/memoria/suspension",
		globals.KConfig.MemoryAddress, globals.KConfig.MemoryPort)

	req := PedidoKernel{PID: pid}
	body, _ := json.Marshal(req)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Error("Error enviando suspensión a Memoria: %v", err)
		return err
	}
	defer resp.Body.Close()

	var r RespuestaMemoriaSWAP
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		logger.Error("Error decodificando respuesta de Memoria: %v", err)
		return err
	}
	if !r.Exito {
		logger.Warn("Memoria rechazó suspensión PID=%d: %s", pid, r.Mensaje)
		return fmt.Errorf("memoria: %s", r.Mensaje)
	}
	return nil
}

func DesuspensionMemoria(pid int) error {
	url := fmt.Sprintf("http://%s:%d/memoria/desuspension", globals.KConfig.MemoryAddress, globals.KConfig.MemoryPort)

	mensaje := PedidoKernel{
		PID: pid,
	}

	var r RespuestaMemoriaSWAP
	err := data.EnviarDatosYRecibirRespuesta(url, mensaje, &r)
	if err != nil {
		logger.Error("Error enviando DESUSPENSION a Memoria: %s", err.Error())
		return err
	}

	if !r.Exito {
		logger.Warn("Memoria rechazó desuspensión PID=%d: %s", pid, r.Mensaje)
		return fmt.Errorf("memoria: %s", r.Mensaje)
	}

	return nil
}
