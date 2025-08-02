package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func ObtenerInstruccion(proceso *g.Proceso, pc int) (respuesta g.RespuestaInstruccion, err error) {
	respuesta = g.RespuestaInstruccion{Exito: nil, Instruccion: ""}

	if proceso == nil {
		logger.Error("Proceso recibido es nil")
		return respuesta, logger.ErrProcessNil
	}
	respuesta.Instruccion = string(proceso.InstruccionesEnBytes[pc])
	return respuesta, nil
}
