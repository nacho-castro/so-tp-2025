## KERNEL 

## FUNCIONALIDAD (SERVER)

0.  Comando de ejecución inicial:

`➜ ~ ./bin/kernel [archivo_pseudocodigo] [tamanio_proceso] [...args]`

1. Leer y cargar en Globals los datos del Archivo de Configuracion

```go
globals.KernelConfig = utils.Config(filepath)
```

2. LISTEN en los puertos HTTP para recibir PUERTOS e IP del IO o CPU

```go
mux.HandleFunc("/kernel/io", utils.RecibirMensajeDeIO)
mux.HandleFunc("/kernel/cpu", utils.RecibirMensajeDeCPU)
```

3. ESCRIBIR EN GLOBALS la IP y PUERTO recibidos por los modulos

```go
	globals.IO = globals.DatosIO{
		Nombre: mensaje.Nombre,
		Ip:     mensaje.Ip,
		Puerto: mensaje.Puerto,
	}

	globals.CPU = globals.DatosCPU{
		Ip:     mensajeRecibido.Ip,
		Puerto: mensajeRecibido.Puerto,
		ID:     mensajeRecibido.ID,
	}
```

4. LISTEN en los puertos HTTP para peticiones de IO o CPU

## 🔌 1. Endpoint expuesto

El Kernel escucha conexiones entrantes desde otros módulos en:

`http://localhost:8081/kernel/io`
`http://localhost:8081/kernel/cpu`

## 📬 2. Formato del mensaje recibido

El cuerpo del mensaje (`body`) debe ser un JSON con una estructura dependiendo de cada Modulo:

CPU:
```json
{
  "id":"1",
  "ip": "127.0.0.1",
  "puerto": 8000
}
```
IO:
```json
{
  "nombre":"impresora",
  "ip": "127.0.0.1",
  "puerto": 8000
}
```

Estos mensajes se decodifican en un struct de Go como los siguientes:
```go
package globals

type Config struct {
	IpMemory           string `json:"ip_memory"`
	PortMemory         int    `json:"port_memory"`
	PortKernel         int    `json:"port_kernel"`
	SchedulerAlgorithm string `json:"scheduler_algorithm"`
	NewAlgorithm       string `json:"new_algorithm"`
	Alpha              string `json:"alpha"`
	SuspensionTime     int    `json:"suspension_time"`
	LogLevel           string `json:"log_level"`
}

// Datos recibidos
type DatosIO struct {
	Nombre string
	Ip     string
	Puerto int
}

type DatosCPU struct {
	Ip     string
	Puerto int
	ID     string
}
```

## 3. Estructura

kernel/ 
├── utils/ # Funciones auxiliares (leer JSON, manejar requests) 
	│ 
	└── utils.go 
├── globals/ 
	│ 
	└── globals.go 
├── config.json # Archivo de configuración 
├── go.mod # Módulo Go 
├── kernel.go # Lógica del módulo Kernel 
└── README.md # Documentación del proyecto
