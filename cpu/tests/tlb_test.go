package tests

import (
	"github.com/sisoputnfrba/tp-golang/cpu/traducciones"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTLBHit(t *testing.T) {
	err := CargarConfigCPU("../configs/CP.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}

	tlb := traducciones.NuevaTLB()
	tlb.AgregarEntrada(10, 24)

	marco, ok := tlb.Buscar(10)

	assert.True(t, ok)
	assert.Equal(t, 24, marco)
}

func TestTLBMiss(t *testing.T) {
	err := CargarConfigCPU("../configs/CP.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}

	tlb := traducciones.NuevaTLB()

	marco, ok := tlb.Buscar(11)

	assert.False(t, ok)
	assert.Equal(t, -1, marco)
}

func TestReemplazoPorFifo(t *testing.T) {
	err := CargarConfigCPU("../configs/CP.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}

	tlb := traducciones.NuevaTLB()
	tlb.AgregarEntrada(1, 1)
	tlb.AgregarEntrada(2, 2)
	tlb.AgregarEntrada(3, 3)
	tlb.AgregarEntrada(4, 4)
	tlb.AgregarEntrada(5, 5)

	//Agrego entrada y por FIFO se elimina (1,1)
	tlb.AgregarEntrada(6, 6)

	marco, ok := tlb.Buscar(1)

	assert.False(t, ok)
	assert.Equal(t, -1, marco)
}

func TestReemplazoPorLRU(t *testing.T) {
	err := CargarConfigCPU("../configs/PLANI.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}

	tlb := traducciones.NuevaTLB()
	tlb.AgregarEntrada(1, 1)
	tlb.AgregarEntrada(2, 2)
	tlb.AgregarEntrada(3, 3)
	tlb.AgregarEntrada(4, 4)
	tlb.AgregarEntrada(5, 5)

	tlb.Buscar(1) //Por LRU se actualiza

	tlb.AgregarEntrada(6, 6) //Remplaza la entrada 2,2

	marco, ok := tlb.Buscar(2)

	assert.False(t, ok)
	assert.Equal(t, -1, marco)
}
