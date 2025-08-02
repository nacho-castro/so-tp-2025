package main

import (
	"bufio"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"net/http"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/kernel/planificadores"
	"github.com/sisoputnfrba/tp-golang/kernel/syscalls"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func main() {
	// ----------------------------------------------------
	// ---------- PARTE CARGA DE PARAMETROS ---------------
	// ----------------------------------------------------
	if len(os.Args) < 4 {
		fmt.Println("Faltan parámetros: archivo_pseudocodigo tamaño_proceso")
		os.Exit(1)
	}

	archivoPseudocodigo := os.Args[1]
	tamanioStr := os.Args[2]
	config := os.Args[3]

	tamanioProceso, err := strconv.Atoi(tamanioStr)
	if err != nil {
		fmt.Printf("Tamaño del proceso inválido: %s\n", tamanioStr)
		os.Exit(1)
	}

	err = logger.ConfigureLogger("kernel.log", "INFO")
	if err != nil {
		fmt.Println("No se pudo crear el logger -", err)
		os.Exit(1)
	}
	logger.Debug("Logger creado")

	logger.Info("Comenzó la ejecución del Kernel")

	// Cargo config directamente sin pasarla por parámetro
	globals.KConfig = globals.CargarConfig(config)

	err = logger.SetLevel(globals.KConfig.LogLevel)
	if err != nil {
		logger.Fatal("No se pudo establecer el nivel de log: %v", err)
	}

	//Inicilizar todas las colas vacias, tipo de dato punteros a PCB y TCB(hilos)
	algoritmos.ColaNuevo = algoritmos.Cola[*pcb.PCB]{}
	algoritmos.ColaBloqueado = algoritmos.Cola[*pcb.PCB]{}
	algoritmos.ColaSalida = algoritmos.Cola[*pcb.PCB]{}
	algoritmos.ColaEjecutando = algoritmos.Cola[*pcb.PCB]{}
	algoritmos.ColaReady = algoritmos.Cola[*pcb.PCB]{}
	algoritmos.ColaBloqueadoSuspendido = algoritmos.Cola[*pcb.PCB]{}
	algoritmos.ColaSuspendidoReady = algoritmos.Cola[*pcb.PCB]{}
	algoritmos.PedidosIO = algoritmos.Cola[*algoritmos.PedidoIO]{}

	// Inicializar recursos compartidos
	Utils.InicializarMutexes()
	Utils.InicializarCanales()

	// ----------------------------------------------------
	// ---------- ENVIAR PSEUDOCODIGO A MEMORIA -----------
	// ----------------------------------------------------
	//1. Crear primer proceso desde los argumentos del main
	planificadores.CrearPrimerProceso(archivoPseudocodigo, tamanioProceso)

	// ----------------------------------------------------
	// ---------- INICAR CORTO Y MEDIANO PLAZO ------------
	// ----------------------------------------------------
	go planificadores.PlanificarCortoPlazo()
	planificadores.PlanificadorMedianoPlazo()

	// ------------------------------------------------------
	// ---------- ESCUCHO REQUESTS DE CPU E IO (Puertos) ----
	// ------------------------------------------------------
	mux := http.NewServeMux()
	mux.HandleFunc("/kernel/io", comunicacion.RecibirMensajeDeIO)
	mux.HandleFunc("/kernel/cpu", comunicacion.RecibirMensajeDeCPU)
	mux.HandleFunc("/kernel/fin_io", comunicacion.RecibirFinDeIO)
	mux.HandleFunc("/kernel/desconexion_io", comunicacion.RecibirFinDeIO)

	// ------------------------------------------------------
	// --------------------- SYSCALLS -----------------------
	// ------------------------------------------------------
	mux.HandleFunc("/kernel/contexto_interrumpido", syscalls.ContextoInterrumpido)
	mux.HandleFunc("/kernel/init_proceso", syscalls.InitProcess)
	mux.HandleFunc("/kernel/exit", syscalls.Exit)
	mux.HandleFunc("/kernel/dump_memory", syscalls.DumpMemory)
	mux.HandleFunc("/kernel/syscallIO", syscalls.Io)
	//para recibir mensajes de confirmacion de SWAP
	mux.HandleFunc("/kernel/suspension_completed", syscalls.ConfirmarSuspensionHandler)

	fmt.Printf("Servidor escuchando en http://localhost:%d/kernel\n", globals.KConfig.KernelPort)

	// ------------------------------------------------------
	// ---------- INICIAR PLANIFICADOR DE LARGO PLAZO  ------
	// ------------------------------------------------------
	// Esperar que el usuario presione Enter
	go iniciarLargoPlazo()

	address := fmt.Sprintf(":%d", globals.KConfig.KernelPort)
	err = http.ListenAndServe(address, mux)
	if err != nil {
		panic(err)
	}

	fmt.Printf("FIN DE EJECUCION")
}

func iniciarLargoPlazo() {
	fmt.Println("Presione ENTER para iniciar el Planificador de Largo Plazo...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	planificadores.PlanificadorLargoPlazo()
}
