package instrucciones

import (
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/traducciones"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"strconv"
	"strings"
)

type Instruccion func(arguments []string) error

type MensajeDump struct {
	PID int    `json:"pid"`
	PC  int    `json:"pc"`
	ID  string `json:"id"`
}

type MensajeIO struct {
	PID    int    `json:"pid"`
	PC     int    `json:"pc"`
	Tiempo int    `json:"tiempo"`
	Nombre string `json:"nombre"`
	ID     string `json:"id"`
}

type MensajeInitProc struct {
	PID      int    `json:"pid"`
	PC       int    `json:"pc"`
	Filename string `json:"filename"` //filename
	Tamanio  int    `json:"tamanio_memoria"`
}

type MensajeExit struct {
	PID int    `json:"pid"`
	PC  int    `json:"pc"`
	ID  string `json:"id"`
}

// Una instruccion es una funcion que recibe un puntero a una struct con el contexto de ejecucion del proceso que se esta
// ejecutando (Pc, variables, registros, etc) y una lista de strings que son los argumentos

var InstruccionSet = map[string]Instruccion{
	// INSTRUCCIONES DE CPU
	"NOOP":  noopInstruccion,
	"GOTO":  gotoInstruccion,
	"WRITE": writeMemInstruccion,
	"READ":  readMemInstruccion,
	// SYSCALLS
	"DUMP_MEMORY": dumpMemoryInstruccion,
	"IO":          ioInstruccion,
	"EXIT":        exitInstruccion,
	"INIT_PROC":   iniciarProcesoInstruccion,
}

// INSTRUCCIONES DE CPU
func noopInstruccion(arguments []string) error {
	if err := checkArguments(arguments, 0); err != nil {
		logger.Error("Error en los argumentos de la instrucción: %s", err)
		return err
	}
	logger.Info("## PID: %d - EJECUTADA: NOOP - %s", globals.PIDActual, strings.Join(arguments, " "))
	return nil
}

func gotoInstruccion(arguments []string) error {
	if err := checkArguments(arguments, 1); err != nil {
		logger.Error("Error en los argumentos de la instrucción: %s", err)
		return err
	}

	nuevoPC, err := strconv.Atoi(arguments[0])
	if err != nil {
		logger.Error("Error al convertir el valor de PC en la instrucción GOTO: %s", err)
		return err
	}

	globals.PCActual = nuevoPC
	globals.SaltarIncrementoPC = true

	logger.Info("## PID: %d - EJECUTADA: GOTO - %s", globals.PIDActual, strings.Join(arguments, " "))
	return nil
}

func writeMemInstruccion(arguments []string) error {
	if err := checkArguments(arguments, 2); err != nil {
		logger.Error("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	dirLogica, err := strconv.Atoi(arguments[0])
	if err != nil {
		logger.Error("Error al convertir la direccion logica: %s", err)
		return err
	}
	datos := arguments[1]

	if traducciones.Cache.EstaActiva() {
		if err := traducciones.EscribirEnCache(datos, dirLogica); err != nil {
			logger.Error("Error escribiendo en Cache: %s", err)
			return err
		}
		logger.Info("## PID: %d - EJECUTADA: WRITE - %s", globals.PIDActual, strings.Join(arguments, " "))
	} else {
		logger.Info("## Cache Inactiva")

		dirFisica := traducciones.Traducir(dirLogica)

		if err := traducciones.EscribirEnMemoria(dirFisica, datos); err != nil {
			logger.Error("Error escribiendo en Memoria: %s", err)
			return err
		}
		logger.Info("## PID: %d - EJECUTADA: WRITE - %s", globals.PIDActual, strings.Join(arguments, " "))
	}
	return nil
}

func readMemInstruccion(arguments []string) error {
	if err := checkArguments(arguments, 2); err != nil {
		logger.Error("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	dirLogica, err := strconv.Atoi(arguments[0])
	if err != nil {
		logger.Error("Error al convertir la direccion logica: %s", err)
		return err
	}
	tamanio, err := strconv.Atoi(arguments[1])
	if err != nil {
		logger.Error("Error al convertir el tamanio: %s", err)
		return err
	}
	nroPagina := dirLogica / globals.TamanioPagina
	dirFisica := traducciones.Traducir(dirLogica)

	if traducciones.Cache.EstaActiva() {
		//logger.Info("Cache Activa")

		valorLeido, hit := traducciones.Cache.Buscar(nroPagina)
		if hit {
			//logger.Info("PID: %d - CACHE HIT - Página: %d", globals.PIDActual, nroPagina)
			return nil
		}

		//logger.Info("PID: %d - CACHE MISS - Página: %d", globals.PIDActual, nroPagina)

		valorLeido, err := traducciones.LeerEnMemoria(dirFisica, tamanio)
		if err != nil {
			logger.Error("Error leyendo de memoria: %v", err)
			return err
		}

		traducciones.Cache.Agregar(dirLogica, valorLeido, false)
		logger.Info("PID: %d - CACHE ADD - Página: %d", globals.PIDActual, nroPagina)
		logger.Info("PID: %d - Accion: LEER - Dirección Física: %d - Valor: %s", globals.PIDActual, dirFisica, valorLeido)
	} else {
		//logger.Info("Cache Inactiva")

		valorLeido, err := traducciones.LeerEnMemoria(dirFisica, tamanio)
		if err != nil {
			logger.Error("Error leyendo de memoria: %v", err)
			return err
		}
		logger.Info("PID: %d - Accion: LEER - Dirección Física: %d - Valor: %s", globals.PIDActual, dirFisica, valorLeido)
	}
	logger.Info("## PID: %d - EJECUTADA: READ - %s", globals.PIDActual, strings.Join(arguments, " "))
	return nil
}

// SYSCALLS
func dumpMemoryInstruccion(arguments []string) error {
	if err := checkArguments(arguments, 0); err != nil {
		logger.Info("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	mensaje := MensajeDump{
		PID: globals.PIDActual,
		PC:  globals.PCActual,
	}

	url := fmt.Sprintf("http://%s:%d/kernel/dump_memory", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	if err := data.EnviarDatos(url, mensaje); err != nil {
		logger.Info("Error al hacer syscall de DUMP_MEMORY a Kernel: %s", err)
		return err
	}

	logger.Info("## PID: %d - EJECUTADA: DUMP_MEMORY - %s", globals.PIDActual, strings.Join(arguments, " "))
	return nil
}

func ioInstruccion(arguments []string) error {
	if err := checkArguments(arguments, 2); err != nil {
		logger.Error("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	nombreIO := arguments[0]
	tiempoIO, err := strconv.Atoi(arguments[1])
	globals.PCActual++
	if err != nil {
		logger.Error("Error al convertir el tiempo de IO: %s", err)
		return err
	}

	mensaje := MensajeIO{
		PID:    globals.PIDActual,
		PC:     globals.PCActual,
		Tiempo: tiempoIO,
		Nombre: nombreIO,
		ID:     globals.ID,
	}

	url := fmt.Sprintf("http://%s:%d/kernel/syscallIO", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	if err := data.EnviarDatos(url, mensaje); err != nil {
		logger.Error("Error al hacer syscall IO a Kernel: %s", err)
		return err
	}

	traducciones.Cache.LimpiarCache()
	traducciones.Tlb.Limpiar()

	logger.Info("## PID: %d - EJECUTADA: IO - %s", globals.PIDActual, strings.Join(arguments, " "))
	return globals.ErrSyscallBloqueante
}

func exitInstruccion(arguments []string) error {
	globals.MutexPID.Lock()
	defer globals.MutexPID.Unlock()

	if err := checkArguments(arguments, 0); err != nil {
		logger.Info("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	mensaje := MensajeExit{
		PID: globals.PIDActual,
		PC:  globals.PCActual,
		ID:  globals.ID,
	}

	url := fmt.Sprintf("http://%s:%d/kernel/exit", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	if err := data.EnviarDatos(url, mensaje); err != nil {
		logger.Error("Error al hacer syscall EXIT a Kernel: %s", err)
		return err
	}

	traducciones.Cache.LimpiarCache()
	traducciones.Tlb.Limpiar()

	logger.Info("## PID: %d - EJECUTADA: EXIT - %s", globals.PIDActual, strings.Join(arguments, " "))
	globals.PIDActual = -1
	return nil
}

func iniciarProcesoInstruccion(arguments []string) error {
	if err := checkArguments(arguments, 2); err != nil {
		logger.Error("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	filename := arguments[0]
	tamanio, err := strconv.Atoi(arguments[1])
	if err != nil {
		logger.Error("Error al convertir el tamaño del proceso: %s", err)
		return err
	}

	mensaje := MensajeInitProc{
		PID:      globals.PIDActual,
		PC:       globals.PCActual,
		Filename: filename,
		Tamanio:  tamanio,
	}

	url := fmt.Sprintf("http://%s:%d/kernel/init_proceso", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	if err := data.EnviarDatos(url, mensaje); err != nil {
		logger.Info("Error al hacer syscall INIT_PROC a Kernel: %s", err)
		return err
	}

	logger.Info("## PID: %d - EJECUTADA: INIT_PROC - %s", globals.PIDActual, strings.Join(arguments, " "))
	return nil
}

func checkArguments(args []string, correctNumberOfArgs int) error {
	if len(args) != correctNumberOfArgs {
		return errors.New("se recibió una cantidad de argumentos no válida")
	}
	return nil
}
