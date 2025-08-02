package traducciones

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"io"
)

// Estructura del mensaje hacia Memoria
type MensajeTabla struct {
	PID            int   `json:"pid"`
	IndicesEntrada []int `json:"indices_entrada"`
}

// Respuesta de Memoria
type RespuestaTabla struct {
	NumeroMarco int `json:"numero_marco"`
}

type MensajeEscritura struct {
	PID             int    `json:"pid"`
	DireccionFisica int    `json:"direccion_fisica"`
	DatosAEscribir  string `json:"datos_a_escribir"`
}

type MensajeLectura struct {
	PID              int `json:"pid"`
	DireccionFisica  int `json:"direccion_fisica"`
	TamanioARecorrer int `json:"tamanio_a_recorrer"`
}

type RespuestaLectura struct {
	Exito      error  `json:"exito"`
	ValorLeido string `json:"valor_leido"`
}

var Tlb *TLB

func InitTLB() {
	Tlb = NuevaTLB()
}

// Función principal de traducción
func Traducir(dirLogica int) int {
	tamPagina := globals.TamanioPagina
	entradasPorNivel := globals.EntradasPorNivel
	niveles := globals.CantidadNiveles

	nroPagina := dirLogica / tamPagina
	desplazamiento := dirLogica % tamPagina

	// Consulto la TLB
	if marco, ok := Tlb.Buscar(nroPagina); ok {
		return (marco * tamPagina) + desplazamiento
	}

	// La página no está en la TLB, voy a Memoria
	entradas := DescomponerPagina(nroPagina, niveles, entradasPorNivel)
	marco, err := accederTabla(globals.PIDActual, entradas)
	if err != nil {
		logger.Error("No se pudo acceder a la tabla de páginas: %s", err.Error())
		return -1
	}
	if marco == -1 {
		logger.Error("No se pudo traducir la dirección lógica %d", dirLogica)
		return -1
	}

	// Agrego entrada a la TLB
	Tlb.AgregarEntrada(nroPagina, marco)
	logger.Info("PID: %d - OBTENER MARCO - Pagina: %d - Marco: %d", globals.PIDActual, nroPagina, marco)

	return (marco * tamPagina) + desplazamiento
}

// Descompone el número de página en los índices para cada nivel (cortesía de PP)
func DescomponerPagina(nroPagina int, niveles int, entradasPorNivel int) []int {
	entradas := make([]int, niveles)
	divisor := 1

	for i := niveles - 1; i >= 0; i-- {
		entradas[i] = (nroPagina / divisor) % entradasPorNivel
		divisor *= entradasPorNivel
	}

	return entradas
}

// Hace una petición HTTP a Memoria para resolver una entrada de tabla
func accederTabla(pid int, indices []int) (int, error) {
	url := fmt.Sprintf("http://%s:%d/memoria/tabla",
		globals.ClientConfig.IpMemory,
		globals.ClientConfig.PortMemory,
	)

	mensaje := MensajeTabla{
		PID:            pid,
		IndicesEntrada: indices,
	}

	resp, err := data.EnviarDatosConRespuesta(url, mensaje)
	if err != nil {
		return -1, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Error al cerrar: %v", err)
		}
	}(resp.Body)

	var respuesta RespuestaTabla
	if err := json.NewDecoder(resp.Body).Decode(&respuesta); err != nil {
		return -1, err
	}
	//logger.Info("Marco Recibido: %d", respuesta.NumeroMarco)

	return respuesta.NumeroMarco, nil
}

func LeerEnMemoria(dirFisica int, tamanio int) (string, error) {
	msg := MensajeLectura{
		PID:              globals.PIDActual,
		DireccionFisica:  dirFisica,
		TamanioARecorrer: tamanio,
	}

	url := fmt.Sprintf("http://%s:%d/memoria/lectura",
		globals.ClientConfig.IpMemory,
		globals.ClientConfig.PortMemory,
	)

	resp, err := data.EnviarDatosConRespuesta(url, msg)
	if err != nil {
		logger.Error("Error enviando Direccion Fisica y Tamanio: %s", err.Error())
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Error al cerrar: %v", err)
		}
	}(resp.Body)

	var respuesta RespuestaLectura
	if err := json.NewDecoder(resp.Body).Decode(&respuesta); err != nil {
		logger.Error("Error decodificando respuesta de memoria: %s", err.Error())
		return "", err
	}

	//logger.Info("Direccion Fisica: %d y Tamanio: %d enviados correctamente a memoria", dirFisica, tamanio)
	//logger.Info("############################Valor leído: %s", respuesta.ValorLeido)

	return respuesta.ValorLeido, nil
}

func EscribirEnMemoria(dirFisica int, datos string) error {
	msg := MensajeEscritura{
		PID:             globals.PIDActual,
		DireccionFisica: dirFisica,
		DatosAEscribir:  datos,
	}

	url := fmt.Sprintf("http://%s:%d/memoria/escritura",
		globals.ClientConfig.IpMemory, globals.ClientConfig.PortMemory)

	err := data.EnviarDatos(url, msg)
	if err != nil {
		logger.Error("Error enviando Direccion Fisica y Datos: %s", err.Error())
		return err
	}

	logger.Info("PID: %d - Accion: ESCRIBIR - Dirección Física: %d - Datos: %s", globals.PIDActual, dirFisica, datos)
	return nil
}
