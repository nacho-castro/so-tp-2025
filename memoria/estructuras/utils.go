package estructuras

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"time"
)

func CalcularEjecutarSleep(tiempo time.Duration) {
	time.Sleep(tiempo)
}

func CalcularCantidadEntradas(tamanio int) (resultado int, err error) {
	err = nil
	resultado = 0
	if tamanio < 0 {
		return resultado, fmt.Errorf("el tamanio pedido de espacio es negativo %v", logger.ErrBadRequest)
	}

	resultado = tamanio / MemoryConfig.PagSize
	if tamanio%MemoryConfig.PagSize > 0 {
		resultado++
	}
	return
}

func RecortarNulosFinales(data []byte) []byte {
	i := len(data) - 1
	for i >= 0 && data[i] == 0 {
		i--
	}
	return data[:i+1]
}
