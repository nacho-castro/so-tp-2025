package tests

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/traducciones"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func CargarConfigCPU(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var config globals.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}
	globals.ClientConfig = &config
	return nil
}

func TestTraducirDireccionTLB(t *testing.T) {
	err := CargarConfigCPU("../configs/CP.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}

	globals.TamanioPagina = 256
	globals.EntradasPorNivel = 4
	globals.CantidadNiveles = 2

	dirLogica := 4 * globals.TamanioPagina

	//Inicializo y cargo la TLB en mmu.go con la pag 4 y marco 7

	//Traduzco
	dirFisica := traducciones.Traducir(dirLogica)

	dirEsperada := 7*globals.TamanioPagina + 0
	assert.Equal(t, dirEsperada, dirFisica)
}

func TestDescomponerPagina(t *testing.T) {
	err := CargarConfigCPU("../configs/CP.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)

	}

	entradas := traducciones.DescomponerPagina(5, 3, 4)

	entradasEsperadas := []int{0, 1, 1}

	assert.Equal(t, entradasEsperadas, entradas)
}

func TestTraduccionMemoria(t *testing.T) {
	//TODO: le pide a memoria el marco
}

func TestLeerEnMemoria(t *testing.T) {
	//TODO
}

func TestEscribirEnMemoria(t *testing.T) {
	//TODO
}
